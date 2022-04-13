package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func sendAndCheckMessage(user userInfo, text string) {
	msg := sendMessage(user.Id, text)
	if DEBUG {
		fmt.Printf("\033[1;33m[ДЕБАГ] msg = %t \033[0m \n", msg)
	}
	if msg {
		fmt.Printf("[СООБЩЕНИЕ] Бот ответил на сообщение от %s %s\n", user.FirstName, user.LastName)
	} else {
		fmt.Printf("[СООБЩЕНИЕ] Бот не смог ответить на сообщение от %s %s\n", user.FirstName, user.LastName)
	}
}

func getInfo(id int) (userInfo, error) {
	var info userInfo
	r, err := http.Get("https://graph.facebook.com/" + fmt.Sprint(id) + "?access_token=" + TOKEN)
	if err != nil {
		return userInfo{}, err
	}
	err_json := json.NewDecoder(r.Body).Decode(&info)
	if err_json != nil {
		return userInfo{}, err
	}
	return info, nil
}

func sendMessage(sender_id int, text string) bool {
	message := ResponseMessage{MessagingType: "RESPONSE", Recipient: struct {
		Id int "json:\"id,string\""
	}{sender_id}, Message: struct {
		Text string "json:\"text\""
	}{text}}
	// Серилизуем структуру
	body, _ := json.Marshal(message)
	// Отправляем json
	resp, _ := http.Post("https://graph.facebook.com/v13.0/me/messages?access_token="+TOKEN, "application/json", bytes.NewBuffer(body))
	facebookResp := decodeMessageResponse(resp)
	if DEBUG {
		fmt.Printf("\033[1;33m[ДЕБАГ] facebookResp=%v; body = %s;  \033[0m \n", facebookResp, string(body))
	}
	if facebookResp.RecipientId != 0 && facebookResp.MessageId != "" {
		return true
	}
	return false
}

func decodeRequest(r *http.Request) (FacebookRequest, error) {
	var frequest FacebookRequest
	err := json.NewDecoder(r.Body).Decode(&frequest)
	return frequest, err
}

func decodeMessageResponse(r *http.Response) FacebookResponse {
	var fresponse FacebookResponse
	json.NewDecoder(r.Body).Decode(&fresponse)
	return fresponse
}
