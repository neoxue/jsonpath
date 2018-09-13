package jsonpath

type filterSlice struct {
	t *pathtoken
}

// TODO action set and unset
func (f *filterSlice) eval(action string, cv interface{}, optionalValue interface{}) ([]interface{}, bool) {
	ah := newaccessins(cv)
	length := 0
	switch cv.(type) {
	case []interface{}:
		length = len(cv.([]interface{}))
	}
	ks, _ := f.getIndexes(f.t.v.([]int), length)
	switch action {
	case actionFind:
		return ah.getByList(ks)
	case actionSet:
		return nil, false
	case actionUnset:
		return nil, false
	default:
		return nil, true
	}
}

func (f *filterSlice) getIndexes(v []int, length int) ([]int, bool) {
	start := v[0]
	end := v[1]
	step := v[2]
	if step < 1 {
		step = 1
	}
	if end <= 0 {
		end += length
	}
	if start < 0 {
		start += length
	}
	k := start
	ks := []int{k}

	for true {
		if k+step < end {
			k += step
			ks = append(ks, k)
		} else {
			break
		}
	}
	return ks, true
}
