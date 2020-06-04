package array

import (
	"strings"
)

func CopyStringMap(source map[string]string, target map[string]string) map[string]string {
	if target == nil {
		target = make(map[string]string)
	}

	for k, v := range source {
		target[k] = v
	}
	return target
}

func MapIntIntKeysToSlice(mymap map[interface{}]interface{}) []int {
	keys := make([]int, len(mymap))
	i := 0
	for k := range mymap {
		keys[i] = k.(int)
		i++
	}
	return keys
}

func ReplaceAtIndex(input string, replacement string, index int) string {
	return strings.Join([]string{input[:index], replacement, input[index+1:]}, "")
}

func StringArrToInterface(arr []string) []interface{} {
	args := make([]interface{}, len(arr))
	i := 0
	for _, v := range arr {
		args[i] = v
		i++
	}
	return args
}

func SliceDifference(inslice1 []string, inslice2 []string, both bool) []string {
	var diff []string
	slice1 := make([]string, len(inslice1))
	slice2 := make([]string, len(inslice2))
	copy(slice1, inslice1)
	copy(slice2, inslice2)

	loopN := 1
	if both {
		loopN = 2
	}
	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < loopN; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
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
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

func InsliceString(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
