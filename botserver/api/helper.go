package api

func contains(s []string, e string) int {
	for index, a := range s {
		if a == e {
			return index
		}
	}
	return -1
}

func containsInt(s []int, e int) int {
	for index, a := range s {
		if a == e {
			return index
		}
	}
	return -1
}

func containsInt64(s []int64, e int64) int {
	for index, a := range s {
		if a == e {
			return index
		}
	}
	return -1
}
