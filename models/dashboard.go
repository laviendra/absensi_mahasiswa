package models

type Dashboard struct {
	TotalMahasiswa   int `json:"total_mahasiswa"`
	TotalDosen       int `json:"total_dosen"`
	TotalKelas       int `json:"total_kelas"`
	PertemuanHariIni int `json:"pertemuan_hari_ini"`
}