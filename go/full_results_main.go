package full_results_main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/brendanwallace/hotspot/simulate"
	"golang.org/x/exp/rand"
)

// Identical simulation run some number of times
// This can contain some unused/optional fields
type SingleResult struct {
	// Mu          float64
	// Std         float64
	// AlphaC      float64
	// AlphaR      float64
	// FinalRs     []int
	// DifEqFinalR float64

	Parameters        simulate.Parameters
	DifEqResults      simulate.DifEqResults
	DifferenceResults simulate.DifferenceResults
	SimulationResults simulate.Results
}

// Simulation results varied systematically along some axis
type SeriesResult struct {
	Description  string
	TrialResults []SingleResult
}

// Top level collections of series results
type FullResult struct {
	Homogeneous   []SeriesResult
	Heterogeneous []SeriesResult
	// Map from R0 values to the list of series results
	RiskStructured map[string][]SeriesResult
}

const N = 1000
const TRIALS = 100

func main() {
	results := FullResult{
		homo_alpha(),
		hetero_alpha(),
		map[string][]SeriesResult{
			"1.1": risk_structure(1.1),
			"1.2": risk_structure(1.2),
			"1.5": risk_structure(1.5),
			"2.0": risk_structure(2.0),
			"3.0": risk_structure(3.0),
			"4.0": risk_structure(4.0),
			"8.0": risk_structure(8.0)},
	}
	write(results)
}

func write(results FullResult) {
	fileName := fmt.Sprintf("full_results_small.json")

	// Output to appropriately named file
	file, jsonErr := json.MarshalIndent(results, "", "\t")
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	writeFileErr := ioutil.WriteFile("data/"+fileName, file, 0644)
	if writeFileErr != nil {
		log.Fatal(writeFileErr)
	}
}

func risk_structure(R0 float64) []SeriesResult {

	// vary alpha c from 8 -> 0
	// match alpha r to keep R_0 = 8 (using a=2, b=6 for riskyness)

	rand.Seed(uint64(time.Now().UnixNano()))
	params := simulate.Parameters{
		AlphaC:        0,
		AlphaR:        0,
		DiseaseLength: 1,
		N:             N,
		Trials:        TRIALS,
	}

	seriesResults := []SeriesResult{}

	for _, run := range []struct {
		description string
		a           float64
		b           float64
	}{
		{"0.75 high", 0.3, 0.1},
		{"0.75 medium", 3, 1},
		{"0.75 low", 6, 2},
		{"0.5 high", 0.1, 0.1},
		{"0.5 medium", 1, 1},
		{"0.5 low", 2, 2},
		{"0.25 high", 0.1, 0.3},
		{"0.25 medium", 1, 3},
		{"0.25 low", 2, 6},
		{"0.1 high", .1, 1.0},
		{"0.1 medium", 1, 10},
		{"0.1 low", 2, 20},
	} {
		resultName := fmt.Sprintf("%v %v", R0, run.description)
		fmt.Println("\nstarting series ", resultName)

		seriesResult := SeriesResult{resultName, []SingleResult{}}

		params.RiskDist = &simulate.RiskDistribution{
			A: run.a,
			B: run.b,
		}

		meanP := run.a / (run.a + run.b)
		for alphaC := R0; alphaC >= 0; alphaC -= 0.2 {
			fmt.Printf("\r%f/%f", alphaC, R0)

			// TODO: turn this into a top level function with unit tests:
			alphaR := (R0 - alphaC) / meanP / meanP
			params.AlphaC = alphaC / N
			params.AlphaR = alphaR / N
			seriesResult.TrialResults = append(
				seriesResult.TrialResults,
				SingleResult{
					Parameters:        params,
					DifEqResults:      simulate.RunDifEq(params),
					DifferenceResults: simulate.RunDifference(params),
					SimulationResults: simulate.RunSimulation(params)})

		}
		seriesResults = append(seriesResults, seriesResult)

	}
	return seriesResults

}

func hetero_alpha() []SeriesResult {

	rand.Seed(uint64(time.Now().UnixNano()))

	params := simulate.Parameters{
		AlphaC:        0,
		AlphaR:        0,
		DiseaseLength: 1,
		N:             1000,
		Trials:        TRIALS,
	}

	seriesResults := []SeriesResult{}

	for _, mu := range []float64{1.5, 2.0, 4.0, 8.0} {

		resultName := fmt.Sprintf("hetero_alphac_%v", mu)
		fmt.Println("\nstarting series", resultName)

		seriesResult := SeriesResult{resultName, []SingleResult{}}

		stdMax := 30.0
		for std := 0.1; std <= stdMax; std += 0.2 {
			fmt.Printf("\r%f/%f", std, stdMax)
			params.AlphaDist = &simulate.AlphaDistribution{Mu: mu, Std: std}

			seriesResult.TrialResults = append(
				seriesResult.TrialResults,
				SingleResult{
					Parameters: params,
					//DifEqResults:      simulate.RunDifEq(params),
					SimulationResults: simulate.RunSimulation(params)})
		}

		seriesResults = append(seriesResults, seriesResult)
	}

	return seriesResults

}

func homo_alpha() []SeriesResult {

	rand.Seed(uint64(time.Now().UnixNano()))

	params := simulate.Parameters{
		AlphaC:        0,
		AlphaR:        0,
		DiseaseLength: 1,
		N:             1000,
		Trials:        TRIALS,
	}

	resultName := "homo_alphac"
	fmt.Println("\nstarting series", resultName)

	seriesResult := SeriesResult{resultName, []SingleResult{}}
	for alphaC := 0.001; alphaC <= 0.008; alphaC += 0.001 {
		params.AlphaC = alphaC

		seriesResult.TrialResults = append(
			seriesResult.TrialResults,
			SingleResult{
				Parameters: params,
				//DifEqResults:      simulate.RunDifEq(params),
				SimulationResults: simulate.RunSimulation(params)})
	}

	return []SeriesResult{seriesResult}
}
