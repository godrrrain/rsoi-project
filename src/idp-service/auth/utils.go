package auth

func Contains(s, substr string) bool {
	parts := Split(s, " ")
	for _, part := range parts {
		if part == substr {
			return true
		}
	}
	return false
}

func Split(s, sep string) []string {
	var result []string
	var current string
	for _, c := range s {
		if string(c) == sep {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
