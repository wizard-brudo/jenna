package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

var DEBUG = false
var test = false
var UserOrders = map[string][]Order{}
var numberOrder = 2

func isLetter(r rune) bool {
	return (r >= 1024 && r <= 1279 || r >= 256 && r <= 383 || r >= 65 && r <= 90 || r >= 48 && r <= 57 || r >= 97 && r <= 122 || r == 32 || r == 44)
}

func isDigit(r rune) bool {
	return (r >= 48 && r <= 57)
}

func createOrdersFile(name string) bool {
	f := excelize.NewFile()
	// Set value of a cell.
	f.SetColWidth("Sheet1", "A", "H", 25)
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

func isIn(text string, words []string) bool {
	for _, v := range words {
		if v == text {
			return true
		}
	}
	return false
}

// Функция возращает завтрашений день
func getTommorow() string {
	return fmt.Sprint(time.Now().Day()+1) + "/" + fmt.Sprint(int(time.Now().Month()))
}

// Функция отдаёт сообщения с потверждением заказа
func (o *Order) BotOrderMsg() string {
	var ru = `Вы потверждаете заказ?
Если да то напишите да если нет то напишите команду заново
-----------------------------
Товар:` + o.Product.Data + `
Адрес:` + o.Adress.Data + `
Телефон:` + o.Telephone.Data + `
Количество:` + o.Quantity.Data
	var ro = `
-----------------------------
Confirmați comanda?
Dacă da, atunci scrieți da dacă nu, apoi scrieți comanda ordin din nou
-----------------------------
Produs:` + o.Product.Data + `
Adresa:` + o.Adress.Data + `
Telefon:` + o.Telephone.Data + `
Cantitate:` + o.Quantity.Data
	if o.PostCode.Data == "" {
		return ru + ro
	}
	return ru + "\nПочтовый индекс:" + o.PostCode.Data + ro + "\nCod poștal:" + o.PostCode.Data
}

func webhook(w http.ResponseWriter, r *http.Request) {
	// Если это потверждение вебхука то потверждаем
	if r.Body == http.NoBody && r.URL.Query().Get("hub.mode") != "" {
		verifyWebhook(w, r)
	} else {
		// Иначе читаем сообщения
		rq, err_rq := decodeRequest(r)
		if err_rq != nil {
			fmt.Println("\033[1;31m[ОШИБКА] " + err_rq.Error() + "\033[0m")
		}
		if DEBUG {
			fmt.Printf("\033[1;33m[ДЕБАГ] rq=%v \033[0m\n", rq)
		}
		receivedMessage := rq.Entry[0].Messaging[0].Message.Text
		senderId := rq.Entry[0].Messaging[0].Sender.Id
		// Получаем информацию о пользователе
		userInfo, err := mes.getInfo(senderId)
		if senderId != mes.PageId && receivedMessage != "" {
			mes.Dialogues[fmt.Sprint(userInfo.Id)] = append(mes.Dialogues[fmt.Sprint(userInfo.Id)], receivedMessage)
		}
		if err != nil {
			fmt.Println("\033[1;31m[ОШИБКА] " + err.Error() + "\033[0m")
		}
		// Выводим то что мы получили
		if len(rq.Entry[0].Messaging[0].Message.Attachments) == 0 && receivedMessage != "" && senderId != mes.PageId {
			fmt.Printf("[СООБЩЕНИЕ] %s %s написал(а) %s \n", userInfo.FirstName, userInfo.LastName, receivedMessage)
		} else if len(rq.Entry[0].Messaging[0].Message.Attachments) > 0 {
			fmt.Printf("[СООБЩЕНИЕ] %s %s отправил(а) %d медиа файла, ссылки:\n", userInfo.FirstName, userInfo.LastName, len(rq.Entry[0].Messaging[0].Message.Attachments))
			for _, v := range rq.Entry[0].Messaging[0].Message.Attachments {
				fmt.Println(v.Type, v.Payload.Url)
			}
		}
		// Отвечаем что получили вебхук
		w.Write([]byte("200 OK"))
		if strings.Contains(receivedMessage, "/заказ") || strings.Contains(receivedMessage, "/ordin") {
			sliced_message := strings.Split(receivedMessage, ";")
			//Проверяем валидность заказа
			// Проверка на то что написали именно команду /заказ
			if len(strings.Split(sliced_message[0], " ")) > 1 && (len(sliced_message) == 4 || len(sliced_message) == 5) {
				// Готовим заказ
				order := Order{
					Product:   parameter{strings.Replace(sliced_message[0], strings.Split(sliced_message[0], " ")[0], "", 1), "string", false},
					Quantity:  parameter{sliced_message[3], "int", false},
					Telephone: parameter{sliced_message[1], "int", false},
					Adress:    parameter{sliced_message[2], "string", false},
					Time:      time.Now().Format("01-02-2006 15:04:05"),
					PostCode:  parameter{"", "int", true},
					User:      userInfo,
				}
				var isCorrectOrder bool
				switch len(sliced_message) {
				case 4:
					isCorrectOrder = order.makeOrder()
				case 5:
					order.PostCode = parameter{sliced_message[4], "int", false}
					isCorrectOrder = order.makeOrder()
				}
				if isCorrectOrder {
					mes.sendAndSaveMessage(userInfo.Id, order.BotOrderMsg())
				} else {
					mes.sendAndSaveMessage(userInfo.Id, "Вы ввели неверные данные в одном из полей\nAți introdus date incorecte într-unul dintre câmpuri")
				}
			} else {
				mes.sendAndSaveMessage(userInfo.Id, "Данные неправильно разделенны либо пусты\nDatele sunt separate incorect sau goale")
			}

		} else if isIn(receivedMessage, []string{"Да", "да", "дА", "Da", "da", "dA", "DA", "ДА"}) && (len(mes.Dialogues[fmt.Sprint(userInfo.Id)]) > 1) && strings.Contains(mes.Dialogues[fmt.Sprint(userInfo.Id)][len(mes.Dialogues[fmt.Sprint(userInfo.Id)])-2], "Вы потверждаете заказ?") {
			mes.sendAndSaveMessage(userInfo.Id, "Заказ принят\nComanda este acceptata")
			// Записываем данные в таблицу
			f, err := excelize.OpenFile("fileOrders.xlsx")
			if err != nil {
				log.Fatal(err)
			}
			f.SetCellValue("Sheet1", "A"+fmt.Sprint(numberOrder), userInfo.FirstName)
			f.SetCellValue("Sheet1", "B"+fmt.Sprint(numberOrder), userInfo.LastName)
			f.SetCellValue("Sheet1", "C"+fmt.Sprint(numberOrder), UserOrders[fmt.Sprint(userInfo.Id)][len(UserOrders[fmt.Sprint(userInfo.Id)])-1].Product.Data)
			f.SetCellValue("Sheet1", "D"+fmt.Sprint(numberOrder), UserOrders[fmt.Sprint(userInfo.Id)][len(UserOrders[fmt.Sprint(userInfo.Id)])-1].Quantity.Data)
			f.SetCellValue("Sheet1", "E"+fmt.Sprint(numberOrder), UserOrders[fmt.Sprint(userInfo.Id)][len(UserOrders[fmt.Sprint(userInfo.Id)])-1].Telephone.Data)
			f.SetCellValue("Sheet1", "F"+fmt.Sprint(numberOrder), UserOrders[fmt.Sprint(userInfo.Id)][len(UserOrders[fmt.Sprint(userInfo.Id)])-1].Adress.Data)
			f.SetCellValue("Sheet1", "G"+fmt.Sprint(numberOrder), UserOrders[fmt.Sprint(userInfo.Id)][len(UserOrders[fmt.Sprint(userInfo.Id)])-1].Time)
			// Пишим индекс только если он есть
			if UserOrders[fmt.Sprint(userInfo.Id)][len(UserOrders[fmt.Sprint(userInfo.Id)])-1].PostCode.Data != "" {
				f.SetCellValue("Sheet1", "H"+fmt.Sprint(numberOrder), UserOrders[fmt.Sprint(userInfo.Id)][len(UserOrders[fmt.Sprint(userInfo.Id)])-1].PostCode.Data)
			}
			f.Save()
			f.Close()
			numberOrder++
			fmt.Printf("\033[32m[УСПЕХ] заказ принят от %s %s\033[0m \n", userInfo.FirstName, userInfo.LastName)
		} else if strings.Contains(receivedMessage, "Доставк") || strings.Contains(receivedMessage, "доставк") || strings.Contains(receivedMessage, "precur") || strings.Contains(receivedMessage, "Precur") || strings.Contains(receivedMessage, "заказать") || strings.Contains(receivedMessage, "livr") || strings.Contains(receivedMessage, "Livr") || strings.Contains(receivedMessage, "a comanda") {
			mes.sendAndSaveMessage(userInfo.Id, `Осуществляем доставку завтра `+getTommorow()+`
Доставка по Кишиневу стоит 50 леев.
Для заказа напишите, команду /заказ имя товара;телефон;домашний адресс;количество товара
-------------------------------------
Livram maine `+getTommorow()+`
Livrarea în Chișinău costă 50 de lei.
Pentru a comanda, scrieți, comanda /ordin numele articolului; telefon; adresa de domiciliu; cantitatea articolului`)
		} else if receivedMessage == "Расскажите мне о вашей компании." {
			mes.sendAndSaveMessage(userInfo.Id, `Торговая компания, г.Кишинев. Оптово-розничная торговля.
-------------------
Societate comerciala, Chisinau. Comerț cu ridicata și cu amănuntul.`)
		} else if receivedMessage == "Могу ли я совершить покупку?" {
			mes.sendAndSaveMessage(userInfo.Id, `Да канечно можете,для этого напишите /заказ имя товара;телефон;домашний адресс;количество товара
----------------
Da, bineinteles ca poti, pentru asta scrie /ordin numele produsului; telefon; adresa de domiciliu; cantitatea de marfa`)
		} else if receivedMessage == "Мне это интересно. Могли бы вы предоставить больше информации?" {
			mes.sendAndSaveMessage(userInfo.Id, `Чтобы узнать больше информации напишите команду /инфо интересная информация
----------------------------------
Pentru mai multe informații scrieți comanda /info informații interesante`)
		} else if strings.Contains(receivedMessage, "/инфо") || strings.Contains(receivedMessage, "/info") {
			sliced_message := strings.Split(receivedMessage, " ")
			var msg string
			if sliced_message[1] == "syl" {
				msg = `Это современная команда специалистов в сфере продаж и продвижения широкого ассортимента товаров на рынке Республики Молдовы.
Страны наших производителей: CША, Китай, Россия , Польша, Украина.Миссия компании – приводить в восторг наших клиентов качеством услуг и строить работу с каждым клиентом на долгосрочную перспективу сотрудничества.
Накопительный опыт мы вкладываем в развитие и совершенствование нашей работы.
Постоянное развитие для нас – залог успеха.Мы строим бизнес на отношениях.
--------------------------------------------------------------------------
Este vorba de o echipa moderna de specialisti in domeniul vanzarii si promovarii unei game largi de produse pe piata Republicii Moldova.
Țările producătorilor noștri: SUA, China, Rusia, Polonia, Ucraina Misiunea companiei este de a încânta clienții noștri cu calitatea serviciilor și de a construi lucrul cu fiecare client pe o perspectivă de cooperare pe termen lung.
Investim experiența acumulată în dezvoltarea și îmbunătățirea muncii noastre.
Dezvoltarea continuă pentru noi este cheia succesului, construim afaceri pe baza relațiilor.`
			}
			mes.sendAndSaveMessage(userInfo.Id, msg)
		}
	}
}

func verifyWebhook(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("hub.mode") == "subscribe" && r.URL.Query().Get("hub.verify_token") == mes.VerifyToken {
		w.Write([]byte(r.URL.Query().Get("hub.challenge")))
		fmt.Println("\033[32m[УСПЕХ] веб хук потверждён\033[0m")
		return
	}
}
