package controllers

import (
	"net/http"
	"strconv"
	"time"

	"absensi-mahasiswa/database"
	"absensi-mahasiswa/models"

	"github.com/gin-gonic/gin"
)

// GetPertemuanAdmin: admin lihat SEMUA pertemuan dari SEMUA dosen,
// lengkap rekap hadir/terlambat/tidak hadir/izin/sakit per pertemuan.
// Bisa difilter ?tanggal=YYYY-MM-DD.
func GetPertemuanAdmin(c *gin.Context) {

	tanggal := c.Query("tanggal")

	query := `
		SELECT p.id, p.tanggal, p.status,
			d.nama, mk.nama, k.nama,
			SUM(CASE WHEN ak.status_kehadiran = 'Hadir' THEN 1 ELSE 0 END),
			SUM(CASE WHEN ak.status_kehadiran = 'Terlambat' THEN 1 ELSE 0 END),
			SUM(CASE WHEN ak.status_kehadiran = 'Tidak Hadir' THEN 1 ELSE 0 END),
			SUM(CASE WHEN ak.status_kehadiran = 'Izin' THEN 1 ELSE 0 END),
			SUM(CASE WHEN ak.status_kehadiran = 'Sakit' THEN 1 ELSE 0 END)
		FROM pertemuan p
		JOIN jadwal j ON j.id = p.jadwal_id
		JOIN dosen d ON d.id = j.dosen_id
		JOIN mata_kuliah mk ON mk.id = j.mata_kuliah_id
		JOIN kelas k ON k.id = j.kelas_id
		LEFT JOIN absensi_kelas ak ON ak.pertemuan_id = p.id
	`

	var args []interface{}

	if tanggal != "" {
		query += " WHERE p.tanggal = ?"
		args = append(args, tanggal)
	}

	query += `
		GROUP BY p.id, p.tanggal, p.status, d.nama, mk.nama, k.nama
		ORDER BY p.tanggal DESC, p.id DESC
	`

	rows, err := database.DB.Query(query, args...)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var data []models.PertemuanAdmin

	for rows.Next() {

		var p models.PertemuanAdmin

		rows.Scan(
			&p.ID, &p.Tanggal, &p.Status,
			&p.NamaDosen, &p.NamaMataKuliah, &p.NamaKelas,
			&p.Hadir, &p.Terlambat, &p.TidakHadir, &p.Izin, &p.Sakit,
		)

		data = append(data, p)
	}

	c.JSON(http.StatusOK, data)
}

// GetAbsensiPertemuanAdmin: admin lihat daftar mahasiswa & status kehadiran
// di satu pertemuan tertentu, tanpa batasan harus dosen pemiliknya.
func GetAbsensiPertemuanAdmin(c *gin.Context) {

	pertemuanID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pertemuan tidak valid"})
		return
	}

	var kelasID int

	err = database.DB.QueryRow(`
		SELECT j.kelas_id
		FROM pertemuan p
		JOIN jadwal j ON j.id = p.jadwal_id
		WHERE p.id = ?
	`, pertemuanID).Scan(&kelasID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pertemuan tidak ditemukan"})
		return
	}

	rows, err := database.DB.Query(`
		SELECT m.id, m.nama,
			COALESCE(ak.jam_hadir, ''),
			COALESCE(ak.status_kehadiran, 'Tidak Hadir')
		FROM mahasiswa m
		LEFT JOIN absensi_kelas ak
			ON ak.mahasiswa_id = m.id AND ak.pertemuan_id = ?
		WHERE m.kelas_id = ? AND m.status = 'aktif'
		ORDER BY m.nama ASC
	`, pertemuanID, kelasID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var data []models.AbsensiKelasItem

	for rows.Next() {

		var item models.AbsensiKelasItem

		rows.Scan(&item.MahasiswaID, &item.Nama, &item.JamHadir, &item.StatusKehadiran)

		data = append(data, item)
	}

	c.JSON(http.StatusOK, gin.H{"pertemuan_id": pertemuanID, "data": data})
}

// UpdateAbsensiKelasAdmin: admin bisa koreksi kehadiran satu mahasiswa
// di pertemuan mana pun, termasuk yang sudah ditutup dosen (buat perbaikan).
func UpdateAbsensiKelasAdmin(c *gin.Context) {

	pertemuanID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pertemuan tidak valid"})
		return
	}

	var req models.SimpanAbsensiKelasRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jamHadir := ""

	if req.StatusKehadiran == "Hadir" || req.StatusKehadiran == "Terlambat" {
		jamHadir = time.Now().Format("15:04:05")
	}

	_, err = database.DB.Exec(`
		INSERT INTO absensi_kelas (pertemuan_id, mahasiswa_id, jam_hadir, status_kehadiran)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE jam_hadir = VALUES(jam_hadir), status_kehadiran = VALUES(status_kehadiran)
	`, pertemuanID, req.MahasiswaID, jamHadir, req.StatusKehadiran)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kehadiran berhasil dikoreksi"})
}
