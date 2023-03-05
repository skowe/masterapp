package models

type User struct {
	UID      int    `json:"uid,omitempty"`
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`

	// Šifra se prima u plaintext formatu kada stiže od klijent aplikacije
	// Šifra se po primanju hešuje i čuva u bazi
	// Tekst kojim se šifra hešuje je sačuvan u bazi ali se nigde ne šalje
	Password string `json:"password,omitempty"`
}
