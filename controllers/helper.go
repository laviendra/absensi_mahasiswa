package controllers

import "strings"

// fkErrorOr mengubah error SQL foreign key jadi pesan yang gampang dibaca user,
// error lain tetap ditampilkan apa adanya.
func fkErrorOr(err error, pesanFK string) string {

	if strings.Contains(err.Error(), "foreign key constraint") {
		return pesanFK
	}

	if strings.Contains(err.Error(), "Duplicate entry") {
		return "Data ini sudah ada"
	}

	return err.Error()
}