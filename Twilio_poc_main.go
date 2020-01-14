package main

import (
	"fmt"
	"strings"
	"net/http"
	"net/url"
	"encoding/json"
	//"io/ioutil"
	"log"
	"github.com/gorilla/mux"
)

type sendMessageRequest struct {
	Message string `json:"message"`
	To      string `json:"to"`
}

type response struct {
	Status     string `json:"status"`
	StatusText string `json:"statusText"`
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome sms sending service!")
}

func sendSms(w http.ResponseWriter, r *http.Request) {

	//Twilio credentials
	accountSid := "AC0e9c4ba71a13276c2d05018fb452fff6"
	authToken := "c248ed0ab16317542bf12475058fc8e9"
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/"+accountSid+"/Messages.json"
	fromPhoneNumber:= "+17867892028"

	var smsReq sendMessageRequest
	err := json.NewDecoder(r.Body).Decode(&smsReq)
	
	if err != nil {
		fmt.Fprintf(w, "Kindly enter valid data to send sms")
	}
	fmt.Println("Client Request: %+v", smsReq)

	//Fetching sms data from sms send request
	msgData := url.Values{}
 	msgData.Set("To",smsReq.To)
	msgData.Set("From",fromPhoneNumber)
	msgData.Set("Body",smsReq.Message)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var resp1 response

	resp, _ := client.Do(req)
	if (resp.StatusCode >= 200 && resp.StatusCode < 300) {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if (err == nil) {
			fmt.Println(data["sid"])
			resp1 = response{Status: "Success", StatusText: "SMS Sent!"}
			w.WriteHeader(http.StatusOK)
		}
	} else {
		fmt.Println(resp.Status);
		resp1 = response{Status: "Failed", StatusText: resp.Status}	
		w.WriteHeader(resp.StatusCode)	
	}	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp1)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/sendsms", sendSms).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
