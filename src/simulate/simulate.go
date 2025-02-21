package simulate

import (
	//"errors"
	randv1 "golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"math"
	"math/rand/v2"
	"time"
)

const INITIAL_INFECTED = 1

// If EXTINCTION_SHORTCUT = true, then we should end the simulation
// if R >= EXTINCTION_CUTOFF because we know we didn't go extinct.
const EXTINCTION_CUTOFF = 50
const EXTINCTION_SHORTCUT = true

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
	Status        Status
	daysInfected  int
	RiskTolerance float64
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

	beta := distuv.Beta{
		Alpha: param.RiskDist.A,
		Beta:  param.RiskDist.B,
		Src:   randv1.NewSource(uint64(time.Now().UnixNano())),
	}

	for p := range population {
		population[p] = &Person{
			Status:        SUSCEPTIBLE,
			daysInfected:  0,
			RiskTolerance: beta.Rand(),
		}
	}
}

// The probability of getting infected when making numContacts contacts, if
// the infection rate per contact is beta.
func infectionProbability(beta float64, numContacts float64) float64 {
	return 1 - math.Pow(1.0-beta, numContacts)
}

// Disease spreads within a subpopulation (possibly the whole population)
// with contact rate * disease spread rate of beta.
func spreadWithin(population []*Person, beta float64) {
	var numInfected float64 = 0
	for _, person := range population {
		if person.Status == INFECTED && person.daysInfected > 0 {
			numInfected += 1
		}
	}

	var infectionProbability float64 = infectionProbability(beta, numInfected)
	for o, other := range population {
		if other.Status == SUSCEPTIBLE {
			if rand.Float64() < infectionProbability {
				population[o].Status = INFECTED
			}
		}
	}
}

func RunSimulation(param Parameters) RunSet {

	// Saves the parameters used for this simulation along with the top level
	// results we care about.
	runSet := RunSet{
		Parameters: param,
		Runs:       make([]Run, 0),
	}

	// Conduct param.Trials discrete trials of the epidemic
	for i := 0; i < param.Trials; i++ {
		//fmt.Printf("\r%v/%v", i, param.Trials)

		// Set up the population for the trial.
		var population []*Person = make([]*Person, param.N)
		Is := []float64{}
		initializePopulation(population, param)

		// Infect initial people:
		for infect := 0; infect < INITIAL_INFECTED; infect++ {
			population[infect].Status = INFECTED
		}

		// Set up timing measurements
		var time int = 0
		var peakTime float64 = 0

		// Time loop of the trial
		// The simulation continues until no-one is infected.
		maxInfected := 0
		for infected := 1; infected > 0; infected = countStatus(population, INFECTED) {
			// measure peak number of infections & timing
			if infected > maxInfected {
				maxInfected = infected
				peakTime = float64(time)
			}

			Is = append(Is, float64(infected))

			// Shortcut out if we only care about probability of extinction.
			if EXTINCTION_SHORTCUT {
				recovered := countStatus(population, RECOVERED)
				if recovered+infected >= EXTINCTION_CUTOFF {
					for p := range population {
						if population[p].Status == INFECTED {
							population[p].Status = RECOVERED
						}
					}
					break
				}
			}

			// risky behavioral spread
			riskTakers := make([]*Person, 0)
			for _, Person := range population {
				if rand.Float64() < Person.RiskTolerance {
					riskTakers = append(riskTakers, Person)
				}
			}

			spreadWithin(riskTakers, param.BetaR)

			// community spread
			spreadWithin(population, param.BetaC)

			// recovery
			for p := range population {
				if population[p].Status == INFECTED {
					if population[p].daysInfected >= param.DiseaseLength {
						population[p].Status = RECOVERED
					} else {
						population[p].daysInfected++
					}
				}
			}
			time += 1
		}

		// The epidemic has run its course, so now we save the things we want
		// to save.
		runSet.Runs = append(runSet.Runs, Run{
			FinalR:   float64(countStatus(population, RECOVERED)),
			MaxI:     float64(maxInfected),
			Duration: computeOutbreakDuration(Is, param),
			PeakTime: peakTime,
			// Is:       Is,
		})

	}
	return runSet
}
