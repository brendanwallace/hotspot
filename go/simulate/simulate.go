package simulate

import (
	//"errors"
	"fmt"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"math"
	"time"
)

type SimulationParameters struct {
	RiskynessDistribution RiskynessDistribution
	AlphaC float64
	AlphaR float64
	DiseaseLength int
	R0c float64
	R0r float64
	R0 float64
	N int
	Trials int
}

// Contains all the necessary parameters to run a simulation (so that these
// are recorded), as well as all the outcomes we mean to measure during the run.
type SimulationResults struct {
	Parameters SimulationParameters
	Infecteds [][]int
	FinalRs []int
	InfectionEvents [][]*InfectionEvent
}


// Infection Status enum.
type Status int
const (
	SUSCEPTIBLE = iota
	INFECTED
	RECOVERED
)


type RiskynessDistribution struct {
	Code string

	// parameters for the beta distribution
	A, B float64
}

// this one is a special case so we use a constant to check for it:
const NORISK = "norisk"

// This is the go way to make a constant slice.
func getRiskynessDistributions() []RiskynessDistribution {
	return []RiskynessDistribution{
		RiskynessDistribution{NORISK, -1, -1},
		RiskynessDistribution{"uniform", 1, 1},
		RiskynessDistribution{"uleft", 0.1, 0.2},
		RiskynessDistribution{"umid", 0.1, 0.1},
		RiskynessDistribution{"uright", 0.2, 0.1},
		RiskynessDistribution{"humpleft", 1.5, 3},
		RiskynessDistribution{"humpmid", 3, 3},
		RiskynessDistribution{"humpright", 3, 1.5},
	}
}


// Agent class for the simulation. Each of these represents a Person
// in the population, who can be S, I or R.
type Person struct {
	Status Status
	daysInfected int
	Riskyness float64
	// Info about the time of this Person's infection.
	// Nil if they were never infected.
	InfectionEvent *InfectionEvent
}

// There are two meaningful ways to measure time in the simulation: number
// of timesteps in, and the population state at the time.
type EventTime struct {
	Steps int
	Infected int
	Succeptible int
}

type InfectionEvent struct {
	EventTime EventTime
	InfectorRiskyness float64
	InfecteeRiskyness float64
	WasRiskyEvent bool
	SecondaryInfections int
}


func countStatus(population []*Person, status Status) int {
	count := 0
	for _, Person := range population {
		if Person.Status == status {
			count++
		}
	}
	return count
}


func initializePopulation(population []*Person, param SimulationParameters) {


	var riskyParam func() float64

	if r := param.RiskynessDistribution; r.Code == NORISK {
		riskyParam = func() float64 {
			return 0
		}
	} else {// All of these use a beta distribution

		// TODO - consider the inverse of the CDF of beta distribution to sample
		// deterministically rather than probabalistically.
		// Probably not a big deal that we're seeding a new random source here -
		// seems way better than plumbing a *rand.Source all the way through.
		// Not sure what the downside would even be to making two rand.Source's.
		beta := distuv.Beta{r.A, r.B, rand.NewSource(uint64(time.Now().UnixNano()))}

		riskyParam = func() float64 {
			return beta.Rand()
		}
	}

	for p := range population {
		population[p] = &Person{SUSCEPTIBLE, 0, riskyParam(), nil}
	}
}



// Computes the various R0 values for a SimulationParameters struct, returns
// a copy with these values set in it.
func computeR0(param SimulationParameters) SimulationParameters {
	// this one doesn't really change
	param.R0c = param.AlphaC * float64(param.N * param.DiseaseLength)

	if r := param.RiskynessDistribution; r.Code == NORISK {
		param.R0r = 0
	} else {
		// risky contacts if you take the risk is:
		// (number people * E[riskyness])
		// so expected number of risky contacts is:
		// E[riskyness] * (number people * E[riskyness])
		// and total R0r is this times alpha times disease length
		expectedRiskyness := distuv.Beta{r.A, r.B, nil}.Mean()
		param.R0r = (math.Pow(expectedRiskyness, 2) *
			param.AlphaR * float64(param.N * param.DiseaseLength))
	}

	// total is always just the sum
	param.R0 = param.R0c + param.R0r
	return param
}


// Disease spreads within a subpopulation (possibly the whole population though)
// with contact rate * disease spread rate of alpha.
// Records any infection events into the appropriate Person struct.
func spreadWithin(
	population []*Person, alpha float64, isRisky bool, eventTime EventTime) {
	for p, person := range population {
		if (person.Status == INFECTED) {
			for o, other := range population {
				if (o != p) && (other.Status == SUSCEPTIBLE) {
					if rand.Float64() < alpha {
						population[o].Status = INFECTED
						population[o].InfectionEvent = &InfectionEvent{
							EventTime: eventTime,
							InfectorRiskyness: person.Riskyness,
							InfecteeRiskyness: other.Riskyness,
							WasRiskyEvent: isRisky,
						}
						population[p].InfectionEvent.SecondaryInfections++
					}
				}
			}
		}
	}
}


func simulate(param SimulationParameters) SimulationResults {

	// Saves the parameters used for this simulation along with the top level
	// results we care about.
	results := &SimulationResults{
		param, make([][]int, param.Trials),
		make([]int, param.Trials),
		make([][]*InfectionEvent, param.Trials),
	}

	// Conduct param.Trials discrete trials of the epidemic
	for i := 0; i < param.Trials; i++ {
		fmt.Printf("\r%v/%v", i, param.Trials)

		// Set up the population for the trial.
		var population []*Person = make([]*Person, param.N)
		initializePopulation(population, param)


		// Measure the number of infecteds over the course of the run:
		infecteds := []int{}

		// Infect person 0 (should maybe vary this?)
		infect := rand.Int() % param.N
		population[infect].Status = INFECTED
		// Set some fields to -1, so we know this one was artificial
		population[infect].InfectionEvent = &InfectionEvent{
			EventTime{0, 0, param.N}, -1, population[infect].Riskyness, false, 0}

		// Time loop of the trial
		// The simulation continues until no-one is infected, but we increment
		// counter `t` to keep track of the number of steps as well.
		infected, t := countStatus(population, INFECTED), 0
		for ; infected > 0; t++ {

			// risky behavioral spread
			riskTakers := make([]*Person, 0)
			for _, Person := range population {
				if rand.Float64() < Person.Riskyness {
					riskTakers = append(riskTakers, Person)
				}
			}
			eventTime := EventTime{t, infected, countStatus(population, SUSCEPTIBLE)}

			spreadWithin(riskTakers, param.AlphaR, true, eventTime)

			// community spread
			spreadWithin(population, param.AlphaC, false, eventTime)

			// recovery
			for p := range population {
				if population[p].Status == INFECTED {
					population[p].daysInfected++
					if population[p].daysInfected >= param.DiseaseLength {
						population[p].Status = RECOVERED
					}
				}
			}
			infected = countStatus(population, INFECTED)
			infecteds = append(infecteds, infected)
		}

		// The epidemic has run its course, so now we save the things we want
		// to save.
		results.FinalRs[i] = countStatus(population, RECOVERED)
		results.Infecteds[i] = infecteds
		results.InfectionEvents[i] = []*InfectionEvent{}
		for _, person := range population {
			if person.InfectionEvent != nil {
				results.InfectionEvents[i] = append(
					results.InfectionEvents[i], person.InfectionEvent)
			}
		}

	}
	fmt.Printf("\r%v/%v\n", param.Trials, param.Trials)
	return *results
}


func (params SimulationParameters) fileDescription() string {
	return fmt.Sprintf("%v,T=%v,N=%v",
	params.RiskynessDistribution.Code,
	params.Trials,
	params.N)
}


func (params SimulationParameters) fileDescriptionLong() string {
	return fmt.Sprintf("%v,T=%v,N=%v,ac=%v,ar=%v,dl=%v",
		params.RiskynessDistribution.Code,
		params.Trials,
		params.N,
		params.AlphaC,
		params.AlphaR,
		params.DiseaseLength)
}