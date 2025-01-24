package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/brendanwallace/hotspot/simulate"
	"golang.org/x/exp/rand"
)

const N = 1000
const TRIALS = 10
const DISEASE_PERIOD int = 1
const GAMMA float64 = 1.0 / float64(DISEASE_PERIOD)
const RUN_TYPE simulate.RunType = simulate.Simulation
const DATA_LOCATION = "../data/"

// Give read/write to user and read to all others.
const DATA_FILE_PERMISSIONS = 0644

func main() {
	// Start profiling
	f, err := os.Create("hotspot.prof")
	if err != nil {
		fmt.Println(err)
		return

	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	title := fmt.Sprintf("%s,D=%d,T=%d", RUN_TYPE, DISEASE_PERIOD, TRIALS)
	fmt.Println(title)
	var results = runR0Series(RUN_TYPE)
	write(results, title)
}

func write(results interface{}, filename string) {
	fileName := fmt.Sprintf("%s.json", filename)

	// Output to appropriately named file
	file, jsonErr := json.MarshalIndent(results, "", "\t")
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	writeFileErr := ioutil.WriteFile(DATA_LOCATION+fileName, file, DATA_FILE_PERMISSIONS)
	if writeFileErr != nil {
		log.Fatal(writeFileErr)
	}
}

func routeRun(runType simulate.RunType, param simulate.Parameters) simulate.RunSet {
	switch runType {
	case simulate.Simulation:
		return simulate.RunSimulation(param)
	case simulate.DifEq:
		return simulate.RunDifEq(param)
	case simulate.Difference:
		return simulate.RunDifference(param)
	default:
		return simulate.RunSet{}
	}
}

// Consider 9 different risk distributions
// For each one, vary R0 from 0.0 -> 8.0
func runR0Series(runType simulate.RunType) []simulate.R0Series {

	// vary R0 0 -> EndR0

	const EndR0 = 8.0
	const R0Step = 0.1

	rand.Seed(uint64(time.Now().UnixNano()))

	// These are typically [0.0, 0.25, 0.5, 0.75]
	hotspotFractions := []float64{0.0, 0.25, 0.5, 0.75}
	// For the main text, we show [0.5, 0.25, 0.125]. For the SI we show this wider range.
	riskMeans := []float64{0.75, 0.5, 0.25, 0.125, 0.06, 0.03}
	// simulate.LowVar, simulate.MediumVar, simulate.HighVar
	riskVariances := []simulate.RiskVariance{simulate.LowVar, simulate.MediumVar, simulate.HighVar}
	allSeries := []simulate.R0Series{}

	for hsf, hotspotFraction := range hotspotFractions {
		for rm, riskMean := range riskMeans {
			for rv, riskVariance := range riskVariances {
				series := simulate.R0Series{
					RunType:         runType,
					RiskMean:        riskMean,
					RiskVariance:    riskVariance,
					HotspotFraction: hotspotFraction,
					RunSets:         make([]simulate.RunSet, 0),
				}

				for R0 := 0.0; R0 <= EndR0; R0 += R0Step {

					fmt.Printf("\r hotspotfraction=%v/%v riskmean=%v/%v riskvar=%v/%v R0=%f",
						hsf+1, len(hotspotFractions),
						rm+1, len(riskMeans),
						rv+1, len(riskVariances),
						R0,
					)

					var alphaR float64
					if riskMean == 0 {
						alphaR = 0
					} else {
						alphaR = GAMMA * (R0 * hotspotFraction / riskMean / riskMean) / N
					}

					params := simulate.Parameters{
						AlphaC:        GAMMA * (R0 * (1 - hotspotFraction)) / N,
						AlphaR:        alphaR,
						DiseaseLength: DISEASE_PERIOD,
						N:             N,
						R0:            R0,
						Trials:        TRIALS,
						RiskDist:      simulate.RiskDist(riskMean, riskVariance),
					}

					series.RunSets = append(series.RunSets, routeRun(runType, params))

				}
				allSeries = append(allSeries, series)
			}
		}
	}
	return allSeries
}
