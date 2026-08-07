package controllers

import (
	"net/http"
	"strconv"
	"time"

	"absensi-mahasiswa/database"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

// query rekap: buat tiap mahasiswa di kelas jadwal ini, hitung total pertemuan
// yang udah berlangsung buat jadwal itu, dan berapa kali masing-masing status.
const queryRekapJadwal = `
	SELECT m.nim, m.nama,
		COUNT(DISTINCT p.id),
		SUM(CASE WHEN COALESCE(ak.status_kehadiran,'Tidak Hadir') = 'Hadir' THEN 1 ELSE 0 END),
		SUM(CASE WHEN COALESCE(ak.status_kehadiran,'Tidak Hadir') = 'Terlambat' THEN 1 ELSE 0 END),
		SUM(CASE WHEN COALESCE(ak.status_kehadiran,'Tidak Hadir') = 'Izin' THEN 1 ELSE 0 END),
		SUM(CASE WHEN COALESCE(ak.status_kehadiran,'Tidak Hadir') = 'Sakit' THEN 1 ELSE 0 END),
		SUM(CASE WHEN COALESCE(ak.status_kehadiran,'Tidak Hadir') = 'Tidak Hadir' THEN 1 ELSE 0 END)
	FROM jadwal j
	JOIN mahasiswa m ON m.kelas_id = j.kelas_id AND m.status = 'aktif'
	LEFT JOIN pertemuan p ON p.jadwal_id = j.id
	LEFT JOIN absensi_kelas ak ON ak.pertemuan_id = p.id AND ak.mahasiswa_id = m.id
	WHERE j.id = ?
	GROUP BY m.id, m.nim, m.nama
	ORDER BY m.nama ASC
`

func ambilInfoJadwal(jadwalID int) (namaMatkul, namaKelas, namaDosen string, err error) {

	err = database.DB.QueryRow(`
		SELECT mk.nama, k.nama, d.nama
		FROM jadwal j
		JOIN mata_kuliah mk ON mk.id = j.mata_kuliah_id
		JOIN kelas k ON k.id = j.kelas_id
		JOIN dosen d ON d.id = j.dosen_id
		WHERE j.id = ?
	`, jadwalID).Scan(&namaMatkul, &namaKelas, &namaDosen)

	return
}

func ExportRekapKelasPDF(c *gin.Context) {

	jadwalID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID jadwal tidak valid"})
		return
	}

	namaMatkul, namaKelas, namaDosen, err := ambilInfoJadwal(jadwalID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jadwal tidak ditemukan"})
		return
	}

	rows, err := database.DB.Query(queryRekapJadwal, jadwalID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 8, "REKAP KEHADIRAN MATA KULIAH")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 6, "Mata Kuliah : "+namaMatkul)
	pdf.Ln(6)
	pdf.Cell(0, 6, "Kelas       : "+namaKelas)
	pdf.Ln(6)
	pdf.Cell(0, 6, "Dosen       : "+namaDosen)
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(10, 8, "No", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 8, "NIM", "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, 8, "Nama", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 8, "Pertemuan", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 8, "Hadir", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 8, "Terlambat", "1", 0, "C", false, 0, "")
	pdf.CellFormat(18, 8, "Izin", "1", 0, "C", false, 0, "")
	pdf.CellFormat(18, 8, "Sakit", "1", 0, "C", false, 0, "")
	pdf.CellFormat(25, 8, "Tdk Hadir", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 8, "%", "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 10)

	no := 1

	for rows.Next() {

		var nim, nama string
		var total, hadir, terlambat, izin, sakit, tidakHadir int

		err := rows.Scan(&nim, &nama, &total, &hadir, &terlambat, &izin, &sakit, &tidakHadir)

		if err != nil {
			continue
		}

		persen := 0

		if total > 0 {
			persen = int((float64(hadir+terlambat) / float64(total)) * 100)
		}

		pdf.CellFormat(10, 8, strconv.Itoa(no), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 8, nim, "1", 0, "C", false, 0, "")
		pdf.CellFormat(60, 8, nama, "1", 0, "L", false, 0, "")
		pdf.CellFormat(25, 8, strconv.Itoa(total), "1", 0, "C", false, 0, "")
		pdf.CellFormat(20, 8, strconv.Itoa(hadir), "1", 0, "C", false, 0, "")
		pdf.CellFormat(25, 8, strconv.Itoa(terlambat), "1", 0, "C", false, 0, "")
		pdf.CellFormat(18, 8, strconv.Itoa(izin), "1", 0, "C", false, 0, "")
		pdf.CellFormat(18, 8, strconv.Itoa(sakit), "1", 0, "C", false, 0, "")
		pdf.CellFormat(25, 8, strconv.Itoa(tidakHadir), "1", 0, "C", false, 0, "")
		pdf.CellFormat(20, 8, strconv.Itoa(persen)+"%", "1", 1, "C", false, 0, "")

		no++
	}

	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 9)
	pdf.Cell(0, 6, "Dicetak pada: "+time.Now().Format("02-01-2006 15:04"))

	c.Header("Content-Disposition", "attachment; filename=rekap_kehadiran.pdf")
	c.Header("Content-Type", "application/pdf")

	if err := pdf.Output(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func ExportRekapKelasExcel(c *gin.Context) {

	jadwalID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID jadwal tidak valid"})
		return
	}

	namaMatkul, namaKelas, namaDosen, err := ambilInfoJadwal(jadwalID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jadwal tidak ditemukan"})
		return
	}

	rows, err := database.DB.Query(queryRekapJadwal, jadwalID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	file := excelize.NewFile()
	sheet := "Rekap"

	file.SetSheetName("Sheet1", sheet)

	file.SetCellValue(sheet, "A1", "Mata Kuliah")
	file.SetCellValue(sheet, "B1", namaMatkul)
	file.SetCellValue(sheet, "A2", "Kelas")
	file.SetCellValue(sheet, "B2", namaKelas)
	file.SetCellValue(sheet, "A3", "Dosen")
	file.SetCellValue(sheet, "B3", namaDosen)

	headerRow := 5

	header := []string{"No", "NIM", "Nama", "Total Pertemuan", "Hadir", "Terlambat", "Izin", "Sakit", "Tidak Hadir", "Persentase"}
	kolom := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}

	for i, h := range header {
		file.SetCellValue(sheet, kolom[i]+strconv.Itoa(headerRow), h)
	}

	baris := headerRow + 1
	no := 1

	for rows.Next() {

		var nim, nama string
		var total, hadir, terlambat, izin, sakit, tidakHadir int

		err := rows.Scan(&nim, &nama, &total, &hadir, &terlambat, &izin, &sakit, &tidakHadir)

		if err != nil {
			continue
		}

		persen := 0

		if total > 0 {
			persen = int((float64(hadir+terlambat) / float64(total)) * 100)
		}

		file.SetCellValue(sheet, "A"+strconv.Itoa(baris), no)
		file.SetCellValue(sheet, "B"+strconv.Itoa(baris), nim)
		file.SetCellValue(sheet, "C"+strconv.Itoa(baris), nama)
		file.SetCellValue(sheet, "D"+strconv.Itoa(baris), total)
		file.SetCellValue(sheet, "E"+strconv.Itoa(baris), hadir)
		file.SetCellValue(sheet, "F"+strconv.Itoa(baris), terlambat)
		file.SetCellValue(sheet, "G"+strconv.Itoa(baris), izin)
		file.SetCellValue(sheet, "H"+strconv.Itoa(baris), sakit)
		file.SetCellValue(sheet, "I"+strconv.Itoa(baris), tidakHadir)
		file.SetCellValue(sheet, "J"+strconv.Itoa(baris), strconv.Itoa(persen)+"%")

		baris++
		no++
	}

	file.SetColWidth(sheet, "A", "A", 6)
	file.SetColWidth(sheet, "B", "B", 15)
	file.SetColWidth(sheet, "C", "C", 25)
	file.SetColWidth(sheet, "D", "J", 14)

	c.Header("Content-Disposition", "attachment; filename=rekap_kehadiran.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	if err := file.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}