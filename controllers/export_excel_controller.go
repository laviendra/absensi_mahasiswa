package controllers

import (
	"strconv"

	"absensi-mahasiswa/database"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func ExportExcel(c *gin.Context) {

	bulan, err := strconv.Atoi(c.Query("bulan"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "bulan tidak valid",
		})
		return
	}

	tahun, err := strconv.Atoi(c.Query("tahun"))
	if err != nil {
		c.JSON(400, gin.H{
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
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	defer rows.Close()


	file := excelize.NewFile()

	sheet := "Rekap Absensi"

	file.SetSheetName(
		"Sheet1",
		sheet,
	)


	// Header
	file.SetCellValue(sheet, "A1", "No")
	file.SetCellValue(sheet, "B1", "Nama Mahasiswa")
	file.SetCellValue(sheet, "C1", "Hadir")
	file.SetCellValue(sheet, "D1", "Terlambat")
	file.SetCellValue(sheet, "E1", "Tidak Hadir")


	baris := 2
	no := 1


	for rows.Next() {

		var nama string
		var hadir int
		var terlambat int
		var tidakHadir int


		rows.Scan(
			&nama,
			&hadir,
			&terlambat,
			&tidakHadir,
		)


		file.SetCellValue(sheet, "A"+strconv.Itoa(baris), no)
		file.SetCellValue(sheet, "B"+strconv.Itoa(baris), nama)
		file.SetCellValue(sheet, "C"+strconv.Itoa(baris), hadir)
		file.SetCellValue(sheet, "D"+strconv.Itoa(baris), terlambat)
		file.SetCellValue(sheet, "E"+strconv.Itoa(baris), tidakHadir)


		baris++
		no++
	}


	// lebar kolom
	file.SetColWidth(sheet, "A", "A", 8)
	file.SetColWidth(sheet, "B", "B", 25)
	file.SetColWidth(sheet, "C", "E", 15)


	c.Header(
		"Content-Disposition",
		"attachment; filename=rekap_absensi.xlsx",
	)

	c.Header(
		"Content-Type",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	)


	err = file.Write(c.Writer)

	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}
}