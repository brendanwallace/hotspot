package simulate

import (
	"math"
	"testing"
)

const tolerance = 0.0001
const N = 1000
const initialInfecteds = INITIAL_INFECTEDS

var defaultParameters Parameters = Parameters{
	RiskDist: &RiskDistribution{1, 1},
	// BetaDist:     nil,
	BetaC:         0,
	BetaR:         0,
	DiseaseLength: 1,
	N:             N,
	Trials:        1,
}

func TestRiskValue(t *testing.T) {
	for _, test := range []struct {
		b       int
		buckets int
		want    float64
	}{
		{b: 0, buckets: 10, want: 0.05},
		{9, 10, 0.95},
	} {
		got := riskValue(test.b, test.buckets)
		if math.Abs(got-test.want) > tolerance {
			t.Fatalf("riskValue(%v, %v) = %v; want %v",
				test.b, test.buckets, got, test.want)
		}

	}
}

func TestInitializePopulation(t *testing.T) {
	for _, p := range []struct {
		riskDist RiskDistribution
	}{
		{RiskDistribution{1, 1}},
		{RiskDistribution{2, 2}},
		{RiskDistribution{0.1, 0.1}},
		{RiskDistribution{1, 3}},
		{RiskDistribution{2, 6}},
		{RiskDistribution{0.1, 0.3}},
		{RiskDistribution{3, 1}},
		{RiskDistribution{6, 2}},
		{RiskDistribution{0.3, 0.1}},
	} {
		param := defaultParameters
		param.RiskDist = &p.riskDist
		S, I, R := InitializePopulations(param)
		totalS, totalI, totalR := 0.0, 0.0, 0.0
		for b := 0; b < BUCKETS; b++ {
			totalS += S[b]
			totalI += I[b]
			totalR += R[b]
		}
		if math.Abs(totalS-(N-initialInfecteds)) > tolerance {
			t.Fatalf("totalS %v != (N - initialInfecteds) %v, riskDist: %v",
				totalS, N-initialInfecteds, p.riskDist)
		}
		if math.Abs(totalI-initialInfecteds) > tolerance {
			t.Fatalf("totalI %v != initialInfecteds %v, riskDist: %v",
				totalI, initialInfecteds, p.riskDist)
		}
	}
}

func TestRunCommunity(t *testing.T) {

	tol := 0.1
	param := defaultParameters

	for _, test := range []struct {
		betaC float64
		want  float64
	}{
		{0.0, initialInfecteds},
		// {2.0, 796.8121},
		{8.0, 999.6636},
	} {
		param.BetaC = test.betaC / N
		results := RunDifEq(param).Runs[0]
		if math.Abs(results.FinalR-test.want) > tol {
			t.Fatalf("FinalR %v != %v; betaC = %v", results.FinalR, test.want, test.betaC)
		}
	}
}

func TestRunRisk(t *testing.T) {
	tol := 0.1
	param := defaultParameters

	for _, test := range []struct {
		betaC float64
		want  float64
	}{
		{0.0, initialInfecteds},
		// {2.0, 796.8121},
		{8.0, 999.6636},
	} {
		param.BetaC = test.betaC / N
		results := RunDifEq(param).Runs[0]
		if math.Abs(results.FinalR-test.want) > tol {
			t.Fatalf("FinalR %v != %v; betaC = %v", results.FinalR, test.want, test.betaC)
		}
	}
}

func TestNewInfectionsDifference(t *testing.T) {
	tol := 0.0001

	for _, test := range []struct {
		St     float64
		I0t    float64
		I1t    float64
		p      float64
		betaC  float64
		alphaR float64

		want float64
	}{
		{1000, 1, 0, 0, 0.001, 0.0, 1},
		{1000, 1, 1, 1, 0.0, 0.001, 1},
		{1000, 1, 1, 1, 0.0, 0.002, 2},
		{1000, 1, 1, 1, 1.0, 0, 1000},
	} {
		result := newInfectionsDifference(test.St, test.I0t, test.I1t, test.p, test.betaC, test.alphaR)
		if math.Abs(result-test.want) > tol {
			t.Fatalf("newInfectionsDifference(%v) = %v; want %v)", test, result, test.want)
		}
	}
}
