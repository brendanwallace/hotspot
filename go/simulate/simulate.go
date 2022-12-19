package simulate

import (
	//"errors"
	"fmt"
	"time"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

const INITIAL_INFECTED = 1

// Measured outcomes from the run(s) as well as parameters.
type Results struct {
	Parameters Parameters
	FinalRs    []int
	MaxIs      []int

	// Tracks how many infections each individual caused
	SecondaryInfectionCounts [][]int
	// Is                  [][]int
	// RiskyInfections     [][]int
	// CommunityInfections [][]int
}

// Infection Status enum.
type Status int

const (
	SUSCEPTIBLE = iota
	INFECTED
	RECOVERED
)

// Agent class for the simulation. Each of these represents a Person
// in the population, who can be S, I or R.
type Person struct {
	Status       Status
	daysInfected int
	Riskyness    float64
	AlphaC       float64
	AlphaR       float64

	SecondaryInfectionCount int
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

func initializePopulation(population []*Person, param Parameters) {
	// Defaults that can get overwritten:
	var risky func() float64 = func() float64 { return 0 }
	var alphaC func() float64 = func() float64 { return param.AlphaC }
	var alphaR func() float64 = func() float64 { return param.AlphaR }

	if param.RiskDist != nil { // All of these use a beta distribution

		// TODO - consider the inverse of the CDF of beta distribution to sample
		// deterministically rather than probabalistically.
		// Probably not a big deal that we're seeding a new random source here -
		// seems way better than plumbing a *rand.Source all the way through.
		// Not sure what the downside would even be to making two rand.Source's.
		beta := distuv.Beta{
			Alpha: param.RiskDist.A,
			Beta:  param.RiskDist.B,
			Src:   rand.NewSource(uint64(time.Now().UnixNano())),
		}

		risky = func() float64 {
			return beta.Rand()
		}
	}

	if param.AlphaDist != nil {
		// we have fixed mu and sigma^2 we want but we need to compute parameters
		// a(lpha) and b(eta) for the gamma distribution
		// given a/b = mu and a/b^2 = sigma^2 we can solve for a and b:

		mu, std := param.AlphaDist.Mu, param.AlphaDist.Std

		a := mu * mu / (std * std)
		b := mu / (std * std)
		gamma := distuv.Gamma{
			Alpha: a,
			Beta:  b,
			Src:   rand.NewSource(uint64(time.Now().UnixNano()))}

		alphaC = func() float64 {
			// g := gamma.Rand()
			// fmt.Printf("%v\n", g)
			// return g
			return gamma.Rand() / float64(param.N)
		}
	}

	for p := range population {
		population[p] = &Person{
			Status:                  SUSCEPTIBLE,
			daysInfected:            0,
			Riskyness:               risky(),
			AlphaC:                  alphaC(),
			AlphaR:                  alphaR(),
			SecondaryInfectionCount: 0,
		}
	}
}

// Disease spreads within a subpopulation (possibly the whole population though)
// with contact rate * disease spread rate of alpha.
// Records any infection events into the appropriate Person struct.
func spreadWithin(
	population []*Person, isRisky bool) {
	for p, person := range population {
		if person.Status == INFECTED && person.daysInfected > 0 {
			for o, other := range population {
				if (o != p) && (other.Status == SUSCEPTIBLE) {
					alpha := person.AlphaC
					if isRisky {
						alpha = person.AlphaR
					}
					if rand.Float64() < alpha {
						population[o].Status = INFECTED
						population[p].SecondaryInfectionCount =
							population[p].SecondaryInfectionCount + 1
					}
				}
			}
		}
	}
}

func RunSimulation(param Parameters) Results {

	// Saves the parameters used for this simulation along with the top level
	// results we care about.
	results := &Results{
		Parameters: param,
		FinalRs:    make([]int, param.Trials),
		MaxIs:      make([]int, param.Trials),
		// Is:         make([][]int, param.Trials),
		SecondaryInfectionCounts: make([][]int, param.Trials),
	}

	// Conduct param.Trials discrete trials of the epidemic
	for i := 0; i < param.Trials; i++ {
		//fmt.Printf("\r%v/%v", i, param.Trials)

		// Set up the population for the trial.
		var population []*Person = make([]*Person, param.N)
		Is := []int{}
		initializePopulation(population, param)

		// Infect initial people:
		for infect := 0; infect < INITIAL_INFECTED; infect++ {
			population[infect].Status = INFECTED
		}

		// Time loop of the trial
		// The simulation continues until no-one is infected.
		maxInfected := -1
		for infected := 1; infected > 0; infected = countStatus(population, INFECTED) {
			if infected > maxInfected {
				maxInfected = infected
			}
			Is = append(Is, infected)

			// risky behavioral spread
			riskTakers := make([]*Person, 0)
			for _, Person := range population {
				if rand.Float64() < Person.Riskyness {
					riskTakers = append(riskTakers, Person)
				}
			}

			spreadWithin(riskTakers, true)

			// community spread
			spreadWithin(population, false)

			// recovery
			for p := range population {
				if population[p].Status == INFECTED {
					if population[p].daysInfected >= param.DiseaseLength {
						population[p].Status = RECOVERED
					}
					population[p].daysInfected++
					// if rand.Float64() < param.Beta {
					// 	population[p].Status = RECOVERED
					// }
				}
			}
		}

		// The epidemic has run its course, so now we save the things we want
		// to save.
		results.FinalRs[i] = countStatus(population, RECOVERED)
		results.MaxIs[i] = maxInfected
		infectionCounts := []int{}
		for _, person := range population {
			if person.Status == RECOVERED {
				infectionCounts = append(infectionCounts, person.SecondaryInfectionCount)
			}
		}
		results.SecondaryInfectionCounts[i] = infectionCounts
		// results.Is[i] = Is

	}
	// fmt.Printf("\r%v/%v\n", param.Trials, param.Trials)
	return *results
}

func (param Parameters) FileDescriptionLong() string {
	return fmt.Sprintf("T=%v,N=%v,ac=%v,ar=%v,dl=%v",
		// param.A,
		// param.B,
		param.Trials,
		param.N,
		param.AlphaC,
		param.AlphaR,
		param.DiseaseLength)
}

func (param Parameters) FileDescriptionExtinction() string {
	return fmt.Sprintf("extinction")
}
