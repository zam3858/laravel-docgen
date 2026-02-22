package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/zam3858/laravel-docgen/internal/model"
)

func (g *Generator) writeSequence(outDir string, diagrams []model.SequenceDiagram) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("creating sequence output %s: %w", outDir, err)
	}
	for _, d := range diagrams {
		if err := writeJSON(filepath.Join(outDir, slug(d.Title)+".json"), d, g.opts.Validate); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(path string, v any, validate bool) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json %s: %w", path, err)
	}
	if validate && !json.Valid(b) {
		return fmt.Errorf("generated invalid json: %s", path)
	}
	b = append(b, '\n')
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return fmt.Errorf("write json %s: %w", path, err)
	}
	return nil
}

var nonSlugChars = regexp.MustCompile(`[^a-z0-9-]+`)
var multiDashRe = regexp.MustCompile(`-+`)

func slug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "/", "-")
	s = nonSlugChars.ReplaceAllString(s, "")
	s = multiDashRe.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return "diagram"
	}
	return s
}
