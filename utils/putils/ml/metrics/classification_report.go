package metrics

import (
	"fmt"
	"strings"
)

func (cm *ConfusionMatrix) PrintClassificationReport() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("%30s\n", "Classification Report"))
	s.WriteString(fmt.Sprintln())

	s.WriteString(fmt.Sprintf("\n%-15s %-10s %-10s %-10s %-10s\n", "", "precision", "recall", "f1-score", "support"))

	totals := map[string]float64{"true": 0, "predicted": 0, "correct": 0}
	macroAvg := map[string]float64{"precision": 0, "recall": 0, "f1-score": 0}

	for i, label := range cm.labels {
		truePos := cm.matrix[i][i]
		falsePos, falseNeg := 0, 0
		for j := 0; j < len(cm.labels); j++ {
			if i != j {
				falsePos += cm.matrix[j][i]
				falseNeg += cm.matrix[i][j]
			}
		}

		precision := float64(truePos) / float64(truePos+falsePos)
		recall := float64(truePos) / float64(truePos+falseNeg)
		f1Score := 2 * precision * recall / (precision + recall)
		support := truePos + falseNeg

		fmt.Printf("%-15s %-10.2f %-10.2f %-10.2f %-10d\n", label, precision, recall, f1Score, support)

		totals["true"] += float64(support)
		totals["predicted"] += float64(truePos + falsePos)
		totals["correct"] += float64(truePos)

		macroAvg["precision"] += precision
		macroAvg["recall"] += recall
		macroAvg["f1-score"] += f1Score
	}

	accuracy := totals["correct"] / totals["true"]
	s.WriteString(fmt.Sprintf("\n%-26s %-10s %-10.2f %-10d", "accuracy", "", accuracy, int(totals["true"])))

	s.WriteString(fmt.Sprintf("\n%-15s %-10.2f %-10.2f %-10.2f %-10d\n", "macro avg",
		macroAvg["precision"]/float64(len(cm.labels)),
		macroAvg["recall"]/float64(len(cm.labels)),
		macroAvg["f1-score"]/float64(len(cm.labels)),
		int(totals["true"])))

	precisionWeightedAvg := totals["correct"] / totals["predicted"]
	recallWeightedAvg := totals["correct"] / totals["true"]
	f1ScoreWeightedAvg := 2 * precisionWeightedAvg * recallWeightedAvg / (precisionWeightedAvg + recallWeightedAvg)

	s.WriteString(fmt.Sprintf("%-15s %-10.2f %-10.2f %-10.2f %-10d\n", "weighted avg",
		precisionWeightedAvg, recallWeightedAvg, f1ScoreWeightedAvg, int(totals["true"])))

	s.WriteString(fmt.Sprintln())

	return s.String()
}
