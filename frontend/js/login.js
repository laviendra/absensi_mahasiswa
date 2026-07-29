function login(){

    let username = document.getElementById("username").value;
    let password = document.getElementById("password").value;


    fetch("http://localhost:8080/login", {

        method: "POST",

        headers:{
            "Content-Type":"application/json"
        },

        body: JSON.stringify({

            username: username,
            password: password

        })

    })


    .then(response => response.json())


    .then(data => {


        if(data.token){

            localStorage.setItem(
                "token",
                data.token
            );


            window.location.href =
            "dashboard.html";


        }else{

            document.getElementById("message").innerHTML =
            "Login gagal";

        }


    })


    .catch(error => {

        console.log(error);

        document.getElementById("message").innerHTML =
        "Server tidak terhubung";

    });

}