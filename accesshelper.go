package jsonpath

import (
	"github.com/sirupsen/logrus"
	"strconv"
)

type access struct {
	t    byte // 0, object(map) 1, array
	cMap map[string]interface{}
	cArr []interface{}
}

const (
	ACCESS_OBJECT = 0
	ACCESS_ARRAY  = 1
)

func isCollection(collection interface{}) bool {
	switch collection.(type) {
	case map[string]interface{}:
		return true
	case []interface{}:
		return true
	}
	return false
}
func newaccess(collection interface{}) (*access, bool) {
	switch collection.(type) {
	case map[string]interface{}:
		return &access{t: 0, cMap: collection.(map[string]interface{})}, true
	case []interface{}:
		return &access{t: 1, cArr: collection.([]interface{})}, true
	}
	return nil, false
}

func (ah *access) keyExists(key string) bool {
	var (
		err error
		vi  int
	)
	switch ah.t {
	case ACCESS_OBJECT:
		if _, ok := ah.cMap[key]; ok {
			return true
		} else {
			return false
		}
	case ACCESS_ARRAY:
		if vi, err = strconv.Atoi(key); err != nil {
			logrus.WithFields(logrus.Fields{"package": "jsonpath", "action": "keyExists", "key": key}).Warn(err)
			return false
		} else {
			if vi+1 <= len(ah.cArr) {
				return true
			}
		}
	}
	return false
}

func (ah *access) getValue(key string) (interface{}, bool) {
	var (
		err error
		vi  int
	)
	switch ah.t {
	case ACCESS_OBJECT:
		if v, ok := ah.cMap[key]; ok {
			return v, true
		} else {
			return nil, false
		}
	case ACCESS_ARRAY:
		if vi, err = strconv.Atoi(key); err != nil {
			logrus.WithFields(logrus.Fields{"package": "jsonpath", "action": "getValue", "key": key, "access.t": ah.t, "access.map": ah.cMap, "access.arr": ah.cArr}).Warn(err)
			return nil, false
		} else {
			if vi+1 <= len(ah.cArr) {
				return ah.cArr[vi], true
			}
		}
	}
	return nil, false
}

func (ah *access) setValue(key string, value interface{}) bool {
	switch ah.t {
	case ACCESS_OBJECT:
		ah.cMap[key] = value
		return true
	case ACCESS_ARRAY:
		if vi, err := strconv.Atoi(key); err != nil {
			logrus.WithFields(logrus.Fields{"package": "jsonpath", "action": "setValue", "key": key}).Warn(err)
			return false
		} else {
			if vi+1 <= len(ah.cArr) {
				ah.cArr[vi] = value
				return true
			}
			// warning: do not use access setValue to expand collection length
			logrus.WithFields(logrus.Fields{"package": "jsonpath", "action": "set", "key": key, "access.t": ah.t, "access.map": ah.cMap, "access.arr": ah.cArr}).Warn(" could not expand arr length")
			return false
		}
	}
	return false
}

/*
TODO verify unset array
*/
func (ah *access) unsetValue(key string) bool {
	switch ah.t {
	case ACCESS_OBJECT:
		delete(ah.cMap, key)
		return true
	case ACCESS_ARRAY:
		if vi, err := strconv.Atoi(key); err != nil {
			logrus.WithFields(logrus.Fields{"package": "jsonpath", "action": "unsetValue", "key": key}).Warn(err)
			return false
		} else {
			if vi+1 < len(ah.cArr) {
				ah.cArr = append(ah.cArr[:vi], ah.cArr[vi+1:]...)
				return true
			}
			if vi+1 == len(ah.cArr) {
				ah.cArr = ah.cArr[:vi]
				return true
			}
			// warning: do not use access setValue to expand collection length
			logrus.WithFields(logrus.Fields{"package": "jsonpath", "action": "unset", "key": key, "ah": ah}).Warn(" could not expand arr length")
			return false
		}
	}
	return false
}

func (ah *access) arrayValues() []interface{} {
	switch ah.t {
	case ACCESS_OBJECT:
		var arr []interface{}
		for _, v := range ah.cMap {
			arr = append(arr, v)
		}
		return arr
	case ACCESS_ARRAY:
		return ah.cArr
	}
	return nil
}

func (ah *access) setAll(v interface{}) bool {
	switch ah.t {
	case ACCESS_OBJECT:
		for i, _ := range ah.cMap {
			ah.cMap[i] = v
		}
	case ACCESS_ARRAY:
		for i, _ := range ah.cArr {
			ah.cArr[i] = v
		}
	}
	return false
}
func (ah *access) unsetAll() bool {
	switch ah.t {
	case ACCESS_OBJECT:
		for i, _ := range ah.cMap {
			delete(ah.cMap, i)
		}
	case ACCESS_ARRAY:
		ah.cArr = nil
	}
	return false
}
