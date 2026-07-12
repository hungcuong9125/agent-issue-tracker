package ait

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// gitignoreEntry is the pattern ait adds to a project's .gitignore so the
// local issue database in .ait/ is not committed.
const gitignoreEntry = ".ait/"

// ensureGitignoreForDB adds .ait/ to the project's .gitignore. init calls it
// on every run, so a lost entry is restored. It is deliberately conservative:
// it only acts when the database is at the default location (.ait/ait.db)
// inside a git repository, leaving custom --db paths and non-git directories
// untouched. It reports whether the entry was added and — so init can surface
// the gap in its output — whether it was skipped because there is no .git
// directory.
func ensureGitignoreForDB(dbPath string) (updated bool, noGit bool, err error) {
	root, err := ProjectRoot()
	if err != nil {
		return false, false, err
	}

	// Only manage .gitignore for the default database location.
	if dbPath != filepath.Join(root, ".ait", "ait.db") {
		return false, false, nil
	}

	// Nothing to ignore if this isn't a git repository.
	if _, err := os.Stat(filepath.Join(root, ".git")); err != nil {
		return false, true, nil
	}

	updated, err = ensureGitignore(root)
	return updated, false, err
}

// ensureGitignore makes sure the .gitignore at root ignores ait's data
// directory, creating the file if it does not exist. It returns true if the
// file was modified. An entry already present in any common form (".ait",
// ".ait/", "/.ait", "/.ait/") is left untouched, so a user's own edits are
// never overwritten.
func ensureGitignore(root string) (bool, error) {
	path := filepath.Join(root, ".gitignore")

	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return false, err
	}

	if gitignoreHasEntry(string(data)) {
		return false, nil
	}

	var b strings.Builder
	b.Write(data)
	// Keep a clean separation from any existing trailing content.
	if len(data) > 0 && !strings.HasSuffix(string(data), "\n") {
		b.WriteByte('\n')
	}
	b.WriteString(gitignoreEntry + "\n")

	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		return false, err
	}
	return true, nil
}

// gitignoreHasEntry reports whether the .gitignore contents already ignore the
// .ait/ directory, accepting the common spellings of the pattern.
func gitignoreHasEntry(contents string) bool {
	for _, line := range strings.Split(contents, "\n") {
		switch strings.TrimSpace(line) {
		case ".ait", ".ait/", "/.ait", "/.ait/":
			return true
		}
	}
	return false
}
