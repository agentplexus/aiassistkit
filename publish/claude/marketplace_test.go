package claude

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPublisher_Name(t *testing.T) {
	p, err := NewPublisher("test-token")
	if err != nil {
		t.Fatalf("NewPublisher() error = %v", err)
	}
	if got := p.Name(); got != "claude" {
		t.Errorf("Name() = %q, want %q", got, "claude")
	}
}

func TestPublisher_Validate(t *testing.T) {
	p, err := NewPublisher("test-token")
	if err != nil {
		t.Fatalf("NewPublisher() error = %v", err)
	}

	// Create temp directory with valid plugin structure
	tmpDir, err := os.MkdirTemp("", "claude-plugin-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test validation fails with empty directory
	err = p.Validate(tmpDir)
	if err == nil {
		t.Error("Validate() should fail for empty directory")
	}

	// Create required files
	claudePluginDir := filepath.Join(tmpDir, ".claude-plugin")
	if err := os.MkdirAll(claudePluginDir, 0755); err != nil {
		t.Fatalf("Failed to create .claude-plugin dir: %v", err)
	}

	pluginJSON := `{"name": "test-plugin", "version": "1.0.0"}`
	if err := os.WriteFile(filepath.Join(claudePluginDir, "plugin.json"), []byte(pluginJSON), 0600); err != nil {
		t.Fatalf("Failed to create plugin.json: %v", err)
	}

	readme := "# Test Plugin\n\nA test plugin."
	if err := os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte(readme), 0600); err != nil {
		t.Fatalf("Failed to create README.md: %v", err)
	}

	// Test validation passes
	err = p.Validate(tmpDir)
	if err != nil {
		t.Errorf("Validate() failed for valid plugin: %v", err)
	}
}

func TestExtractDescription(t *testing.T) {
	tests := []struct {
		name     string
		readme   string
		wantDesc string
	}{
		{
			name:     "simple readme",
			readme:   "# My Plugin\n\nThis is a test plugin.",
			wantDesc: "This is a test plugin.",
		},
		{
			name:     "multi-line description",
			readme:   "# My Plugin\n\nLine one.\nLine two.\nLine three.",
			wantDesc: "Line one. Line two. Line three.",
		},
		{
			name:     "description with next heading",
			readme:   "# My Plugin\n\nShort description.\n\n## Installation\n\nInstall steps.",
			wantDesc: "Short description.",
		},
		{
			name:     "empty readme",
			readme:   "",
			wantDesc: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractDescription(tt.readme)
			if got != tt.wantDesc {
				t.Errorf("extractDescription() = %q, want %q", got, tt.wantDesc)
			}
		})
	}
}

func TestMarketplaceConfig(t *testing.T) {
	if MarketplaceOwner != "anthropics" {
		t.Errorf("MarketplaceOwner = %q, want %q", MarketplaceOwner, "anthropics")
	}
	if MarketplaceRepo != "claude-plugins-official" {
		t.Errorf("MarketplaceRepo = %q, want %q", MarketplaceRepo, "claude-plugins-official")
	}
	if ExternalPluginsPath != "external_plugins" {
		t.Errorf("ExternalPluginsPath = %q, want %q", ExternalPluginsPath, "external_plugins")
	}
}
