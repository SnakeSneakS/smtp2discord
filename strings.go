package main

// TruncateAndSplit splits a string into a slice of strings with each element being at most x characters long.
// If the input string exceeds x characters, "..." is prefixed to the subsequent parts.
func TruncateAndSplit(input string, x int) []string {
	if x <= 3 {
		panic("x must be greater than 3 to accommodate '...'")
	}
	result := make([]string, 0)
	for len(input) > x {
		if len(result) == 0 {
			// Add the first part without "..."
			result = append(result, input[:x])
			input = input[x:]
		} else {
			// Add subsequent parts prefixed with "..."
			prefixLength := x - 3 // Reserve space for "..."
			result = append(result, "..."+input[:prefixLength])
			input = input[prefixLength:]
		}
	}
	// Add the remaining part
	if len(input) > 0 {
		if len(result) > 0 {
			result = append(result, "..."+input)
		} else {
			result = append(result, input)
		}
	}

	return result
}
