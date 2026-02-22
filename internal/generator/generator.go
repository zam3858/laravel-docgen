package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zam3858/laravel-docgen/internal/model"
)

type Options struct {
	Validate bool
}

type Generator struct {
	opts Options
}

func New(opts Options) *Generator {
	return &Generator{opts: opts}
}

func (g *Generator) Generate(outputRoot string, analysis model.Analysis, kind string) error {
	if err := os.MkdirAll(outputRoot, 0o755); err != nil {
		return fmt.Errorf("creating output root %s: %w", outputRoot, err)
	}

	switch kind {
	case "sequence":
		return g.writeSequence(filepath.Join(outputRoot, "sequence"), analysis.Sequence)
	case "objects":
		return g.writeObjects(filepath.Join(outputRoot, "objects"), analysis.Objects)
	case "usecase":
		return g.writeUseCase(filepath.Join(outputRoot, "usecase"), analysis.UseCase)
	case "all":
		if err := g.writeSequence(filepath.Join(outputRoot, "sequence"), analysis.Sequence); err != nil {
			return err
		}
		if err := g.writeObjects(filepath.Join(outputRoot, "objects"), analysis.Objects); err != nil {
			return err
		}
		return g.writeUseCase(filepath.Join(outputRoot, "usecase"), analysis.UseCase)
	default:
		return fmt.Errorf("unsupported diagram type: %s", kind)
	}
}
