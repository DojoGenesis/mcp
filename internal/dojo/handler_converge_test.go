package dojo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// isolateGit keeps the developer's git setup out of these tests. A global
// core.hooksPath would otherwise fire on the fixture commits, and a repository
// above the temp directory would turn the "not a repository" fixtures into real ones.
func isolateGit(t *testing.T) {
	t.Helper()
	t.Setenv("GIT_CONFIG_GLOBAL", filepath.Join(t.TempDir(), "no-global-config"))
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_CEILING_DIRECTORIES", os.TempDir())
}

func gitIn(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-c", "user.name=test", "-c", "user.email=test@example.com"}, args...)...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

// newRepo returns a fresh repository holding one commit and no uncommitted files.
func newRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	gitIn(t, dir, "init", "-q")
	gitIn(t, dir, "commit", "-q", "--allow-empty", "-m", "first")
	return dir
}

func writeUntracked(t *testing.T, dir string, count int) {
	t.Helper()
	for i := 0; i < count; i++ {
		if err := os.WriteFile(filepath.Join(dir, fmt.Sprintf("file%d.txt", i)), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func converge(t *testing.T, root string) *mcp.CallToolResult {
	t.Helper()
	h, err := NewHandler("", t.TempDir(), nil, root)
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}
	return callTool(t, h.handleConverge, nil)
}

func TestConvergeDidNotRunWhenGitFails(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T) string // returns the workspace root
	}{
		{"not a repository", func(t *testing.T) string { return t.TempDir() }},
		{"missing directory", func(t *testing.T) string { return filepath.Join(t.TempDir(), "gone") }},
		// A worktree whose target was deleted: the drift report's "Unreadable repos".
		{"broken worktree pointer", func(t *testing.T) string {
			dir := t.TempDir()
			pointer := "gitdir: " + filepath.Join(dir, "gone") + "\n"
			if err := os.WriteFile(filepath.Join(dir, ".git"), []byte(pointer), 0o644); err != nil {
				t.Fatal(err)
			}
			return dir
		}},
		{"git not on PATH", func(t *testing.T) string {
			dir := newRepo(t)
			t.Setenv("PATH", t.TempDir())
			return dir
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isolateGit(t)
			// Start inside a clean repository: a handler reading the process directory would say GREEN.
			t.Chdir(newRepo(t))
			root := tt.setup(t)

			res := converge(t, root)
			text := resultText(t, res)

			if !res.IsError {
				t.Errorf("IsError = false, want true\n%s", text)
			}
			for _, want := range []string{"# Convergence Check: DID NOT RUN", root, "not a clean result"} {
				if !strings.Contains(text, want) {
					t.Errorf("output missing %q\n%s", want, text)
				}
			}
			for _, bad := range []string{"GREEN", "Dirty files"} {
				if strings.Contains(text, bad) {
					t.Errorf("output must not contain %q\n%s", bad, text)
				}
			}
		})
	}
}

func TestConvergeDidNotRunWithoutWorkspaceRoot(t *testing.T) {
	isolateGit(t)
	t.Chdir(newRepo(t))

	res := converge(t, "")
	text := resultText(t, res)

	if !res.IsError {
		t.Errorf("IsError = false, want true\n%s", text)
	}
	for _, want := range []string{"# Convergence Check: DID NOT RUN", "DOJO_WORKSPACE_ROOT"} {
		if !strings.Contains(text, want) {
			t.Errorf("output missing %q\n%s", want, text)
		}
	}
	if strings.Contains(text, "GREEN") {
		t.Errorf("output must not contain GREEN\n%s", text)
	}
}

func TestConvergeSignalComesFromWorkspaceRoot(t *testing.T) {
	tests := []struct {
		dirty   int
		signal  string
		verdict string
	}{
		{0, "GREEN", "**GREEN** — 0 dirty files in"},
		{9, "GREEN", "**GREEN** — 9 dirty files in"},
		{10, "YELLOW", "**YELLOW** — 10 dirty files."},
		{24, "YELLOW", "**YELLOW** — 24 dirty files."},
		{25, "RED", "**RED** — 25 dirty files."},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d dirty", tt.dirty), func(t *testing.T) {
			isolateGit(t)
			// The process directory is a different repository, with a count no case expects.
			elsewhere := newRepo(t)
			writeUntracked(t, elsewhere, 3)
			t.Chdir(elsewhere)

			root := newRepo(t)
			writeUntracked(t, root, tt.dirty)

			res := converge(t, root)
			text := resultText(t, res)

			if res.IsError {
				t.Errorf("IsError = true, want false\n%s", text)
			}
			for _, want := range []string{
				"# Convergence Check: " + tt.signal,
				tt.verdict,
				fmt.Sprintf("| Dirty files | %d |", tt.dirty),
				"| Commits (last 7 days) | 1 |",
				root,
				"not every repo in the workspace",
			} {
				if !strings.Contains(text, want) {
					t.Errorf("output missing %q\n%s", want, text)
				}
			}
			if strings.Contains(text, "workspace is clean") {
				t.Errorf("one repository must not be reported as a clean workspace\n%s", text)
			}
		})
	}
}

// git log fails in a repository with no commits, while git status works.
func TestConvergeCommitCountUnknownWhenLogFails(t *testing.T) {
	isolateGit(t)
	root := t.TempDir()
	gitIn(t, root, "init", "-q")

	res := converge(t, root)
	text := resultText(t, res)

	if res.IsError {
		t.Errorf("IsError = true, want false\n%s", text)
	}
	for _, want := range []string{"# Convergence Check: GREEN", "| Commits (last 7 days) | unknown |", "`git log` failed"} {
		if !strings.Contains(text, want) {
			t.Errorf("output missing %q\n%s", want, text)
		}
	}
	if strings.Contains(text, "| Commits (last 7 days) | 0 |") {
		t.Errorf("a failed git log must not read as zero commits\n%s", text)
	}
}
