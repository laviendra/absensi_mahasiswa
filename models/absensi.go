package models

type Absensi struct {
	ID          int    `json:"id"`
	MahasiswaID int    `json:"mahasiswa_id"`
	Tanggal     string `json:"tanggal"`
	JamMasuk    string `json:"jam_masuk"`
	JamPulang   string `json:"jam_pulang"`
	Status      string `json:"status"`
	StatusKehadiran string `json:"status_kehadiran"`
}