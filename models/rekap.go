package models

type Rekap struct {
	Nama       string `json:"nama"`
	Hadir      int    `json:"hadir"`
	Terlambat int    `json:"terlambat"`
	TidakHadir int    `json:"tidak_hadir"`
}