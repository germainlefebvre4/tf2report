package main

import (
	"fmt"
	"os"

	"github.com/germainlefebvre4/tf2report/pkg/config"
	"github.com/germainlefebvre4/tf2report/pkg/report"
	"github.com/germainlefebvre4/tf2report/pkg/terraform"
	"github.com/spf13/cobra"
)

var (
	cfgFile       string
	planPath      string
	outputFormat  string
	resourceTypes []string
	actions       []string
	verbose       bool
	appVersion    = "dev"
	buildTime     = "unknown"
	buildCommit   = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "tf2report",
	Short: "Summarize Terraform plan changes into human-readable reports",
	Long: `tf2report analyzes Terraform plan files and generates resource change reports.

It parses Terraform plan JSON files, extracts resource changes, and generates
summaries in multiple output formats including Markdown, plain text, and JSON.`,
	RunE:    run,
	Version: appVersion,
}

func init() {
	// Custom version template to include build time
	versionTemplate := fmt.Sprintf("tf2report\nVersion: %s\nCommit: %s\nBuild time: %s\n", appVersion, buildCommit, buildTime)
	rootCmd.SetVersionTemplate(versionTemplate)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path (default: ./tf2report.yaml)")
	rootCmd.Flags().StringVarP(&planPath, "plan", "p", "", "path to Terraform plan JSON file")
	rootCmd.Flags().StringVarP(&outputFormat, "format", "f", "", "output format (markdown, text, json)")
	rootCmd.Flags().StringSliceVarP(&resourceTypes, "type", "t", []string{}, "filter by resource type (can be specified multiple times)")
	rootCmd.Flags().StringSliceVarP(&actions, "action", "a", []string{}, "filter by action (create, update, delete, replace)")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func run(cmd *cobra.Command, args []string) error {
	// Load configuration
	loader := config.NewLoader()
	cfg, err := loader.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Override config with command-line flags
	if planPath != "" {
		cfg.TerraformPlanPath = planPath
	}
	if outputFormat != "" {
		cfg.OutputFormat = outputFormat
	}
	if len(resourceTypes) > 0 {
		cfg.Filters.ResourceTypes = resourceTypes
	}
	if len(actions) > 0 {
		cfg.Filters.Actions = actions
	}
	if verbose {
		cfg.Verbosity = "debug"
	}

	// Validate plan file path
	if cfg.TerraformPlanPath == "" {
		return fmt.Errorf("terraform plan path not specified (use --plan or config file)")
	}

	// Check if file exists
	if _, err := os.Stat(cfg.TerraformPlanPath); os.IsNotExist(err) {
		return fmt.Errorf("plan file not found: %s", cfg.TerraformPlanPath)
	}

	// Parse plan file
	parser := terraform.NewParser()
	plan, err := parser.ParseFile(cfg.TerraformPlanPath)
	if err != nil {
		return fmt.Errorf("failed to parse plan: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Parsed plan with %d resource changes\n", len(plan.ResourceChanges))
	}

	// Apply filters
	changes := plan.ResourceChanges
	if len(cfg.Filters.ResourceTypes) > 0 {
		changes = terraform.FilterByType(changes, cfg.Filters.ResourceTypes)
		if verbose {
			fmt.Fprintf(os.Stderr, "Filtered by type: %d resources remain\n", len(changes))
		}
	}
	if len(cfg.Filters.Actions) > 0 {
		actionTypes := cfg.Filters.GetActionTypes()
		changes = terraform.FilterByAction(changes, actionTypes)
		if verbose {
			fmt.Fprintf(os.Stderr, "Filtered by action: %d resources remain\n", len(changes))
		}
	}

	// Update plan with filtered changes
	plan.ResourceChanges = changes

	// Generate summary
	summary := terraform.NewSummary(plan.ResourceChanges)

	// Parse output format
	format, err := report.ParseFormat(cfg.OutputFormat)
	if err != nil {
		return err
	}

	// Create formatter
	formatter, err := report.NewFormatter(format)
	if err != nil {
		return err
	}

	// Generate report
	if err := formatter.Format(plan, summary, os.Stdout); err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
