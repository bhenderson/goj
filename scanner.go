package goj

import ()

func dotSplit(s string) (strs []string) {
	start := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\\':
			i++
		case '.':
			if i != 0 {
				strs = append(strs, s[start:i])
			}
			if s[i+1] == '.' {
				strs = append(strs, s[i:i+2])
				i++
			}
			start = i + 1
		}
	}
	strs = append(strs, s[start:])
	return
}
