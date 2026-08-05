package models

type Mahasiswa struct {
	ID          int    `json:"id"`
	Nama        string `json:"nama"`
	Nim         string `json:"nim"`
	KelasID     *int   `json:"kelas_id"`
	JurusanID   *int   `json:"jurusan_id"`
	NamaKelas   string `json:"nama_kelas,omitempty"`
	NamaJurusan string `json:"nama_jurusan,omitempty"`
	Status      string `json:"status"`
}