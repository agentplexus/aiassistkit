// Package codex provides helpers for OpenAI Codex hook commands.
package codex

import (
	"encoding/json"
	"fmt"
	"io"
)

const (
	// DefaultAssistantName is the product name used in generated commit trailers.
	DefaultAssistantName = "OpenAI Codex"

	// DefaultUnknownModel is used when Codex hook input omits the active model.
	DefaultUnknownModel = "unknown-model"
)

// HookInput is the subset of Codex hook input used by the generated-with hook.
type HookInput struct {
	HookEventName string `json:"hook_event_name,omitempty"`
	Model         string `json:"model,omitempty"`
}

// HookOutput is the Codex hook output shape for adding developer context.
type HookOutput struct {
	HookSpecificOutput HookSpecificOutput `json:"hookSpecificOutput"`
}

// HookSpecificOutput is the Codex hook-specific output payload.
type HookSpecificOutput struct {
	HookEventName     string `json:"hookEventName"`
	AdditionalContext string `json:"additionalContext"`
}

// GeneratedWithTrailer returns the commit trailer for the given model.
func GeneratedWithTrailer(assistantName, model string) string {
	if assistantName == "" {
		assistantName = DefaultAssistantName
	}
	if model == "" {
		model = DefaultUnknownModel
	}
	return fmt.Sprintf("Generated-with: %s %s", assistantName, model)
}

// GeneratedWithContext returns developer context that asks Codex to include
// detailed commit messages and a model-specific Generated-with trailer.
func GeneratedWithContext(assistantName, model string) string {
	return fmt.Sprintf(
		"When creating or amending commits, use a detailed Conventional Commit message with an explanatory body and a Tests section when tests were run. Include this trailer when the active model is known:\n\n%s",
		GeneratedWithTrailer(assistantName, model),
	)
}

// NewGeneratedWithOutput builds the SessionStart hook output for the model.
func NewGeneratedWithOutput(assistantName, model string) HookOutput {
	return HookOutput{
		HookSpecificOutput: HookSpecificOutput{
			HookEventName:     "SessionStart",
			AdditionalContext: GeneratedWithContext(assistantName, model),
		},
	}
}

// RunGeneratedWith reads Codex hook input from r and writes Codex hook output to w.
func RunGeneratedWith(r io.Reader, w io.Writer, assistantName, unknownModel string) error {
	if unknownModel == "" {
		unknownModel = DefaultUnknownModel
	}

	var input HookInput
	if err := json.NewDecoder(r).Decode(&input); err != nil {
		return fmt.Errorf("decode codex hook input: %w", err)
	}

	model := input.Model
	if model == "" {
		model = unknownModel
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(NewGeneratedWithOutput(assistantName, model)); err != nil {
		return fmt.Errorf("encode codex hook output: %w", err)
	}
	return nil
}
