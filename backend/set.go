package main

type set[elem comparable] map[elem]struct{}

func (s set[elem]) contains(e elem) bool {
	_, ok := s[e]
	return ok
}
