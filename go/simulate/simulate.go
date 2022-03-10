package simulate

import (
	//"errors"
	"fmt"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"math"
	"time"
)

// The parameters to run a simulation.
type Parameters struct {
	// riskyness distribution parameters: 
	A, B float64
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
}

// Measured outcomes from the run(s) as well as parameters.
type Results struct {
	Parameters Parameters
	Infecteds [][]int
	Susceptibles [][]int
	FinalRs []int
	InfectionEvents [][]*InfectionEvent
	// chance of being infected by risk taker if one takes a risk
	// at each time step, for each trial
	RiskyRisks [][]float64
	// expected number of risky secondary infections (per primary infection)
	// if p_r = 1, at each time step, for each trial
	ERtrs [][]float64
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


func initializePopulation(population []*Person, param Parameters) {
	var riskyParam func() float64

	if param.A == -1 && param.B == -1 {
		// special case: use no risk here.
		riskyParam = func() float64 {
			return 0
		}
	} else {// All of these use a beta distribution

		// TODO - consider the inverse of the CDF of beta distribution to sample
		// deterministically rather than probabalistically.
		// Probably not a big deal that we're seeding a new random source here -
		// seems way better than plumbing a *rand.Source all the way through.
		// Not sure what the downside would even be to making two rand.Source's.
		beta := distuv.Beta{param.A, param.B,
			rand.NewSource(uint64(time.Now().UnixNano()))}

		riskyParam = func() float64 {
			return beta.Rand()
		}
	}

	for p := range population {
		population[p] = &Person{SUSCEPTIBLE, 0, riskyParam(), nil}
	}
}



// Computes the various R0 values for a Parameters struct, returns
// a copy with these values set in it.
func ComputeR0(param Parameters) Parameters {
	// this one doesn't really change
	param.R0c = param.AlphaC * float64(param.N) * float64(param.DiseaseLength)
	//param.R0c = param.AlphaC * float64(param.N) / param.Beta


	if param.A == -1 && param.B == -1 {
		// special case: use no risk here.
		param.R0r = 0
	} else {
		// risky contacts if you take the risk is:
		// (number people * E[riskyness])
		// so expected number of risky contacts is:
		// E[riskyness] * (number people * E[riskyness])
		// and total R0r is this times alpha times disease length
		expectedRiskyness := distuv.Beta{param.A, param.B, nil}.Mean()
		param.R0r = (math.Pow(expectedRiskyness, 2) *
			param.AlphaR * float64(param.N) * float64(param.DiseaseLength))
			//param.AlphaR * float64(param.N * param.DiseaseLength))

	}

	// total is always just the sum
	param.R0 = param.R0c + param.R0r
	return param
}


// Computes the chance of becoming infected when taking the "risky" action in
// a day.
func computeRiskyRisk(population []*Person, param Parameters) float64 {
	c := 1.0
	for _, person := range population {
		if person.Status == INFECTED {
			c *= (1 - person.Riskyness * param.AlphaR)
		}
	}
	return 1 - c
}


func computeERtr(population []*Person, param Parameters) float64 {
	// Effective instantaneous Rt due to Riskyness in the population
	// Only considers pr in the susceptible population
	sumSRiskyness := 0.0
	numSusceptible := 0.0
	for _, person := range population {
		if person.Status == SUSCEPTIBLE {
			sumSRiskyness += person.Riskyness
			numSusceptible += 1.0
		}
	}
	// avoids dividing by zero:
	if numSusceptible == 0.0 {
		return 0.0
	}
	//return sumSRiskyness * param.AlphaR * (sumIRiskyness / numInfected) / param.Beta
	return sumSRiskyness * param.AlphaR * (sumSRiskyness / numSusceptible) * float64(param.DiseaseLength)
}


// Disease spreads within a subpopulation (possibly the whole population though)
// with contact rate * disease spread rate of alpha.
// Records any infection events into the appropriate Person struct.
func spreadWithin(
	population []*Person, alpha float64, isRisky bool, eventTime EventTime) {
	for p, person := range population {
		if (person.Status == INFECTED && person.daysInfected > 0) {
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


func Run(param Parameters) Results {

	// Saves the parameters used for this simulation along with the top level
	// results we care about.
	results := &Results{
		Parameters: param,
		Infecteds: make([][]int, param.Trials),
		Susceptibles: make([][]int, param.Trials),
		FinalRs: make([]int, param.Trials),
		InfectionEvents: make([][]*InfectionEvent, param.Trials),
		RiskyRisks: make([][]float64, param.Trials),
		ERtrs: make([][]float64, param.Trials),
	}

	// Conduct param.Trials discrete trials of the epidemic
	for i := 0; i < param.Trials; i++ {
		fmt.Printf("\r%v/%v", i, param.Trials)

		// Set up the population for the trial.
		var population []*Person = make([]*Person, param.N)
		initializePopulation(population, param)


		// Set up slices for each of the metrics to measure over the
		// course of the run
		infecteds := []int{}
		susceptibles := []int{}
		riskyRisks := []float64{}
		eRtrs := []float64{}


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
					// if rand.Float64() < param.Beta {
					// 	population[p].Status = RECOVERED
					// }
				}
			}
			// Incidentally, we use `infected` to track when the simulation
			// is complete too, so we have to update a variable here:
			infected = countStatus(population, INFECTED)
			infecteds = append(infecteds, infected)
			susceptibles = append(susceptibles, countStatus(population, SUSCEPTIBLE))
			riskyRisks = append(riskyRisks, computeRiskyRisk(population, param))
			eRtrs = append(eRtrs, computeERtr(population, param))
		}

		// The epidemic has run its course, so now we save the things we want
		// to save.
		results.FinalRs[i] = countStatus(population, RECOVERED)
		results.Infecteds[i] = infecteds
		results.Susceptibles[i] = susceptibles
		results.RiskyRisks[i] = riskyRisks
		results.ERtrs[i] = eRtrs
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



func (params Parameters) FileDescriptionLong() string {
	return fmt.Sprintf("A=%v,B=%v,T=%v,N=%v,ac=%v,ar=%v,dl=%v",
		params.A,
		params.B,
		params.Trials,
		params.N,
		params.AlphaC,
		params.AlphaR,
		params.DiseaseLength)
}