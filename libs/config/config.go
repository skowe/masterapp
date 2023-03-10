package config

import "os"

// Promenljive definisane u okviru ovog paketa će biti postavljene na svoje podrazumevane vrednosti
// ako vrednosti nisu postavljene u okruženju.
// ostale vrednosti iz isečka koje sluze za postacku će biti pročitane iz okruženja ili postavljene na prazan string

// Objekat konfiguracije iz okruzenaja prati singleton šablon kako bi se obezbedilo da konfiguracija jednom ucitana ostane ne promenljiva

var env map[string]string = nil

func Configure(envVars []string) map[string]string {

	if env == nil {
		env = make(map[string]string)
		for _, varName := range envVars {
			env[varName] = setVar(varName)
		}
	}
	return env
}

// Logika za postavku vrednosti za promenljive Okruženja
func setVar(varName string) string {
	val, ok := os.LookupEnv(varName)

	if ok {
		return val
	}

	switch varName {
	case DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASS:
		return DefaultEnvValues[varName]
	default:
		return ""
	}
}

func FreeConf() {
	env = nil
}
