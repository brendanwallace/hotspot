package simulate

import (
	"testing"
	"math"
)

const floatTolerance = 0.001

var p0I Person = Person{INFECTED, 0, 0.0, nil}
var p05I Person = Person{INFECTED, 0, 0.5, nil}
var p1I Person = Person{INFECTED, 0, 1.0, nil}
var p1R Person = Person{RECOVERED, 0, 1.0, nil}
var p0S Person = Person{SUSCEPTIBLE, 0, 0.0, nil}
var p05S Person = Person{SUSCEPTIBLE, 0, 0.5, nil}
var p1S Person = Person{SUSCEPTIBLE, 0, 1.0, nil}

func TestComputeRiskyRisk(t *testing.T) {


	param := Parameters{}
	param.AlphaR = 0.1



	tests := []struct {
		name string
		population []*Person
		param Parameters
		want float64
	}{
		{"one person p0", []*Person{&p0I}, param, 0},
		{"one person p1", []*Person{&p1I}, param, 0.1},
		// 1 - (0.9)(0.9) = 0.19
		{"multiple p1s and p0s", []*Person{&p1I, &p1I, &p0I, &p0I}, param, 0.19},
		{"recovereds and susceptibles", []*Person{&p1R, &p1S, &p1S}, param, 0},
		// 1 - (0.95)(0.95) = 0.0975
		{"p05I", []*Person{&p1R, &p05I, &p05I}, param, 0.0975},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := computeRiskyRisk(tc.population, tc.param)
			if math.Abs(tc.want - got) > floatTolerance {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}

}



func TestComputeERtr(t *testing.T) {


	param := Parameters{}
	param.AlphaR = 0.1
	param.DiseaseLength = 10


	tests := []struct {
		name string
		population []*Person
		param Parameters
		want float64
	}{
		{"one person p0", []*Person{&p0S, &p1I}, param, 0},
		{"one person p1", []*Person{&p1S, &p1I}, param, 1},
		{"multiple", []*Person{&p1I, &p1S, &p0S, &p05S}, param, 1.5*(0.5)},
		{"multiple half pr I", []*Person{&p05I, &p1S, &p0S, &p05S}, param, 0.75},
		{"noone infected", []*Person{&p1S}, param, 1.0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := computeERtr(tc.population, tc.param)
			if math.IsNaN(got) || math.Abs(tc.want - got) > floatTolerance {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}

}