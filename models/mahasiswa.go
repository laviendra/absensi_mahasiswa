package models

type Mahasiswa struct {
	ID      int    `json:"id"`
	Nama    string `json:"nama"`
	Jurusan string `json:"jurusan"`
}