package neuralnet

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
)

type ObjectDetection struct {
	modelPath string
	pythonCmd string
}

func NewObjectDetection() *ObjectDetection {
	return &ObjectDetection{
		pythonCmd: "python3",
	}
}

func (od *ObjectDetection) LoadModel(modelPath string) error {
	absPath, err := filepath.Abs(modelPath)
	if err != nil {
		return err
	}

	od.modelPath = absPath
	return nil
}

func (od *ObjectDetection) Detect(imagePath string) (*DetectionResult, error) {
	if od.modelPath == "" {
		return nil, fmt.Errorf("no model loaded")
	}

	cmd := exec.Command(od.pythonCmd, "-m", "python.embeddings", "--command", "predict", "--args", fmt.Sprintf(`{"model_type":"object_detection","model_path":"%s","image_path":"%s"}`, od.modelPath, imagePath))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("detection failed: %v, output: %s", err, output)
	}

	var result struct {
		Result *DetectionResult `json:"result"`
		Error  string           `json:"error"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse result: %v", err)
	}

	if result.Error != "" {
		return nil, fmt.Errorf("python error: %s", result.Error)
	}

	return result.Result, nil
}

type DetectionResult struct {
	Objects []DetectedObject `json:"objects"`
	Count   int              `json:"count"`
}

type DetectedObject struct {
	Class      string      `json:"class"`
	Confidence float64     `json:"confidence"`
	Box        BoundingBox `json:"box"`
}

type BoundingBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}
