package helpers

func IsInvalid(value string) bool {
	if value == "" {
		return true
	}
	return false
}

func ByteEmpty(s []byte) bool {
	for _, v := range s {
		if v != 0 {
			return false
		}
	}
	return true
}

func Paginate(pageNum int, pageSize int, sliceLength int) (int, int) {
	start := (pageNum - 1) * pageSize

	if start > sliceLength {
		start = sliceLength
	}

	end := pageNum * pageSize
	if end > sliceLength {
		end = sliceLength
	}

	return start, end
}
