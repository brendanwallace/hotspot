package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/brendanwallace/riskySIR/simulate"
	"golang.org/x/exp/rand"
)

// Identical simulation run some number of times
// This can contain some unused/optional fields
type Result struct {
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

const N = 1000
const TRIALS = 1000

func main() {

	var results [][]Result

	// fmt.Println("\n100 0")

	// results = extinction_results(1.0, 0.0)
	// write(results, "extinction_results_100_0")

	fmt.Println("\n75 25")

	results = extinction_results(.75, .25)
	write(results, "extinction_results_75_25")

	fmt.Println("\n50 50")

	results = extinction_results(.50, .50)
	write(results, "extinction_results_50_50")

	fmt.Println("\n25 75")

	results = extinction_results(.25, .75)
	write(results, "extinction_results_25_75")

}

func write(results [][]Result, filename string) {
	fileName := fmt.Sprintf("%s.json", filename)

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

// Spread contribution is set to 50/50
// Consider 9 different risk distributions
// For each one, vary R0 from 0.0 -> 8.0
func extinction_results(C float64, R float64) [][]Result {

	// vary R0 0 -> 8

	const EndR0 = 8.0

	rand.Seed(uint64(time.Now().UnixNano()))
	params := simulate.Parameters{
		AlphaC:        0,
		AlphaR:        0,
		DiseaseLength: 1,
		N:             N,
		Trials:        TRIALS,
	}

	results := [][]Result{}

	for _, run := range []struct {
		description string
		a           float64
		b           float64
	}{
		{"0.5 high", 0.1, 0.1},
		{"0.5 medium", 1, 1},
		{"0.5 low", 2, 2},
		{"0.25 high", 0.1, 0.3},
		{"0.25 medium", 1, 3},
		{"0.25 low", 2, 6},
		{"0.125 high", .1, 0.7},
		{"0.125 medium", 1, 7},
		{"0.125 low", 2, 14},
	} {
		resultName := fmt.Sprintf("%v", run.description)
		fmt.Println("\nstarting series ", resultName)

		series := []Result{}

		params.RiskDist = &simulate.RiskDistribution{
			A: run.a,
			B: run.b,
		}

		meanP := run.a / (run.a + run.b)
		for R0 := 0.0; R0 <= EndR0; R0 += 0.1 {
			fmt.Printf("\r%f/%f", R0, EndR0)

			params.AlphaC = (R0 * C) / N
			params.AlphaR = (R0 * R / meanP / meanP) / N
			series = append(
				series,
				Result{
					Parameters:        params,
					DifEqResults:      simulate.RunDifEq(params),
					DifferenceResults: simulate.RunDifference(params),
					SimulationResults: simulate.RunSimulation(params)})

		}
		results = append(results, series)

	}
	return results

}
