let dosenToken = localStorage.getItem("dosenToken");
let pertemuanAktif = null;

if (dosenToken) {
    tampilkanJadwal();
}

function loginDosen() {

    let username = document.getElementById("username").value;
    let password = document.getElementById("password").value;

    fetch(API_BASE + "/login-dosen", {

        method: "POST",

        headers: {
            "Content-Type": "application/json"
        },

        body: JSON.stringify({ username, password })

    })

    .then(response => response.json())

    .then(data => {

        if (data.token) {

            localStorage.setItem("dosenToken", data.token);
            localStorage.setItem("dosenNama", data.nama);
            dosenToken = data.token;

            tampilkanJadwal();

        } else {

            document.getElementById("pesanLogin").innerText = data.message || "Login gagal";

        }

    })

    .catch(() => {

        document.getElementById("pesanLogin").innerText = "Server tidak terhubung";

    });

}

function logoutDosen() {

    localStorage.removeItem("dosenToken");
    localStorage.removeItem("dosenNama");

    location.reload();

}

function tampilkanJadwal() {

    document.getElementById("areaLogin").classList.add("d-none");
    document.getElementById("areaJadwal").classList.remove("d-none");
    document.getElementById("areaAbsensi").classList.add("d-none");
    document.getElementById("btnLogout").classList.remove("d-none");

    fetch(API_BASE + "/dosen-area/jadwal-saya", {

        headers: {
            "Authorization": "Bearer " + dosenToken
        }

    })

    .then(response => response.json())

    .then(data => {

        let tabel = "";

        (data || []).forEach(j => {

            tabel += `

            <tr>
                <td>${j.hari}</td>
                <td>${j.jam_mulai}</td>
                <td>${j.nama_mata_kuliah}</td>
                <td>${j.nama_kelas}</td>
                <td>
                    <button class="btn btn-success btn-sm" onclick="ambilAbsen(${j.id})">
                        Ambil Absen Hari Ini
                    </button>
                </td>
            </tr>

            `;

        });

        document.getElementById("dataJadwal").innerHTML =
            tabel || `<tr><td colspan="5" class="text-center text-muted">Belum ada jadwal</td></tr>`;

    });

}

function ambilAbsen(jadwalId) {

    fetch(API_BASE + "/dosen-area/jadwal/" + jadwalId + "/buka-pertemuan", {

        method: "POST",

        headers: {
            "Authorization": "Bearer " + dosenToken
        }

    })

    .then(response => response.json())

    .then(data => {

        if (!data.pertemuan_id) {
            alert(data.error || "Gagal membuka pertemuan");
            return;
        }

        pertemuanAktif = data.pertemuan_id;

        loadAbsensiKelas();

    });

}

function loadAbsensiKelas() {

    document.getElementById("areaJadwal").classList.add("d-none");
    document.getElementById("areaAbsensi").classList.remove("d-none");

    fetch(API_BASE + "/dosen-area/pertemuan/" + pertemuanAktif + "/absensi", {

        headers: {
            "Authorization": "Bearer " + dosenToken
        }

    })

    .then(response => response.json())

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
                    <select
                        class="form-select form-select-sm"
                        onchange="simpanAbsen(${item.mahasiswa_id}, this.value)">
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

        document.getElementById("dataAbsensiKelas").innerHTML = tabel;

    });

}

function simpanAbsen(mahasiswaId, status) {

    fetch(API_BASE + "/dosen-area/pertemuan/" + pertemuanAktif + "/absensi", {

        method: "POST",

        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + dosenToken
        },

        body: JSON.stringify({
            mahasiswa_id: mahasiswaId,
            status_kehadiran: status
        })

    })

    .then(response => response.json())

    .then(data => {

        if (data.error) {
            alert(data.error);
        }

    });

}

function kembaliKeJadwal() {
    tampilkanJadwal();
}

function tutupPertemuan() {

    if (!confirm("Tutup pertemuan ini? Kehadiran tidak bisa diubah lagi setelah ditutup.")) {
        return;
    }

    fetch(API_BASE + "/dosen-area/pertemuan/" + pertemuanAktif + "/tutup", {

        method: "PUT",

        headers: {
            "Authorization": "Bearer " + dosenToken
        }

    })

    .then(response => response.json())

    .then(data => {

        alert(data.message || data.error || "Terjadi kesalahan");

        tampilkanJadwal();

    });

}