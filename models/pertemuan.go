package models

type Jadwal struct {
	ID           int    `json:"id"`
	DosenID      int    `json:"dosen_id"`
	KelasID      int    `json:"kelas_id"`
	MataKuliahID int    `json:"mata_kuliah_id"`
	Hari         string `json:"hari"`
	JamMulai     string `json:"jam_mulai"`

	// hasil join buat ditampilkan, bukan disimpan langsung
	NamaDosen      string `json:"nama_dosen,omitempty"`
	NamaKelas      string `json:"nama_kelas,omitempty"`
	NamaMataKuliah string `json:"nama_mata_kuliah,omitempty"`
}