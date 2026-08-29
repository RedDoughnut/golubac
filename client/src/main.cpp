#include <iostream>
#include <fstream>
#include <string>

#include <cpr/cpr.h>

#include <nlohmann/json.hpp>


/*
 * Used when the client sends a message
 */
class SentMessage {
  public:
    /*
     * Pass empty string to show an absence of a value
     */
    SentMessage(const std::u8string& text, const std::string& image_file_path) {
        has_text = true;
        has_image = true;
        if(text.empty()) { has_text = false; }
        if(image_file_path.empty()) { has_image = false; }

        this->text = text;
        this->image_file_path = image_file_path;

    }

    const std::string& getImageFilePath() const {
        return image_file_path;
    }
    const std::u8string& getText() const {
        return text;
    }
    const bool& hasText() const {
        return has_text;
    }
    const bool& hasImage() const {
        return has_image;
    }
  private:
    std::u8string text;
    bool has_text;
    std::string image_file_path;
    bool has_image;
};

/*
 * Used for storing sent and recieved messages
 */
struct ChatMessage {
    std::u8string text;
    std::string image_file_path;
    uint64_t unix_timestamp;
};

// TODO: chat class


enum class HTTPStatus {
    INFORMATIONAL,
    SUCCESS,
    REDIRECTION,
    CLIENT_ERROR,
    SERVER_ERROR,
    NONE
};
HTTPStatus responseStatusCodeType(cpr::Response response) {
    auto code = response.status_code;
    if       (100 <= code && code < 200) {
        return HTTPStatus::INFORMATIONAL;
    } else if(200 <= code && code < 300) {
        return HTTPStatus::SUCCESS;
    } else if(300 <= code && code < 400) {
        return HTTPStatus::REDIRECTION;
    } else if(400 <= code && code < 500) {
        return HTTPStatus::CLIENT_ERROR;
    } else if(500 <= code && code < 600) {
        return HTTPStatus::SERVER_ERROR;
    } else {
        return HTTPStatus::NONE;
    }
}

class ChatClient {
  public:
    std::string username;
    std::string display_name;
    std::string session_token{"0"};
    std::string refresh_token{"0"};
    std::string server_ip;

    ChatClient(const std::string& server_ip) {
        this->server_ip = server_ip;
    }

    bool registerUser(const std::string& username, const std::string& email, const std::string& display_name, const std::string& password) { // make username be alphanumeric and email text checking
        nlohmann::json register_json;
        register_json["email"] = email;
        register_json["username"] = username;
        register_json["displayname"] = display_name;
        register_json["password"] = password;

        auto response = cpr::Post(
            cpr::Url{server_ip + "/register"},
            cpr::Header{{"Content-Type", "application/json"}},
            cpr::Body{register_json.dump()}
        );

        if(!checkResponse(response, "register")) {
            std::cout << "register failed" << std::endl;
            return false;
        }

        std::cout << "register successful" << std::endl;
        this->username = username;
        this->display_name = display_name;

        nlohmann::json response_json = nlohmann::json::parse(response.text);
        refresh_token = response_json["refreshtoken"];

        std::ofstream refresh_token_file_stream("../data/refresh-token.txt");
        refresh_token_file_stream << refresh_token;

        return true;
    }
    /*
     * Logs in as a user
     */
    bool login(const std::string& username, const std::string& password) {
        nlohmann::json body_json;
        body_json["username"] = username;
        body_json["password"] = password;
        auto response = cpr::Post(
            cpr::Url{server_ip + "/login"},
            cpr::Header{{"Content-Type", "application/json"}},
            cpr::Body{body_json.dump()}
        );

        if (!checkResponse(response, "login")) {
            return false;
        }

        nlohmann::json response_body = nlohmann::json::parse(response.text);
        refresh_token = response_body["refreshtoken"];

        std::ofstream refresh_token_file_stream("../data/refresh-token.txt");
        refresh_token_file_stream << refresh_token;
        return true;
    }

    bool refreshSession() {
        nlohmann::json refresh_session_json;
        refresh_session_json["refreshtoken"] = refresh_token;
        auto response = cpr::Post(
            cpr::Url{server_ip + "/refresh-session-token"},
            cpr::Header{{"Content-Type", "application/json"}},
            cpr::Body{refresh_session_json.dump()}
        );
        if (response.error.code != cpr::ErrorCode::OK) {
            std::cerr << "refreshing session token" << " failed (connection error): " << response.error.message << "\n";
            return false;
        }

        if(response.error.message == "invalid-refresh-token" && response.status_code == 401) {
            std::cerr << "invalid refresh token, try logging out and loggin in" << std::endl;
            return false;
        }

        if(responseStatusCodeType(response) != HTTPStatus::SUCCESS) {
            std::cerr << "refreshing session token" << " failed (HTTP " << response.status_code << "): " << response.text << std::endl;
            return false;
        }

        nlohmann::json response_json = nlohmann::json::parse(response.text);
        session_token = response_json["sessiontoken"];

        return true;
    }

    /*
     * Sends a text message *
     * *should not be used in favour of sendMessage()*
     */
    bool sendTextMessage(const std::string& recipient, const std::string& message) { // TODO: add message class that can hold an image or other things
        nlohmann::json j;
        j["sessiontoken"] = session_token;
        j["recipient"] = recipient;
        j["message"] = message;

        cpr::Body body{j.dump()};

        auto response = cpr::Post(
            cpr::Url{server_ip + "/send-text-message"},
            cpr::Header{{"Content-Type", "application/json"}},
            body
        );
        if (response.error.code != cpr::ErrorCode::OK) {
            std::cerr << "sending message" << " failed (connection error): " << response.error.message << "\n";
            return false;
        }
        if(responseStatusCodeType(response) != HTTPStatus::SUCCESS) {
            try {
                nlohmann::json response_json = nlohmann::json::parse(response.text);
                if (response_json.contains("error") && response_json["error"] == "invalid_token" && response.status_code == 401) {
                    refreshSession();
                    std::cerr << "session token refreshed, message not sent" << std::endl;
                    return false;
                }
            } catch (const nlohmann::json::parse_error& e) {
                std::cerr << "send-message failed (HTTP " << response.status_code << "): " << response.text << std::endl;
            }
            return false;
        }


        return true;
    }
    
    bool sendMessage(const std::string& recipient, const SentMessage& message) {
        const auto& u8_text = message.getText();
        std::string text_utf8(reinterpret_cast<const char*>(u8_text.data()), u8_text.size());

        cpr::Multipart multipart_parts{
            {"recipient", recipient},
            {"sessiontoken", session_token},
            {"text", text_utf8},
            {"has_text", message.hasText() ? "true" : "false"},
            {"has_image", message.hasImage() ? "true" : "false"}
        };
    
        if (message.hasImage()) {
            multipart_parts.parts.push_back(cpr::Part{"image_data", cpr::File{message.getImageFilePath()}});
        }
    
        auto response = cpr::Post(
            cpr::Url{server_ip + "/send-message"},
            multipart_parts
        );
            
        if (response.error.code != cpr::ErrorCode::OK) {
            std::cerr << "sending message" << " failed (connection error): " << response.error.message << "\n";
            return false;
        }
        if(responseStatusCodeType(response) != HTTPStatus::SUCCESS) {
            try {
                nlohmann::json response_json = nlohmann::json::parse(response.text);
                if (response_json.contains("error") && response_json["error"] == "invalid_token" && response.status_code == 401) {
                    refreshSession();
                    std::cerr << "session token refreshed, message not sent" << std::endl;
                    return false;
                }
            } catch (const nlohmann::json::parse_error& e) {
                std::cerr << "send-message failed (HTTP " << response.status_code << "): " << response.text << std::endl;
            }
            return false;
        }


        return true;
    }
   
    // TODO: Refactor
    bool checkIncomingMessages() { // TODO: make this return a ChatMessage
        auto response = cpr::Get(
            cpr::Url{server_ip + "/recieve-message"},
            cpr::Parameters{{"name", username}}
        );
        if(response.status_code == 400) return false;
        if(!checkResponse(response, "recieve-message")) return false;
        std::cout << response.text << std::endl;

        return true;

    }
  private:
    // TODO: Refactor
    bool checkResponse(const cpr::Response& response, const std::string& context) {
        if (response.error.code != cpr::ErrorCode::OK) {
            std::cerr << context << " failed (connection error): " << response.error.message << "\n";
            return false;
        }
        if (response.status_code < 200 || response.status_code >= 300) {
            std::cerr << context << " failed (HTTP " << response.status_code << "): " << response.text << "\n";
            return false;
        }
        return true;
    }

};


int main() {
    std::string server_ip;
    std::ifstream ip_text_stream("../data/server_ip.txt");
    ip_text_stream >> server_ip;

    ChatClient chat_client1(server_ip);
    // register/login -> sendMessage -> checkIncomingMessages -> refreshToken

    std::string input;
    std::string message;
    std::string recipient;
    while(true){
        std::cout << "type 'q' to quit, 'r' to recieve polled messages and 's' to send a message" << std::endl;
        std::getline(std::cin, input);
        switch (input[0]) {
            case 'q':
                std::cout << "quitting" << std::endl;
                return 0;

            case 'r':
                chat_client1.checkIncomingMessages();
                break;

            case 's':
                std::cout << "enter the recipient: " << std::flush;
                std::getline(std::cin, recipient);
                std::cout << "enter the message: " << std::endl;
                std::getline(std::cin, message);
                chat_client1.sendTextMessage(recipient, message);
                break;
        }
    }
    return 0;
}
