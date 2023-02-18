package config

// Imena promenljivih okruženja treba da budu konstantna
const (
	DB_HOST = "DB_HOST"
	DB_PORT = "DB_PORT"
	DB_NAME = "DB_TEST"
	DB_USER = "DB_USER"
	DB_PASS = "DB_PASS"
)

// Realno ce sve biti izmenjeno ali su mi potrebne podrazumevane vrednosti kada sam lenj za testiranje
var DefaultEnvValues = map[string]string{
	DB_HOST: "localhost",

	// Podrazumevana baza je MySQL ili MariaDB pa se koristi njihov podrazumevani port
	DB_PORT: "3306",
	DB_NAME: "test",

	// Ovo bi trebalo da se zameni u aplikaciji
	// root nalog ima previše prava da bi se koristio za aplikaciju
	DB_USER: "root",
	DB_PASS: "password",
}

// Koristi da vrati niz podrazumevanih promenljivih koje se obicno koriste da se uspostavi konekcija sa bazom podataka
// Ako se imena menjaju u aplikaciji trebala bi da se napise posebna logika za komunikaciju sa bazom koja se ne oslanja na
// libs modul ali ostatak config paketa bi i dalje mogao da se koristi za brzu i izolovanu postavku
func GetDbEnvNames() []string {
	return []string{DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASS}
}
