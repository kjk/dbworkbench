package main

import "testing"

func TestSanitizeUsage(t *testing.T) {
	s := `foo

bar`
	sGot := string(sanitizeUsage([]byte(s)))
	sExp := `foo
bar

`
	if sGot != sExp {
		t.Fatalf("sanitizeUsage() failed, for '%s' got:\n'%s'\nexpected:\n'%s'\n", s, sGot, sExp)
	}
}
