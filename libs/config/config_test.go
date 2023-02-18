package config

import (
	"testing"
)

const (
	ARG   = "ARG"
	UNSET = "UNSET"
)

func TestDefaultConfigure(t *testing.T) {
	env := GetDbEnvNames()

	got := Configure(env)

	for k, v := range DefaultEnvValues {
		if got[k] == v {
			continue
		}

		t.Fatalf("\nGot: %s\nExpected: %s\n", got, DefaultEnvValues)
	}
}

func TestSetDefaultParams(t *testing.T) {
	env := GetDbEnvNames()

	got := Configure(env)

	for k, v := range DefaultEnvValues {
		if got[k] == v {
			t.Fatalf("\nGot: %s\nExpected: %s\n", got, DefaultEnvValues)
		}
	}
}

func TestCustomVariableSetup(t *testing.T) {
	env := []string{ARG, UNSET}

	got := Configure(env)

	expected := map[string]string{
		ARG:   "value",
		UNSET: "",
	}

	for k, v := range got {
		if v == expected[k] {
			continue
		}
		t.Fatalf("\nGot: %s\nExpected: %s\n", got, expected)
	}
	t.Logf("\nGot: %s\nExpected: %s\n", got, expected)
}
