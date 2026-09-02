function getLogInData(){
    let usernameInput = document.getElementById("username");
    let passwordInput = document.getElementById("pwd");
    
    let username = usernameInput.value;
    let password = passwordInput.value;

    if (username === "" || password === "") {
        document.getElementById("error").innerText = "Please enter your credidentials";
        return;
    }

    if(checkLogInData()){
        window.location.href = "messaging-screen.html"
    }else{
        document.getElementById("error").innerText = "Invalid credidentials";
    }
}

function getSignUpData(){
    let fname = document.getElementById("fname").value;
    let lname = document.getElementById("lname").value;
    let email = document.getElementById("mail").value;
    let username = document.getElementById("username").value;
    let pwd1 = document.getElementById("pwd1").value;
    let pwd2 = document.getElementById("pwd2").value;

    if(!(pwd1 === pwd2)){
        document.getElementById("error").innerText = "Passwords mismatch";
    }else{
        if(sendSignUpData()){
            window.location.href = "authentication-screen.html";
        }
    }
}

function sendSignUpData(){
    return true;
}

function checkLogInData(){
    return true;
}

function checkCode(code){
    return true;
}

function getCode(){
    let code = document.getElementById("code").value;

    if(checkCode(code)){
        window.location.href = "login-screen.html";
    }else{
        document.getElementById("error").innerText = "Invalid code";
    }
}

function openConversation(){
    
}