package config

import (
	"github.com/go-logr/logr/testr"
	stdOS "os"
	"path"
	"site/test_helpers/asserts"
	"strings"
	"testing"
)

type callsStore map[string]*functionCall

type functionsCallsStorage struct {
	callsByFunctionName callsStore
}

type fnArgs []interface{}
type fnCallArgs []fnArgs
type functionCall struct {
	calls fnCallArgs
}

func (storage *functionsCallsStorage) For(name string) *functionCall {
	call, ok := storage.callsByFunctionName[name]
	if !ok {
		call = &functionCall{
			calls: fnCallArgs{},
		}
		storage.callsByFunctionName[name] = call
	}
	return call
}

func (fc *functionCall) Add(args ...interface{}) {
	fc.calls = append(fc.calls, args)
}

func (fc *functionCall) Count() int {
	return len(fc.calls)
}

func (fc *functionCall) ArgsAt(index int) []interface{} {
	return fc.calls[index]
}

func newMockFunctionCalls() *functionsCallsStorage {
	return &functionsCallsStorage{
		callsByFunctionName: make(callsStore),
	}
}

type mockedFileInfo struct {
	stdOS.FileInfo
}

type mockedOS struct {
	reportErr bool
	calls     *functionsCallsStorage
}

func (mockedOS) IsNotExist(err error) bool { return stdOS.IsNotExist(err) }

func (m mockedOS) Stat(name string) (stdOS.FileInfo, error) {
	m.calls.For("Stat").Add(name)
	if m.reportErr {
		return nil, stdOS.ErrNotExist
	}
	return mockedFileInfo{}, nil
}
func (m mockedOS) Exit(code int) { stdOS.Exit(code) }

type mockedIOutil struct {
	calls *functionsCallsStorage
}

func (m mockedIOutil) WriteFile(filename string, data []byte, perm stdOS.FileMode) error {
	m.calls.For("WriteFile").Add(filename, data, perm)
	return nil
}

func TestCreateConfigIfNotExists(t *testing.T) {
	oldOs := os
	mos := &mockedOS{
		reportErr: true,
		calls:     newMockFunctionCalls(),
	}
	os = mos
	oldIoUtil := ioutil
	moutil := &mockedIOutil{
		calls: newMockFunctionCalls(),
	}
	ioutil = moutil
	configFilename := path.Join("testdata", "config.cfg")
	defer func() {
		os = oldOs
		ioutil = oldIoUtil
	}()

	LoadConfiguration(testr.New(t).V(1), configFilename)

	asserts.Equals(t, 1, mos.calls.For("Stat").Count(), "stat")
	asserts.Equals(t, 1, moutil.calls.For("WriteFile").Count(), "writeFile")
	asserts.Equals(t, configFilename, moutil.calls.For("WriteFile").ArgsAt(0)[0], "file name")
	asserts.Equals(t, []byte(strings.TrimSpace(defaultConfig)), moutil.calls.For("WriteFile").ArgsAt(0)[1], "content")
}

func TestNotCreateConfigIfExists(t *testing.T) {
	oldOs := os
	mos := &mockedOS{
		calls: newMockFunctionCalls(),
	}
	os = mos
	oldIoUtil := ioutil
	moutil := &mockedIOutil{
		calls: newMockFunctionCalls(),
	}
	ioutil = moutil
	configFilename := path.Join("testdata", "config.cfg")
	defer func() {
		os = oldOs
		ioutil = oldIoUtil
	}()

	LoadConfiguration(testr.New(t).V(1), configFilename)

	asserts.Equals(t, 1, mos.calls.For("Stat").Count(), "stat")
	asserts.Equals(t, 0, moutil.calls.For("WriteFile").Count(), "write file")
}

func TestReadConfig(t *testing.T) {
	configFilename := path.Join("testdata", "config.cfg")
	LoadConfiguration(testr.New(t).V(1), configFilename)

	asserts.Equals(t, server{Port: 3001, AllowedDomains: []string{"http://127.0.0.1:3000", "http://localhost:3000"}}, Server, "server")
}
