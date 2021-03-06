package matchers

import (
	"fmt"
	"math"
)

type BeNumericallyMatcher struct {
	Comparator string
	CompareTo  []interface{}
}

func (matcher *BeNumericallyMatcher) Match(actual interface{}) (success bool, message string, err error) {
	if len(matcher.CompareTo) == 0 || len(matcher.CompareTo) > 2 {
		return false, "", fmt.Errorf("BeNumerically requires 1 or 2 CompareTo arguments.  Got:%s", formatObject(matcher.CompareTo))
	}
	if !isNumber(actual) {
		return false, "", fmt.Errorf("Expected a number, got:%s", formatObject(actual))
	}
	if !isNumber(matcher.CompareTo[0]) {
		return false, "", fmt.Errorf("Expected a number, got:%s", formatObject(matcher.CompareTo[0]))
	}
	if len(matcher.CompareTo) == 2 && !isNumber(matcher.CompareTo[1]) {
		return false, "", fmt.Errorf("Expected a number, got:%s", formatObject(matcher.CompareTo[0]))
	}
	switch matcher.Comparator {
	case "==", "~", ">", ">=", "<", "<=":
	default:
		return false, "", fmt.Errorf("Unknown comparator: %s", matcher.Comparator)
	}

	if isInteger(actual) {
		success = matcher.matchIntegers(toInteger(actual), toInteger(matcher.CompareTo[0]))
	} else if isUnsignedInteger(actual) {
		success = matcher.matchUnsignedIntegers(toUnsignedInteger(actual), toUnsignedInteger(matcher.CompareTo[0]))
	} else if isFloat(actual) {
		var secondOperand float64 = 1e-8
		if len(matcher.CompareTo) == 2 {
			secondOperand = toFloat(matcher.CompareTo[1])
		}
		success = matcher.matchFloats(toFloat(actual), toFloat(matcher.CompareTo[0]), secondOperand)
	} else {
		return false, "", fmt.Errorf("Failed to compare:%s\n%s:%s", formatObject(actual), matcher.Comparator, formatObject(matcher.CompareTo[0]))
	}

	if success {
		return true, formatMessage(actual, fmt.Sprintf("not to be %s", matcher.Comparator), matcher.CompareTo[0]), nil
	} else {
		return false, formatMessage(actual, fmt.Sprintf("to be %s", matcher.Comparator), matcher.CompareTo[0]), nil
	}
}

func (matcher *BeNumericallyMatcher) matchIntegers(actual, compareTo int64) (success bool) {
	switch matcher.Comparator {
	case "==", "~":
		return (actual == compareTo)
	case ">":
		return (actual > compareTo)
	case ">=":
		return (actual >= compareTo)
	case "<":
		return (actual < compareTo)
	case "<=":
		return (actual <= compareTo)
	}
	return false
}

func (matcher *BeNumericallyMatcher) matchUnsignedIntegers(actual, compareTo uint64) (success bool) {
	switch matcher.Comparator {
	case "==", "~":
		return (actual == compareTo)
	case ">":
		return (actual > compareTo)
	case ">=":
		return (actual >= compareTo)
	case "<":
		return (actual < compareTo)
	case "<=":
		return (actual <= compareTo)
	}
	return false
}

func (matcher *BeNumericallyMatcher) matchFloats(actual, compareTo, threshold float64) (success bool) {
	switch matcher.Comparator {
	case "~":
		return math.Abs(actual-compareTo) <= threshold
	case "==":
		return (actual == compareTo)
	case ">":
		return (actual > compareTo)
	case ">=":
		return (actual >= compareTo)
	case "<":
		return (actual < compareTo)
	case "<=":
		return (actual <= compareTo)
	}
	return false
}
