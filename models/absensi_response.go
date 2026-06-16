package models

type AbsensiResponse struct {
	ID        int    `json:"id"`
	Nama      string `json:"nama"`
	Tanggal   string `json:"tanggal"`
	JamMasuk  string `json:"jam_masuk"`
	JamPulang string `json:"jam_pulang"`
	Status    string `json:"status"`
}