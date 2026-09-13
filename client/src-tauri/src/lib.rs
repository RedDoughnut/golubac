// Learn more about Tauri commands at https://tauri.app/develop/calling-rust/
use tauri::State;
use tokio::sync::Mutex;

use reqwest;
use serde::{Deserialize, Serialize};

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    let chat_client = ChatClient::new("127.0.0.1".to_string());

    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .manage(Mutex::new(chat_client))
        .invoke_handler(tauri::generate_handler![command_register, command_login])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

#[tauri::command]
async fn command_register(
    state: State<'_, Mutex<ChatClient>>,
    username: String,
    email: String,
    display_name: String,
    password: String,
) -> Result<(), String> {
    let mut client = state.lock().await;
    client
        .register(username, email, display_name, password)
        .await
        .map_err(|e| e.to_string())
}

#[tauri::command]
async fn command_login(
    state: State<'_, Mutex<ChatClient>>,
    username: String,
    password: String,
) -> Result<(), String> {
    let mut client = state.lock().await;
    client
        .login(username, password)
        .await
        .map_err(|e| e.to_string())
}

#[derive(Serialize)]
struct RegisterPayload {
    email: String,
    username: String,
    displayname: String,
    password: String,
}
#[derive(Deserialize)]
struct RegisterResponse {
    refreshtoken: String,
}

#[derive(Serialize)]
struct LoginPayload {
    username: String,
    password: String,
}
#[derive(Deserialize)]
struct LoginResponse {
    refreshtoken: String,
}

#[derive(Serialize)]
struct RefreshSessionPayload {
    refreshtoken: String,
}
#[derive(Deserialize)]
struct RefreshSessionResponse {
    sessiontoken: String,
}

#[derive(Serialize)]
struct SendTextMessagePayload {
    to: String,
    message: String,
}
// TODO: send images

struct ChatClient {
    username: String,
    display_name: String,
    session_token: String,
    refresh_token: String,
    server_ip: String,
    http_client: reqwest::Client,
}

impl ChatClient {
    fn new(server_ip: String) -> ChatClient {
        // TODO: make this do register and return an Option mabye
        Self {
            username: String::new(),
            display_name: String::new(),
            session_token: String::new(),
            refresh_token: String::new(),
            server_ip,
            http_client: reqwest::Client::new(),
        }
    }

    

    pub async fn register(
        &mut self,
        username: String,
        email: String,
        display_name: String,
        password: String,
    ) -> Result<(), Box<dyn std::error::Error>> {
        let payload = RegisterPayload {
            username: username.clone(),
            email,
            displayname: display_name.clone(),
            password,
        };

        let url = format!("https://{}/register", self.server_ip);

        let response = self.http_client.post(&url).json(&payload).send().await?;

        if !response.status().is_success() {
            let error_text = response.text().await?;
            eprintln!("register failed: {}", error_text);
            return Err("register failed".into());
        }

        self.refresh_token = response.json::<RegisterResponse>().await?.refreshtoken;

        println!("Register successful!");
        Ok(())
    }

    pub async fn login(
        &mut self,
        username: String,
        password: String,
    ) -> Result<(), Box<dyn std::error::Error>> {
        let payload = LoginPayload {
            username: username.clone(),
            password: password.clone(),
        };

        let url = format!("https://{}/login", self.server_ip);

        let response = self.http_client.post(&url).json(&payload).send().await?;

        if !response.status().is_success() {
            let error_text = response.text().await?;
            eprintln!("login failed: {}", error_text);
            return Err("login failed".into());
        }

        self.refresh_token = response.json::<LoginResponse>().await?.refreshtoken;

        Ok(())
    }
    

    async fn refresh_session(&mut self) -> Result<(), Box<dyn std::error::Error>> {
        // TODO: implement saving the refresh token on device
        let payload = RefreshSessionPayload {
            refreshtoken: self.refresh_token.clone(),
        };

        let url = format!("https://{}/refresh-session-token", self.server_ip);

        let response = self.http_client.post(&url).json(&payload).send().await?;

        if !response.status().is_success() {
            let error_text = response.text().await?;
            eprintln!("refreshing session failed: {}", error_text);
            return Err("refreshing session failed".into());
        }

        self.session_token = response
            .json::<RefreshSessionResponse>()
            .await?
            .sessiontoken;

        Ok(())
    }
}
