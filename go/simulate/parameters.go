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

// TODO -- function to convert from risk variance & mean to a and b.

// The parameters to carry out a set of runs.
type Parameters struct {
	// Model dynamics:
	// Number of individuals:
	N int
	// cchance of being infected per contact:
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

	// Stuff not really in use anymore:
	Caution      bool
	AlphaDist    *AlphaDistribution
	Intervention *Intervention
}

func AlphaR(R0 float64, R0c float64, meanP float64, N float64) float64 {
	AlphaR := (R0 - R0c) / meanP / meanP / N
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
	RunType RunType
	// always capture these:
	FinalR float64
	MaxI   float64

	// these are optional:
	Is                  []float64
	Rts                 []float64
	EffectiveAlphas     []float64
	IRisks              []float64
	SRisks              []float64
	RiskyInfections     []float64
	CommunityInfections []float64
}

// One or multiple Runs with identical Parameters
type RunSet struct {
	Parameters Parameters
	Runs       []Run
}

// An R0 Series fixes a bunch of values and varies R0 systematically
type R0Series struct {
	RiskMean             float64
	RiskVariance         RiskVariance
	ProblemPlaceFraction float64
	RunSets              []RunSet
}
