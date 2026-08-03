let token = localStorage.getItem("token");

if(!token){
    window.location.href="index.html";
}

function loadKelasOptions(){

    fetch("http://localhost:8080/kelas", {
        headers:{
            "Authorization":"Bearer " + token
        }
    })

    .then(response => response.json())

    .then(data => {

        let opt = `<option value="">Pilih Kelas</option>`;

        (data || []).forEach(k => {
            opt += `<option value="${k.id}">${k.nama} - ${k.nama_jurusan}</option>`;
        });

        document.getElementById("kelas").innerHTML = opt;

    });

}

function loadMahasiswa(){

    let tampilkanSemua = document.getElementById("tampilkanSemua").checked;
    let url = "http://localhost:8080/mahasiswa";

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

loadKelasOptions();
loadMahasiswa();

function bukaModalTambah(){

    document.getElementById("idMahasiswa").value = "";
    document.getElementById("nama").value = "";
    document.getElementById("kelas").value = "";

    document.querySelector("#modalMahasiswa .modal-title").innerText = "Tambah Mahasiswa";

}

function simpanMahasiswa(){

    let id = document.getElementById("idMahasiswa").value;
    let nama = document.getElementById("nama").value;
    let kelasId = document.getElementById("kelas").value;

    if(kelasId === ""){
        alert("Kelas wajib dipilih");
        return;
    }

    let token = localStorage.getItem("token");


    let method = "POST";
    let url = "http://localhost:8080/mahasiswa";


    if(id != ""){

        method = "PUT";
        url = "http://localhost:8080/mahasiswa/" + id;

    }


    fetch(url, {

        method: method,

        headers:{
            "Content-Type":"application/json",
            "Authorization":"Bearer " + token
        },

        body: JSON.stringify({

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
        "http://localhost:8080/mahasiswa?status=semua",
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

        document.getElementById("nama").value = mhs.nama;

        document.getElementById("kelas").value = mhs.kelas_id || "";

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
            "http://localhost:8080/mahasiswa/" + id,
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
        "http://localhost:8080/mahasiswa/" + id + "/aktifkan",
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