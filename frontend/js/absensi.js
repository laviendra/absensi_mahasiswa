let token = localStorage.getItem("token");

if (!token) {
    window.location.href = "index.html";
}

const inputTanggal = document.getElementById("tanggalFilter");

// default tanggal = hari ini
const hariIni = new Date().toISOString().slice(0, 10);
inputTanggal.value = hariIni;

// isi dropdown mahasiswa (dipakai buat tombol Absen Masuk)
fetch("http://localhost:8080/mahasiswa", {

    headers: {
        "Authorization": "Bearer " + token
    }

})

.then(response => response.json())

.then(data => {

    let option = "";

    data.forEach(mhs => {

        option += `
        <option value="${mhs.id}">
        ${mhs.nama}
        </option>
        `;

    });

    document.getElementById("mahasiswa").innerHTML = option;

});

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

            let tombol = "-";

            if (item.jam_masuk && !item.jam_pulang) {

                tombol = `
                <button
                    class="btn btn-warning btn-sm"
                    onclick="absenPulang(${item.absensi_id})">
                    Pulang
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
                <td>${item.jam_pulang || "-"}</td>
                <td><span class="badge bg-${warna}">${item.status_kehadiran}</span></td>
                <td>${tombol}</td>

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

function absenMasuk() {

    let id = document.getElementById("mahasiswa").value;

    fetch("http://localhost:8080/absensi/masuk", {

        method: "POST",

        headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + token
        },

        body: JSON.stringify({

            mahasiswa_id: Number(id)

        })

    })

    .then(response => response.json())

    .then(data => {

        alert(data.message);

        loadAbsensi(inputTanggal.value);

    });

}

function absenPulang(id) {

    const konfirmasi = confirm("Yakin ingin melakukan absensi pulang?");

    if (!konfirmasi) {
        return;
    }

    fetch("http://localhost:8080/absensi/pulang/" + id, {

        method: "PUT",

        headers: {
            "Authorization": "Bearer " + token
        }

    })

    .then(response => response.json())

    .then(data => {

        alert(data.message);

        if (data.message === "Absensi pulang berhasil") {
            loadAbsensi(inputTanggal.value);
        }

    })

    .catch(error => {

        console.log(error);
        alert("Terjadi kesalahan.");

    });

}