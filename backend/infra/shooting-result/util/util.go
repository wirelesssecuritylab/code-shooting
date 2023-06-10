package util

func GetIndexesByNames(names []string, rows [][]string) []int {
	indexes := make([]int, len(names), len(names))
	for i := range names {
		found := false
		for j := range rows[0] {
			if names[i] == rows[0][j] {
				found = true
				indexes[i] = j
			}
		}
		if !found {
			indexes[i] = -1
		}
	}

	return indexes
}

func IsIndexesValid(indexes []int, names []string) bool {
	if len(indexes) != len(names) {
		return false
	}
	for _, index := range indexes {
		if index == -1 {
			return false
		}
	}
	return true
}

func IsRowValid(indexes []int, lenOfRow int) bool {
	for _, index := range indexes {
		if index >= lenOfRow {
			return false
		}
	}

	return true
}
