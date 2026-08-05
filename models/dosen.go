package models

type Dosen struct {
	ID            int    `json:"id"`
	Nama          string `json:"nama"`
	Username      string `json:"username"`
	Password      string `json:"password,omitempty"`
	MataKuliahIDs []int  `json:"mata_kuliah_ids"`
}