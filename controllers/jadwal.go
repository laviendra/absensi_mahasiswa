package controllers

import (
	"net/http"
	"strconv"

	"absensi-mahasiswa/database"
	"absensi-mahasiswa/models"

	"github.com/gin-gonic/gin"
)

func GetJadwal(c *gin.Context) {

	rows, err := database.DB.Query(`
		SELECT j.id, j.dosen_id, j.kelas_id, j.mata_kuliah_id, j.hari, j.jam_mulai,
			d.nama, k.nama, mk.nama
		FROM jadwal j
		JOIN dosen d ON d.id = j.dosen_id
		JOIN kelas k ON k.id = j.kelas_id
		JOIN mata_kuliah mk ON mk.id = j.mata_kuliah_id
		ORDER BY j.hari, j.jam_mulai
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var data []models.Jadwal

	for rows.Next() {

		var j models.Jadwal

		rows.Scan(
			&j.ID, &j.DosenID, &j.KelasID, &j.MataKuliahID, &j.Hari, &j.JamMulai,
			&j.NamaDosen, &j.NamaKelas, &j.NamaMataKuliah,
		)

		data = append(data, j)
	}

	c.JSON(http.StatusOK, data)
}

func CreateJadwal(c *gin.Context) {

	var j models.Jadwal

	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if j.DosenID == 0 || j.KelasID == 0 || j.MataKuliahID == 0 || j.Hari == "" || j.JamMulai == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Semua field jadwal wajib diisi"})
		return
	}

	result, err := database.DB.Exec(
		"INSERT INTO jadwal (dosen_id, kelas_id, mata_kuliah_id, hari, jam_mulai) VALUES (?, ?, ?, ?, ?)",
		j.DosenID, j.KelasID, j.MataKuliahID, j.Hari, j.JamMulai,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()

	c.JSON(http.StatusCreated, gin.H{"message": "Jadwal berhasil ditambahkan", "id": id})
}

func DeleteJadwal(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	result, err := database.DB.Exec("DELETE FROM jadwal WHERE id = ?", id)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": fkErrorOr(err, "Jadwal ini sudah punya riwayat pertemuan, tidak bisa dihapus")})
		return
	}

	rowAffected, _ := result.RowsAffected()

	if rowAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jadwal tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Jadwal berhasil dihapus"})
}

// GetJadwalSaya dipakai DOSEN buat lihat jadwal ngajarnya sendiri.
// dosen_id diambil dari token (lewat middleware), bukan dari input,
// biar dosen A nggak bisa lihat jadwal dosen B.
func GetJadwalSaya(c *gin.Context) {

	dosenID := c.GetInt("dosen_id")

	rows, err := database.DB.Query(`
		SELECT j.id, j.dosen_id, j.kelas_id, j.mata_kuliah_id, j.hari, j.jam_mulai,
			k.nama, mk.nama
		FROM jadwal j
		JOIN kelas k ON k.id = j.kelas_id
		JOIN mata_kuliah mk ON mk.id = j.mata_kuliah_id
		WHERE j.dosen_id = ?
		ORDER BY j.hari, j.jam_mulai
	`, dosenID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var data []models.Jadwal

	for rows.Next() {

		var j models.Jadwal

		rows.Scan(
			&j.ID, &j.DosenID, &j.KelasID, &j.MataKuliahID, &j.Hari, &j.JamMulai,
			&j.NamaKelas, &j.NamaMataKuliah,
		)

		data = append(data, j)
	}

	c.JSON(http.StatusOK, data)
}