package main

type set[elem comparable] map[elem]struct{}

func (s set[elem]) insert(e elem) {
	s[e] = struct{}{}
}

func (s set[elem]) contains(e elem) bool {
	_, ok := s[e]
	return ok
}

func (s set[elem]) toArray() []elem {
	arr := []elem{}
	for k := range s {
		arr = append(arr, k)
	}
	return arr
}
