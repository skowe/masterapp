package config

import "os"

// Promenljive definisane u okviru ovog paketa će biti postavljene na svoje podrazumevane vrednosti
// ako vrednosti nisu postavljene u okruženju.
// ostale vrednosti iz isečka koje sluze za postacku će biti pročitane iz okruženja ili postavljene na prazan string
func Configure(envVars []string) map[string]string {
	env := make(map[string]string)

	for _, varName := range envVars {
		env[varName] = setVar(varName)
	}
	return env
}

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
