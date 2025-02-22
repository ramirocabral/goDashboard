package utils

func RoundFloat(value float64, places int) float64 {
	shift := 1
	for i := 0; i < places; i++ {
		shift *= 10
	}
	return float64(int(value*float64(shift))) / float64(shift)
}
