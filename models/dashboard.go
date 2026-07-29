package models

type Dashboard struct {
	TotalMahasiswa int `json:"total_mahasiswa"`
	HadirHariIni int `json:"hadir_hari_ini"`
	Terlambat int `json:"terlambat"`
	BelumAbsen int `json:"belum_absen"`
}