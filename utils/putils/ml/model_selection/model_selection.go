package modelselection

import (
	"math/rand"
)

func TrainTestSplit(dataset []interface{}, testSize float64) (train, test []interface{}) {
	for _, data := range dataset {
		if rand.Float64() > testSize {
			train = append(train, data)
		} else {
			test = append(test, data)
		}
	}
	return train, test
}
