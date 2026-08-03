package models

type Dosen struct {
	ID       int    `json:"id"`
	Nama     string `json:"nama"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}