function login(){

    let identifier = document.getElementById("username").value;
    let password = document.getElementById("password").value;


    fetch(API_BASE + "/login-universal", {

        method: "POST",

        headers:{
            "Content-Type":"application/json"
        },

        body: JSON.stringify({

            identifier: identifier,
            password: password

        })

    })


    .then(response => response.json())


    .then(data => {


        if(!data.token){

            document.getElementById("message").innerHTML =
            data.message || "Login gagal";

            return;

        }


        if(data.role === "admin"){

            localStorage.setItem("token", data.token);
            window.location.href = "dashboard.html";

        } else if(data.role === "dosen"){

            localStorage.setItem("dosenToken", data.token);
            localStorage.setItem("dosenNama", data.nama);
            window.location.href = "dosen.html";

        } else if(data.role === "mahasiswa"){

            localStorage.setItem("mahasiswaToken", data.token);
            localStorage.setItem("mahasiswaNama", data.nama);
            window.location.href = "mahasiswa-login.html";

        }


    })


    .catch(error => {

        console.log(error);

        document.getElementById("message").innerHTML =
        "Server tidak terhubung";

    });

}