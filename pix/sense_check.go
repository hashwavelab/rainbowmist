package pix

func PriceSenseCheck(price float64) bool {
	if price < 0.0000000001 || price > 1000000000.0 {
		return false
	}
	return true
}
