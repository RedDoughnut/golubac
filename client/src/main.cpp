#include <cpr/response.h>
#include <iostream>
#include <fstream>


#include <cpr/cpr.h>

#include <nlohmann/json.hpp>

std::string ip;

/*
 * Loads ip from data/ip.txt   TODO: change the file location
 */
void loadIP() {
    std::ifstream ip_text_stream("../data/ip.txt");
    ip_text_stream >> ip;
}

/*
 *
 */
cpr::Response login(std::string username) {
    return cpr::Get(
        cpr::Url{ip + "/login"},
        cpr::Parameters{{"name", username.c_str()}}
    );
}

/*
 * Sends a text message
 */
cpr::Response sendMessage(std::string recipient, std::string message) { // TODO: add message class that can hold an image or other things
    nlohmann::json j;
    j["recipient"] = recipient;
    j["message"] = message;

    cpr::Body body{j.dump()};

    return cpr::Post(
        cpr::Url{ip + "/send-message"},
        body
    );
}


int main() {
    loadIP();
    login("me");
    sendMessage("him", "hi");
    
    return 0;
}
