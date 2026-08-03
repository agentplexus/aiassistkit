package main

import (
	"os"

	codexhooks "github.com/plexusone/assistantkit/hooks/codex"
	"github.com/spf13/cobra"
)

var (
	codexGeneratedWithAssistantName string
	codexGeneratedWithUnknownModel  string
)

var codexCmd = &cobra.Command{
	Use:   "codex",
	Short: "Utilities for OpenAI Codex",
}

var codexHooksCmd = &cobra.Command{
	Use:   "hooks",
	Short: "Codex hook commands",
}

var codexHooksGeneratedWithCmd = &cobra.Command{
	Use:   "generated-with",
	Short: "Emit SessionStart context for Generated-with commit trailers",
	Long: `Read Codex hook JSON from stdin and emit SessionStart hook output that
instructs Codex to include a Generated-with trailer containing the active model.

Example:
  assistantkit codex hooks generated-with`,
	RunE: runCodexHooksGeneratedWith,
}

func init() {
	rootCmd.AddCommand(codexCmd)
	codexCmd.AddCommand(codexHooksCmd)
	codexHooksCmd.AddCommand(codexHooksGeneratedWithCmd)

	codexHooksGeneratedWithCmd.Flags().StringVar(&codexGeneratedWithAssistantName, "assistant-name", codexhooks.DefaultAssistantName, "Assistant name for the Generated-with trailer")
	codexHooksGeneratedWithCmd.Flags().StringVar(&codexGeneratedWithUnknownModel, "unknown-model", codexhooks.DefaultUnknownModel, "Model name to use when Codex hook input omits model")
}

func runCodexHooksGeneratedWith(cmd *cobra.Command, args []string) error {
	return codexhooks.RunGeneratedWith(os.Stdin, os.Stdout, codexGeneratedWithAssistantName, codexGeneratedWithUnknownModel)
}
