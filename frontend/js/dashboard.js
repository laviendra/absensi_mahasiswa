let token = localStorage.getItem("token");


if(!token){

    window.location.href = "index.html";

}


document.getElementById("tanggalHariIni").innerText =
    new Date().toLocaleDateString("id-ID", {
        weekday: "long",
        year: "numeric",
        month: "long",
        day: "numeric"
    });


fetch("http://localhost:8080/dashboard", {

    method:"GET",

    headers:{
        "Authorization":"Bearer " + token
    }

})


.then(response => response.json())


.then(data => {


    document.getElementById("totalMahasiswa").innerHTML =
    data.total_mahasiswa;


    document.getElementById("totalDosen").innerHTML =
    data.total_dosen;


    document.getElementById("totalKelas").innerHTML =
    data.total_kelas;


    document.getElementById("pertemuanHariIni").innerHTML =
    data.pertemuan_hari_ini;


});