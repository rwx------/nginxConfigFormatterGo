package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestDecomposeLine(t *testing.T) {
	s1 := "aa;\n lskdfjl;\n # aa"
	t1 := []string{"aa;", " lskdfjl;", " # aa"}
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

	s3 := "{\nrewrite \"^(.*[;]{2,+})$\" /test.html;\n\n}\n"
	t3 := []string{`{`, `rewrite "^(.*[;]{2,+})$" /test.html;`, ``, `}`, ``}
	r3, f3 := decomposeLine(s3)
	if f3 != true || !stringSliceEqual(r3, t3) {
		t.Error(testFailedMessageBool(false, true, f3))
		t.Error(testFailedMessageString2Slice(s3, t3, r3))
	}
}

func TestCheackEveryChar(t *testing.T) {
	s1 := `{ rewrite "^/a/([\d]{2,}).html" /b/$1; } # here have qutoes(")`
	t1 := " {\n rewrite \"^/a/([\\d]___TEMPLATE_OPENING_TAG___2,___TEMPLATE_CLOSING_TAG___).html\" /b/$1;\n \n}\n \n# here have qutoes(\")"
	r1 := cheackEveryChar(s1)
	if r1 != t1 {
		t.Error(testFailedMessageString(s1, t1, r1))
	}

	s2 := `{ rewrite '^/a/([\d]{2,}).html' /b/$1; } # here have qutoes(')`
	t2 := " {\n rewrite '^/a/([\\d]___TEMPLATE_OPENING_TAG___2,___TEMPLATE_CLOSING_TAG___).html' /b/$1;\n \n}\n \n# here have qutoes(')"
	r2 := cheackEveryChar(s2)
	if r2 != t2 {
		t.Error(testFailedMessageString(s2, t2, r2))
	}

	s3 := `{ rewrite "^/a/([\d]{2,}).html" /b/$1; } # here have qutoes(") { test1 }`
	t3 := " {\n rewrite \"^/a/([\\d]___TEMPLATE_OPENING_TAG___2,___TEMPLATE_CLOSING_TAG___).html\" /b/$1;\n \n}\n \n# here have qutoes(\") { test1 }"
	r3 := cheackEveryChar(s3)
	if r3 != t3 {
		t.Error(testFailedMessageString(s3, t3, r3))
	}

	s4 := `{ rewrite "^/a/([\d]{2,}).html" /b/$1; } # here no qutoes { test1 }`
	t4 := " {\n rewrite \"^/a/([\\d]___TEMPLATE_OPENING_TAG___2,___TEMPLATE_CLOSING_TAG___).html\" /b/$1;\n \n}\n \n# here no qutoes { test1 }"
	r4 := cheackEveryChar(s4)
	if r4 != t4 {
		t.Error(testFailedMessageString(s4, t4, r4))
	}
}

func TestStripBracketTemplateTags(t *testing.T) {
	s1 := "{ rewrite \"^/a/([\\d]___TEMPLATE_OPENING_TAG___2,___TEMPLATE_CLOSING_TAG___).html\" /b/$1; }"
	t1 := `{ rewrite "^/a/([\d]{2,}).html" /b/$1; }`
	r1 := stripBracketTemplateTags(s1)
	if r1 != t1 {
		t.Error(testFailedMessageString(s1, t1, r1))
	}
}

func TestPerformIndentation(t *testing.T) {
	s1 := []string{
		"http {",
		"server {",
		"listen 80;",
		"",
		"# It's my domain:  liaoyongfu.com",
		"server_name www.liaoyongfu.com;",
		"",
		"location /nginx_status {",
		"stub_status on;",
		"allow 127.0.0.1;",
		"deny all;",
		"}",
		"}",
		"}",
	}
	b1 := 4
	t1 := []string{
		"http {",
		"    server {",
		"        listen 80;",
		"",
		"        # It's my domain:  liaoyongfu.com",
		"        server_name www.liaoyongfu.com;",
		"",
		"        location /nginx_status {",
		"            stub_status on;",
		"            allow 127.0.0.1;",
		"            deny all;",
		"        }",
		"    }",
		"}",
	}
	r1 := performIndentation(s1, b1)

	if !stringSliceEqual(t1, r1) {
		t.Error(testFailedMessageSlice(s1, t1, r1))
	}

	b2 := 2
	t2 := []string{
		"http {",
		"  server {",
		"    listen 80;",
		"",
		"    # It's my domain:  liaoyongfu.com",
		"    server_name www.liaoyongfu.com;",
		"",
		"    location /nginx_status {",
		"      stub_status on;",
		"      allow 127.0.0.1;",
		"      deny all;",
		"    }",
		"  }",
		"}",
	}
	r2 := performIndentation(s1, b2)

	if !stringSliceEqual(t2, r2) {
		t.Error(testFailedMessageSlice(s1, t2, r2))
	}
}

func TestJoinOpeningBracket(t *testing.T) {
	s1 := []string{
		"http {",
		"server",
		"",
		"",
		"{",
		"",
		"",
		"listen 80;",
		"",
		"# It's my domain:  {liaoyongfu.com}",
		"server_name www.liaoyongfu.com;",
		"",
		"location /nginx_status {",
		"stub_status on;",
		"allow 127.0.0.1;",
		"deny all;",
		"}",
		"}",
		"",
		"",
		"}",
	}
	t1 := []string{
		"http {",
		"server {",
		"listen 80;",
		"",
		"# It's my domain:  {liaoyongfu.com}",
		"server_name www.liaoyongfu.com;",
		"",
		"location /nginx_status {",
		"stub_status on;",
		"allow 127.0.0.1;",
		"deny all;",
		"}",
		"}",
		"}",
	}

	r1 := joinOpeningBracket(s1)

	if !stringSliceEqual(t1, r1) {
		t.Error(testFailedMessageSlice(s1, t1, r1))
	}
}

func TestCleanLines(t *testing.T) {
	s1 := []string{
		"http {",
		"server",
		"",
		"",
		"{",
		"",
		"",
		"listen 80;",
		"",
		"# location { allow all; }",
		"location { allow all; } # It's my domain:  {liaoyongfu.com}",
		"server_name www.liaoyongfu.com;",
		"",
		"location /nginx_status {stub_status on;allow 127.0.0.1;",
		"deny all;}",
		"",
		"}}",
	}
	t1 := []string{
		"http {",
		"server",
		"",
		"",
		"{",
		"",
		"",
		"listen 80;",
		"",
		"# location { allow all; }",
		"location  {",
		"allow all;",
		"",
		"}",
		"",
		"# It's my domain:  {liaoyongfu.com}",
		"server_name www.liaoyongfu.com;",
		"",
		"location /nginx_status  {",
		"stub_status on;",
		"allow 127.0.0.1;",
		"deny all;",
		"",
		"}",
		"",
		"",
		"}",
		"",
		"}",
		"",
	}
	r1 := cleanLines(s1)
	if !stringSliceEqual(t1, r1) {
		t.Error(testFailedMessageSlice(s1, t1, r1))
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
		eq := strings.Compare(a[i], b[i])
		if eq != 0 {
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
