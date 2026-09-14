import { invoke } from '@tauri-apps/api/core';

export async function login() {
    let username = document.getElementById("username").value;
    let password = document.getElementById("pwd").value;

    try {
      await invoke("command_login", {
        username: username,
        password: password,
      });

      console.log("Login successful");
      window.location.href = "messages.html";
    } catch (error) {
      document.getElementById("error").innerText = "Invalid credentials";
    }
}
window.login = login;

export async function signup() {
    // username, email, display_name, password
    let username = document.getElementById("username").value;
    let password = document.getElementById("pwd").value;
    let display_name = document.getElementById("displayname").value;
    let email = document.getElementById("email").value;

    try {
      await invoke("command_register", {
        username: username,
        password: password,
        displayName: display_name,
        email: email,
      });
      
      console.log("Register successful");
      window.location.href = "messages.html";
    } catch (error) {
      document.getElementById("error").innerText = String(error);
    }
}
window.signup = signup;
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
