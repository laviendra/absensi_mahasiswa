let token = localStorage.getItem("token");


if(!token){

    window.location.href="index.html";

}



fetch("http://localhost:8080/mahasiswa",{

    method:"GET",

    headers:{
        "Authorization":"Bearer " + token
    }

})


.then(response => response.json())


.then(data => {


    let tabel = "";


    data.forEach((mhs,index)=>{


        tabel += `

        <tr>

        <td>${index+1}</td>

        <td>${mhs.nama}</td>

        <td>${mhs.jurusan}</td>

        <td>

        <button 
        class="btn btn-warning btn-sm"
        onclick="editMahasiswa(${mhs.id})">
        Edit
        </button>


        <button 
        class="btn btn-danger btn-sm"
        onclick="hapusMahasiswa(${mhs.id})">
        Hapus
        </button>

        </td>

        </tr>

        `;


    });



    document.getElementById(
        "dataMahasiswa"
    ).innerHTML = tabel;


});

function simpanMahasiswa(){

    let id = document.getElementById("idMahasiswa").value;
    let nama = document.getElementById("nama").value;
    let jurusan = document.getElementById("jurusan").value;

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
            jurusan: jurusan

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
        "http://localhost:8080/mahasiswa",
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

        document.getElementById("jurusan").value = mhs.jurusan;


        let modal = new bootstrap.Modal(
            document.getElementById("modalMahasiswa")
        );


        modal.show();


    });

}

function hapusMahasiswa(id){

    let token = localStorage.getItem("token");


    if(confirm("Yakin hapus mahasiswa ini?")){


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