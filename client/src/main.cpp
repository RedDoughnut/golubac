#include <cpr/response.h>
#include <iostream>
#include <fstream>


#include <cpr/cpr.h>

#include <nlohmann/json.hpp>
#include <string>



class ChatClient {
  public:
    std::string username;
    bool logged_in{false};

    ChatClient(const std::string& username, const std::string& ip) {
        this->ip = ip;
        this->username = username;
        logged_in = login();


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
        j["recipient"] = recipient;
        j["message"] = message;

        cpr::Body body{j.dump()};

        auto response = cpr::Post(
            cpr::Url{ip + "/send-message"},
            body
        );

        return checkResponse(response, "send-message");
    }

  private:
    std::string ip;

    /*
     * Logs in as a user
     */
    bool login() {
        auto response = cpr::Get(
            cpr::Url{ip + "/login"},
            cpr::Parameters{{"name", username.c_str()}}
        );

        return checkResponse(response, "login");

    }

    bool checkResponse(cpr::Response response, const std::string& context) {
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

    return 0;
}
