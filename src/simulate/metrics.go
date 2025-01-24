// Computes key metrics of an outbreak
package simulate

// Outbreak occurs if at least 5% are infected.
const OUTBREAK_THRESHOLD = 0.05

// Assumes that we exceed outbreak threshold and then go back below it
// once each.
func computeOutbreakDuration(Is []float64, param Parameters) float64 {

	duration := 0.0
	outbreakStarted := false

	var outbreakThreshold float64 = OUTBREAK_THRESHOLD * float64(param.N)
	for _, infected := range Is {
		// Count the number of times we meet or exceed the threshold until
		// dipping below.
		if float64(infected) >= outbreakThreshold {
			outbreakStarted = true
			duration += 1.0
		} else if outbreakStarted {
			break
		}
	}
	if param.RunType == DifEq {
		return duration * DT
	}
	return duration
}

func computePeakTime(Is []float64, param Parameters) float64 {
	peakTime := 0.0
	peakInfected := 0.0
	for t, infected := range Is {

		if infected > peakInfected {
			peakTime = float64(t)
			peakInfected = infected
		}
	}
	if param.RunType == DifEq {
		return peakTime * DT
	}
	return peakTime
}
