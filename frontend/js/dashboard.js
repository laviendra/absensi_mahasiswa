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


    document.getElementById("hadir").innerHTML =
    data.hadir_hari_ini;


    document.getElementById("terlambat").innerHTML =
    data.terlambat;


    document.getElementById("belumAbsen").innerHTML =
    data.belum_absen;


});