package simulate

import "math"

type RiskDistribution struct {
	// riskyness distribution parameters:
	A, B float64
}

type AlphaDistribution struct {
	Mu, Std float64
}

type Intervention struct {
	Start    float64
	Duration float64
}

type RunType string

const (
	Unknown    RunType = ""
	Simulation RunType = "simulation"
	DifEq      RunType = "difeq"
	Difference RunType = "difference"
)

type RiskVariance string

const (
	LowVar    RiskVariance = "low"
	MediumVar RiskVariance = "medium"
	HighVar   RiskVariance = "high"
)

// function to convert from risk variance & mean to a and b.
func RiskDist(riskMean float64, riskVariance RiskVariance) *RiskDistribution {
	a := 1.0
	b := (1 - riskMean) / riskMean
	var factor float64
	if riskVariance == LowVar {
		factor = 2.0
	} else if riskVariance == MediumVar {
		factor = 1.0
	} else if riskVariance == HighVar {
		factor = 0.1
	}
	a, b = factor*a, factor*b
	return &RiskDistribution{A: a, B: b}
}

// The parameters to carry out a set of runs.
type Parameters struct {
	// Model dynamics:
	// Number of individuals:
	N int
	// chance of being infected per contact:
	AlphaC, AlphaR float64
	// disease lasts for this long before the individual recovers:
	DiseaseLength int
	// if not nil, contains information to construct riskyness distribution:
	RiskDist *RiskDistribution

	// More meta/computed stuff:
	RunType RunType
	R0      float64
	// RiskynessMean     float64
	// RiskynessVariance float64 // This could either be an enum (HIGH, MEDIUM, LOW), or need to standardize them

	// Number of identical simulations to run:
	Trials int
}

func AlphaR(R0 float64, R0c float64, meanP float64, N float64) float64 {
	AlphaR := (R0 - R0c) / meanP / meanP / N
	if math.IsNaN(AlphaR) {
		AlphaR = 0
	}
	return AlphaR
}

func CautionAlphaR(infectedFraction float64, alphaR float64) float64 {
	return math.Max(0, (0.5-infectedFraction)*alphaR)
}

// RESULTS

// Plots we want to make:
// 1. extinction probability
// 2. final Rs
// 3. max Is
// 4. dynamics over time
type Run struct {
	// always capture these:
	FinalR float64
	MaxI   float64

	// these are optional:
	Ts                  []float64 `json:",omitempty"`
	Is                  []float64 `json:",omitempty"`
	Rs                  []float64 `json:",omitempty"`
	Rts                 []float64 `json:",omitempty"`
	EffectiveAlphas     []float64 `json:",omitempty"`
	IRisks              []float64 `json:",omitempty"`
	SRisks              []float64 `json:",omitempty"`
	RiskyInfections     []float64 `json:",omitempty"`
	CommunityInfections []float64 `json:",omitempty"`
}

// One or multiple Runs with identical Parameters
type RunSet struct {
	Parameters Parameters
	Runs       []Run
}

// An R0 Series fixes a bunch of values and varies R0 systematically
type R0Series struct {
	RunType         RunType
	RiskMean        float64
	RiskVariance    RiskVariance
	HotspotFraction float64
	RunSets         []RunSet
}

// func (param Parameters) FileDescriptionLong() string {
// 	return fmt.Sprintf("T=%v,N=%v,ac=%v,ar=%v,dl=%v",
// 		// param.A,
// 		// param.B,
// 		param.Trials,
// 		param.N,
// 		param.AlphaC,
// 		param.AlphaR,
// 		param.DiseaseLength)
// }

// func (param Parameters) FileDescriptionExtinction() string {
// 	return fmt.Sprintf("extinction")
// }
