package models

type AbsensiKelasItem struct {
	MahasiswaID     int    `json:"mahasiswa_id"`
	Nama            string `json:"nama"`
	JamHadir        string `json:"jam_hadir"`
	StatusKehadiran string `json:"status_kehadiran"`
}

type SimpanAbsensiKelasRequest struct {
	MahasiswaID     int    `json:"mahasiswa_id" binding:"required"`
	StatusKehadiran string `json:"status_kehadiran" binding:"required"`
}