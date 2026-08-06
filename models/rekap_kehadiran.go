package models

type RekapKehadiranItem struct {
	PertemuanID     int    `json:"pertemuan_id"`
	Tanggal         string `json:"tanggal"`
	StatusPertemuan string `json:"status_pertemuan"`
	NamaMataKuliah  string `json:"nama_mata_kuliah"`
	NamaDosen       string `json:"nama_dosen"`
	JamHadir        string `json:"jam_hadir"`
	StatusKehadiran string `json:"status_kehadiran"`
}