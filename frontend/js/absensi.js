let token = localStorage.getItem("token");

if (!token) {
    window.location.href = "index.html";
}

const inputTanggal = document.getElementById("tanggalFilter");

// default tanggal = hari ini
const hariIni = new Date().toISOString().slice(0, 10);
inputTanggal.value = hariIni;

function loadAbsensi(tanggal) {

    fetch("http://localhost:8080/absensi/filter?tanggal=" + tanggal, {

        headers: {
            "Authorization": "Bearer " + token
        }

    })

    .then(response => response.json())

    .then(hasil => {

        let data = hasil.data || [];

        let tabel = "";

        data.forEach((item, index) => {

            let aksi = "-";

            // absen masuk cuma boleh dilakukan untuk tanggal hari ini,
            // data tanggal lain sifatnya cuma riwayat (read only)
            if (tanggal === hariIni && !item.absensi_id) {

                aksi = `
                <button
                    class="btn btn-success btn-sm"
                    onclick="absenMasuk(${item.mahasiswa_id})">
                    Masuk
                </button>
                `;
            }

            let warna = "secondary";

            if (item.status_kehadiran === "Hadir") warna = "success";
            if (item.status_kehadiran === "Terlambat") warna = "warning";
            if (item.status_kehadiran === "Tidak Hadir") warna = "danger";

            tabel += `

            <tr>

                <td>${index + 1}</td>
                <td>${item.nama}</td>
                <td>${item.jam_masuk || "-"}</td>
                <td><span class="badge bg-${warna}">${item.status_kehadiran}</span></td>
                <td>${aksi}</td>

            </tr>

            `;

        });

        document.getElementById("dataAbsensi").innerHTML = tabel;

    });

}

// load data untuk tanggal hari ini saat halaman dibuka
loadAbsensi(inputTanggal.value);

// reload tabel setiap kali tanggal diganti
inputTanggal.addEventListener("change", () => {
    loadAbsensi(inputTanggal.value);
});

function absenMasuk(mahasiswaId) {

    fetch("http://localhost:8080/absensi/masuk", {

        method: "POST",

        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + token
        },

        body: JSON.stringify({

            mahasiswa_id: Number(mahasiswaId)

        })

    })

    .then(response => response.json())

    .then(data => {

        alert(data.message || data.error || "Terjadi kesalahan");

        loadAbsensi(inputTanggal.value);

    });

}