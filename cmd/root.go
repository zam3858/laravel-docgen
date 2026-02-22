package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zam3858/laravel-docgen/internal/analyzer"
	"github.com/zam3858/laravel-docgen/internal/cache"
	"github.com/zam3858/laravel-docgen/internal/generator"
	"github.com/zam3858/laravel-docgen/internal/parser"
)

var errUsage = errors.New("usage")

func Execute() error {
	if len(os.Args) < 2 {
		printUsage()
		return errUsage
	}

	switch os.Args[1] {
	case "analyze":
		return runAnalyze(os.Args[2:])
	case "generate":
		return runGenerate(os.Args[2:])
	case "-h", "--help", "help":
		printUsage()
		return nil
	default:
		printUsage()
		return fmt.Errorf("unknown command: %s", os.Args[1])
	}
}

func runAnalyze(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("analyze requires an input path")
	}
	inputPath := args[0]
	opts, err := parseGlobalFlags(args[1:])
	if err != nil {
		return err
	}
	return runPipeline(inputPath, opts, "all")
}

func runGenerate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("generate requires a diagram type")
	}
	kind := strings.ToLower(args[0])
	switch kind {
	case "sequence", "objects", "usecase", "all":
	default:
		return fmt.Errorf("unsupported diagram type: %s", kind)
	}

	fs := flag.NewFlagSet("generate", flag.ContinueOnError)
	inputPath := fs.String("input", "", "Path to Laravel project root")
	outputPath := fs.String("output", "./docs", "Output directory")
	verbose := fs.Bool("verbose", false, "Enable verbose logging")
	incremental := fs.Bool("incremental", false, "Skip unchanged files")
	validate := fs.Bool("validate", false, "Validate generated JSON")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	if *inputPath == "" {
		return fmt.Errorf("--input is required")
	}

	return runPipeline(*inputPath, runOptions{
		outputPath:  *outputPath,
		verbose:     *verbose,
		incremental: *incremental,
		validate:    *validate,
	}, kind)
}

type runOptions struct {
	outputPath  string
	verbose     bool
	incremental bool
	validate    bool
}

func parseGlobalFlags(args []string) (runOptions, error) {
	fs := flag.NewFlagSet("analyze", flag.ContinueOnError)
	outputPath := fs.String("output", "./docs", "Output directory")
	verbose := fs.Bool("verbose", false, "Enable verbose logging")
	incremental := fs.Bool("incremental", false, "Skip unchanged files")
	validate := fs.Bool("validate", false, "Validate generated JSON")
	if err := fs.Parse(args); err != nil {
		return runOptions{}, err
	}
	return runOptions{
		outputPath:  *outputPath,
		verbose:     *verbose,
		incremental: *incremental,
		validate:    *validate,
	}, nil
}

func runPipeline(inputPath string, opts runOptions, kind string) error {
	absInput, err := filepath.Abs(inputPath)
	if err != nil {
		return fmt.Errorf("resolving input path: %w", err)
	}
	absOutput, err := filepath.Abs(opts.outputPath)
	if err != nil {
		return fmt.Errorf("resolving output path: %w", err)
	}

	if _, err := os.Stat(absInput); err != nil {
		return fmt.Errorf("invalid input path %s: %w", absInput, err)
	}

	c, err := cache.New(filepath.Join(absOutput, ".cache.json"))
	if err != nil {
		return fmt.Errorf("creating cache: %w", err)
	}

	p := parser.NewParser(parser.Options{Incremental: opts.incremental, Cache: c})
	project, warnings, err := p.Parse(absInput)
	if err != nil {
		return err
	}

	if opts.verbose {
		for _, w := range warnings {
			fmt.Fprintf(os.Stderr, "warning: %s\n", w)
		}
	}

	analysis := analyzer.New().Analyze(project)
	if err := generator.New(generator.Options{Validate: opts.validate}).Generate(absOutput, analysis, kind); err != nil {
		return err
	}

	if opts.incremental {
		if err := c.Save(); err != nil {
			return fmt.Errorf("saving cache: %w", err)
		}
	}

	return nil
}

func printUsage() {
	fmt.Println("laravel-docgen")
	fmt.Println("  analyze <path> [--output ./docs] [--verbose] [--incremental] [--validate]")
	fmt.Println("  generate <sequence|objects|usecase|all> --input <path> [--output ./docs] [--verbose] [--incremental] [--validate]")
}
