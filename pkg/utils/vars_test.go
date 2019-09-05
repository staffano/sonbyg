package utils

import (
	"testing"
)

func TestCopyVars(t *testing.T) {
	s1 := createS1()
	s2 := createS2()

	v := Variables{}
	v["S1"] = s1
	v["S2"] = s2

	v2 := v.Copy("S1", "S2")
	// v["S1"] contains *S1, therefore we need to fiddle with the pointers
	if !v2["S1"].(*S1).Equal(*(v["S1"].(*S1))) {
		t.Errorf("Value S1 not copied")
	}
	if !v2["S2"].(*S2).Equal(*(v["S2"].(*S2))) {
		t.Errorf("Value S2 not copied")
	}
}

func TestResolveVars(t *testing.T) {
	v := Variables{}
	v["A"] = "World"
	v["B"] = "Hello ${A}!"
	r := v.Resolve("${B}")
	if r != "Hello World!" {
		t.Errorf("resolved value should be 'Hello World!', is '%s'", r)
	}
}

func TestRecursiveVarsResolve(t *testing.T) {
	v := Variables{}
	v["A"] = "World"
	v["B"] = "Hello ${A}!"
	v["C"] = "I say '${B}'"
	r := v.Resolve("${C}")
	if r != "I say 'Hello World!'" {
		t.Errorf("resolved value should be [I say 'Hello World!'], is '%s'", r)
	}
}

func TestResolveNonStrings(t *testing.T) {
	v := Variables{}
	v["A"] = 666
	v["B"] = 123.23
	v["C"] = []int{5, 7, 9}
	v["D"] = "${A}:${B}:${C}"
	r := v.Resolve("${D}")
	if r != "666:123.23:[5 7 9]" {
		t.Errorf("Wrong result for resolve: %v", r)
	}
}

func TestResolveStructs(t *testing.T) {
	v := Variables{}
	s1 := *createS1()
	s1.String1 = "${HELLO}"
	v["A"] = s1
	v["HELLO"] = "World"
	r := v.Resolve("${A}")
	if r != "{World Olsson [1 2 3 4]}" {
		t.Errorf("Wrong result for resolve: %v", r)
	}
}

func TestResolveAll(t *testing.T) {
	v := Variables{}
	v["A"] = "World"
	v["B"] = "Hello ${A}!"
	v["C"] = "I say '${B}'"
	v.ResolveAll()
	if v["C"].(string) != "I say 'Hello World!'" {
		t.Errorf("resolved value should be [I say 'Hello World!'], is '%s'", v["C"])
	}
}
