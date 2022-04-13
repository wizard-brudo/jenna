package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type parameter struct {
	Data string
	Type string
}

var DEBUG = false
var dialogues = map[string][]string{}
var numberOrder = 2

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "debug" {
			DEBUG = true
			fmt.Println("\033[1;33m[ИНФОРМАЦИЯ] Дебаг включён \033[0m")
		} else if os.Args[1] == "set" && os.Args[2] == "numberOrders" {
			if os.Args[2] == "numberOrders" {
				numberOrder, _ = strconv.Atoi(os.Args[3])
			}
			fmt.Println("\033[1;33m[ИНФОРМАЦИЯ] Значение изменено \033[0m")
		} else if os.Args[1] == "help" {
			fmt.Printf(`Доступные команды:
debug-Включить дебаг
set numberOrder [количество заказов]-устоновить начальное значение для количества заказов
set verifyToken [токен]-устоновить начальное значение для количества заказов
`)
			os.Exit(0)
		} else if os.Args[1] == "version" {
			fmt.Println(0.1)
			os.Exit(0)
		}
	}
	fmt.Println("\033[32m[ИНФОРМАЦИЯ] Бот запустился \033[0m")
	var fileOrdersExists bool
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.Name() == "fileOrders.xlsx" {
			fileOrdersExists = true
		}
	}
	if !fileOrdersExists {
		createOrdersFile("fileOrders")
		fmt.Println("\033[1;33m[ИНФОРМАЦИЯ] файл с заказами создан \033[0m")
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("templates/index.html")
		tmpl.Execute(w, "")
	})
	http.HandleFunc("/privacy", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles("templates/privacy.html")
		tmpl.Execute(w, "")
	})
	http.HandleFunc("/webhook", webhook)

	log.Fatal(http.ListenAndServe(":8083", nil))
}

func createOrdersFile(name string) bool {
	f := excelize.NewFile()
	// Set value of a cell.
	f.SetColWidth("Sheet1", "A", "H", 30)
	f.SetCellValue("Sheet1", "A1", "ИМЯ")
	f.SetCellValue("Sheet1", "B1", "ФАМИЛИЯ")
	f.SetCellValue("Sheet1", "C1", "ТОВАР")
	f.SetCellValue("Sheet1", "D1", "КОЛИЧЕСТВО")
	f.SetCellValue("Sheet1", "E1", "ТЕЛЕФОН")
	f.SetCellValue("Sheet1", "F1", "АДРЕСС")
	f.SetCellValue("Sheet1", "G1", "ВРЕМЯ")
	f.SetCellValue("Sheet1", "H1", "ПОЧТОВЫЙ ИНДЕКС")
	// Set active sheet of the workbook.
	if err := f.SaveAs(name + ".xlsx"); err != nil {
		return false
	}
	return true
}

func checkOrder(leng int, data ...parameter) bool {
	if len(data) != leng {
		return false
	}
	for _, v := range data {
		if v.Data != "" {
			if v.Type == "int" {
				isValidInt, _ := regexp.Compile("[0-9]")
				if !isValidInt.MatchString(v.Data) {
					return false
				}
			} else if v.Type == "string" {
				isValidInt, _ := regexp.Compile("[A-Za-zА-Яа-я]")
				if !isValidInt.MatchString(v.Data) {
					return false
				}
			}
		} else {
			return false
		}
	}
	return true
}

func trimMany(stringi ...string) []string {
	var str []string
	for _, v := range stringi {
		str = append(str, strings.TrimSpace(v))
	}
	return str
}

func BotOrderMsg(params ...string) string {
	return "Вы потверждаете заказ? \nЕсли да то напишите да если нет то напишите команду заново\n----------------------\nТовар:" + params[0] + "\nАдрес:" + params[1] + "\nТелефон:" + params[2] + "\nКоличество:" + params[3]
}

func webhook(w http.ResponseWriter, r *http.Request) {
	// Если это потверждение вебхука то потверждаем
	if r.Body == http.NoBody && r.URL.Query().Get("hub.mode") != "" {
		verifyWebhook(w, r)
	} else {
		// Иначе читаем сообщения
		rq, err_rq := decodeRequest(r)
		if err_rq != nil {
			fmt.Println(err_rq)
		}
		if DEBUG {
			fmt.Printf("\033[1;33m[ДЕБАГ] rq=%v \033[0m\n", rq)
		}
		receivedMessage := rq.Entry[0].Messaging[0].Message.Text
		senderId := rq.Entry[0].Messaging[0].Sender.Id
		// Получаем информацию о пользователе
		userInfo, err := getInfo(senderId)
		if senderId != PAGE_ID && receivedMessage != "" {
			dialogues[fmt.Sprint(userInfo.Id)] = append(dialogues[fmt.Sprint(userInfo.Id)], receivedMessage)
		}
		if err != nil {
			fmt.Println(err)
		}
		// Выводим то что мы получили
		if len(rq.Entry[0].Messaging[0].Message.Attachments) == 0 && receivedMessage != "" && senderId != PAGE_ID {
			fmt.Printf("[СООБЩЕНИЕ] %s %s написал(а) %s \n", userInfo.FirstName, userInfo.LastName, receivedMessage)
		} else if len(rq.Entry[0].Messaging[0].Message.Attachments) > 0 {
			fmt.Printf("[СООБЩЕНИЕ] %s %s отправил(а) %d медиа файла, ссылки:\n", userInfo.FirstName, userInfo.LastName, len(rq.Entry[0].Messaging[0].Message.Attachments))
			for _, v := range rq.Entry[0].Messaging[0].Message.Attachments {
				fmt.Println(v.Type, v.Payload.Url)
			}
		}
		// Отвечаем что получили вебхук
		w.Write([]byte("200 OK"))
		if strings.Contains(receivedMessage, "/заказ") {
			sliced_message := strings.Split(receivedMessage, ";")
			parameters := []parameter{
				{strings.Split(sliced_message[0], " ")[1], "string"},
				{sliced_message[1], "int"},
				{sliced_message[2], "string"},
				{sliced_message[3], "int"},
			}
			//Проверяем валидность заказа
			// Проверка на то что написали именно команду /заказ
			if len(strings.Split(sliced_message[0], " ")) > 1 {
				if len(sliced_message) == 4 {
					// Проверяем заказ
					if checkOrder(4, parameters...) {
						trimmed := trimMany(strings.Replace(sliced_message[0], "/заказ", "", 1), sliced_message[1], sliced_message[2], sliced_message[3])
						msg := BotOrderMsg(trimmed[0], trimmed[2], trimmed[1], trimmed[3])
						sendAndCheckMessage(userInfo, msg)
						dialogues[fmt.Sprint(userInfo.Id)] = append(dialogues[fmt.Sprint(userInfo.Id)], msg)
					} else {
						sendAndCheckMessage(userInfo, "Вы ввели неверные данные в одном из полей")
						dialogues[fmt.Sprint(userInfo.Id)] = append(dialogues[fmt.Sprint(userInfo.Id)], "Вы ввели неверные данные в одном из полей")
					}
				} else if len(sliced_message) == 5 {
					parameters = append(parameters, parameter{Data: sliced_message[4], Type: "int"})
					if checkOrder(5, parameters...) {
						// Очищаем от пробелов
						trimmed := trimMany(strings.Replace(sliced_message[0], "/заказ", "", 1), sliced_message[1], sliced_message[2], sliced_message[3], sliced_message[4])
						msg := BotOrderMsg(trimmed[0], trimmed[2], trimmed[1], trimmed[3]) + "\nПочтовый код:" + trimmed[4]
						sendAndCheckMessage(userInfo, msg)
						dialogues[fmt.Sprint(userInfo.Id)] = append(dialogues[fmt.Sprint(userInfo.Id)], msg)
					} else {
						sendAndCheckMessage(userInfo, "Вы ввели неверные данные в одном из полей")
						dialogues[fmt.Sprint(userInfo.Id)] = append(dialogues[fmt.Sprint(userInfo.Id)], "Вы ввели неверные данные в одном из полей")
					}
				} else {
					sendAndCheckMessage(userInfo, "Данные неправильно разделенны либо пусты")
					dialogues[fmt.Sprint(userInfo.Id)] = append(dialogues[fmt.Sprint(userInfo.Id)], "Данные неправильно разделенны либо пусты")
				}
			}

		} else if (receivedMessage == "Да" || receivedMessage == "да" || receivedMessage == "дА") && (len(dialogues[fmt.Sprint(userInfo.Id)]) > 1) && strings.Contains(dialogues[fmt.Sprint(userInfo.Id)][len(dialogues[fmt.Sprint(userInfo.Id)])-2], "Вы потверждаете заказ?") {
			sendAndCheckMessage(userInfo, "Заказ принят")
			product := strings.TrimSpace(strings.Split(strings.Split(dialogues[fmt.Sprint(userInfo.Id)][len(dialogues[fmt.Sprint(userInfo.Id)])-2], ":")[1], "\n")[0])
			address := strings.TrimSpace(strings.Split(strings.Split(dialogues[fmt.Sprint(userInfo.Id)][len(dialogues[fmt.Sprint(userInfo.Id)])-2], ":")[2], "\n")[0])
			telephone := strings.TrimSpace(strings.Split(strings.Split(dialogues[fmt.Sprint(userInfo.Id)][len(dialogues[fmt.Sprint(userInfo.Id)])-2], ":")[3], "\n")[0])
			amount := strings.TrimSpace(strings.Split(strings.Split(dialogues[fmt.Sprint(userInfo.Id)][len(dialogues[fmt.Sprint(userInfo.Id)])-2], ":")[4], "\n")[0])
			f, _ := excelize.OpenFile("fileOrders.xlsx")
			f.SetCellValue("Sheet1", "A"+fmt.Sprint(numberOrder), userInfo.FirstName)
			f.SetCellValue("Sheet1", "B"+fmt.Sprint(numberOrder), userInfo.LastName)
			f.SetCellValue("Sheet1", "C"+fmt.Sprint(numberOrder), product)
			f.SetCellValue("Sheet1", "D"+fmt.Sprint(numberOrder), amount)
			f.SetCellValue("Sheet1", "E"+fmt.Sprint(numberOrder), telephone)
			f.SetCellValue("Sheet1", "F"+fmt.Sprint(numberOrder), address)
			f.SetCellValue("Sheet1", "G"+fmt.Sprint(numberOrder), time.Now().Format("01-02-2006 15:04:05"))
			if len(strings.Split(dialogues[fmt.Sprint(userInfo.Id)][len(dialogues[fmt.Sprint(userInfo.Id)])-2], ":")) == 5 {
				f.SetCellValue("Sheet1", "H"+fmt.Sprint(numberOrder), "----")
			} else {
				f.SetCellValue("Sheet1", "H"+fmt.Sprint(numberOrder), strings.TrimSpace(strings.Split(strings.Split(dialogues[fmt.Sprint(userInfo.Id)][len(dialogues[fmt.Sprint(userInfo.Id)])-2], ":")[5], "\n")[0]))
			}
			f.Save()
			f.Close()
			numberOrder++
			fmt.Printf("\033[32m[УСПЕХ] заказ принят от %s %s\033[0m \n", userInfo.FirstName, userInfo.LastName)
			dialogues[fmt.Sprint(userInfo.Id)] = append(dialogues[fmt.Sprint(userInfo.Id)], "Заказ принят")
		}
	}
}

func verifyWebhook(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("hub.mode") == "subscribe" && r.URL.Query().Get("hub.verify_token") == VERIFY_TOKEN {
		w.Write([]byte(r.URL.Query().Get("hub.challenge")))
		fmt.Println("\033[32m[УСПЕХ] веб хук потверждён\033[0m")
		return
	}
}
