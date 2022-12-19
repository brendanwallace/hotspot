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

// The parameters to run a simulation.
type Parameters struct {
	// if not nil, contains information to construct riskyness distribution:
	RiskDist *RiskDistribution
	// if not nil, contains information to construct infectiousness distribution:
	AlphaDist *AlphaDistribution
	// disease parameters - chance of being infected per contact:
	AlphaC, AlphaR float64
	// disease lasts for this long before the individual recovers:
	DiseaseLength int
	// Beta float64
	// // computed instantaneous values:
	// R0c, R0r, R0 float64
	// Number of individuals:
	N int
	// Number of identical simulations to run:
	Trials       int
	Intervention *Intervention
	Caution      bool
}

func AlphaR(R0 float64, R0c float64, meanP float64, N float64) float64 {
	AlphaR := (R0 - R0c) / meanP / meanP / N
	return AlphaR
}

func CautionAlphaR(infectedFraction float64, alphaR float64) float64 {
	return math.Max(0, (0.5-infectedFraction)*alphaR)
}
