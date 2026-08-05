let token = localStorage.getItem("token");

if(!token){
    window.location.href="index.html";
}

function loadJurusanOptions(){

    fetch(API_BASE + "/jurusan", {
        headers:{
            "Authorization":"Bearer " + token
        }
    })

    .then(response => response.json())

    .then(data => {

        let opt = `<option value="">Pilih Jurusan</option>`;

        (data || []).forEach(j => {
            opt += `<option value="${j.id}">${j.nama}</option>`;
        });

        document.getElementById("jurusanMahasiswa").innerHTML = opt;

    });

}

function loadKelasByJurusan(preselectKelasId){

    let jurusanId = document.getElementById("jurusanMahasiswa").value;
    let kelasSelect = document.getElementById("kelas");

    if(jurusanId === ""){
        kelasSelect.innerHTML = `<option value="">Pilih jurusan dulu</option>`;
        kelasSelect.disabled = true;
        return;
    }

    kelasSelect.disabled = false;

    fetch(API_BASE + "/kelas?jurusan_id=" + jurusanId, {
        headers:{
            "Authorization":"Bearer " + token
        }
    })

    .then(response => response.json())

    .then(data => {

        let opt = `<option value="">Pilih Kelas</option>`;

        (data || []).forEach(k => {
            opt += `<option value="${k.id}">${k.nama}</option>`;
        });

        kelasSelect.innerHTML = opt;

        if(preselectKelasId){
            kelasSelect.value = preselectKelasId;
        }

    });

}

function loadMahasiswa(){

    let tampilkanSemua = document.getElementById("tampilkanSemua").checked;
    let url = API_BASE + "/mahasiswa";

    if(tampilkanSemua){
        url += "?status=semua";
    }

    fetch(url, {

        method:"GET",

        headers:{
            "Authorization":"Bearer " + token
        }

    })

    .then(response => response.json())

    .then(data => {

        let tabel = "";

        (data || []).forEach((mhs,index)=>{

            let badgeStatus = mhs.status === "aktif" ? "success" : "secondary";

            let aksi = "";

            if(mhs.status === "aktif"){

                aksi = `
                <button 
                class="btn btn-warning btn-sm"
                onclick="editMahasiswa(${mhs.id})">
                Edit
                </button>

                <button 
                class="btn btn-danger btn-sm"
                onclick="nonaktifkanMahasiswa(${mhs.id})">
                Nonaktifkan
                </button>
                `;

            } else {

                aksi = `
                <button 
                class="btn btn-success btn-sm"
                onclick="aktifkanMahasiswa(${mhs.id})">
                Aktifkan
                </button>
                `;
            }

            tabel += `

            <tr>

            <td>${index+1}</td>

            <td>${mhs.nim || "-"}</td>

            <td>${mhs.nama}</td>

            <td>${mhs.nama_kelas || "-"}</td>

            <td>${mhs.nama_jurusan || "-"}</td>

            <td><span class="badge bg-${badgeStatus}">${mhs.status}</span></td>

            <td>${aksi}</td>

            </tr>

            `;

        });

        document.getElementById("dataMahasiswa").innerHTML = tabel;

    });

}

loadJurusanOptions();
loadMahasiswa();
loadKelasCard();

function loadKelasCard() {

    fetch(API_BASE + "/kelas", {
        headers: {
            "Authorization": "Bearer " + token
        }
    })

    .then(response => response.json())

    .then(data => {

        let html = "";

        (data || []).forEach(k => {

            html += `
            <div class="col-md-3">
                <div class="card shadow-sm" role="button" onclick="tampilkanMahasiswaKelas(${k.id}, '${k.nama}')">
                    <div class="card-body text-center">
                        <h5 class="mb-1">${k.nama}</h5>
                        <p class="text-muted small mb-0">${k.nama_jurusan}</p>
                    </div>
                </div>
            </div>
            `;

        });

        document.getElementById("daftarKelasCard").innerHTML =
            html || `<p class="text-muted">Belum ada data kelas. Tambahkan dulu di halaman Data Master.</p>`;

    });

}

function tampilkanMahasiswaKelas(kelasId, namaKelas) {

    document.getElementById("daftarKelasCard").classList.add("d-none");
    document.getElementById("detailKelas").classList.remove("d-none");
    document.getElementById("judulDetailKelas").innerText = "Mahasiswa Kelas " + namaKelas;

    fetch(API_BASE + "/mahasiswa?kelas_id=" + kelasId, {
        headers: {
            "Authorization": "Bearer " + token
        }
    })

    .then(response => response.json())

    .then(data => {

        let tabel = "";

        (data || []).forEach((mhs, index) => {
            tabel += `
            <tr>
                <td>${index + 1}</td>
                <td>${mhs.nim || "-"}</td>
                <td>${mhs.nama}</td>
                <td><span class="badge bg-success">${mhs.status}</span></td>
            </tr>
            `;
        });

        document.getElementById("dataMahasiswaKelas").innerHTML =
            tabel || `<tr><td colspan="4" class="text-center text-muted">Belum ada mahasiswa di kelas ini</td></tr>`;

    });

}

function kembaliKeDaftarKelas() {
    document.getElementById("detailKelas").classList.add("d-none");
    document.getElementById("daftarKelasCard").classList.remove("d-none");
}

function bukaModalTambah(){

    document.getElementById("idMahasiswa").value = "";
    document.getElementById("nim").value = "";
    document.getElementById("nama").value = "";
    document.getElementById("jurusanMahasiswa").value = "";
    document.getElementById("kelas").innerHTML = `<option value="">Pilih jurusan dulu</option>`;
    document.getElementById("kelas").disabled = true;

    document.querySelector("#modalMahasiswa .modal-title").innerText = "Tambah Mahasiswa";

}

function simpanMahasiswa(){

    let id = document.getElementById("idMahasiswa").value;
    let nim = document.getElementById("nim").value;
    let nama = document.getElementById("nama").value;
    let kelasId = document.getElementById("kelas").value;

    if(nim === ""){
        alert("NIM wajib diisi");
        return;
    }

    if(kelasId === ""){
        alert("Kelas wajib dipilih");
        return;
    }

    let token = localStorage.getItem("token");


    let method = "POST";
    let url = API_BASE + "/mahasiswa";


    if(id != ""){

        method = "PUT";
        url = API_BASE + "/mahasiswa/" + id;

    }


    fetch(url, {

        method: method,

        headers:{
            "Content-Type":"application/json",
            "Authorization":"Bearer " + token
        },

        body: JSON.stringify({

            nim: nim,
            nama: nama,
            kelas_id: Number(kelasId)

        })

    })


    .then(response=>response.json())


    .then(data=>{

        alert(data.message || data.error || "Terjadi kesalahan");

        location.reload();

    });

}

function editMahasiswa(id){

    let token = localStorage.getItem("token");


    fetch(
        API_BASE + "/mahasiswa?status=semua",
        {
            headers:{
                "Authorization":"Bearer " + token
            }
        }
    )


    .then(response=>response.json())


    .then(data=>{


        let mhs = data.find(
            item => item.id == id
        );


        document.getElementById("idMahasiswa").value = mhs.id;

        document.getElementById("nim").value = mhs.nim || "";

        document.getElementById("nama").value = mhs.nama;

        document.getElementById("jurusanMahasiswa").value = mhs.jurusan_id || "";

        loadKelasByJurusan(mhs.kelas_id);

        document.querySelector("#modalMahasiswa .modal-title").innerText = "Edit Mahasiswa";


        let modal = new bootstrap.Modal(
            document.getElementById("modalMahasiswa")
        );


        modal.show();


    });

}

function nonaktifkanMahasiswa(id){

    let token = localStorage.getItem("token");


    if(confirm("Nonaktifkan mahasiswa ini? Data absensinya tetap tersimpan.")){


        fetch(
            API_BASE + "/mahasiswa/" + id,
            {
                method:"DELETE",

                headers:{
                    "Authorization":"Bearer " + token
                }
            }
        )


        .then(response=>response.json())


        .then(data=>{

            alert(data.message || data.error || "Terjadi kesalahan");

            location.reload();

        });


    }

}

function aktifkanMahasiswa(id){

    let token = localStorage.getItem("token");

    fetch(
        API_BASE + "/mahasiswa/" + id + "/aktifkan",
        {
            method:"PUT",

            headers:{
                "Authorization":"Bearer " + token
            }
        }
    )

    .then(response=>response.json())

    .then(data=>{

        alert(data.message || data.error || "Terjadi kesalahan");

        location.reload();

    });

}