package adapters

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/knights-analytics/hugot"
	"github.com/knights-analytics/hugot/options"
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

// Close implements [ports.VectorGenerator].
func (v *VectorGenerator) Close() {
	v.session.Destroy()
}

func NewVectorGenerator() (ports.VectorGenerator, error) {
	modelDir, err := extractModel()
	if err != nil {
		return nil, fmt.Errorf("model Path: %w", err)
	}

	onnxLibPath := os.Getenv("ONNX_LIB_PATH")
	if onnxLibPath == "" {
		onnxLibPath = "/usr/lib"
	}
	session, err := hugot.NewORTSession(
		options.WithOnnxLibraryPath(onnxLibPath),
		options.WithInterOpNumThreads(10),
		options.WithIntraOpNumThreads(10),
		options.WithCPUMemArena(true),
		options.WithMemPattern(true),
	)
	if err != nil {
		return nil, fmt.Errorf("init: %w", err)
	}
	// //defer session.Destroy()

	// downloadOptions := hugot.NewDownloadOptions()

	// downloadOptions.OnnxFilePath = "onnx/model_qint8_avx512.onnx"

	// modelPath, err := hugot.DownloadModel("sentence-transformers/all-MiniLM-L6-v2", "./models/", downloadOptions)
	// //check(err)
	// if err != nil {
	// 	return nil, fmt.Errorf("download: %w", err)
	// }

	// Load the quantized 3-class model from the extracted temp dir
	config := hugot.FeatureExtractionConfig{
		ModelPath:    modelDir,
		Name:         "sentence-transformers_all-MiniLM-L6-v2",
		OnnxFilename: "onnx/model_qint8_avx512.onnx",
	}

	commandPipeline, err := hugot.NewPipeline(session, config)
	if err != nil {
		return nil, fmt.Errorf("pipeline creating: %w", err)
	}

	return &VectorGenerator{
		session:  session,
		pipeline: commandPipeline,
		modelDir: modelDir,
	}, nil
}

func (v *VectorGenerator) Generate(ctx context.Context, queryTerm string) (domain.Vector, error) {
	output, err := v.pipeline.Run([]string{queryTerm})
	if err != nil {
		return nil, fmt.Errorf("generate: %w", err)
	}

	return domain.Vector(lo.Map(output.GetOutput()[0].([]float32), func(v float32, _ int) float64 {
		return float64(v)
	})), nil
}

//go:embed model_data/*
var modelFS embed.FS

// extractModel extracts the embedded model files to a temporary directory
// and returns the path to the extracted model directory.
// The caller is responsible for cleaning up with os.RemoveAll.
func extractModel() (string, error) {
	tmpDir, err := os.MkdirTemp("", "hugot-model-*")
	if err != nil {
		return "", fmt.Errorf("creating temp dir: %w", err)
	}

	err = fs.WalkDir(modelFS, "model_data", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Get relative path under model_data/
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
