package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/zam3858/laravel-docgen/internal/analyzer"
	"github.com/zam3858/laravel-docgen/internal/cache"
	"github.com/zam3858/laravel-docgen/internal/generator"
	"github.com/zam3858/laravel-docgen/internal/parser"
)

func TestPipelineGeneratesDeterministicJSON(t *testing.T) {
	fixtureRoot := filepath.Join("fixtures", "laravel")
	out := t.TempDir()

	c, err := cache.New(filepath.Join(out, ".cache.json"))
	if err != nil {
		t.Fatalf("cache: %v", err)
	}

	run := func() map[string][]byte {
		t.Helper()
		p := parser.NewParser(parser.Options{Incremental: false, Cache: c})
		project, warnings, err := p.Parse(fixtureRoot)
		if err != nil {
			t.Fatalf("parse: %v", err)
		}
		if len(warnings) > 0 {
			t.Logf("warnings: %v", warnings)
		}

		analysis := analyzer.New().Analyze(project)
		if err := generator.New(generator.Options{Validate: true}).Generate(out, analysis, "all"); err != nil {
			t.Fatalf("generate: %v", err)
		}

		files := []string{
			filepath.Join(out, "sequence", "post-login.json"),
			filepath.Join(out, "objects", "user-model.json"),
			filepath.Join(out, "usecase", "application-use-cases.json"),
		}
		res := make(map[string][]byte, len(files))
		for _, f := range files {
			b, err := os.ReadFile(f)
			if err != nil {
				t.Fatalf("read %s: %v", f, err)
			}
			res[f] = b
		}
		return res
	}

	first := run()
	second := run()

	for file, fb := range first {
		sb := second[file]
		if string(fb) != string(sb) {
			t.Fatalf("non-deterministic output for %s", file)
		}
	}
}
