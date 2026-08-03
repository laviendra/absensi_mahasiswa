package models

type Mahasiswa struct {
	ID          int    `json:"id"`
	Nama        string `json:"nama"`
	KelasID     *int   `json:"kelas_id"`
	NamaKelas   string `json:"nama_kelas,omitempty"`
	NamaJurusan string `json:"nama_jurusan,omitempty"`
	Status      string `json:"status"`
}