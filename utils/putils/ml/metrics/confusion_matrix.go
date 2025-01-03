package metrics

import (
	"fmt"
	"strings"
)

type ConfusionMatrix struct {
	matrix [][]int
	labels []string
}

func NewConfusionMatrix(actual, predicted []string, labels []string) *ConfusionMatrix {
	n := len(labels)
	matrix := make([][]int, n)
	for i := range matrix {
		matrix[i] = make([]int, n)
	}

	labelIndices := make(map[string]int)
	for i, label := range labels {
		labelIndices[label] = i
	}

	for i := range actual {
		matrix[labelIndices[actual[i]]][labelIndices[predicted[i]]]++
	}

	return &ConfusionMatrix{
		matrix: matrix,
		labels: labels,
	}
}

func (cm *ConfusionMatrix) PrintConfusionMatrix() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("%30s\n", "Confusion Matrix"))
	s.WriteString(fmt.Sprintln())
	// Print header
	s.WriteString(fmt.Sprintf("%-15s", ""))
	for _, label := range cm.labels {
		s.WriteString(fmt.Sprintf("%-15s", label))
	}
	s.WriteString(fmt.Sprintln())

	// Print rows
	for i, row := range cm.matrix {
		s.WriteString(fmt.Sprintf("%-15s", cm.labels[i]))
		for _, value := range row {
			s.WriteString(fmt.Sprintf("%-15d", value))
		}
		s.WriteString(fmt.Sprintln())
	}
	s.WriteString(fmt.Sprintln())

	return s.String()
}
