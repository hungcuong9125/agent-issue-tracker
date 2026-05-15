package ait

import (
	"os"
	"os/exec"
	"strings"
)

const editorCreateHelp = `# Enter issue details. First line is the title, everything after
# a blank line is the description. Lines starting with # are ignored.
# Save and close the editor to create the issue; leave it empty to abort.
`

const editorUpdateHelp = `# Edit issue details. First line is the title, everything after
# a blank line is the description. Lines starting with # are ignored.
# Save and close the editor to apply the changes; leave it empty to abort.
`

// EditIssueMessage opens $EDITOR (falling back to vi) with a template and
// parses the result into a title and description, similar to git commit.
// If initialTitle/initialDescription are non-empty they are pre-populated for
// the user to edit.
func EditIssueMessage(initialTitle, initialDescription string) (title, description string, err error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	tmp, err := os.CreateTemp("", "ait-edit-*.md")
	if err != nil {
		return "", "", &CLIError{Code: "io", Message: "cannot create temp file: " + err.Error(), ExitCode: 73}
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	var body strings.Builder
	if initialTitle == "" && initialDescription == "" {
		body.WriteString("\n")
		body.WriteString(editorCreateHelp)
	} else {
		body.WriteString(initialTitle)
		body.WriteString("\n\n")
		body.WriteString(initialDescription)
		if initialDescription != "" {
			body.WriteString("\n")
		}
		body.WriteString("\n")
		body.WriteString(editorUpdateHelp)
	}

	if _, err := tmp.WriteString(body.String()); err != nil {
		tmp.Close()
		return "", "", &CLIError{Code: "io", Message: "cannot write temp file: " + err.Error(), ExitCode: 73}
	}
	tmp.Close()

	cmd := exec.Command(editor, tmpPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", "", &CLIError{Code: "io", Message: "editor exited with error: " + err.Error(), ExitCode: 1}
	}

	data, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", "", &CLIError{Code: "io", Message: "cannot read temp file: " + err.Error(), ExitCode: 73}
	}

	return parseEditorMessage(string(data))
}

// parseEditorMessage extracts a title and description from editor output.
// The first non-comment, non-blank line is the title. Everything after the
// next blank line (skipping comments) is the description.
func parseEditorMessage(raw string) (title, description string, err error) {
	var titleLine string
	var descLines []string
	pastTitle := false
	pastBlank := false

	for _, line := range strings.Split(raw, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}

		if !pastTitle {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				titleLine = trimmed
				pastTitle = true
			}
			continue
		}

		if !pastBlank {
			if strings.TrimSpace(line) == "" {
				pastBlank = true
			}
			continue
		}

		descLines = append(descLines, line)
	}

	if titleLine == "" {
		return "", "", &CLIError{Code: "validation", Message: "aborted: empty title", ExitCode: 1}
	}

	return titleLine, strings.TrimSpace(strings.Join(descLines, "\n")), nil
}
