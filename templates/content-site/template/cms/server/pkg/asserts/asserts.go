package asserts

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/go-test/deep"
)

// Equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act any, message ...string) {
	tb.Helper()
	equals(tb, exp, act, message, nil)
}

// EqualsIgnoreOrder fails the test if exp is not equal to act.
func EqualsIgnoreOrder(tb testing.TB, exp, act any, message ...string) {
	tb.Helper()
	equals(tb, exp, act, message, []byte{deep.FLAG_IGNORE_SLICE_ORDER})
}

// ErrorsEqual fails the test if expected error is not equal to actual error.
func ErrorsEqual(tb testing.TB, expected, actual error) {
	tb.Helper()
	var expStr, actStr string
	if expected != nil {
		expStr = expected.Error()
	}
	if actual != nil {
		actStr = actual.Error()
	}
	equals(tb, expStr, actStr, []string{"error"}, nil)
}

// NoError fails the test if err is not nil.
func NoError(tb testing.TB, err error) {
	tb.Helper()
	equals(tb, nil, err, []string{"error is nil"}, nil)
}

// Contains checks if a string contains a substring.
//
// tb: The testing.TB interface for logging and reporting test failures.
// str: The string to be checked.
// substr: The substring to be searched for.
//
// This function does not return anything.
func Contains(tb testing.TB, str, substr string) {
	tb.Helper()
	contains(tb, str, substr, false)
}

// NotContains checks if the given string does not contain a specific substring.
//
// Parameters:
// - tb: The testing.TB interface for logging and reporting test failures.
// - str: the string to be checked.
// - substr: the substring to check if it is not contained in the given string.
//
// This function does not return anything.
func NotContains(tb testing.TB, str, substr string) {
	tb.Helper()
	contains(tb, str, substr, true)
}

func contains(tb testing.TB, str, substr string, inverse bool) {
	tb.Helper()
	const (
		containsSubstring    = "contains substring"
		notContainsSubstring = "does not contain substring"
	)
	var successString string
	var failString string
	if inverse {
		successString = notContainsSubstring
		failString = containsSubstring
	} else {
		successString = containsSubstring
		failString = notContainsSubstring
	}
	if strings.Contains(str, substr) != inverse {
		tb.Logf("\t\033[32;1m✔\033[37;0m %s\033[39m\n", successString)
	} else {
		tb.Logf(
			"\t\033[31m✖ %s\033[39m,"+
				"\n\n\texpected: \n\n\t\t\033[37;0m%s\033[39m, "+
				"\n\n\tcontained in: \n\n\t\t\033[37;0m%s\033[39m\n\n",
			failString,
			substr,
			str,
		)
		tb.FailNow()
	}
}

// JSONEquals fails the test if exp is not equal to act.
func JSONEquals(tb testing.TB, exp, act string, messages ...string) {
	tb.Helper()

	var actMap any
	NoError(tb, json.Unmarshal([]byte(act), &actMap))

	var expMap any
	NoError(tb, json.Unmarshal([]byte(exp), &expMap))

	equals(tb, expMap, actMap, messages, nil)
}

func equals(tb testing.TB, expected, actual any, messages []string, flags []byte) {
	tb.Helper()
	convFlags := make([]any, len(flags))
	for i, f := range flags {
		convFlags[i] = f
	}
	message := strings.Join(messages, " ")
	if diff := deep.Equal(expected, actual, convFlags...); diff != nil {
		tb.Logf(
			"\t\033[31m✖ %s"+
				"\n\n\t\033[36mexp\033[39m: %#v"+
				"\n\n\t\033[36mgot\033[39m: %#v"+
				"\n\n\t\033[36mdiff\033[39m:\n\t\t%s\033[39m\n\n",
			message,
			expected,
			actual,
			strings.Join(diff, "\n\t\t"),
		)
		tb.FailNow()
	} else {
		tb.Logf("\t\033[32;1m✔ \033[37;0m%s\033[39m\n", message)
	}
}
