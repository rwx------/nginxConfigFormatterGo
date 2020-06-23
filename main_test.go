package main

import (
	"fmt"
	"testing"
)

func TestDecomposeLine(t *testing.T) {
	s1 := "aa; lskdfjl; "
	t1 := []string{"aa;", " lskdfjl;", " "}
	r1, f1 := decomposeLine(s1)
	if f1 != true || !stringSliceEqual(r1, t1) {
		t.Error(testFailedMessageBool(false, true, f1))
		t.Error(testFailedMessageString2Slice(s1, t1, r1))
	}

	s2 := "aa;"
	t2 := []string{"aa;"}
	r2, f2 := decomposeLine(`aa;`)
	if f2 != false || !stringSliceEqual(r2, t2) {
		t.Error(testFailedMessageBool(false, false, f2))
		t.Error(testFailedMessageString2Slice(s2, t2, r2))
	}

	s3 := `{rewrite "^(.*\'\"[;]{2,+})$" /test.html;}`
	t3 := []string{` {`, `rewrite "^(.*\'\"[;]{2,+})$" /test.html;`, ``, `}`, ``}
	r3, f3 := decomposeLine(s3)
	if f3 != true || !stringSliceEqual(r3, t3) {
		t.Error(testFailedMessageBool(false, true, f3))
		t.Error(testFailedMessageString2Slice(s3, t3, r3))
	}
}

func TestAddNewLineString(t *testing.T) {
	s1 := "aa; lskdfjl; "
	t1 := "aa;\n lskdfjl;\n "
	r1 := addNewLineString(s1)
	if r1 != t1 {
		t.Error(testFailedMessageString(s1, t1, r1))
	}

	s2 := "aa;"
	t2 := "aa;"
	r2 := addNewLineString(s2)
	if r2 != t2 {
		t.Error(testFailedMessageString(s2, t2, r2))
	}

	s3 := `{rewrite "^(.*'\"[;]{2,+})$" /test.html;}`
	t3 := " {\nrewrite \"^(.*'\\\"[;]{2,+})$\" /test.html;\n\n}\n"
	r3 := addNewLineString(s3)
	if r3 != t3 {
		t.Error(testFailedMessageString(s3, t3, r3))
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

func testFailedMessageString(s, t, r string) string {
	err := fmt.Sprintf(`%vSource: %#v%v`, "\n", s, "\n")
	err += fmt.Sprintf(`Target: %#v%v`, t, "\n")
	err += fmt.Sprintf(`Result: %#v%v`, r, "\n")
	err += fmt.Sprintf(`Result is not match the Target`)
	return err
}

func testFailedMessageBool(s, t, r bool) string {
	err := fmt.Sprintf(`%vSource: %#v%v`, "\n", s, "\n")
	err += fmt.Sprintf(`Target: %#v%v`, t, "\n")
	err += fmt.Sprintf(`Result: %#v%v`, r, "\n")
	err += fmt.Sprintf(`Result is not match the Target`)
	return err
}

func testFailedMessageSlice(s, t, r []string) string {
	err := fmt.Sprintf(`%vSource: %#v%v`, "\n", s, "\n")
	err += fmt.Sprintf(`Target: %#v%v`, t, "\n")
	err += fmt.Sprintf(`Result: %#v%v`, r, "\n")
	err += fmt.Sprintf(`Result is not match the Target`)
	return err
}

func testFailedMessageString2Slice(s string, t, r []string) string {
	err := fmt.Sprintf(`%vSource: %#v%v`, "\n", s, "\n")
	err += fmt.Sprintf(`Target: %#v%v`, t, "\n")
	err += fmt.Sprintf(`Result: %#v%v`, r, "\n")
	err += fmt.Sprintf(`Result is not match the Target`)
	return err
}
