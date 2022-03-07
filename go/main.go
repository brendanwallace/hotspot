package main

import (
	//"errors"
	"github.com/brendanwallace/risky_nonrisky/simulate"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"golang.org/x/exp/rand"
	"simulate"
	"time"
)


func parseRiskyDistribution(flag string) (RiskynessDistribution, error) {
	for _, dist := range getRiskynessDistributions() {
		if flag == dist.Code {
			return dist, nil
		}
	}
	return RiskynessDistribution{}, fmt.Errorf("unsupported distribution %v\n", flag)
}

var distributionFlag = flag.String("D", "uniform", "riskyness distribution")
var NFlag = flag.Int("N", 1000, "number of people in simulation")
var trialsFlag = flag.Int("T", 1000, "times to run the simulation")
var alphaRFlag = flag.Float64("ar", 0.0004,
	"infectiousness parameter of risky behavior (default 0.0004)")
var alphaCFlag = flag.Float64("ac", 0.0001,
	"infectiousness parameter of community spread (default 0.0001)")

func main() {
	flag.Parse()
	rand.Seed(uint64(time.Now().UnixNano()))


	distribution, distFlagErr := parseRiskyDistribution(*distributionFlag)
	if distFlagErr != nil {
		fmt.Printf("%v\n", distFlagErr)
	}

	// Set up the parameters of the simulation
	params := SimulationParameters{
		RiskynessDistribution: distribution,
		AlphaC: *alphaCFlag,
		AlphaR: *alphaRFlag,
		DiseaseLength: 10,
		// these get added by computeR0 function:
		R0c: -1,
		R0r: -1,
		R0: -1,
		N: *NFlag,
		Trials: *trialsFlag,
	}
	params = computeR0(params)
	
	/////////////////////////////////////
	// Run the simulation
	/////////////////////////////////////
	var results SimulationResults = simulate(params)

	// Output to appropriately named file
	file, jsonErr := json.MarshalIndent(results, "", "\t")
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	filename := fmt.Sprintf("%v.json", params.fileDescription())

	fmt.Println(filename)

	writeFileErr := ioutil.WriteFile("data/"+filename, file, 0644)
	if writeFileErr != nil {
		log.Fatal(writeFileErr)
	}
}