package api

import (
	"fmt"
	"sort"
)

func GetMaxPrefix(ss []string) (string, error) {
	if len(ss) == 0 {
		return "", fmt.Errorf("Slice ss is empty.")
	}
	if !sort.StringsAreSorted(ss) {
		return "", fmt.Errorf("Slice ss must be sorted.")
	}
	if ss[0] == "" {
		return "", fmt.Errorf("There is an empty string.")
	}
	first := ss[0]
	end := ss[len(ss)-1]
	var i int
	for ; i < len(first); i++ {
		if first[i] != end[i] {
			break
		}
	}
	if i == 0 {
		return "", fmt.Errorf("Have no same prefix.")
	}
	return first[:i], nil
}
