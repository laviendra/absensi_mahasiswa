package models

type PertemuanAdmin struct {
	ID             int    `json:"id"`
	Tanggal        string `json:"tanggal"`
	Status         string `json:"status"`
	NamaDosen      string `json:"nama_dosen"`
	NamaMataKuliah string `json:"nama_mata_kuliah"`
	NamaKelas      string `json:"nama_kelas"`
	Hadir          int    `json:"hadir"`
	Terlambat      int    `json:"terlambat"`
	TidakHadir     int    `json:"tidak_hadir"`
	Izin           int    `json:"izin"`
	Sakit          int    `json:"sakit"`
}