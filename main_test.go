package main

import (
	"os"
	"testing"
	"time"
)

var user_id = 4978729188841810

func TestSendMesseage(t *testing.T) {
	test = true
	t.Log("Тестирование отправки сообщения")
	test := mes.sendMessage(user_id, "test")
	if !test {
		t.Error("Не удалось отправить сообщение test", "Полученно", test)
		return
	}
	t.Log("Тестирование отправки большого сообщения")
	bigMesseage := mes.sendMessage(user_id, `
	Lorem ipsum dolor sit amet, consectetur adipiscing elit. 
	Sed euismod sed tellus id pulvinar. 
	Cras feugiat ante tincidunt, fringilla mauris nec, consequat libero. Duis vitae tortor lobortis, cursus libero nec, aliquet quam. Cras ac tortor sed orci ornare vestibulum. 
	Ut eget augue enim. Donec euismod orci et libero sollicitudin viverra. Nulla ac purus aliquam, rutrum ante et, varius elit. 
	Etiam accumsan volutpat lorem eu euismod. 
	Ut pharetra mi orci, sit amet pharetra nunc mollis at. 
	Pellentesque ut magna ut justo tristique rutrum. 
	Nunc id porttitor enim, vel ultricies leo. 
	Quisque rhoncus hendrerit dui, vel ornare libero commodo et. 
	Vestibulum euismod felis et scelerisque pharetra.
	Praesent vel lacus tincidunt leo fringilla luctus. 
	In in elit sollicitudin, pulvinar dolor et, aliquam leo. 
	Nunc mollis iaculis risus, ac dictum arcu ullamcorper ac.
	Duis in porta odio. Vivamus tincidunt ipsum id erat pulvinar porttitor.
	Donec sed purus a sem porta porttitor.
	Fusce et orci eget elit luctus cursus.
	Etiam vestibulum vel urna vitae vestibulum. 
	`)
	if !bigMesseage {
		t.Error("Не удалось отправить большое сообщение", "Полученно", bigMesseage)
		return
	}
	t.Log("Тестирование отправки сообщения неверному id")
	notCorrectId := mes.sendMessage(-12341314, "test")
	if notCorrectId {
		t.Error("Удалось отправить сообщение некорректному id", "Полученно", notCorrectId)
		return
	}
	t.Log("Тестирование отправки и сохранения сообщения")
	saveMessage := mes.sendAndSaveMessage(user_id, "test")
	if !saveMessage {
		t.Error("Удалось отправить сообщение некорректному id", "Полученно", saveMessage)
		return
	}
}

func TestGetInfo(t *testing.T) {
	info, err := mes.getInfo(user_id)
	t.Log("Тестирование получения информации о пользователе")
	if err != nil && info.FirstName == "" && info.LastName == "" {
		t.Error("Не удалось получить информацию о пользователе", "Полученно", info)
		return
	}
}

func TestCreateOrdersFile(t *testing.T) {
	t.Log("Тестирование создания файла с заказами")
	file := createOrdersFile("test")
	if !file {
		t.Error("Не удалось создать тестовый excel файл", "Полученно", file)
		return
	}
	t.Log("Тестирование удаления файла с заказами")
	remove_err := os.Remove("test.xlsx")
	if remove_err != nil {
		t.Error("Не удалось удалить тестовый excel файл", "Полученно", remove_err)
		return
	}
}

func TestIsIn(t *testing.T) {
	t.Log("Тестирование сравнения строки с массивом")
	b := isIn("da", []string{"Да", "да", "дА", "Da", "da", "dA", "DA", "ДА"})
	if !b {
		t.Error("Неправильный результат для сравнения строки с массивом", "Полученно", b, "Ожидалось false")
		return
	}
}

func TestCheckOrder(t *testing.T) {
	userInfo := userInfo{
		"Valetin",
		"Petrenjo",
		user_id,
		"-",
	}
	order := Order{
		Product:   parameter{"Ночник галактика", "string", false},
		Quantity:  parameter{"1", "int", false},
		Telephone: parameter{"67687399", "int", false},
		Adress:    parameter{"Кишинёв", "string", false},
		Time:      time.Now().Format("01-02-2006 15:04:05"),
		PostCode:  parameter{"", "int", true},
		User:      userInfo,
	}
	t.Log("Тестирование проверки заказа")
	checkOrder := order.checkOrder()
	if !checkOrder {
		t.Error("Неожиданный результат для проверки корректного заказа", "Полученно", checkOrder, "Ожидалось true")
		return
	}
	NotCorrectOrders := []Order{
		{
			Product:   parameter{"Ночник галактика", "string", false},
			Quantity:  parameter{"trei", "int", false},
			Telephone: parameter{"unudoidoi", "int", false},
			Adress:    parameter{"Кишинёв 3a", "string", false},
			Time:      time.Now().Format("01-02-2006 15:04:05"),
			PostCode:  parameter{"dve tisichi shesti", "int", false},
			User:      userInfo,
		},
		{
			Product:   parameter{"asdasddsa", "string", false},
			Quantity:  parameter{"asdasdsad", "int", false},
			Telephone: parameter{"asdsadasd", "int", false},
			Adress:    parameter{"asdasdsadsad", "string", false},
			Time:      time.Now().Format("01-02-2006 15:04:05"),
			PostCode:  parameter{"asdasdasd", "int", false},
			User:      userInfo,
		},
		{
			Product:   parameter{"asdasddsa", "string", false},
			Quantity:  parameter{"asdasdsad", "int", false},
			Telephone: parameter{"asdsadasd", "int", false},
			Adress:    parameter{"asdasdsadsad", "string", false},
			Time:      time.Now().Format("01-02-2006 15:04:05"),
			PostCode:  parameter{"asdasdasd", "int", false},
			User:      userInfo,
		},
		{
			Product:   parameter{"Nocnick galaktika", "string", false},
			Quantity:  parameter{"1sad", "int", false},
			Telephone: parameter{"+37367676399", "int", false},
			Adress:    parameter{"Chisinau 3a", "string", false},
			Time:      time.Now().Format("01-02-2006 15:04:05"),
			PostCode:  parameter{"2006", "int", false},
			User:      userInfo,
		},
	}
	t.Log("Тестирование проверки некорректных заказов")
	for _, order := range NotCorrectOrders {
		checkNotCorrectOrder := order.checkOrder()
		if checkNotCorrectOrder {
			t.Error("Неожиданный результат для проверки некорректного заказа ", "Полученно", checkNotCorrectOrder, "Ожидалось false")
			return
		}
	}
}

func TestOrder(t *testing.T) {
	test = true
	userInfo := userInfo{
		"Valetin",
		"Petrenjo",
		user_id,
		"-",
	}
	order := Order{
		Product:   parameter{"Ночник галактика", "string", false},
		Quantity:  parameter{"1", "int", false},
		Telephone: parameter{"67687399", "int", false},
		Adress:    parameter{"Кишинёв", "string", false},
		Time:      time.Now().Format("01-02-2006 15:04:05"),
		PostCode:  parameter{"", "int", true},
		User:      userInfo,
	}
	NotCorrectOrder := Order{
		Product:   parameter{"Ночник галактика", "string", false},
		Quantity:  parameter{"trei", "int", false},
		Telephone: parameter{"unudoidoi", "int", false},
		Adress:    parameter{"Кишинёв 3a", "string", false},
		Time:      time.Now().Format("01-02-2006 15:04:05"),
		PostCode:  parameter{"dve tisichi shesti", "int", false},
		User:      userInfo,
	}
	t.Log("Тестирование создания корректного заказа")
	checkMakeCorrectOrder := order.makeOrder()
	if !checkMakeCorrectOrder {
		t.Error("Неожиданный результат для создания корректного заказа ", "Полученно", checkMakeCorrectOrder, "Ожидалось true")
		return
	}
	t.Log("Тестирование создания некорректного заказа")
	checkMakeNotCorrectOrder := NotCorrectOrder.makeOrder()
	if checkMakeNotCorrectOrder {
		t.Error("Неожиданный результат для создания некорректного заказа ", "Полученно", checkMakeNotCorrectOrder, "Ожидалось false")
		return
	}

}
