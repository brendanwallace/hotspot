package simulate

import (
	"testing"
)

func TestInfectionProbability(t *testing.T) {

	for _, test := range []struct {
		numContacts float64
		beta        float64
		want        float64
	}{
		{numContacts: 1, beta: 0.5, want: 0.5},
		{numContacts: 2, beta: 0.5, want: 0.75},
	} {
		got := infectionProbability(test.beta, test.numContacts)
		if got != test.want {
			t.Fatalf("infectionProbability(%v %v) = %v; want %v",
				test.numContacts, test.beta, got, test.want)
		}

	}
}
