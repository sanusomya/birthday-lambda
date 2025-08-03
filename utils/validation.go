package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func CorrectMonth(month string) string{
	return strings.ToLower(month)[0:3]
}

func CheckDates(date int8, month string) bool {
	calenderMap := map[string]int8{
		"jan":  31,
		"feb":  28,
		"mar":  31,
		"apr":  30,
		"may":  31,
		"jun": 30,
		"jul": 31,
		"aug":  31,
		"sep":  30,
		"oct":  31,
		"nov":  30,
		"dec":  31,
	}
	val, exists := calenderMap[month]
	if exists {
		if date < 0 || date > val {
			return false
		} else {
			return true
		}
	}
	return false

}

func ValidName(name string) bool {
	length := len(name) < 9
	namePattern := `^[a-zA-Z]+$`
	nameRegex := regexp.MustCompile(namePattern)
	return length && nameRegex.MatchString(name)
}

func ValidMobile(mobile int64) bool {
	numPattern := `^[0-9]+$`
	numRegex := regexp.MustCompile(numPattern)
	number := fmt.Sprintf("%v", mobile)
	return len(number) == 10 && numRegex.MatchString(number)
}
