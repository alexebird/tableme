package tableme

import (
	"strconv"
)

func StringifyBool(b bool) string {
	return strconv.FormatBool(b)
}

func StringifyIntPtr(i *int64) string {
	return strconv.FormatInt(*i, 10)
}

func StringifyStringPtr(s *string) string {
	return WithEmptyStringDefault(s)
}

func StringifyString(s string) string {
	return s
}

func WithEmptyStringDefault(val *string) string {
	return WithDefault(val, "")
}

func WithDefault(val *string, defaultVal string) string {
	if val != nil {
		return *val
	} else {
		return defaultVal
	}
}
