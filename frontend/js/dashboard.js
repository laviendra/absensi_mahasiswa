let token = localStorage.getItem("token");


if(!token){

    window.location.href = "index.html";

}



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


})



function logout(){

    localStorage.removeItem("token");

    window.location.href="index.html";

}