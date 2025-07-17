package neuralnet

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
)

type ImageClassification struct {
	modelPath string
	pythonCmd string
}

func NewImageClassification() *ImageClassification {
	return &ImageClassification{
		pythonCmd: "python3",
	}
}

func (ic *ImageClassification) LoadModel(modelPath string) error {
	absPath, err := filepath.Abs(modelPath)
	if err != nil {
		return err
	}

	ic.modelPath = absPath
	return nil
}

func (ic *ImageClassification) Classify(imagePath string) (*ClassificationResult, error) {
	if ic.modelPath == "" {
		return nil, fmt.Errorf("no model loaded")
	}

	cmd := exec.Command(ic.pythonCmd, "-m", "python.embeddings", "--command", "predict", "--args", fmt.Sprintf(`{"model_type":"image_classification","model_path":"%s","image_path":"%s"}`, ic.modelPath, imagePath))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("classification failed: %v, output: %s", err, output)
	}

	var result struct {
		Result *ClassificationResult `json:"result"`
		Error  string                `json:"error"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse result: %v", err)
	}

	if result.Error != "" {
		return nil, fmt.Errorf("python error: %s", result.Error)
	}

	return result.Result, nil
}

type ClassificationResult struct {
	Predictions []float64    `json:"predictions"`
	Confidence  float64      `json:"confidence"`
	Class       string       `json:"class"`
	TopK        []ClassScore `json:"top_k,omitempty"`
}

type ClassScore struct {
	Class      string  `json:"class"`
	Confidence float64 `json:"confidence"`
}
