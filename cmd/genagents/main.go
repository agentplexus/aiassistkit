// Command genagents generates platform-specific agent files from canonical specs.
//
// Usage:
//
//	genagents -spec=plugins/spec/agents -output=.claude/agents -format=claude
//	genagents -spec=plugins/spec/agents -output=plugins/kiro/agents -format=kiro
//	genagents -spec=plugins/spec/agents -targets=claude:.claude/agents,kiro:plugins/kiro/agents
//
// Multi-agent-spec format (reads deployment.json for targets):
//
//	genagents -project=examples/stats-agent-team
//	genagents -project=examples/stats-agent-team -priority=p1
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/agentplexus/assistantkit/agents"
	"github.com/agentplexus/assistantkit/agents/agentkit"
	"github.com/agentplexus/assistantkit/agents/awsagentcore"
	"github.com/agentplexus/assistantkit/agents/core"

	// Import adapters to register them
	_ "github.com/agentplexus/assistantkit/agents/claude"
	_ "github.com/agentplexus/assistantkit/agents/kiro"
)

func main() {
	specDir := flag.String("spec", "plugins/spec/agents", "Directory containing canonical agent specs (.md files)")
	outputDir := flag.String("output", "", "Output directory for generated agents")
	format := flag.String("format", "claude", "Output format (claude, kiro, agentkit, aws-agentcore)")
	targets := flag.String("targets", "", "Multiple targets as format:dir pairs (e.g., claude:.claude/agents,kiro:plugins/kiro/agents)")
	project := flag.String("project", "", "Multi-agent-spec project directory (reads deployment.json)")
	priority := flag.String("priority", "", "Filter by priority (p1, p2, p3) - only with -project")
	verbose := flag.Bool("verbose", false, "Verbose output")
	flag.Parse()

	// Handle multi-agent-spec project mode
	if *project != "" {
		if err := runProjectMode(*project, *priority, *verbose); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Read canonical agents from spec directory
	agentList, err := agents.ReadCanonicalDir(*specDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading spec directory %s: %v\n", *specDir, err)
		os.Exit(1)
	}

	if len(agentList) == 0 {
		fmt.Fprintf(os.Stderr, "No agents found in %s\n", *specDir)
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("Found %d agents in %s\n", len(agentList), *specDir)
		for _, agent := range agentList {
			fmt.Printf("  - %s: %s\n", agent.Name, agent.Description)
		}
	}

	// Handle multiple targets
	if *targets != "" {
		targetPairs := strings.Split(*targets, ",")
		for _, pair := range targetPairs {
			parts := strings.SplitN(pair, ":", 2)
			if len(parts) != 2 {
				fmt.Fprintf(os.Stderr, "Invalid target format: %s (expected format:dir)\n", pair)
				os.Exit(1)
			}
			targetFormat := strings.TrimSpace(parts[0])
			targetDir := strings.TrimSpace(parts[1])

			if err := generateAgents(agentList, targetFormat, targetDir, *verbose); err != nil {
				fmt.Fprintf(os.Stderr, "Error generating %s agents: %v\n", targetFormat, err)
				os.Exit(1)
			}
		}
		return
	}

	// Handle single target
	if *outputDir == "" {
		fmt.Fprintf(os.Stderr, "Error: -output, -targets, or -project required\n")
		flag.Usage()
		os.Exit(1)
	}

	if err := generateAgents(agentList, *format, *outputDir, *verbose); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating agents: %v\n", err)
		os.Exit(1)
	}
}

func generateAgents(agentList []*core.Agent, format, outputDir string, verbose bool) error {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Get the adapter
	adapter, ok := core.GetAdapter(format)
	if !ok {
		available := core.AdapterNames()
		return fmt.Errorf("unknown format %q (available: %s)", format, strings.Join(available, ", "))
	}

	// Write each agent
	for _, agent := range agentList {
		filename := agent.Name + adapter.FileExtension()
		path := filepath.Join(outputDir, filename)

		if err := adapter.WriteFile(agent, path); err != nil {
			return fmt.Errorf("failed to write %s: %w", path, err)
		}

		if verbose {
			fmt.Printf("Generated %s\n", path)
		}
	}

	fmt.Printf("Generated %d %s agents in %s\n", len(agentList), format, outputDir)
	return nil
}

// Deployment represents deployment.json from multi-agent-spec format.
type Deployment struct {
	Schema  string   `json:"$schema"`
	Team    string   `json:"team"`
	Targets []Target `json:"targets"`
}

// Target represents a deployment target.
type Target struct {
	Name     string                 `json:"name"`
	Platform string                 `json:"platform"`
	Priority string                 `json:"priority"`
	Output   string                 `json:"output"`
	Config   map[string]interface{} `json:"config"`
}

// runProjectMode processes a multi-agent-spec project directory.
func runProjectMode(projectDir, priorityFilter string, verbose bool) error {
	// Read deployment.json
	deploymentPath := filepath.Join(projectDir, "deployment.json")
	deploymentData, err := os.ReadFile(deploymentPath)
	if err != nil {
		return fmt.Errorf("failed to read deployment.json: %w", err)
	}

	var deployment Deployment
	if err := json.Unmarshal(deploymentData, &deployment); err != nil {
		return fmt.Errorf("failed to parse deployment.json: %w", err)
	}

	if verbose {
		fmt.Printf("Processing project: %s\n", deployment.Team)
		fmt.Printf("Found %d deployment targets\n", len(deployment.Targets))
	}

	// Read agents from agents/ directory
	agentsDir := filepath.Join(projectDir, "agents")
	agentList, err := agents.ReadCanonicalDir(agentsDir)
	if err != nil {
		return fmt.Errorf("failed to read agents: %w", err)
	}

	if len(agentList) == 0 {
		return fmt.Errorf("no agents found in %s", agentsDir)
	}

	if verbose {
		fmt.Printf("Found %d agents:\n", len(agentList))
		for _, agent := range agentList {
			fmt.Printf("  - %s\n", agent.Name)
		}
	}

	// Process each target
	for _, target := range deployment.Targets {
		// Filter by priority if specified
		if priorityFilter != "" && target.Priority != priorityFilter {
			if verbose {
				fmt.Printf("Skipping %s (priority %s, filter %s)\n", target.Name, target.Priority, priorityFilter)
			}
			continue
		}

		outputDir := filepath.Join(projectDir, target.Output)

		if verbose {
			fmt.Printf("\nProcessing target: %s (%s)\n", target.Name, target.Platform)
			fmt.Printf("  Output: %s\n", outputDir)
		}

		if err := generateForPlatform(deployment.Team, agentList, target, outputDir, verbose); err != nil {
			return fmt.Errorf("failed to generate %s: %w", target.Name, err)
		}
	}

	return nil
}

// generateForPlatform generates output for a specific platform.
func generateForPlatform(teamName string, agentList []*core.Agent, target Target, outputDir string, verbose bool) error {
	switch target.Platform {
	case "claude-code":
		return generateAgents(agentList, "claude", outputDir, verbose)

	case "kiro-cli":
		return generateAgents(agentList, "kiro", outputDir, verbose)

	case "agentkit-local":
		// Generate full agentkit config
		configPath := filepath.Join(outputDir, "config.json")
		if err := agentkit.WriteFullConfig(agentList, configPath); err != nil {
			return err
		}
		fmt.Printf("Generated agentkit config: %s\n", configPath)
		return nil

	case "aws-agentcore":
		// Generate CDK project
		config := &awsagentcore.AgentCoreConfig{
			StackName: toPascalCase(teamName) + "Stack",
		}
		// Apply config from deployment.json if present
		if region, ok := target.Config["region"].(string); ok {
			config.Region = region
		}
		if model, ok := target.Config["foundationModel"].(string); ok {
			config.FoundationModel = model
		}
		if runtime, ok := target.Config["lambdaRuntime"].(string); ok {
			config.LambdaRuntime = runtime
		}

		if err := awsagentcore.WriteCDKProject(teamName, agentList, outputDir, config); err != nil {
			return err
		}
		fmt.Printf("Generated CDK project in %s\n", outputDir)
		return nil

	case "aws-eks", "azure-aks", "gcp-gke", "kubernetes":
		// TODO: Implement Helm chart generation
		fmt.Printf("Kubernetes deployment not yet implemented for %s\n", target.Platform)
		return nil

	default:
		return fmt.Errorf("unsupported platform: %s", target.Platform)
	}
}

// toPascalCase converts a hyphenated string to PascalCase.
func toPascalCase(s string) string {
	parts := strings.Split(s, "-")
	var result strings.Builder
	for _, part := range parts {
		if len(part) > 0 {
			result.WriteString(strings.ToUpper(part[:1]))
			result.WriteString(part[1:])
		}
	}
	return result.String()
}
