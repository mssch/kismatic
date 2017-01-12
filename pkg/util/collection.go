package util

// Subset returns true if the first array is completely contained in the second array
func Subset(first, second []string) bool {
	set := make(map[string]int)
	for _, value := range second {
		set[value]++
	}

	for _, value := range first {
		if count, found := set[value]; !found {
			return false
		} else if count < 1 {
			return false
		}
	}

	return true
}
