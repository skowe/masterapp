package models

// ApiError poruka se šalje kao odgovor klijentu kada servis iz nekog razloga
// ne može da odgovori na zahtev očekivanom porukom.

// Funkcije koje vraćaju potrebne podatke treba da vrate pokazivač na ApiError objekat umesto error objekta
type ApiError struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	ErrorType    string `json:"error_type"`
	Info         string `json:"info,omitempty"`
}
