package models

type AbsensiHarian struct {
	AbsensiID   *int   `json:"absensi_id"`
	MahasiswaID int    `json:"mahasiswa_id"`
	Nama        string `json:"nama"`
	JamMasuk    string `json:"jam_masuk"`
	JamPulang   string `json:"jam_pulang"`
	Status      string `json:"status_kehadiran"`
}