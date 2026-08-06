let mahasiswaToken = localStorage.getItem("mahasiswaToken");

if (!mahasiswaToken) {
    window.location.href = "index.html";
}

tampilkanRekap();

function logoutMahasiswa() {

    localStorage.removeItem("mahasiswaToken");
    localStorage.removeItem("mahasiswaNama");

    window.location.href = "index.html";

}

function tampilkanRekap() {

    fetch(API_BASE + "/mahasiswa-area/rekap", {

        headers: {
            "Authorization": "Bearer " + mahasiswaToken
        }

    })

    .then(response => response.json())

    .then(hasil => {

        let data = hasil || [];

        // ringkasan persentase kehadiran per mata kuliah
        let ringkasan = {};

        data.forEach(item => {

            if (!ringkasan[item.nama_mata_kuliah]) {
                ringkasan[item.nama_mata_kuliah] = { total: 0, hadir: 0 };
            }

            ringkasan[item.nama_mata_kuliah].total++;

            if (item.status_kehadiran === "Hadir" || item.status_kehadiran === "Terlambat") {
                ringkasan[item.nama_mata_kuliah].hadir++;
            }

        });

        let htmlRingkasan = "";

        Object.keys(ringkasan).forEach(matkul => {

            let r = ringkasan[matkul];
            let persen = r.total > 0 ? Math.round((r.hadir / r.total) * 100) : 0;
            let warna = persen >= 75 ? "success" : "danger";

            htmlRingkasan += `
            <div class="col-md-4">
                <div class="card shadow-sm">
                    <div class="card-body">
                        <h6 class="card-title">${matkul}</h6>
                        <div class="d-flex justify-content-between align-items-center">
                            <span class="text-muted small">${r.hadir}/${r.total} pertemuan</span>
                            <span class="badge bg-${warna}">${persen}%</span>
                        </div>
                    </div>
                </div>
            </div>
            `;

        });

        document.getElementById("ringkasanMatkul").innerHTML =
            htmlRingkasan || `<p class="text-muted">Belum ada riwayat pertemuan</p>`;

        // tabel riwayat lengkap
        let tabel = "";

        data.forEach(item => {

            let warna = "secondary";

            if (item.status_kehadiran === "Hadir") warna = "success";
            if (item.status_kehadiran === "Terlambat") warna = "warning";
            if (item.status_kehadiran === "Tidak Hadir") warna = "danger";
            if (item.status_kehadiran === "Izin" || item.status_kehadiran === "Sakit") warna = "info";

            tabel += `
            <tr>
                <td>${item.tanggal}</td>
                <td>${item.nama_mata_kuliah}</td>
                <td>${item.nama_dosen}</td>
                <td>${item.jam_hadir || "-"}</td>
                <td><span class="badge bg-${warna}">${item.status_kehadiran}</span></td>
            </tr>
            `;

        });

        document.getElementById("dataRekap").innerHTML =
            tabel || `<tr><td colspan="5" class="text-center text-muted">Belum ada riwayat pertemuan</td></tr>`;

    });

}