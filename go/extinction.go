package main

import (
	"github.com/brendanwallace/riskySIR/simulate"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"golang.org/x/exp/rand"
	"time"
)

type TrialResult struct {
	Mu float64
	Std float64
	AlphaC float64
	AlphaR float64
	FinalRs []int
}

type ExtinctionResult struct {
	Description string
	TrialResults []TrialResult
}

func main() {
	risk_structure()
}

const N = 1000

func risk_structure() {


	// vary alpha c from 8 -> 0
	// match alpha r to keep R_0 = 8 (using a=2, b=6 for riskyness)

	rand.Seed(uint64(time.Now().UnixNano()))
	trials := 1000
	params := simulate.Parameters{
		AlphaC: 0,
		AlphaR: 0,
		DiseaseLength: 1,
		N: N,
		Trials: trials,
		ExtinctionShortcircuit: true,
	}


	R0 := 8.0


	for _, run := range []struct {
		description string
		a float64
		b float64
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
	} {
		resultName := fmt.Sprintf("extinctions_risk_p075_%v", run.description)
		fileName := resultName + ".json"
		fmt.Println("\nstarting simulation. will save output as:", fileName)

		result := ExtinctionResult{resultName, []TrialResult{}}

		params.RiskDist = &simulate.RiskDistribution{
			A: run.a,
			B: run.b,
		}


		for alphaC := R0; alphaC >= 0; alphaC -= 0.2 {
			fmt.Printf("\r%f/%f", alphaC, R0)

			alphaR := (R0 - alphaC)*(16.0/9.0)
			params.AlphaC = alphaC / N
			params.AlphaR = alphaR / N
			res := simulate.Run(params)
			result.TrialResults = append(
				result.TrialResults, TrialResult{0, 0, alphaC, alphaR, res.FinalRs})

		}


		// Output to appropriately named file
		file, jsonErr := json.MarshalIndent(result, "", "\t")
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}


		writeFileErr := ioutil.WriteFile("data/"+fileName, file, 0644)
		if writeFileErr != nil {
			log.Fatal(writeFileErr)
		}
	}
}

func hetero_alpha() {


	rand.Seed(uint64(time.Now().UnixNano()))

	trials := 1000

	params := simulate.Parameters{
		AlphaC: 0,
		AlphaR: 0,
		DiseaseLength: 1,
		R0c: -1,
		R0r: -1,
		R0: -1,
		N: 1000,
		Trials: trials,
		ExtinctionShortcircuit: true,
	}


	for _, mu := range []float64{2.0, 4.0, 6.0, 8.0} {



		resultName := fmt.Sprintf("extinctions_hetero_alphac_%v", mu)
		fileName := resultName + ".json"
		fmt.Println("\nstarting simulation. will save output as:", fileName)

		results := ExtinctionResult{resultName, []TrialResult{}}


		stdMax := 30.0
		for std := 0.1; std <= stdMax; std += 0.2 {
			fmt.Printf("\r%f/%f", std, stdMax)
			params.AlphaDist = &simulate.AlphaDistribution{mu, std}
			//results := simulate.Run(params)
			res := simulate.Run(params)
			results.TrialResults = append(
				results.TrialResults, TrialResult{mu, std, 0, 0, res.FinalRs})
		}

		//var results simulate.Results = simulate.Run(params)
		params = simulate.ComputeR0(params)
		

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
}
