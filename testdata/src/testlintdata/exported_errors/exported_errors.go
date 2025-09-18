package exported_errors

import (
	"errors"
	"fmt"
)

// Экспортные ошибки уровня пакета
var ErrNotFound = errors.New("not found")
var ErrInvalidInput = errors.New("invalid input")

func ExportedErrorsCorrect() error {
	// Правильное использование %w для экспортных ошибок
	return fmt.Errorf("operation failed: %w", ErrNotFound)
}

func ExportedErrorsWrong() error {
	// Неправильное использование %v для экспортных ошибок
	return fmt.Errorf("operation failed: %v", ErrNotFound) // want "exported package-level error 'ErrNotFound' should use %w instead of %v for error wrapping"
}

func MultipleExportedErrors() error {
	// Множественные экспортные ошибки
	return fmt.Errorf("errors: %w and %w", ErrNotFound, ErrInvalidInput)
}

func MixedExportedWrong() error {
	// Смешанное неправильное использование
	return fmt.Errorf("errors: %v and %w", ErrNotFound, ErrInvalidInput) // want "exported package-level error 'ErrNotFound' should use %w instead of %v for error wrapping"
}
