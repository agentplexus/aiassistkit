package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewWorkflow(t *testing.T) {
	wf := NewWorkflow("aws-working-backwards", "product")

	if wf.Name != "aws-working-backwards" {
		t.Errorf("expected name 'aws-working-backwards', got '%s'", wf.Name)
	}
	if wf.Level != "product" {
		t.Errorf("expected level 'product', got '%s'", wf.Level)
	}
	if wf.Trigger != "Using VisionSpec," {
		t.Errorf("expected trigger 'Using VisionSpec,', got '%s'", wf.Trigger)
	}
}

func TestLoadSourceRepo(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()

	// Create workflows directory structure
	workflowDir := filepath.Join(tmpDir, "workflows", "test-workflow", "product")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		t.Fatalf("failed to create workflow dir: %v", err)
	}

	// Create core-workflow.md
	workflowFile := filepath.Join(workflowDir, "core-workflow.md")
	content := []byte("# Test Workflow\n\n## Phase 1\n\nDo something.")
	if err := os.WriteFile(workflowFile, content, 0644); err != nil {
		t.Fatalf("failed to write workflow file: %v", err)
	}

	// Create rule-details directory
	ruleDetailsDir := filepath.Join(tmpDir, "rule-details")
	if err := os.MkdirAll(ruleDetailsDir, 0755); err != nil {
		t.Fatalf("failed to create rule-details dir: %v", err)
	}

	// Create templates directory
	templatesDir := filepath.Join(tmpDir, "templates", "test-workflow")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		t.Fatalf("failed to create templates dir: %v", err)
	}

	// Create rubrics directory
	rubricsDir := filepath.Join(tmpDir, "rubrics", "test-workflow")
	if err := os.MkdirAll(rubricsDir, 0755); err != nil {
		t.Fatalf("failed to create rubrics dir: %v", err)
	}

	// Load the repo
	repo, err := LoadSourceRepo(tmpDir)
	if err != nil {
		t.Fatalf("failed to load source repo: %v", err)
	}

	// Verify workflows
	workflows := repo.ListWorkflows()
	if len(workflows) != 1 {
		t.Errorf("expected 1 workflow, got %d", len(workflows))
	}

	// Get the workflow
	wf, err := repo.GetWorkflow("test-workflow/product")
	if err != nil {
		t.Fatalf("failed to get workflow: %v", err)
	}

	if wf.Name != "test-workflow" {
		t.Errorf("expected name 'test-workflow', got '%s'", wf.Name)
	}
	if wf.Level != "product" {
		t.Errorf("expected level 'product', got '%s'", wf.Level)
	}
	if wf.EntryPoint != workflowFile {
		t.Errorf("expected entry point '%s', got '%s'", workflowFile, wf.EntryPoint)
	}
}

func TestLoadSourceRepo_NotFound(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := LoadSourceRepo(filepath.Join(tmpDir, "nonexistent"))
	if err == nil {
		t.Error("expected error for nonexistent repo")
	}
}

func TestSourceRepo_GetWorkflow_NotFound(t *testing.T) {
	repo := &SourceRepo{
		Workflows: make(map[string]*Workflow),
	}

	_, err := repo.GetWorkflow("nonexistent/product")
	if err == nil {
		t.Error("expected error for nonexistent workflow")
	}

	if _, ok := err.(*NotFoundError); !ok {
		t.Errorf("expected NotFoundError, got %T", err)
	}
}
