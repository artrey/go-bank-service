package mcc_test

import (
	"fmt"
	"github.com/artrey/go-bank-service/pkg/mcc"
)

func ExampleMCC_ToCategory() {
	fmt.Println(mcc.MCC("5411").ToCategory())
	fmt.Println(mcc.MCC("5533").ToCategory())
	fmt.Println(mcc.MCC("5812").ToCategory())
	fmt.Println(mcc.MCC("5912").ToCategory())
	fmt.Println(mcc.MCC("7788").ToCategory())
	// Output:
	// Супермаркеты
	// Автоуслуги
	// Рестораны
	// Аптеки
	// Категория не указана
}
