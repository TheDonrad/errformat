package unexported_errors

import (
	"fmt"
)

func UnexportedErrorsCorrect() error {
	err := someFunction()
	if err != nil {
		// Правильное использование %v для неэкспортных ошибок
		return fmt.Errorf("processing failed: %v", err)
	}
	return nil
}

func UnexportedErrorsWrong() error {
	err := someFunction()
	if err != nil {
		// Неправильное использование %w для неэкспортных ошибок
		return fmt.Errorf("processing failed: %w", err) // want "non-exported error 'err' should use %v instead of %w for error formatting"
	}
	return nil
}

func LocalVariableErrors() error {
	localErr := someFunction()
	// Даже если имя начинается с заглавной буквы, это локальная переменная
	LocalErr := someFunction()
	
	return fmt.Errorf("errors: %w %w", localErr, LocalErr) // want "non-exported error 'localErr' should use %v instead of %w for error formatting" "non-exported error 'LocalErr' should use %v instead of %w for error formatting"
}

func someFunction() error {
	return nil
}