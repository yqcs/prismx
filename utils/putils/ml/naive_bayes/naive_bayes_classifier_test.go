package naive_bayes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNaiveBayesClassifier(t *testing.T) {
	// Create a new Naive Bayes Classifier
	threshold := 1.1
	nb := New(threshold)

	// Create a new training set
	trainingSet := map[string][]string{
		"Baseball": {
			"Pitcher",
			"Shortstop",
			"Outfield",
		},
		"Basketball": {
			"Point Guard",
			"Shooting Guard",
			"Small Forward",
			"Power Forward",
			"Center",
		},
		"Soccer": {
			"Goalkeeper",
			"Defender",
			"Midfielder",
			"Forward",
		},
	}

	// Train the classifier
	nb.Fit(trainingSet)

	//then
	assert.Equal(t, nb.Classify("Point guard"), "Basketball")
}
