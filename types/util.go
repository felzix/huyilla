package types

func absInt64(i int64) int64 {
	if i >= 0 {
		return i
	} else {
		return -i
	}
}
