let token = localStorage.getItem("token");

if (!token) {
    window.location.href = "index.html";
}

fetch("http://localhost:8080/mahasiswa",{

headers:{
"Authorization":"Bearer "+token
}

})

.then(response=>response.json())

.then(data=>{

let option="";

data.forEach(mhs=>{

option+=`
<option value="${mhs.id}">
${mhs.nama}
</option>
`;

});

document.getElementById("mahasiswa").innerHTML=option;

});

fetch("http://localhost:8080/absensi", {

    method: "GET",

    headers: {
        "Authorization": "Bearer " + token
    }

})

.then(response => response.json())

.then(data => {

    let tabel = "";

    data.forEach((item, index) => {

        let tombol = "-";

        if(item.jam_pulang == ""){

            tombol = `
            <button
                class="btn btn-warning btn-sm"
                onclick="absenPulang(${item.id})">
                Pulang
            </button>
            `;
        }

        tabel += `

        <tr>

            <td>${index + 1}</td>
            <td>${item.nama}</td>
            <td>${item.tanggal}</td>
            <td>${item.jam_masuk}</td>
            <td>${item.jam_pulang}</td>
            <td>${item.status}</td>
            <td>${tombol}</td>

        </tr>

        `;

    });

    document.getElementById("dataAbsensi").innerHTML = tabel;

});

function absenMasuk(){

let id=document.getElementById("mahasiswa").value;

fetch("http://localhost:8080/absensi/masuk",{

method:"POST",

headers:{
"Content-Type":"application/json",
"Authorization":"Bearer "+token
},

body:JSON.stringify({

mahasiswa_id:Number(id)

})

})

.then(response=>response.json())

.then(data=>{

alert(data.message);

location.reload();

});

}

function absenPulang(id){

    const konfirmasi = confirm("Yakin ingin melakukan absensi pulang?");

    if(!konfirmasi){
        return;
    }

    fetch("http://localhost:8080/absensi/pulang/" + id, {

        method: "PUT",

        headers:{
            "Authorization":"Bearer " + token
        }

    })

    .then(response => response.json())

    .then(data => {

        alert(data.message);

        if(data.message === "Absensi pulang berhasil"){
            location.reload();
        }

    })

    .catch(error => {

        console.log(error);
        alert("Terjadi kesalahan.");

    });

}