package mixed_cases

import (
	"errors"
	"fmt"
)

// Экспортные ошибки уровня пакета
var ErrTimeout = errors.New("timeout")

func MixedCorrectUsage() error {
	localErr := someFunction()
	
	// Правильное использование: %w для экспортной, %v для неэкспортной
	if localErr != nil {
		return fmt.Errorf("local error: %v, timeout: %w", localErr, ErrTimeout)
	}
	return nil
}

func MixedWrongUsage() error {
	localErr := someFunction()
	
	// Неправильное использование: %w для неэкспортной, %v для экспортной
	if localErr != nil {
		return fmt.Errorf("local error: %w, timeout: %v", localErr, ErrTimeout) // want "non-exported error 'localErr' should use %v instead of %w for error formatting" "exported package-level error 'ErrTimeout' should use %w instead of %v for error wrapping"
	}
	return nil
}

func ComplexMixedCase() error {
	err1 := someFunction()
	err2 := someOtherFunction()
	
	// Множественные неэкспортные ошибки с одной экспортной
	return fmt.Errorf("errors: %w, %w, and %w", err1, err2, ErrTimeout) // want "non-exported error 'err1' should use %v instead of %w for error formatting" "non-exported error 'err2' should use %v instead of %w for error formatting"
}

func someFunction() error {
	return nil
}

func someOtherFunction() error {
	return nil
}