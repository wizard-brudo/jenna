package main

import (
	"fmt"
	"reflect"
)

type parameter struct {
	Data         string
	Type         string
	IsMaybeEmpty bool
}

type Order struct {
	Product   parameter
	Quantity  parameter
	Telephone parameter
	Adress    parameter
	Time      string
	PostCode  parameter
	User      userInfo
}

// Функция проходиться по полям структуры и отдаёт эти поля в виде списка параметров
func (o *Order) getFields() []parameter {
	// Получаем структуру
	val := reflect.ValueOf(*o)
	var arr []parameter
	// Проходимся по полям
	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).Type().String() == "main.parameter" && !val.Field(i).Interface().(parameter).IsMaybeEmpty {
			arr = append(arr, val.Field(i).Interface().(parameter))
		}
	}
	return arr
}

// Функция проверяет заказа и добавляет заказ в список заказов
func (o *Order) makeOrder() bool {
	// Проверяем заказ
	if o.checkOrder() {
		// Если он верный то добовляем заказ в список заказов
		UserOrders[fmt.Sprint(o.User.Id)] = append(UserOrders[fmt.Sprint(o.User.Id)], *o)
		return true
	} else {
		return false
	}
}

// Метод проверяет валидность заказа
func (o *Order) checkOrder() bool {
	// Проверяем валидность заказа
	for _, v := range o.getFields() {
		// Смотрим не пустой ли параметер
		if v.Data != "" {
			// Смотрим соответсвует ли параметер заданному типу
			if v.Type == "int" {
				for _, value := range v.Data {
					if !isDigit(value) {
						return false
					}
				}
			} else if v.Type == "string" {
				for _, value := range v.Data {
					if !isLetter(value) {
						return false
					}
				}
			}
		} else {
			return false
		}
	}
	return true
}
