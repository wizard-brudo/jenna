package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "debug" {
			DEBUG = true
			fmt.Println("\033[1;33m[ИНФОРМАЦИЯ] Дебаг включён \033[0m")
		} else if os.Args[1] == "set" {
			if os.Args[2] == "numberOrders" {
				numberOrder, _ = strconv.Atoi(os.Args[3])
			}
			if os.Args[2] == "api" {
				mes.Token = os.Args[3]
				mes.VerifyToken = os.Args[4]
				mes.PageId, _ = strconv.Atoi(os.Args[5])
			}
			fmt.Println("\033[1;33m[ИНФОРМАЦИЯ] Значение изменено \033[0m")
		} else if os.Args[1] == "help" {
			fmt.Printf(`Доступные команды:
debug-Включить дебаг
set numberOrder [количество заказов]-устоновить начальное значение для количества заказов
set api токен проверочный токен-устоновить маркер доступа и проверочный токен для вебхука
version-Версия бота
`)
			os.Exit(0)
		} else if os.Args[1] == "version" {
			fmt.Println(0.2)
			os.Exit(0)
		}
	}
	var fileOrdersExists bool
	var configFileExists bool
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.Name() == "fileOrders.xlsx" {
			fileOrdersExists = true
		} else if f.Name() == "jenna.json" {
			configFileExists = true
		}
	}
	if configFileExists {
		conf := map[string]string{}
		data, _ := os.ReadFile("jenna.json")
		json.Unmarshal(data, &conf)
		mes.Token = conf["token"]
		mes.VerifyToken = conf["verifyToken"]
		mes.PageId, _ = strconv.Atoi(conf["pageId"])
		fmt.Println("\033[1;33m[ИНФОРМАЦИЯ] файл конфигурации прочитан \033[0m")
	}
	if !fileOrdersExists {
		createOrdersFile("fileOrders")
		fmt.Println("\033[1;33m[ИНФОРМАЦИЯ] файл с заказами создан \033[0m")
	}
	if mes.Token == "" && mes.VerifyToken == "" && mes.PageId == 0 {
		fmt.Println("\033[1;31m[ОШИБКА] Нету входных данных \033[0m")
		os.Exit(1)
	}
	fmt.Println("\033[32m[ИНФОРМАЦИЯ] Бот запустился \033[0m")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})
	http.HandleFunc("/privacy", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/privacy.html")
	})
	http.HandleFunc("/webhook", webhook)

	log.Fatal(http.ListenAndServe(":8083", nil))
}
