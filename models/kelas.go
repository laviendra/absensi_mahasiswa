package models

type Kelas struct {
	ID          int    `json:"id"`
	Nama        string `json:"nama"`
	JurusanID   int    `json:"jurusan_id"`
	NamaJurusan string `json:"nama_jurusan,omitempty"`
}