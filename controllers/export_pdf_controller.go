package controllers

import (
	"net/http"
	"strconv"
	"time"

	"absensi-mahasiswa/database"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func ExportPDF(c *gin.Context) {

	bulan, err := strconv.Atoi(c.Query("bulan"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bulan tidak valid",
		})
		return
	}

	tahun, err := strconv.Atoi(c.Query("tahun"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "tahun tidak valid",
		})
		return
	}


	rows, err := database.DB.Query(`
		SELECT 
			m.nama,
			SUM(CASE WHEN a.status_kehadiran='Hadir' THEN 1 ELSE 0 END),
			SUM(CASE WHEN a.status_kehadiran='Terlambat' THEN 1 ELSE 0 END),
			SUM(CASE WHEN a.status_kehadiran='Tidak Hadir' THEN 1 ELSE 0 END)
		FROM mahasiswa m
		LEFT JOIN absensi a
		ON m.id = a.mahasiswa_id
		AND MONTH(a.tanggal)=?
		AND YEAR(a.tanggal)=?
		GROUP BY m.nama
		ORDER BY m.nama ASC
	`,
		bulan,
		tahun,
	)


	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	defer rows.Close()


	// membuat PDF
	pdf := gofpdf.New(
		"P",
		"mm",
		"A4",
		"",
	)


	pdf.AddPage()


	// Judul
	pdf.SetFont(
		"Arial",
		"B",
		16,
	)

	pdf.Cell(
		0,
		10,
		"REKAP ABSENSI MAHASISWA",
	)

	pdf.Ln(8)


	pdf.SetFont(
		"Arial",
		"",
		12,
	)

	pdf.Cell(
		0,
		10,
		"Bulan: "+strconv.Itoa(bulan)+" Tahun: "+strconv.Itoa(tahun),
	)

	pdf.Ln(15)


	// Header tabel
	pdf.SetFont(
		"Arial",
		"B",
		10,
	)


	pdf.CellFormat(10, 10, "No", "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, 10, "Nama Mahasiswa", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 10, "Hadir", "1", 0, "C", false, 0, "")
	pdf.CellFormat(35, 10, "Terlambat", "1", 0, "C", false, 0, "")
	pdf.CellFormat(35, 10, "Tidak Hadir", "1", 1, "C", false, 0, "")


	// Isi tabel
	pdf.SetFont(
		"Arial",
		"",
		10,
	)


	no := 1


	for rows.Next() {

		var nama string
		var hadir int
		var terlambat int
		var tidakHadir int


		err := rows.Scan(
			&nama,
			&hadir,
			&terlambat,
			&tidakHadir,
		)


		if err != nil {
			continue
		}


		pdf.CellFormat(
			10,
			10,
			strconv.Itoa(no),
			"1",
			0,
			"C",
			false,
			0,
			"",
		)

		pdf.CellFormat(
			60,
			10,
			nama,
			"1",
			0,
			"L",
			false,
			0,
			"",
		)

		pdf.CellFormat(
			30,
			10,
			strconv.Itoa(hadir),
			"1",
			0,
			"C",
			false,
			0,
			"",
		)

		pdf.CellFormat(
			35,
			10,
			strconv.Itoa(terlambat),
			"1",
			0,
			"C",
			false,
			0,
			"",
		)

		pdf.CellFormat(
			35,
			10,
			strconv.Itoa(tidakHadir),
			"1",
			1,
			"C",
			false,
			0,
			"",
		)


		no++
	}


	// Footer
	pdf.Ln(15)

	pdf.SetFont(
		"Arial",
		"I",
		10,
	)

	pdf.Cell(
		0,
		10,
		"Dicetak pada: "+time.Now().Format("02-01-2006 15:04"),
	)


	c.Header(
		"Content-Disposition",
		"attachment; filename=rekap_absensi.pdf",
	)

	c.Header(
		"Content-Type",
		"application/pdf",
	)


	err = pdf.Output(c.Writer)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}