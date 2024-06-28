package util

type Set[elem comparable] map[elem]struct{}

func (s Set[elem]) Insert(e elem) {
	s[e] = struct{}{}
}

func (s Set[elem]) Contains(e elem) bool {
	_, ok := s[e]
	return ok
}

func (s Set[elem]) ToArray() []elem {
	arr := []elem{}
	for k := range s {
		arr = append(arr, k)
	}
	return arr
}
