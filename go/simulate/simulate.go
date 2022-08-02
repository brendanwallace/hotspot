package simulate

import (
	//"errors"
	"fmt"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"time"
)

type RiskDistribution struct {
	// riskyness distribution parameters: 
	A, B float64
}

type AlphaDistribution struct {
	Mu, Std float64
}

// The parameters to run a simulation.
type Parameters struct {
	// if not nil, contains information to construct riskyness distribution:
	RiskDist *RiskDistribution
	// if not nil, contains information to construct infectiousness distribution:
	AlphaDist *AlphaDistribution
	// disease parameters - chance of being infected per contact:
	AlphaC, AlphaR float64
	// disease lasts for this long and then the individual recovers:
	DiseaseLength int
	// Beta float64
	// computed instantaneous values:
	R0c, R0r, R0 float64
	// Number of individuals:
	N int
	// Number of identical simulations to run:
	Trials int
	// If true, stop the simulation after hitting N/10 infected individuals
	ExtinctionShortcircuit bool
}

// Measured outcomes from the run(s) as well as parameters.
type Results struct {
	Parameters Parameters
	FinalRs []int
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
	Status Status
	daysInfected int
	Riskyness float64
	AlphaC float64
	AlphaR float64
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

	if param.RiskDist != nil {// All of these use a beta distribution

		// TODO - consider the inverse of the CDF of beta distribution to sample
		// deterministically rather than probabalistically.
		// Probably not a big deal that we're seeding a new random source here -
		// seems way better than plumbing a *rand.Source all the way through.
		// Not sure what the downside would even be to making two rand.Source's.
		beta := distuv.Beta{param.RiskDist.A, param.RiskDist.B,
			rand.NewSource(uint64(time.Now().UnixNano()))}

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
			a, b, rand.NewSource(uint64(time.Now().UnixNano()))}

		alphaC = func() float64 { 
			// g := gamma.Rand()
			// fmt.Printf("%v\n", g)
			// return g
			return gamma.Rand()/float64(param.N)
		}
	}

	for p := range population {
		population[p] = &Person{SUSCEPTIBLE, 0, risky(), alphaC(), alphaR()}
	}
}


// Disease spreads within a subpopulation (possibly the whole population though)
// with contact rate * disease spread rate of alpha.
// Records any infection events into the appropriate Person struct.
func spreadWithin(
	population []*Person, isRisky bool) {
	for p, person := range population {
		if (person.Status == INFECTED && person.daysInfected > 0) {
			for o, other := range population {
				if (o != p) && (other.Status == SUSCEPTIBLE) {
					alpha := person.AlphaC
					if isRisky {
						alpha = person.AlphaR
					}
					if rand.Float64() < alpha {
						population[o].Status = INFECTED
					}
				}
			}
		}
	}
}


func Run(param Parameters) Results {

	// Saves the parameters used for this simulation along with the top level
	// results we care about.
	results := &Results{
		Parameters: param,
		FinalRs: make([]int, param.Trials),
	}

	// Conduct param.Trials discrete trials of the epidemic
	for i := 0; i < param.Trials; i++ {
		//fmt.Printf("\r%v/%v", i, param.Trials)

		// Set up the population for the trial.
		var population []*Person = make([]*Person, param.N)
		initializePopulation(population, param)



		// Infect person 0 (should maybe vary this?)
		infect := rand.Int() % param.N
		population[infect].Status = INFECTED

		// Time loop of the trial
		// The simulation continues until no-one is infected.
		for infected := countStatus(population, INFECTED); infected > 0; {

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