package xstrings

import (
	"fmt"
	"strconv"
	"strings"
)

type constraints interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func SplitNumber[T constraints](s string, seqs string) ([]T, error) {
	fields := strings.FieldsFunc(s, func(r rune) bool { return strings.ContainsRune(seqs, r) })
	vresult := make([]T, 0, len(fields))
	for _, f := range fields {
		if f == "" {
			continue
		}
		var v T
		switch any(v).(type) {
		case int8, int16, int32, int64, int:
			x, err2 := strconv.ParseInt(f, 10, 64)
			if err2 != nil {
				return nil, err2
			}
			v = T(x)
		case uint8, uint16, uint32, uint64, uint:
			x, err2 := strconv.ParseUint(f, 10, 64)
			if err2 != nil {
				return nil, err2
			}
			v = T(x)
		default:
			return nil, fmt.Errorf("unsupported type")
		}
		vresult = append(vresult, v)
	}
	return vresult, nil
}
