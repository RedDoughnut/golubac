#include <chrono>
#include <iostream>
#include <fstream>


#include <cpr/cpr.h>

#include <nlohmann/json.hpp>
#include <string>
#include <thread>



class ChatClient {
  public:
    std::string username;
    std::string token;
    bool logged_in{false};

    ChatClient(const std::string& username, const std::string& ip) {
        this->ip = ip;
        this->username = username;
        login();
    }


    /*
     * Sends a text message
     */
    bool sendMessage(const std::string& recipient, const std::string& message) { // TODO: add message class that can hold an image or other things
        if(!logged_in) {
            std::cout << "not logged in, cant send a message" << std::endl;
            return false;
        }
        nlohmann::json j;
        j["token"] = token;
        j["recipient"] = recipient;
        j["message"] = message;

        cpr::Body body{j.dump()};

        auto response = cpr::Post(
            cpr::Url{ip + "/send-message"},
            cpr::Header{{"Content-Type", "application/json"}},
            body
        );

        return checkResponse(response, "send-message");
    }

    bool checkIncomingMessages() {
        if(!logged_in) {
            std::cout << "cant recieve a message while not logged in" << std::endl;
            return false;
        }
        auto response = cpr::Get(
            cpr::Url{ip + "/recieve-message"},
            cpr::Parameters{{"name", username}}
        );
        if(response.status_code == 400) return false;
        if(!checkResponse(response, "recieve-message")) return false;
        std::cout << response.text << std::endl;
        
        return true;
        
    }
  private:
    std::string ip;

    /*
     * Logs in as a user
     */
    void login() {
        nlohmann::json body_json;
        body_json["username"] = username;
        auto response = cpr::Post(
            cpr::Url{ip + "/login"},
            cpr::Header{{"Content-Type", "application/json"}},
            cpr::Body{body_json.dump()}
        );
        
        if (!checkResponse(response, "login")) {
            logged_in = false;
            return;
        }
        logged_in = true;
        
        nlohmann::json response_body = nlohmann::json::parse(response.text);
        token = response_body["token"];
        
    }

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
    std::string ip;
    std::ifstream ip_text_stream("../data/ip.txt");
    ip_text_stream >> ip;

    ChatClient chat_client1("me", ip);
    chat_client1.sendMessage("nm", "hi");
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
                chat_client1.sendMessage(recipient, message);
                break;
        }
    }
    return 0;
}
