package main

// this is the code for -reverse flag
// reverse a string
func reverse(s string) string {
	rs := []rune(s)

	for i, j := 0, len(rs)-1; i < j; i, j = i+1, j-1 {
		rs[i], rs[j] = rs[j], rs[i]
	}

	return string(rs)
}
