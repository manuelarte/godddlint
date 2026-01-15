package simple

import "slices"

//godddlint:valueObject
type MySlice struct {
	s []int
}

func NewMySlice(s []int) MySlice {
	return MySlice{
		s: s, // want `VOX002: Maps/Slices Not Defensive Copied`
	}
}

func NewMySlice2(s []int) MySlice {
	return MySlice{
		s: slices.Clone(s),
	}
}
