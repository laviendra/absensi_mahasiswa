package models

type MataKuliah struct {
	ID       int    `json:"id"`
	Nama     string `json:"nama"`
	Kode     string `json:"kode"`
	KelasIDs []int  `json:"kelas_ids"`
}