package ait

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureGitignoreCreatesFile(t *testing.T) {
	root := t.TempDir()

	changed, err := ensureGitignore(root)
	if err != nil {
		t.Fatalf("ensureGitignore failed: %v", err)
	}
	if !changed {
		t.Fatal("expected ensureGitignore to report a change")
	}

	got := readFile(t, filepath.Join(root, ".gitignore"))
	if got != ".ait/\n" {
		t.Fatalf("unexpected .gitignore contents: %q", got)
	}
}

func TestEnsureGitignoreAppendsToExisting(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".gitignore")
	if err := os.WriteFile(path, []byte("node_modules/\n*.log\n"), 0o644); err != nil {
		t.Fatalf("seed .gitignore: %v", err)
	}

	changed, err := ensureGitignore(root)
	if err != nil {
		t.Fatalf("ensureGitignore failed: %v", err)
	}
	if !changed {
		t.Fatal("expected ensureGitignore to report a change")
	}

	got := readFile(t, path)
	want := "node_modules/\n*.log\n.ait/\n"
	if got != want {
		t.Fatalf("unexpected .gitignore contents:\n got %q\nwant %q", got, want)
	}
}

func TestEnsureGitignoreAddsTrailingNewlineBeforeAppending(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".gitignore")
	// No trailing newline on the existing content.
	if err := os.WriteFile(path, []byte("*.log"), 0o644); err != nil {
		t.Fatalf("seed .gitignore: %v", err)
	}

	if _, err := ensureGitignore(root); err != nil {
		t.Fatalf("ensureGitignore failed: %v", err)
	}

	got := readFile(t, path)
	want := "*.log\n.ait/\n"
	if got != want {
		t.Fatalf("unexpected .gitignore contents:\n got %q\nwant %q", got, want)
	}
}

func TestEnsureGitignoreIsIdempotent(t *testing.T) {
	variants := []string{".ait", ".ait/", "/.ait", "/.ait/"}
	for _, entry := range variants {
		t.Run(entry, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, ".gitignore")
			original := "build/\n" + entry + "\n"
			if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
				t.Fatalf("seed .gitignore: %v", err)
			}

			changed, err := ensureGitignore(root)
			if err != nil {
				t.Fatalf("ensureGitignore failed: %v", err)
			}
			if changed {
				t.Fatalf("expected no change when %q is already present", entry)
			}

			if got := readFile(t, path); got != original {
				t.Fatalf("file was modified:\n got %q\nwant %q", got, original)
			}
		})
	}
}

func TestGitignoreHasEntryIgnoresComments(t *testing.T) {
	if gitignoreHasEntry("# .ait/ is where the db lives\nbuild/\n") {
		t.Fatal("commented-out mention should not count as an entry")
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(data)
}
