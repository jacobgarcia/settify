// Package helper implements helper functions for the tests.
package helper

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "update .golden files")

// SaveGoldenData saves test data in a golden file.
func SaveGoldenData(t *testing.T, name string, data []byte) {
	t.Helper()
	golden := filepath.Join("testdata", name+".golden")
	if *update {
		err := ioutil.WriteFile(golden, data, 0644)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
	}
}

// GetGoldenData gets data from a golden file.
func GetGoldenData(t *testing.T, name string) []byte {
	t.Helper()
	golden := filepath.Join("testdata", name+".golden")
	expected, err := ioutil.ReadFile(golden)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	return expected
}
