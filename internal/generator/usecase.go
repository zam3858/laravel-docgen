package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zam3858/laravel-docgen/internal/model"
)

func (g *Generator) writeUseCase(outDir string, diagrams []model.UseCaseDiagram) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("creating usecase output %s: %w", outDir, err)
	}
	for _, d := range diagrams {
		if err := writeJSON(filepath.Join(outDir, slug(d.Title)+".json"), d, g.opts.Validate); err != nil {
			return err
		}
	}
	return nil
}
