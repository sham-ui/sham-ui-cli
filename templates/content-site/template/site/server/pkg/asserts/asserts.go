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
	if strings.Contains(str, substr) {
		tb.Logf("\t\033[32;1m✔ \033[37;0mcontains substring\033[39m\n")
	} else {
		tb.Logf(
			"\t\033[31m✖ does not contain substring,"+
				"\n\n\texpected: \n\n\t\t\033[37;0m%s\033[39m, "+
				"\n\n\tcontained in: \n\n\t\t\033[37;0m%s\033[39m\n\n",
			substr,
			str,
		)
		tb.FailNow()
	}
}

// JSONEquals fails the test if exp is not equal to act.
func JSONEquals(tb testing.TB, exp, act string, messages ...string) {
	tb.Helper()
	message := strings.Join(messages, " ")

	var actMap any
	err := json.Unmarshal([]byte(act), &actMap)
	if nil != err {
		tb.Logf(
			"\t\033[31m✖ %s: can't unmarshal act string \033[39m\n\n",
			message,
		)
		tb.FailNow()
	}

	var expMap any
	err = json.Unmarshal([]byte(exp), &expMap)
	if nil != err {
		tb.Logf(
			"\t\033[31m✖ %s: can't unmarshal exp string \033[39m\n\n",
			message,
		)
		tb.FailNow()
	}

	if diff := deep.Equal(expMap, actMap); diff != nil {
		tb.Logf(
			"\t\033[31m✖ %s:JSON"+
				"\n\n\t\033[36mexp\033[39m: %#v"+
				"\n\n\t\033[36mgot\033[39m: %#v"+
				"\n\n\t\033[36mdiff\033[39m:\n\t\t%s\033[39m\n\n",
			message,
			expMap,
			actMap,
			strings.Join(diff, "\n\t\t"),
		)
		tb.FailNow()
	} else {
		tb.Logf("\t\033[32;1m✔ \033[37;0m%s:JSON\033[39m\n", message)
	}
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
