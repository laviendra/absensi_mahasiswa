let token = localStorage.getItem("token");

if (!token) {
    window.location.href = "index.html";
}

let pertemuanDilihat = null;

function loadPertemuan() {

    let tanggal = document.getElementById("tanggalFilter").value;

    let url = API_BASE + "/rekap-kelas/pertemuan";

    if (tanggal) {
        url += "?tanggal=" + tanggal;
    }

    fetch(url, {
        headers: { "Authorization": "Bearer " + token }
    })

    .then(r => r.json())

    .then(data => {

        let tabel = "";

        (data || []).forEach((p, i) => {

            let badgeStatus = p.status === "selesai" ? "secondary" : "success";

            tabel += `
            <tr>
                <td>${i + 1}</td>
                <td>${p.tanggal}</td>
                <td>${p.nama_mata_kuliah}</td>
                <td>${p.nama_kelas}</td>
                <td>${p.nama_dosen}</td>
                <td><span class="badge bg-${badgeStatus}">${p.status}</span></td>
                <td>${p.hadir}</td>
                <td>${p.terlambat}</td>
                <td>${p.izin}</td>
                <td>${p.sakit}</td>
                <td>${p.tidak_hadir}</td>
                <td>
                    <button class="btn btn-primary btn-sm" onclick="lihatDetail(${p.id}, '${p.nama_mata_kuliah}', '${p.nama_kelas}', '${p.tanggal}')">
                        Lihat Detail
                    </button>
                </td>
            </tr>
            `;

        });

        document.getElementById("dataPertemuan").innerHTML =
            tabel || `<tr><td colspan="12" class="text-center text-muted">Belum ada pertemuan</td></tr>`;

    });

}

function resetFilterTanggal() {
    document.getElementById("tanggalFilter").value = "";
    loadPertemuan();
}

document.getElementById("tanggalFilter").addEventListener("change", loadPertemuan);

loadPertemuan();

function lihatDetail(pertemuanId, matkul, kelas, tanggal) {

    pertemuanDilihat = pertemuanId;

    document.getElementById("areaDaftar").classList.add("d-none");
    document.getElementById("areaDetail").classList.remove("d-none");
    document.getElementById("judulDetail").innerText =
        matkul + " - " + kelas + " (" + tanggal + ")";

    fetch(API_BASE + "/rekap-kelas/pertemuan/" + pertemuanId + "/absensi", {
        headers: { "Authorization": "Bearer " + token }
    })

    .then(r => r.json())

    .then(hasil => {

        let data = hasil.data || [];

        let tabel = "";

        data.forEach((item, index) => {

            tabel += `
            <tr>
                <td>${index + 1}</td>
                <td>${item.nama}</td>
                <td>${item.jam_hadir || "-"}</td>
                <td>
                    <select class="form-select form-select-sm" onchange="koreksiAbsen(${item.mahasiswa_id}, this.value)">
                        <option value="Hadir" ${item.status_kehadiran === "Hadir" ? "selected" : ""}>Hadir</option>
                        <option value="Terlambat" ${item.status_kehadiran === "Terlambat" ? "selected" : ""}>Terlambat</option>
                        <option value="Izin" ${item.status_kehadiran === "Izin" ? "selected" : ""}>Izin</option>
                        <option value="Sakit" ${item.status_kehadiran === "Sakit" ? "selected" : ""}>Sakit</option>
                        <option value="Tidak Hadir" ${item.status_kehadiran === "Tidak Hadir" ? "selected" : ""}>Tidak Hadir</option>
                    </select>
                </td>
            </tr>
            `;

        });

        document.getElementById("dataDetailAbsensi").innerHTML = tabel;

    });

}

function koreksiAbsen(mahasiswaId, status) {

    fetch(API_BASE + "/rekap-kelas/pertemuan/" + pertemuanDilihat + "/absensi", {

        method: "PUT",

        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + token
        },

        body: JSON.stringify({
            mahasiswa_id: mahasiswaId,
            status_kehadiran: status
        })

    })

    .then(r => r.json())

    .then(data => {
        if (data.error) alert(data.error);
    });

}

function kembaliKeDaftar() {
    document.getElementById("areaDetail").classList.add("d-none");
    document.getElementById("areaDaftar").classList.remove("d-none");
    loadPertemuan();
}