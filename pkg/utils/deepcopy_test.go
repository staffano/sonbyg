package utils

import (
	"testing"
)

type S1 struct {
	String1 string
	String2 string
	IntA    []int
}

func (s S1) Equal(o S1) bool {
	if s.String1 != o.String1 || s.String2 != o.String2 {
		return false
	}
	if len(s.IntA) != len(o.IntA) {
		return false
	}
	for i := range s.IntA {
		if s.IntA[i] != o.IntA[i] {
			return false
		}
	}
	return true
}

func createS1() *S1 {
	return &S1{String1: "Staffan", String2: "Olsson", IntA: []int{1, 2, 3, 4}}
}

type S2 struct {
	Struct1 S1
	Float1  float64
	Struct2 S1
}

func (s S2) Equal(o S2) bool {
	return s.Struct1.Equal(o.Struct1) && s.Struct2.Equal(o.Struct2) && s.Float1 == o.Float1
}

func createS2() *S2 {
	return &S2{Struct1: *createS1(), Float1: 1.2345, Struct2: *createS1()}
}

func TestDeepCopy(t *testing.T) {
	s1 := createS1()
	s2 := createS2()
	dc1 := DeepCopy(s1).(*S1)
	dc2 := DeepCopy(s2).(*S2)
	if !s1.Equal(*dc1) {
		t.Errorf("Deep Copy failed s1 != dc1  [%v] != [%v]", s1, dc1)
	}
	if !s2.Equal(*dc2) {
		t.Errorf("Deep Copy failed s2 != dc2  [%v] != [%v]", s2, dc2)
	}
}

func TestDeepCopyNonPointer(t *testing.T) {

	a := "Staffan"
	b := DeepCopy(a)
	if b != a {
		t.Errorf("DeepCopy failed b=%v", b)
	}
}
