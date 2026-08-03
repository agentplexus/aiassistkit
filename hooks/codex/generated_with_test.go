package codex

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestGeneratedWithTrailer(t *testing.T) {
	got := GeneratedWithTrailer("OpenAI Codex", "gpt-5.5")
	want := "Generated-with: OpenAI Codex gpt-5.5"
	if got != want {
		t.Fatalf("GeneratedWithTrailer() = %q, want %q", got, want)
	}
}

func TestRunGeneratedWith(t *testing.T) {
	var out bytes.Buffer
	err := RunGeneratedWith(strings.NewReader(`{"hook_event_name":"SessionStart","model":"gpt-5.6-sol"}`), &out, "OpenAI Codex", "")
	if err != nil {
		t.Fatalf("RunGeneratedWith() error = %v", err)
	}

	var got HookOutput
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if got.HookSpecificOutput.HookEventName != "SessionStart" {
		t.Fatalf("hookEventName = %q, want SessionStart", got.HookSpecificOutput.HookEventName)
	}
	if !strings.Contains(got.HookSpecificOutput.AdditionalContext, "Generated-with: OpenAI Codex gpt-5.6-sol") {
		t.Fatalf("additionalContext missing generated-with trailer: %q", got.HookSpecificOutput.AdditionalContext)
	}
}

func TestRunGeneratedWith_UnknownModelFallback(t *testing.T) {
	var out bytes.Buffer
	err := RunGeneratedWith(strings.NewReader(`{"hook_event_name":"SessionStart"}`), &out, "OpenAI Codex", "unknown")
	if err != nil {
		t.Fatalf("RunGeneratedWith() error = %v", err)
	}
	if !strings.Contains(out.String(), "Generated-with: OpenAI Codex unknown") {
		t.Fatalf("output missing unknown model fallback: %s", out.String())
	}
}
