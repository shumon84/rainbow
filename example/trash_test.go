package main

import (
	"testing"
)

type TestCase struct {
	Case   string
	Expect bool
}

var (
	D1Test = []TestCase{
		{"abc.ABC", true},
		{":@", false},
		{":@abc", false},
		{"abc:@", false},
	}

	D2Test = []TestCase{
		{"abc", true},
		{".abc", false},
		{"ab.c", true},
		{".ab.c", false},
	}

	D3Test = []TestCase{
		{"abc", true},
		{"abc.", false},
	}

	D4Test = []TestCase{
		{"ab.c", true},
		{"ab..c", false},
	}

	D5Test = []TestCase{
		{"abc", true},
		{"", false},
	}
)

func verify(t *testing.T, testCases []TestCase) {
	for _, testCase := range testCases {
		result := DomainVerify(testCase.Case)
		if result != testCase.Expect {
			t.Log(testCase)
			t.Error("failed test")
		}
	}
}

func TestDomainVerifyD1(t *testing.T) {
	verify(t, D1Test)
}

func TestDomainVerifyD2(t *testing.T) {
	verify(t, D2Test)
}

func TestDomainVerifyD3(t *testing.T) {
	verify(t, D3Test)
}

func TestDomainVerifyD4(t *testing.T) {
	verify(t, D4Test)
}

func TestDomainVerifyD5(t *testing.T) {
	verify(t, D5Test)
}
