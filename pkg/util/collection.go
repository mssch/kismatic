package util

// StringToInterfaceSlice converts slice of strings to a slice of interfaces
func StringToInterfaceSlice(slice []string) []interface{} {
	new := make([]interface{}, len(slice))
	for i, v := range slice {
		new[i] = v
	}
	return new
}

// Difference find the difference of 2 lists
func Difference(first []interface{}, second []interface{}) []interface{} {
	var diff []interface{}

	// Loop two times, first to find first strings not in second,
	// second loop to find second strings not in first
	for i := 0; i < 2; i++ {
		for _, s1 := range first {
			found := false
			for _, s2 := range second {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			first, second = second, first
		}
	}

	return diff
}

// Subset returns true if the first array is completely contained in the second array
func Subset(first, second []interface{}) bool {
	set := make(map[interface{}]int)
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
