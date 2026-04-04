package adapters

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/knights-analytics/hugot"
	"github.com/knights-analytics/hugot/pipelines"
	"github.com/pdaccess/ws/internal/core/domain"
	"github.com/pdaccess/ws/internal/core/ports"
	"github.com/samber/lo"
)

type VectorGenerator struct {
	session  *hugot.Session
	pipeline *pipelines.FeatureExtractionPipeline
	modelDir string
}

func (v *VectorGenerator) Close() {
	if v.session != nil {
		v.session.Destroy()
	}
	if v.modelDir != "" {
		os.RemoveAll(v.modelDir)
	}
}

type StubVectorGenerator struct{}

func (s *StubVectorGenerator) Close() {}

func (s *StubVectorGenerator) Generate(ctx context.Context, queryTerm string) (domain.Vector, error) {
	return nil, fmt.Errorf("stub vector generator")
}

func NewVectorGenerator() (ports.VectorGenerator, error) {
	modelDir, err := extractModel()
	if err != nil {
		return nil, fmt.Errorf("failed to extract model: %w", err)
	}

	session, err := hugot.NewGoSession()
	if err != nil {
		os.RemoveAll(modelDir)
		return nil, fmt.Errorf("failed to create hugot session: %w", err)
	}

	config := hugot.FeatureExtractionConfig{
		ModelPath: modelDir,
		Name:      "feature-extraction",
	}
	pipeline, err := hugot.NewPipeline(session, config)
	if err != nil {
		session.Destroy()
		os.RemoveAll(modelDir)
		return nil, fmt.Errorf("failed to create feature extraction pipeline: %w", err)
	}

	return &VectorGenerator{
		session:  session,
		pipeline: pipeline,
		modelDir: modelDir,
	}, nil
}

func (v *VectorGenerator) Generate(ctx context.Context, queryTerm string) (domain.Vector, error) {
	result, err := v.pipeline.Run([]string{queryTerm})
	if err != nil {
		return nil, fmt.Errorf("generate: %w", err)
	}

	output := result.GetOutput()
	if len(output) == 0 {
		return nil, fmt.Errorf("no output from pipeline")
	}

	return domain.Vector(lo.Map(output[0].([]float32), func(f float32, _ int) float64 {
		return float64(f)
	})), nil
}

//go:embed model_data/*
var modelFS embed.FS

func extractModel() (string, error) {
	tmpDir, err := os.MkdirTemp("", "hugot-model-*")
	if err != nil {
		return "", fmt.Errorf("creating temp dir: %w", err)
	}

	err = fs.WalkDir(modelFS, "model_data", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel("model_data", path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(tmpDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0o755)
		}

		data, err := modelFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading embedded file %s: %w", path, err)
		}

		return os.WriteFile(targetPath, data, 0o644)
	})
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("extracting model: %w", err)
	}

	return tmpDir, nil
}
