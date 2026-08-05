package controllers

import (
	"net/http"
	"strconv"

	"absensi-mahasiswa/database"
	"absensi-mahasiswa/models"

	"github.com/gin-gonic/gin"
)

// getKelasIDsByMataKuliah ambil daftar id kelas yang pakai mata kuliah tsb
func getKelasIDsByMataKuliah(mataKuliahID int) []int {

	rows, err := database.DB.Query(
		"SELECT kelas_id FROM mata_kuliah_kelas WHERE mata_kuliah_id = ?",
		mataKuliahID,
	)

	ids := []int{}

	if err != nil {
		return ids
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		rows.Scan(&id)
		ids = append(ids, id)
	}

	return ids
}

// simpanRelasiMataKuliahKelas sinkronkan relasi: hapus yang lama, insert yang baru
func simpanRelasiMataKuliahKelas(mataKuliahID int, kelasIDs []int) error {

	_, err := database.DB.Exec("DELETE FROM mata_kuliah_kelas WHERE mata_kuliah_id = ?", mataKuliahID)

	if err != nil {
		return err
	}

	for _, kelasID := range kelasIDs {

		_, err := database.DB.Exec(
			"INSERT INTO mata_kuliah_kelas (mata_kuliah_id, kelas_id) VALUES (?, ?)",
			mataKuliahID, kelasID,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func GetMataKuliah(c *gin.Context) {

	rows, err := database.DB.Query("SELECT id, nama, kode FROM mata_kuliah ORDER BY nama ASC")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var data []models.MataKuliah

	for rows.Next() {

		var mk models.MataKuliah

		rows.Scan(&mk.ID, &mk.Nama, &mk.Kode)

		mk.KelasIDs = getKelasIDsByMataKuliah(mk.ID)

		data = append(data, mk)
	}

	c.JSON(http.StatusOK, data)
}

func CreateMataKuliah(c *gin.Context) {

	var mk models.MataKuliah

	if err := c.ShouldBindJSON(&mk); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if mk.Nama == "" || mk.Kode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan kode mata kuliah wajib diisi"})
		return
	}

	result, err := database.DB.Exec("INSERT INTO mata_kuliah (nama, kode) VALUES (?, ?)", mk.Nama, mk.Kode)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()

	if err := simpanRelasiMataKuliahKelas(int(id), mk.KelasIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Mata kuliah berhasil ditambahkan", "id": id})
}

func UpdateMataKuliah(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var mk models.MataKuliah

	if err := c.ShouldBindJSON(&mk); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if mk.Nama == "" || mk.Kode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama dan kode mata kuliah wajib diisi"})
		return
	}

	result, err := database.DB.Exec("UPDATE mata_kuliah SET nama = ?, kode = ? WHERE id = ?", mk.Nama, mk.Kode, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowAffected, _ := result.RowsAffected()

	if rowAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mata kuliah tidak ditemukan"})
		return
	}

	if err := simpanRelasiMataKuliahKelas(id, mk.KelasIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mata kuliah berhasil diperbarui"})
}

func DeleteMataKuliah(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	result, err := database.DB.Exec("DELETE FROM mata_kuliah WHERE id = ?", id)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": fkErrorOr(err, "Mata kuliah masih dipakai jadwal, tidak bisa dihapus")})
		return
	}

	rowAffected, _ := result.RowsAffected()

	if rowAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mata kuliah tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mata kuliah berhasil dihapus"})
}

// GetKelasByMataKuliah dipakai buat cascading dropdown pas bikin jadwal:
// cuma nampilin kelas yang emang pakai mata kuliah tsb.
func GetKelasByMataKuliah(c *gin.Context) {

	mataKuliahID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	rows, err := database.DB.Query(`
		SELECT k.id, k.nama, k.jurusan_id, j.nama
		FROM kelas k
		JOIN mata_kuliah_kelas mkk ON mkk.kelas_id = k.id
		JOIN jurusan j ON j.id = k.jurusan_id
		WHERE mkk.mata_kuliah_id = ?
		ORDER BY k.nama ASC
	`, mataKuliahID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var data []models.Kelas

	for rows.Next() {

		var k models.Kelas

		rows.Scan(&k.ID, &k.Nama, &k.JurusanID, &k.NamaJurusan)

		data = append(data, k)
	}

	c.JSON(http.StatusOK, data)
}