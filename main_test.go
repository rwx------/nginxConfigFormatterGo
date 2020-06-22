package main

import (
	"testing"
)

func TestDecomposeLine(t *testing.T) {
	s1 := "aa; lskdfjl; "
	t1 := []string{"aa;", " lskdfjl;", " "}
	r1, f1 := decomposeLine(s1)
	if f1 != true || !stringSliceEqual(r1, t1) {
		t.Errorf(`Source: %#v%v`, s1, "\n")
		t.Errorf(`Target: %#v%v`, t1, "\n")
		t.Errorf(`Result: %#v%v`, r1, "\n")
		t.Errorf(`mFlag: %#v%v`, f1, "\n")
		t.Error(`Result is not match the Target`)
	}

	s2 := "aa;"
	t2 := []string{"aa;"}
	r2, f2 := decomposeLine(`aa;`)
	if f2 != false || !stringSliceEqual(r2, t2) {
		t.Errorf(`Source: %#v%v`, s2, "\n")
		t.Errorf(`Target: %#v%v`, t2, "\n")
		t.Errorf(`Result: %#v%v`, r2, "\n")
		t.Errorf(`mFlag: %#v%v`, f2, "\n")
		t.Error(`Result is not match the Target`)
	}

	s3 := `{rewrite "^(.*\'\"[;]{2,+})$" /test.html;}`
	t3 := []string{` {`, `rewrite "^(.*\'\"[;]{2,+})$" /test.html;`, ``, `}`, ``}
	r3, f3 := decomposeLine(s3)
	if f3 != true || !stringSliceEqual(r3, t3) {
		t.Errorf(`Source: %#v%v`, s3, "\n")
		t.Errorf(`Target: %#v%v`, t3, "\n")
		t.Errorf(`Result: %#v%v`, r3, "\n")
		t.Errorf(`mFlag: %#v%v`, f3, "\n")
		t.Error(`Result is not match the Target`)
	}
}

func TestAddNewLineString(t *testing.T) {
	s1 := "aa; lskdfjl; "
	t1 := "aa;\n lskdfjl;\n "
	r1 := addNewLineString(s1)
	if r1 != t1 {
		t.Errorf(`Source: %#v%v`, s1, "\n")
		t.Errorf(`Target: %#v%v`, t1, "\n")
		t.Errorf(`Result: %#v%v`, r1, "\n")
		t.Error(`Result is not match the Target`)
	}

}

func stringSliceEqual(a, b []string) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
