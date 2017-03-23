package util

// Subset returns true if the first slice is completely contained in the second slice
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

// Intersects returns true if any element from the first slice is in the second slice
func Intersects(first, second []string) bool {
	set := make(map[string]int)
	for _, value := range second {
		set[value]++
	}

	for _, value := range first {
		if _, found := set[value]; found {
			return true
		}
	}
	return false
}
