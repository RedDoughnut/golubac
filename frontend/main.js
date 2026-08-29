function getData(){
    let usernameInput = document.getElementById("username");
    let passwordInput = document.getElementById("pwd");
    
    let username = usernameInput.value;
    let password = passwordInput.value;

    console.log("Korisnik:", username);
    console.log("Lozinka:", password);

    if (username === "" || password === "") {
        alert("Please enter your credidentials!");
    }
}