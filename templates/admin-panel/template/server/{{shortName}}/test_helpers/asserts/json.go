package asserts

import (
	"strconv"
	"strings"
)

func RemoveKeysFromMap(paths []string, data map[string]interface{}) map[string]interface{} {
	for _, path := range paths {
		removeKey(path, data)
	}
	return data
}

func removeKey(keyPath string, data interface{}) {
	owner := data
	keys := strings.Split(keyPath, ".")
	for i, key := range keys[:len(keys)-1] {
		switch own := owner.(type) {
		case map[string]interface{}:
			item, ok := own[key]
			if !ok {
				return
			}
			owner = item
		case []interface{}:
			if "$" == key {
				restKeys := strings.Join(keys[i+1:], ".")
				for _, item := range own {
					removeKey(restKeys, item)
				}
				return
			}
			index, err := strconv.Atoi(key)
			if nil != err {
				return
			}
			if index >= len(own) {
				return
			}
			owner = own[index]
		}
	}
	switch own := owner.(type) {
	case map[string]interface{}:
		delete(own, keys[len(keys)-1])
	}

}
