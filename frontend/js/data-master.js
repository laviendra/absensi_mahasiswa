let token = localStorage.getItem("token");

if (!token) {
    window.location.href = "index.html";
}

const API = "http://localhost:8080";

function authHeader(json) {
    let h = { "Authorization": "Bearer " + token };
    if (json) h["Content-Type"] = "application/json";
    return h;
}


// ================= JURUSAN =================

function loadJurusan() {

    fetch(API + "/jurusan", { headers: authHeader() })
    .then(r => r.json())
    .then(data => {

        let tabel = "";

        (data || []).forEach((j, i) => {
            tabel += `
            <tr>
                <td>${i + 1}</td>
                <td>${j.nama}</td>
                <td>
                    <button class="btn btn-warning btn-sm" onclick="editJurusan(${j.id}, '${j.nama}')">Edit</button>
                    <button class="btn btn-danger btn-sm" onclick="hapusJurusan(${j.id})">Hapus</button>
                </td>
            </tr>
            `;
        });

        document.getElementById("dataJurusan").innerHTML =
            tabel || `<tr><td colspan="3" class="text-center text-muted">Belum ada data</td></tr>`;

    });

}

function bukaModalJurusan() {
    document.getElementById("idJurusan").value = "";
    document.getElementById("namaJurusan").value = "";
    new bootstrap.Modal(document.getElementById("modalJurusan")).show();
}

function editJurusan(id, nama) {
    document.getElementById("idJurusan").value = id;
    document.getElementById("namaJurusan").value = nama;
    new bootstrap.Modal(document.getElementById("modalJurusan")).show();
}

function simpanJurusan() {

    let id = document.getElementById("idJurusan").value;
    let nama = document.getElementById("namaJurusan").value;

    let method = id ? "PUT" : "POST";
    let url = id ? API + "/jurusan/" + id : API + "/jurusan";

    fetch(url, {
        method: method,
        headers: authHeader(true),
        body: JSON.stringify({ nama })
    })
    .then(r => r.json())
    .then(data => {
        alert(data.message || data.error || "Terjadi kesalahan");
        bootstrap.Modal.getInstance(document.getElementById("modalJurusan")).hide();
        loadJurusan();
        loadKelasOptions();
    });

}

function hapusJurusan(id) {

    if (!confirm("Hapus jurusan ini?")) return;

    fetch(API + "/jurusan/" + id, { method: "DELETE", headers: authHeader() })
    .then(r => r.json())
    .then(data => {
        alert(data.message || data.error || "Terjadi kesalahan");
        loadJurusan();
    });

}


// ================= KELAS =================

function loadKelasOptions() {

    fetch(API + "/jurusan", { headers: authHeader() })
    .then(r => r.json())
    .then(data => {

        let opt = `<option value="">Pilih Jurusan</option>`;

        (data || []).forEach(j => {
            opt += `<option value="${j.id}">${j.nama}</option>`;
        });

        document.getElementById("jurusanKelas").innerHTML = opt;

    });

}

function loadKelas() {

    fetch(API + "/kelas", { headers: authHeader() })
    .then(r => r.json())
    .then(data => {

        let tabel = "";

        (data || []).forEach((k, i) => {
            tabel += `
            <tr>
                <td>${i + 1}</td>
                <td>${k.nama}</td>
                <td>${k.nama_jurusan}</td>
                <td>
                    <button class="btn btn-warning btn-sm" onclick="editKelas(${k.id}, '${k.nama}', ${k.jurusan_id})">Edit</button>
                    <button class="btn btn-danger btn-sm" onclick="hapusKelas(${k.id})">Hapus</button>
                </td>
            </tr>
            `;
        });

        document.getElementById("dataKelas").innerHTML =
            tabel || `<tr><td colspan="4" class="text-center text-muted">Belum ada data</td></tr>`;

    });

}

function bukaModalKelas() {
    document.getElementById("idKelas").value = "";
    document.getElementById("namaKelas").value = "";
    document.getElementById("jurusanKelas").value = "";
    new bootstrap.Modal(document.getElementById("modalKelas")).show();
}

function editKelas(id, nama, jurusanId) {
    document.getElementById("idKelas").value = id;
    document.getElementById("namaKelas").value = nama;
    document.getElementById("jurusanKelas").value = jurusanId;
    new bootstrap.Modal(document.getElementById("modalKelas")).show();
}

function simpanKelas() {

    let id = document.getElementById("idKelas").value;
    let nama = document.getElementById("namaKelas").value;
    let jurusanId = document.getElementById("jurusanKelas").value;

    if (!jurusanId) {
        alert("Jurusan wajib dipilih");
        return;
    }

    let method = id ? "PUT" : "POST";
    let url = id ? API + "/kelas/" + id : API + "/kelas";

    fetch(url, {
        method: method,
        headers: authHeader(true),
        body: JSON.stringify({ nama, jurusan_id: Number(jurusanId) })
    })
    .then(r => r.json())
    .then(data => {
        alert(data.message || data.error || "Terjadi kesalahan");
        bootstrap.Modal.getInstance(document.getElementById("modalKelas")).hide();
        loadKelas();
        loadKelasJadwalOptions();
    });

}

function hapusKelas(id) {

    if (!confirm("Hapus kelas ini?")) return;

    fetch(API + "/kelas/" + id, { method: "DELETE", headers: authHeader() })
    .then(r => r.json())
    .then(data => {
        alert(data.message || data.error || "Terjadi kesalahan");
        loadKelas();
    });

}


// ================= DOSEN =================

function loadDosen() {

    fetch(API + "/dosen", { headers: authHeader() })
    .then(r => r.json())
    .then(data => {

        let tabel = "";

        (data || []).forEach((d, i) => {
            tabel += `
            <tr>
                <td>${i + 1}</td>
                <td>${d.nama}</td>
                <td>${d.username}</td>
                <td>
                    <button class="btn btn-warning btn-sm" onclick="editDosen(${d.id}, '${d.nama}', '${d.username}')">Edit</button>
                    <button class="btn btn-danger btn-sm" onclick="hapusDosen(${d.id})">Hapus</button>
                </td>
            </tr>
            `;
        });

        document.getElementById("dataDosen").innerHTML =
            tabel || `<tr><td colspan="4" class="text-center text-muted">Belum ada data</td></tr>`;

    });

}

function bukaModalDosen() {
    document.getElementById("idDosen").value = "";
    document.getElementById("namaDosen").value = "";
    document.getElementById("usernameDosen").value = "";
    document.getElementById("passwordDosen").value = "";
    new bootstrap.Modal(document.getElementById("modalDosen")).show();
}

function editDosen(id, nama, username) {
    document.getElementById("idDosen").value = id;
    document.getElementById("namaDosen").value = nama;
    document.getElementById("usernameDosen").value = username;
    document.getElementById("passwordDosen").value = "";
    new bootstrap.Modal(document.getElementById("modalDosen")).show();
}

function simpanDosen() {

    let id = document.getElementById("idDosen").value;
    let nama = document.getElementById("namaDosen").value;
    let username = document.getElementById("usernameDosen").value;
    let password = document.getElementById("passwordDosen").value;

    if (!id && !password) {
        alert("Password wajib diisi untuk dosen baru");
        return;
    }

    let method = id ? "PUT" : "POST";
    let url = id ? API + "/dosen/" + id : API + "/dosen";

    fetch(url, {
        method: method,
        headers: authHeader(true),
        body: JSON.stringify({ nama, username, password })
    })
    .then(r => r.json())
    .then(data => {
        alert(data.message || data.error || "Terjadi kesalahan");
        bootstrap.Modal.getInstance(document.getElementById("modalDosen")).hide();
        loadDosen();
        loadDosenJadwalOptions();
    });

}

function hapusDosen(id) {

    if (!confirm("Hapus dosen ini?")) return;

    fetch(API + "/dosen/" + id, { method: "DELETE", headers: authHeader() })
    .then(r => r.json())
    .then(data => {
        alert(data.message || data.error || "Terjadi kesalahan");
        loadDosen();
    });

}


// ================= MATA KULIAH =================

function loadMatkul() {

    fetch(API + "/mata-kuliah", { headers: authHeader() })
    .then(r => r.json())
    .then(data => {

        let tabel = "";

        (data || []).forEach((mk, i) => {
            tabel += `
            <tr>
                <td>${i + 1}</td>
                <td>${mk.kode}</td>
                <td>${mk.nama}</td>
                <td>
                    <button class="btn btn-warning btn-sm" onclick="editMatkul(${mk.id}, '${mk.kode}', '${mk.nama}')">Edit</button>
                    <button class="btn btn-danger btn-sm" onclick="hapusMatkul(${mk.id})">Hapus</button>
                </td>
            </tr>
            `;
        });

        document.getElementById("dataMatkul").innerHTML =
            tabel || `<tr><td colspan="4" class="text-center text-muted">Belum ada data</td></tr>`;

    });

}

function bukaModalMatkul() {
    document.getElementById("idMatkul").value = "";
    document.getElementById("kodeMatkul").value = "";
    document.getElementById("namaMatkul").value = "";
    new bootstrap.Modal(document.getElementById("modalMatkul")).show();
}

function editMatkul(id, kode, nama) {
    document.getElementById("idMatkul").value = id;
    document.getElementById("kodeMatkul").value = kode;
    document.getElementById("namaMatkul").value = nama;
    new bootstrap.Modal(document.getElementById("modalMatkul")).show();
}

function simpanMatkul() {

    let id = document.getElementById("idMatkul").value;
    let kode = document.getElementById("kodeMatkul").value;
    let nama = document.getElementById("namaMatkul").value;

    let method = id ? "PUT" : "POST";
    let url = id ? API + "/mata-kuliah/" + id : API + "/mata-kuliah";

    fetch(url, {
        method: method,
        headers: authHeader(true),
        body: JSON.stringify({ kode, nama })
    })
    .then(r => r.json())
    .then(data => {
        alert(data.message || data.error || "Terjadi kesalahan");
        bootstrap.Modal.getInstance(document.getElementById("modalMatkul")).hide();
        loadMatkul();
        loadMatkulJadwalOptions();
    });

}

function hapusMatkul(id) {

    if (!confirm("Hapus mata kuliah ini?")) return;

    fetch(API + "/mata-kuliah/" + id, { method: "DELETE", headers: authHeader() })
    .then(r => r.json())
    .then(data => {
        alert(data.message || data.error || "Terjadi kesalahan");
        loadMatkul();
    });

}


// ================= JADWAL =================

function loadDosenJadwalOptions() {

    fetch(API + "/dosen", { headers: authHeader() })
    .then(r => r.json())
    .then(data => {
        let opt = `<option value="">Pilih Dosen</option>`;
        (data || []).forEach(d => opt += `<option value="${d.id}">${d.nama}</option>`);
        document.getElementById("dosenJadwal").innerHTML = opt;
    });

}

function loadMatkulJadwalOptions() {

    fetch(API + "/mata-kuliah", { headers: authHeader() })
    .then(r => r.json())
    .then(data => {
        let opt = `<option value="">Pilih Mata Kuliah</option>`;
        (data || []).forEach(mk => opt += `<option value="${mk.id}">${mk.kode} - ${mk.nama}</option>`);
        document.getElementById("matkulJadwal").innerHTML = opt;
    });

}

function loadKelasJadwalOptions() {

    fetch(API + "/kelas", { headers: authHeader() })
    .then(r => r.json())
    .then(data => {
        let opt = `<option value="">Pilih Kelas</option>`;
        (data || []).forEach(k => opt += `<option value="${k.id}">${k.nama}</option>`);
        document.getElementById("kelasJadwal").innerHTML = opt;
    });

}

function loadJadwal() {

    fetch(API + "/jadwal", { headers: authHeader() })
    .then(r => r.json())
    .then(data => {

        let tabel = "";

        (data || []).forEach((j, i) => {
            tabel += `
            <tr>
                <td>${i + 1}</td>
                <td>${j.hari}</td>
                <td>${j.jam_mulai}</td>
                <td>${j.nama_dosen}</td>
                <td>${j.nama_mata_kuliah}</td>
                <td>${j.nama_kelas}</td>
                <td>
                    <button class="btn btn-danger btn-sm" onclick="hapusJadwal(${j.id})">Hapus</button>
                </td>
            </tr>
            `;
        });

        document.getElementById("dataJadwal").innerHTML =
            tabel || `<tr><td colspan="7" class="text-center text-muted">Belum ada data</td></tr>`;

    });

}

function bukaModalJadwal() {
    document.getElementById("dosenJadwal").value = "";
    document.getElementById("matkulJadwal").value = "";
    document.getElementById("kelasJadwal").value = "";
    document.getElementById("hariJadwal").value = "";
    document.getElementById("jamJadwal").value = "";
    new bootstrap.Modal(document.getElementById("modalJadwal")).show();
}

function simpanJadwal() {

    let dosenId = document.getElementById("dosenJadwal").value;
    let matkulId = document.getElementById("matkulJadwal").value;
    let kelasId = document.getElementById("kelasJadwal").value;
    let hari = document.getElementById("hariJadwal").value;
    let jam = document.getElementById("jamJadwal").value;

    if (!dosenId || !matkulId || !kelasId || !hari || !jam) {
        alert("Semua field wajib diisi");
        return;
    }

    fetch(API + "/jadwal", {
        method: "POST",
        headers: authHeader(true),
        body: JSON.stringify({
            dosen_id: Number(dosenId),
            kelas_id: Number(kelasId),
            mata_kuliah_id: Number(matkulId),
            hari: hari,
            jam_mulai: jam
        })
    })
    .then(r => r.json())
    .then(data => {
        alert(data.message || data.error || "Terjadi kesalahan");
        bootstrap.Modal.getInstance(document.getElementById("modalJadwal")).hide();
        loadJadwal();
    });

}

function hapusJadwal(id) {

    if (!confirm("Hapus jadwal ini?")) return;

    fetch(API + "/jadwal/" + id, { method: "DELETE", headers: authHeader() })
    .then(r => r.json())
    .then(data => {
        alert(data.message || data.error || "Terjadi kesalahan");
        loadJadwal();
    });

}


// ================= INIT =================

loadJurusan();
loadKelasOptions();
loadKelas();
loadDosen();
loadMatkul();
loadDosenJadwalOptions();
loadMatkulJadwalOptions();
loadKelasJadwalOptions();
loadJadwal();