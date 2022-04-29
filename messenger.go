package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Messenger struct {
	Token       string
	PageId      int
	VerifyToken string
	Dialogues   map[string][]string
}

var mes = Messenger{
	"",
	0,
	"",
	map[string][]string{},
}

// Функция отсылает и сохраняет сообщения в json
func (m *Messenger) sendAndSaveMessage(userId int, text string) bool {
	msg := m.sendMessage(userId, text)
	m.Dialogues[fmt.Sprint(userId)] = append(m.Dialogues[fmt.Sprint(userId)], text)
	if DEBUG {
		fmt.Printf("\033[1;33m[ДЕБАГ] msg = %t \033[0m \n", msg)
	}
	if userId != m.PageId {
		if msg {
			if !test {
				fmt.Printf("[ОТВЕТ] Бот ответил на сообщение от %d \n", userId)
			}
			return true
		}
	}
	return false
}

// Функция получает информацию о пользователе
func (m *Messenger) getInfo(id int) (userInfo, error) {
	var info userInfo
	r, err := http.Get("https://graph.facebook.com/" + fmt.Sprint(id) + "?access_token=" + m.Token)
	if err != nil {
		return userInfo{}, err
	}
	err_json := json.NewDecoder(r.Body).Decode(&info)
	if err_json != nil {
		return userInfo{}, err
	}
	return info, nil
}

// Функция отсылает сообщение
func (m *Messenger) sendMessage(sender_id int, text string) bool {
	message := ResponseMessage{MessagingType: "RESPONSE", Recipient: struct {
		Id int "json:\"id,string\""
	}{sender_id}, Message: struct {
		Text string "json:\"text\""
	}{text}}
	// Серилизуем структуру
	body, _ := json.Marshal(message)
	// Отправляем json
	resp, _ := http.Post("https://graph.facebook.com/v13.0/me/messages?access_token="+m.Token, "application/json", bytes.NewBuffer(body))
	// Декодируем ответ
	facebookResp := decodeMessageResponse(resp)
	if DEBUG {
		fmt.Printf("\033[1;33m[ДЕБАГ] facebookResp=%v; body = %s;  \033[0m \n", facebookResp, string(body))
	}
	// И если ответ не пустой то отдаём ответ что сообщение было отправленно
	if facebookResp.RecipientId != 0 && facebookResp.MessageId != "" {
		return true
	}
	return false
}

// Функция декодирует запросы facebook
func decodeRequest(r *http.Request) (FacebookRequest, error) {
	var frequest FacebookRequest
	err := json.NewDecoder(r.Body).Decode(&frequest)
	return frequest, err
}

// Функция декодирует ответы facebook
func decodeMessageResponse(r *http.Response) FacebookResponse {
	var fresponse FacebookResponse
	json.NewDecoder(r.Body).Decode(&fresponse)
	return fresponse
}
