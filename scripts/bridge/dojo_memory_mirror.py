#!/usr/bin/env python3
"""dojo_memory_mirror.py - markdown memory -> CSV for the Postgres Memory Hub.

Walks the top-level *.md memory files, parses frontmatter, and emits a CSV
(id,name,slug,type,description,body,updated,origin_session,links) suitable for
COPY into the `memories` table. stdlib only; runs on Mac (proof-load) or bridge
(ongoing cron sync) identically.

Usage: dojo_memory_mirror.py <src_dir> <out_csv>
"""
import csv, os, re, sys, datetime

# These are indexes / navigation / contracts, not atomic memories -> skip.
SKIP_EXACT = {"MEMORY.md", "README.md", "DESIGN.md", "OPERATIONS.md", "CLAUDE.md"}
SKIP_PREFIX = ("MEMORY-",)

FM_RE = re.compile(r"^---\s*\n(.*?)\n---\s*\n", re.DOTALL)
LINK_RE = re.compile(r"\[\[([^\]]+)\]\]")
DATE_RE = re.compile(r"(\d{4}-\d{2}-\d{2})")


def parse_frontmatter(text):
    m = FM_RE.match(text)
    fm, body = {}, text
    if not m:
        return fm, body
    block, body = m.group(1), text[m.end():]
    cur_key, meta = None, {}
    for line in block.splitlines():
        if not line.strip():
            continue
        if cur_key == "metadata" and re.match(r"^\s+\S", line):
            mm = re.match(r"^\s+([\w-]+)\s*:\s*(.*)$", line)
            if mm:
                meta[mm.group(1).strip()] = mm.group(2).strip().strip('"\'')
            continue
        km = re.match(r"^([\w-]+)\s*:\s*(.*)$", line)
        if km:
            cur_key = km.group(1).strip()
            fm[cur_key] = km.group(2).strip().strip('"\'')
    if meta:
        fm["_meta"] = meta
    return fm, body


def get_type(fm):
    if fm.get("type"):
        return fm["type"]
    meta = fm.get("_meta") or {}
    return meta.get("type") or "unknown"


def main():
    if len(sys.argv) != 3:
        print("usage: dojo_memory_mirror.py <src_dir> <out_csv>", file=sys.stderr)
        sys.exit(2)
    src, out = sys.argv[1], sys.argv[2]
    rows, seen = [], 0
    for fn in sorted(os.listdir(src)):
        if not fn.endswith(".md") or fn in SKIP_EXACT or fn.startswith(SKIP_PREFIX):
            continue
        path = os.path.join(src, fn)
        if not os.path.isfile(path):
            continue
        seen += 1
        text = open(path, encoding="utf-8", errors="replace").read()
        fm, body = parse_frontmatter(text)
        stem = fn[:-3]
        name = fm.get("name") or stem.replace("_", " ").replace("-", " ")
        typ = get_type(fm)
        desc = fm.get("description", "")
        upd = None
        for k in ("updated", "date", "created"):
            if fm.get(k):
                dm = DATE_RE.search(fm[k])
                if dm:
                    upd = dm.group(1) + "T00:00:00Z"
                    break
        if not upd:
            upd = datetime.datetime.fromtimestamp(os.path.getmtime(path), datetime.timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
        links = sorted(set(LINK_RE.findall(body)))
        links_lit = "{" + ",".join('"%s"' % l.replace('"', '').replace('\\', '') for l in links) + "}"
        # id == slug == filename stem (guaranteed unique, satisfies UNIQUE(slug))
        rows.append([stem, name, stem, typ, desc, body, upd, fm.get("origin_session", ""), links_lit])
    with open(out, "w", newline="", encoding="utf-8") as f:
        w = csv.writer(f)
        w.writerow(["id", "name", "slug", "type", "description", "body", "updated", "origin_session", "links"])
        w.writerows(rows)
    print("files_seen=%d rows_written=%d" % (seen, len(rows)))


if __name__ == "__main__":
    main()
