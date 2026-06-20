// Package core provides the canonical workflow definition types.
// Workflows define multi-phase spec creation processes that can be
// deployed across different AI coding assistants.
package core

import (
	"os"
	"path/filepath"
)

// DefaultFileMode is the default permission for workflow files.
const DefaultFileMode = 0644

// DefaultDirMode is the default permission for workflow directories.
const DefaultDirMode = 0755

// Workflow represents a canonical workflow definition.
// Workflows guide AI assistants through multi-phase spec creation.
type Workflow struct {
	// Name is the workflow identifier (e.g., "aws-working-backwards")
	Name string `json:"name" yaml:"name"`

	// Level is the workflow level: "product" or "feature"
	Level string `json:"level" yaml:"level"`

	// Description describes the workflow purpose
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Trigger is the phrase that activates this workflow
	Trigger string `json:"trigger,omitempty" yaml:"trigger,omitempty"`

	// EntryPoint is the path to the core-workflow.md file
	EntryPoint string `json:"entryPoint" yaml:"entryPoint"`

	// RuleDetails is the path to the rule-details directory
	RuleDetails string `json:"ruleDetails,omitempty" yaml:"ruleDetails,omitempty"`

	// Templates is the path to the templates directory
	Templates string `json:"templates,omitempty" yaml:"templates,omitempty"`

	// Rubrics is the path to the rubrics directory
	Rubrics string `json:"rubrics,omitempty" yaml:"rubrics,omitempty"`

	// Phases defines the workflow phases in order
	Phases []Phase `json:"phases,omitempty" yaml:"phases,omitempty"`

	// Extension is an optional extension to apply
	Extension string `json:"extension,omitempty" yaml:"extension,omitempty"`
}

// Phase represents a workflow phase.
type Phase struct {
	// Name is the phase identifier
	Name string `json:"name" yaml:"name"`

	// SpecType is the spec type created in this phase
	SpecType string `json:"specType" yaml:"specType"`

	// Mode is the authoring mode: "interactive", "synthesis", or "reconciliation"
	Mode string `json:"mode" yaml:"mode"`

	// Sources are the spec types used as input (for synthesis mode)
	Sources []string `json:"sources,omitempty" yaml:"sources,omitempty"`

	// Gate indicates if approval is required before proceeding
	Gate bool `json:"gate,omitempty" yaml:"gate,omitempty"`
}

// WorkflowLevel constants.
const (
	LevelProduct = "product"
	LevelFeature = "feature"
)

// AuthoringMode constants.
const (
	ModeInteractive    = "interactive"
	ModeSynthesis      = "synthesis"
	ModeReconciliation = "reconciliation"
)

// NewWorkflow creates a new Workflow with the given name and level.
func NewWorkflow(name, level string) *Workflow {
	return &Workflow{
		Name:    name,
		Level:   level,
		Trigger: "Using VisionSpec,",
	}
}

// SourceRepo represents a spec-workflows repository.
type SourceRepo struct {
	// Path is the local filesystem path to the repo
	Path string `json:"path" yaml:"path"`

	// Workflows maps workflow identifiers to Workflow definitions
	Workflows map[string]*Workflow `json:"workflows,omitempty" yaml:"workflows,omitempty"`
}

// LoadSourceRepo loads a spec-workflows repository from the given path.
func LoadSourceRepo(path string) (*SourceRepo, error) {
	repo := &SourceRepo{
		Path:      path,
		Workflows: make(map[string]*Workflow),
	}

	// Scan workflows directory
	workflowsDir := filepath.Join(path, "workflows")
	entries, err := os.ReadDir(workflowsDir)
	if err != nil {
		return nil, &LoadError{Path: workflowsDir, Err: err}
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		methodologyName := entry.Name()
		methodologyDir := filepath.Join(workflowsDir, methodologyName)

		// Check for product and feature subdirectories
		for _, level := range []string{LevelProduct, LevelFeature} {
			levelDir := filepath.Join(methodologyDir, level)
			entryPoint := filepath.Join(levelDir, "core-workflow.md")

			if _, err := os.Stat(entryPoint); err == nil {
				workflowID := methodologyName + "/" + level
				repo.Workflows[workflowID] = &Workflow{
					Name:        methodologyName,
					Level:       level,
					EntryPoint:  entryPoint,
					RuleDetails: filepath.Join(path, "rule-details"),
					Templates:   filepath.Join(path, "templates", methodologyName),
					Rubrics:     filepath.Join(path, "rubrics", methodologyName),
					Trigger:     "Using VisionSpec,",
				}
			}
		}
	}

	return repo, nil
}

// GetWorkflow returns the workflow with the given ID.
func (r *SourceRepo) GetWorkflow(id string) (*Workflow, error) {
	wf, ok := r.Workflows[id]
	if !ok {
		return nil, &NotFoundError{WorkflowID: id}
	}
	return wf, nil
}

// ListWorkflows returns all available workflow IDs.
func (r *SourceRepo) ListWorkflows() []string {
	var ids []string
	for id := range r.Workflows {
		ids = append(ids, id)
	}
	return ids
}
