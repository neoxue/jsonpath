package jsonpath

import (
	"github.com/sirupsen/logrus"
	"strconv"
)

type accessarray struct {
	v []interface{}
}

func (ah *accessarray) get(key string) ([]interface{}, bool) {
	if key == "*" {
		return ah.v, true
	}

	var (
		vi  int
		err error
	)
	if vi, err = strconv.Atoi(key); err != nil {
		logrus.WithFields(logrus.Fields{"package": "jsonpath", "action": "get arr by string key", "key": key, "access.arr": ah.v}).Debug(err)
		return nil, false
	}
	if vi+1 <= len(ah.v) {
		return []interface{}{ah.v[vi]}, true
	}
	return nil, false
}

func (ah *accessarray) set(key string, data interface{}) bool {
	// TODO to be tested
	if key == "*" {
		for k := range ah.v {
			ah.v[k] = data
		}
		return true
	}
	var (
		vi  int
		err error
	)
	if vi, err = strconv.Atoi(key); err != nil {
		logrus.WithFields(logrus.Fields{"package": "jsonpath", "action": "set arr by string key", "key": key, "access.arr": ah.v}).Debug(err)
		return false
	}
	if vi <= len(ah.v) {
		ah.v[vi] = data
		return true
	}
	// warning: do not use access setValue to expand collection length by 2
	logrus.WithFields(logrus.Fields{"package": "jsonpath", "action": "set", "key": key, "access.arr": ah.v}).Debug(" could not expand arr length by 2")
	return false
}

func (ah *accessarray) unset(key string) bool {
	// TODO to be tested
	if key == "*" {
		ah.v = nil
	}
	var (
		vi  int
		err error
	)
	if vi, err = strconv.Atoi(key); err != nil {
		logrus.WithFields(logrus.Fields{"package": "jsonpath", "action": "set arr by string key", "key": key, "access.arr": ah.v}).Debug(err)
		return false
	}
	if vi+1 <= len(ah.v) {
		ah.v = append(ah.v[:vi], ah.v[vi+1:]...)
		return true
	}
	// warning: do not use access setValue to expand collection length
	logrus.WithFields(logrus.Fields{"package": "jsonpath", "action": "set", "key": key, "access.arr": ah.v}).Debug(" could not expand arr length by 2")
	return false
}

func (ah *accessarray) getByList(keys interface{}) ([]interface{}, bool) {
	var vs = []interface{}{}

	switch keys.(type) {
	case []int:
		for _, vi := range keys.([]int) {
			vs = append(vs, ah.v[vi])
		}
	case string:
		for _, key := range keys.([]string) {
			var vi int
			var err error
			if vi, err = strconv.Atoi(key); err != nil {
				logrus.WithFields(logrus.Fields{"package": "jsonpath", "action": "get arr by string key", "key": key, "access.arr": ah.v}).Debug(err)
				break
			}
			vs = append(vs, ah.v[vi])
		}
	}
	return vs, true
}
