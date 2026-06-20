// Package claude provides the Claude Code workflow adapter.
package claude

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/plexusone/assistantkit/workflows/core"
)

func init() {
	core.Register(&Adapter{})
}

// Adapter converts between canonical Workflow and Claude Code format.
type Adapter struct{}

// Name returns the adapter identifier.
func (a *Adapter) Name() string {
	return "claude"
}

// EntryPointFile returns the entry point filename for Claude Code.
func (a *Adapter) EntryPointFile() string {
	return "CLAUDE.md"
}

// EntryPointDir returns the directory for the entry point.
// Empty string means project root.
func (a *Adapter) EntryPointDir() string {
	return ""
}

// RuleDetailsDir returns the directory for rule details.
func (a *Adapter) RuleDetailsDir() string {
	return ".spec-workflows"
}

// Generate creates Claude Code workflow files from a Workflow.
func (a *Adapter) Generate(workflow *core.Workflow, outputDir string) error {
	// Generate entry point (CLAUDE.md)
	entryPointContent, err := a.GenerateEntryPoint(workflow)
	if err != nil {
		return err
	}

	entryPointPath := filepath.Join(outputDir, a.EntryPointFile())
	if err := os.WriteFile(entryPointPath, entryPointContent, core.DefaultFileMode); err != nil {
		return &core.WriteError{Path: entryPointPath, Err: err}
	}

	// Copy rule details
	ruleDetailsDir := filepath.Join(outputDir, a.RuleDetailsDir())
	if workflow.RuleDetails != "" {
		if err := copyDir(workflow.RuleDetails, filepath.Join(ruleDetailsDir, "rule-details")); err != nil {
			return &core.GenerateError{Platform: "claude", Path: ruleDetailsDir, Err: err}
		}
	}

	// Copy templates
	if workflow.Templates != "" {
		if err := copyDir(workflow.Templates, filepath.Join(ruleDetailsDir, "templates")); err != nil {
			return &core.GenerateError{Platform: "claude", Path: ruleDetailsDir, Err: err}
		}
	}

	// Copy rubrics
	if workflow.Rubrics != "" {
		if err := copyDir(workflow.Rubrics, filepath.Join(ruleDetailsDir, "rubrics")); err != nil {
			return &core.GenerateError{Platform: "claude", Path: ruleDetailsDir, Err: err}
		}
	}

	return nil
}

// GenerateEntryPoint generates the CLAUDE.md content.
func (a *Adapter) GenerateEntryPoint(workflow *core.Workflow) ([]byte, error) {
	var buf bytes.Buffer

	// Write header
	buf.WriteString(fmt.Sprintf("# VisionSpec: %s (%s Level)\n\n", workflow.Name, workflow.Level))

	// Write trigger pattern
	trigger := workflow.Trigger
	if trigger == "" {
		trigger = "Using VisionSpec,"
	}
	buf.WriteString(fmt.Sprintf("Activate this workflow when user says: **\"%s [intent]\"**\n\n", trigger))

	// Write rule details location
	buf.WriteString("## Rule Details Location\n\n")
	buf.WriteString("Check these paths in order and use the first one that exists:\n\n")
	buf.WriteString("1. `.spec-workflows/rule-details/`\n")
	buf.WriteString("2. `.visionspec/rule-details/`\n\n")

	// Read and include the core workflow content
	if workflow.EntryPoint != "" {
		content, err := os.ReadFile(workflow.EntryPoint)
		if err != nil {
			return nil, &core.ParseError{Path: workflow.EntryPoint, Err: err}
		}

		// Skip the header from the source file (we wrote our own)
		// Find the first "## " section and include from there
		lines := bytes.Split(content, []byte("\n"))
		startIdx := 0
		for i, line := range lines {
			if bytes.HasPrefix(line, []byte("## ")) {
				startIdx = i
				break
			}
		}

		buf.Write(bytes.Join(lines[startIdx:], []byte("\n")))
	}

	return buf.Bytes(), nil
}

// copyDir recursively copies a directory.
func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
