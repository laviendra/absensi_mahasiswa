let token = localStorage.getItem("token");

if (!token) {
    window.location.href = "index.html";
}

const API = "http://localhost:8080";

let daftarDosenCache = [];
let daftarMatkulCache = [];

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
        loadKelas();
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

function loadKelasOptionsForKelasForm() {

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

        daftarDosenCache = data || [];

        let tabel = "";

        daftarDosenCache.forEach((d, i) => {
            tabel += `
            <tr>
                <td>${i + 1}</td>
                <td>${d.nama}</td>
                <td>${d.username}</td>
                <td>
                    <button class="btn btn-warning btn-sm" onclick="editDosen(${d.id})">Edit</button>
                    <button class="btn btn-danger btn-sm" onclick="hapusDosen(${d.id})">Hapus</button>
                </td>
            </tr>
            `;
        });

        document.getElementById("dataDosen").innerHTML =
            tabel || `<tr><td colspan="4" class="text-center text-muted">Belum ada data</td></tr>`;

    });

}

// isi checkbox mata kuliah di form Dosen, sekalian tandain yang udah dipilih (edit)
function isiPilihanMatkulDosen(selectedIds) {

    fetch(API + "/mata-kuliah", { headers: authHeader() })
    .then(r => r.json())
    .then(data => {

        daftarMatkulCache = data || [];

        let html = "";

        daftarMatkulCache.forEach(mk => {

            let dicentang = selectedIds.includes(mk.id) ? "checked" : "";

            html += `
            <div class="form-check">
                <input
                    class="form-check-input matkul-dosen-checkbox"
                    type="checkbox"
                    value="${mk.id}"
                    id="cbMatkulDosen${mk.id}"
                    ${dicentang}>
                <label class="form-check-label" for="cbMatkulDosen${mk.id}">
                    ${mk.kode} - ${mk.nama}
                </label>
            </div>
            `;

        });

        document.getElementById("matkulDosenCheckbox").innerHTML =
            html || `<p class="text-muted mb-0 small">Belum ada mata kuliah, tambahkan dulu di tab Mata Kuliah</p>`;

    });

}

function bukaModalDosen() {
    document.getElementById("idDosen").value = "";
    document.getElementById("namaDosen").value = "";
    document.getElementById("usernameDosen").value = "";
    document.getElementById("passwordDosen").value = "";
    isiPilihanMatkulDosen([]);
    new bootstrap.Modal(document.getElementById("modalDosen")).show();
}

function editDosen(id) {

    let d = daftarDosenCache.find(item => item.id === id);

    if (!d) return;

    document.getElementById("idDosen").value = d.id;
    document.getElementById("namaDosen").value = d.nama;
    document.getElementById("usernameDosen").value = d.username;
    document.getElementById("passwordDosen").value = "";

    isiPilihanMatkulDosen(d.mata_kuliah_ids || []);

    new bootstrap.Modal(document.getElementById("modalDosen")).show();

}

function simpanDosen() {

    let id = document.getElementById("idDosen").value;
    let nama = document.getElementById("namaDosen").value;
    let username = document.getElementById("usernameDosen").value;
    let password = document.getElementById("passwordDosen").value;

    let mataKuliahIds = Array.from(
        document.querySelectorAll(".matkul-dosen-checkbox:checked")
    ).map(cb => Number(cb.value));

    if (!id && !password) {
        alert("Password wajib diisi untuk dosen baru");
        return;
    }

    let method = id ? "PUT" : "POST";
    let url = id ? API + "/dosen/" + id : API + "/dosen";

    fetch(url, {
        method: method,
        headers: authHeader(true),
        body: JSON.stringify({ nama, username, password, mata_kuliah_ids: mataKuliahIds })
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

        daftarMatkulCache = data || [];

        let tabel = "";

        daftarMatkulCache.forEach((mk, i) => {
            tabel += `
            <tr>
                <td>${i + 1}</td>
                <td>${mk.kode}</td>
                <td>${mk.nama}</td>
                <td>
                    <button class="btn btn-warning btn-sm" onclick="editMatkul(${mk.id})">Edit</button>
                    <button class="btn btn-danger btn-sm" onclick="hapusMatkul(${mk.id})">Hapus</button>
                </td>
            </tr>
            `;
        });

        document.getElementById("dataMatkul").innerHTML =
            tabel || `<tr><td colspan="4" class="text-center text-muted">Belum ada data</td></tr>`;

    });

}

// isi checkbox kelas di form Mata Kuliah, sekalian tandain yang udah dipilih (edit)
function isiPilihanKelasMatkul(selectedIds) {

    fetch(API + "/kelas", { headers: authHeader() })
    .then(r => r.json())
    .then(data => {

        let html = "";

        (data || []).forEach(k => {

            let dicentang = selectedIds.includes(k.id) ? "checked" : "";

            html += `
            <div class="form-check">
                <input
                    class="form-check-input kelas-matkul-checkbox"
                    type="checkbox"
                    value="${k.id}"
                    id="cbKelasMatkul${k.id}"
                    ${dicentang}>
                <label class="form-check-label" for="cbKelasMatkul${k.id}">
                    ${k.nama} - ${k.nama_jurusan}
                </label>
            </div>
            `;

        });

        document.getElementById("kelasMatkulCheckbox").innerHTML =
            html || `<p class="text-muted mb-0 small">Belum ada kelas, tambahkan dulu di tab Kelas</p>`;

    });

}

function bukaModalMatkul() {
    document.getElementById("idMatkul").value = "";
    document.getElementById("kodeMatkul").value = "";
    document.getElementById("namaMatkul").value = "";
    isiPilihanKelasMatkul([]);
    new bootstrap.Modal(document.getElementById("modalMatkul")).show();
}

function editMatkul(id) {

    let mk = daftarMatkulCache.find(item => item.id === id);

    if (!mk) return;

    document.getElementById("idMatkul").value = mk.id;
    document.getElementById("kodeMatkul").value = mk.kode;
    document.getElementById("namaMatkul").value = mk.nama;

    isiPilihanKelasMatkul(mk.kelas_ids || []);

    new bootstrap.Modal(document.getElementById("modalMatkul")).show();

}

function simpanMatkul() {

    let id = document.getElementById("idMatkul").value;
    let kode = document.getElementById("kodeMatkul").value;
    let nama = document.getElementById("namaMatkul").value;

    let kelasIds = Array.from(
        document.querySelectorAll(".kelas-matkul-checkbox:checked")
    ).map(cb => Number(cb.value));

    let method = id ? "PUT" : "POST";
    let url = id ? API + "/mata-kuliah/" + id : API + "/mata-kuliah";

    fetch(url, {
        method: method,
        headers: authHeader(true),
        body: JSON.stringify({ kode, nama, kelas_ids: kelasIds })
    })
    .then(r => r.json())
    .then(data => {
        alert(data.message || data.error || "Terjadi kesalahan");
        bootstrap.Modal.getInstance(document.getElementById("modalMatkul")).hide();
        loadMatkul();
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


// ================= JADWAL (cascading: Dosen -> Mata Kuliah -> Kelas) =================

function loadDosenJadwalOptions() {

    fetch(API + "/dosen", { headers: authHeader() })
    .then(r => r.json())
    .then(data => {
        let opt = `<option value="">Pilih Dosen</option>`;
        (data || []).forEach(d => opt += `<option value="${d.id}">${d.nama}</option>`);
        document.getElementById("dosenJadwal").innerHTML = opt;
    });

}

function loadMatkulJadwalByDosen() {

    let dosenId = document.getElementById("dosenJadwal").value;
    let matkulSelect = document.getElementById("matkulJadwal");
    let kelasSelect = document.getElementById("kelasJadwal");

    kelasSelect.innerHTML = `<option value="">Pilih mata kuliah dulu</option>`;
    kelasSelect.disabled = true;

    if (!dosenId) {
        matkulSelect.innerHTML = `<option value="">Pilih dosen dulu</option>`;
        matkulSelect.disabled = true;
        return;
    }

    fetch(API + "/dosen/" + dosenId + "/mata-kuliah", { headers: authHeader() })
    .then(r => r.json())
    .then(data => {

        matkulSelect.disabled = false;

        let opt = `<option value="">Pilih Mata Kuliah</option>`;

        (data || []).forEach(mk => opt += `<option value="${mk.id}">${mk.kode} - ${mk.nama}</option>`);

        matkulSelect.innerHTML =
            (data && data.length) ? opt : `<option value="">Dosen ini belum diampu ke mata kuliah manapun</option>`;

    });

}

function loadKelasJadwalByMatkul() {

    let matkulId = document.getElementById("matkulJadwal").value;
    let kelasSelect = document.getElementById("kelasJadwal");

    if (!matkulId) {
        kelasSelect.innerHTML = `<option value="">Pilih mata kuliah dulu</option>`;
        kelasSelect.disabled = true;
        return;
    }

    fetch(API + "/mata-kuliah/" + matkulId + "/kelas", { headers: authHeader() })
    .then(r => r.json())
    .then(data => {

        kelasSelect.disabled = false;

        let opt = `<option value="">Pilih Kelas</option>`;

        (data || []).forEach(k => opt += `<option value="${k.id}">${k.nama}</option>`);

        kelasSelect.innerHTML =
            (data && data.length) ? opt : `<option value="">Mata kuliah ini belum dipakai kelas manapun</option>`;

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

    document.getElementById("matkulJadwal").innerHTML = `<option value="">Pilih dosen dulu</option>`;
    document.getElementById("matkulJadwal").disabled = true;

    document.getElementById("kelasJadwal").innerHTML = `<option value="">Pilih mata kuliah dulu</option>`;
    document.getElementById("kelasJadwal").disabled = true;

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
loadKelasOptionsForKelasForm();
loadKelas();
loadDosen();
loadMatkul();
loadDosenJadwalOptions();
loadJadwal();