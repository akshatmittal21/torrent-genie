package utils

import "testing"

type addTest struct {
	arg1, expected string
}

var addTests = []addTest{
	{"1024", "1 KB"},
	{"1048576", "1 MB"},
	{"1073741824", "1 GB"},
	{"5486", "5.36 KB"},
	{"2684354560", "2.5 GB"},
	{"", ""},
	{"dsdsd", ""},
}

func TestGetSize(t *testing.T) {
	t.Log("TestGetSize")
	for _, test := range addTests {
		if output := GetFileSize(test.arg1); test.expected != output {
			t.Errorf("Output %q not equal to expected %q", output, test.expected)
		}
	}

}
