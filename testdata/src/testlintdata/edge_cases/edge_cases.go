package edge_cases

import (
	"fmt"
)

// Случаи, которые не должны вызывать ошибок
func ValidCases() {
	// Не fmt.Errorf
	fmt.Sprintf("not an error: %v", "something")

	// Без переменных ошибок
	fmt.Errorf("simple error message")

	// С другими типами аргументов
	fmt.Errorf("error with number: %d", 42)
}

func NoArguments() error {
	// Недостаточно аргументов
	return fmt.Errorf("no arguments")
}

func NonStringFormat() {
	// Не строковый первый аргумент
	var format interface{} = "dynamic format: %v"
	err := someFunction()
	fmt.Errorf(format.(string), err) // Не должно анализироваться
}

func someFunction() error {
	return nil
}