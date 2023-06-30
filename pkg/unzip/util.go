package unzip

func isArrayContaining(arr []string, val string) bool {
	for _, j := range arr {
		if j == val {
			return true
		}
	}

	return false
}
