package asserts

import (
	"encoding/json"
	"fmt"
	"github.com/go-test/deep"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// Assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// Ok fails the test if an err is not nil.
func Ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}, message string) {
	if diff := deep.Equal(exp, act); diff != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(
			"\t\033[31m✖ %s:%d:%s\n\n\t\033[36mexp\033[39m: %#v\n\n\t\033[36mgot\033[39m: %#v\n\n\t\033[36mdiff\033[39m:\n\t\t%s\033[39m\n\n",
			filepath.Base(file),
			line,
			message,
			exp,
			act,
			strings.Join(diff, "\n\t\t"),
		)
		tb.FailNow()
	} else {
		fmt.Printf("\t\033[32;1m✔ \033[37;0m%s\033[39m\n", message)
	}
}

func JSONEqualsWithoutSomeKeys(tb testing.TB, keysPaths []string, exp, act, message string) {
	var actMap interface{}
	err := json.Unmarshal([]byte(act), &actMap)
	if nil != err {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(
			"\t\033[31m✖ %s:%d:%s: can't unmarshal act string \033[39m\n\n",
			filepath.Base(file),
			line,
			message,
		)
		tb.FailNow()
	}

	var expMap interface{}
	err = json.Unmarshal([]byte(exp), &expMap)
	if nil != err {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(
			"\t\033[31m✖ %s:%d:%s: can't unmarshal exp string \033[39m\n\n",
			filepath.Base(file),
			line,
			message,
		)
		tb.FailNow()
	}

	RemoveKeysFromMap(keysPaths, actMap.(map[string]interface{}))
	RemoveKeysFromMap(keysPaths, expMap.(map[string]interface{}))

	if diff := deep.Equal(expMap, actMap); diff != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf(
			"\t\033[31m✖ %s:%d:%s:JSON without some keys\n\n\t\033[36mexp\033[39m: %#v\n\n\t\033[36mgot\033[39m: %#v\n\n\t\033[36mdiff\033[39m:\n\t\t%s\033[39m\n\n",
			filepath.Base(file),
			line,
			message,
			expMap,
			actMap,
			strings.Join(diff, "\n\t\t"),
		)
		tb.FailNow()
	} else {
		fmt.Printf("\t\033[32;1m✔ \033[37;0m%s:JSON without some keys\033[39m\n", message)
	}
}
