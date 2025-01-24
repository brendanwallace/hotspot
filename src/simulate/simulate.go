package simulate

import (
	//"errors"

	"time"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

const INITIAL_INFECTED = 1

// If EXTINCTION_SHORTCUT = true, then we should end the simulation
// if R >= EXTINCTION_CUTOFF because we know we didn't go extinct.
const EXTINCTION_CUTOFF = 50
const EXTINCTION_SHORTCUT = false

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

	beta := distuv.Beta{
		Alpha: param.RiskDist.A,
		Beta:  param.RiskDist.B,
		Src:   rand.NewSource(uint64(time.Now().UnixNano())),
	}

	risky = func() float64 {
		return beta.Rand()
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

// Disease spreads within a subpopulation (possibly the whole population)
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
