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
	FreeConf()
}

func TestSetDefaultParams(t *testing.T) {
	env := GetDbEnvNames()

	got := Configure(env)

	for k, v := range DefaultEnvValues {
		if got[k] == v {
			continue
		}
		t.Errorf("\nGot: %s\nExpected: %s\n", got, DefaultEnvValues)
	}
	FreeConf()
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
	FreeConf()
}

func TestSingleton(t *testing.T) {
	instance1 := Configure(GetDbEnvNames())
	instance1["DB_HOST"] = "172.0.2.1"

	instance2 := Configure(GetDbEnvNames())
	if instance1["DB_HOST"] != instance2["DB_HOST"] {
		t.Errorf("Expected value for instance2 to be %s but got %s instead", instance1["DB_HOST"], instance2["DB_HOST"])
	}

	instance2["DB_USER"] = "Skowe"
	if instance1["DB_USER"] != instance2["DB_USER"] {
		t.Errorf("Expected value for instance1 to be %s but got %s instead", instance2["DB_USER"], instance1["DB_USER"])
	}
	FreeConf()
}
