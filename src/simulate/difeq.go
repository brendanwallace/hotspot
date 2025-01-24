package simulate

import (
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

// Generic differential equation parameters:
const DT = 0.01
const INITIAL_INFECTEDS = 1.0
const BUCKETS = 100
const END_THRESHOLD = 0.1

func sum(population []float64) float64 {
	total := 0.0
	for _, p := range population {
		total += p
	}
	return total
}

func firstMoment(population []float64, buckets int) float64 {
	moment := 0.0
	for b := 0; b < buckets; b++ {
		moment += riskValue(b, buckets) * population[b]
	}
	return moment
}

func riskValue(b int, buckets int) float64 {
	return (float64(b) + 0.5) / float64(buckets)
}

func InitializePopulations(param Parameters) ([]float64, []float64, []float64) {

	A, B := 1.0, 1.0
	if param.RiskDist != nil {
		A, B = param.RiskDist.A, param.RiskDist.B
	}
	beta := distuv.Beta{Alpha: A, Beta: B, Src: nil}

	S := make([]float64, BUCKETS)
	I := make([]float64, BUCKETS)
	R := make([]float64, BUCKETS)

	// initialize the susceptible population using the CDF
	for b := 0; b < BUCKETS; b++ {
		// each bucket should have this much mass in it cdf(x+1) - cdf(x)
		cumulative := beta.CDF(float64(b+1)/BUCKETS) - beta.CDF(float64(b)/BUCKETS)
		S[b] = float64(param.N) * cumulative
	}
	// move a total of INITIAL_INFECTEDS from S to I
	for b := 0; b < BUCKETS; b++ {
		I[b] = S[b] * (INITIAL_INFECTEDS / float64(param.N))
		S[b] = S[b] * (1 - (INITIAL_INFECTEDS / float64(param.N)))
	}
	return S, I, R
}

func RunDifEq(param Parameters) RunSet {

	S, I, R := InitializePopulations(param)
	// Gamma in the dif eq is the inverse of disease length:
	gamma := 1 / float64(param.DiseaseLength)
	Ts := []float64{}
	Is := []float64{}
	Rs := []float64{}
	Rts := []float64{}
	EffectiveAlphas := []float64{}
	IRisks, SRisks := []float64{}, []float64{}

	maxInfected := -1.0
	currentTime := 0.0

	for sumI := sum(I); sumI >= END_THRESHOLD; sumI = sum(I) {

		if sumI > maxInfected {
			maxInfected = sumI
		}
		// Need to compute these two updates:
		recoveries := make([]float64, BUCKETS)
		newInfections := make([]float64, BUCKETS)

		// Recoveries are straight forward
		for b := 0; b < BUCKETS; b++ {
			recoveries[b] = I[b] * gamma * DT
		}

		// sumI = sum(I)
		var rInfectionsThisGen, cInfectionsThisGen float64
		momentI := firstMoment(I, BUCKETS)

		// apply interventions like this:
		var alphaR float64 = param.AlphaR
		for b := 0; b < BUCKETS; b++ {
			communityInfections := S[b] * param.AlphaC * sumI * DT
			// if interventionInEffect {
			// 	communityInfections *= (1.0 / 2.0)
			// }

			riskyInfections := S[b] * riskValue(b, BUCKETS) * alphaR * momentI * DT

			newInfections[b] = communityInfections + riskyInfections

			// Save these to return per generation
			cInfectionsThisGen += communityInfections
			rInfectionsThisGen += riskyInfections
		}
		saveData := int(currentTime/DT)%10 == 0
		// Compute Rt
		if saveData {
			sumS := sum(S)
			momentS := firstMoment(S, BUCKETS)
			// alphaR is the version that possibly uses 'caution' and interventions
			Rt := sumS * (param.AlphaC + alphaR*(momentI/sumI)*(momentS/sumS))
			Rts = append(Rts, Rt)

			EffectiveAlpha := (param.AlphaC + alphaR*(momentI/sumI)*(momentS/sumS))
			EffectiveAlphas = append(EffectiveAlphas, EffectiveAlpha)

			IRisks = append(IRisks, (momentI / sumI))
			SRisks = append(SRisks, (momentS / sumS))
			Is = append(Is, sumI)
			Rs = append(Rs, sum(R))
			Ts = append(Ts, currentTime)
		}

		for b := 0; b < BUCKETS; b++ {
			S[b] -= newInfections[b]
			I[b] += newInfections[b]
			I[b] -= recoveries[b]
			R[b] += recoveries[b]
		}
		currentTime += DT
	}

	return RunSet{
		Parameters: param,
		Runs: []Run{
			Run{
				FinalR:          sum(R),
				MaxI:            maxInfected,
				Ts:              Ts,
				Is:              Is,
				Rs:              Rs,
				Rts:             Rts,
				EffectiveAlphas: EffectiveAlphas,
				IRisks:          IRisks,
				SRisks:          SRisks,
				Duration:        computeOutbreakDuration(Is, param),
				PeakTime:        computePeakTime(Is, param),
			},
		},
	}
}

// This is kind of complicated.
// For a differential equation, we can get away with - alpha * S * I; but for
// a difference equation this overshoots especially when S and I are both kind
// of big.
// Instead we have to do something more like S * (1 - alpha)**I
func newInfectionsDifference(St, I0t, I1t, p, alphaC, alphaR float64) float64 {
	return St * (1 - math.Pow((1-alphaC), I0t)*(1-p+p*math.Pow((1-alphaR), I1t)))
}

func RunDifference(param Parameters) RunSet {

	S, I, R := InitializePopulations(param)
	Is := []float64{}
	Rs := []float64{}

	if param.DiseaseLength != 1.0 {
		panic("disease length needs to be exactly 1.0 here")
	}

	maxInfected := -1.0
	for sumI := sum(I); sumI >= END_THRESHOLD; sumI = sum(I) {

		Is = append(Is, sumI)
		Rs = append(Rs, sum(R))
		if sumI > maxInfected {
			maxInfected = sumI
		}

		// New infections are the only thing we have to compute: everyone recovers.
		newInfections := make([]float64, BUCKETS)

		// sumI = sum(I)
		momentI := firstMoment(I, BUCKETS)
		for b := 0; b < BUCKETS; b++ {
			risk := riskValue(b, BUCKETS)

			newInfections[b] = newInfectionsDifference(S[b], sumI, momentI, risk, param.AlphaC, param.AlphaR)
			// newInfections[b] = (S[b] * math.Pow((1-param.AlphaC), sumI) *
			// 	(1 - risk + risk*math.Pow((1-param.AlphaR), momentI)))
		}

		for b := 0; b < BUCKETS; b++ {
			S[b] -= newInfections[b]
			R[b] += I[b]
			I[b] = newInfections[b]
		}
	}

	return RunSet{
		Parameters: param,
		Runs: []Run{
			Run{
				FinalR: sum(R),
				MaxI:   maxInfected,
				Is:     Is,
				Rs:     Rs,
			},
		},
	}
}
