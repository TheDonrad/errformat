# Примеры использования ErrFormat Linter

## Экспортные ошибки уровня пакета

### ✅ Правильно

```go
package mypackage

import (
    "errors"
    "fmt"
)

// Экспортные ошибки уровня пакета
var (
    ErrNotFound     = errors.New("not found")
    ErrInvalidInput = errors.New("invalid input")
    ErrTimeout      = errors.New("timeout")
)

func FindUser(id string) (*User, error) {
    // ✅ Экспортная ошибка уровня пакета должна использовать %w
    return nil, fmt.Errorf("user lookup failed: %w", ErrNotFound)
}

func ProcessData(data []byte) error {
    if len(data) == 0 {
        // ✅ Экспортная ошибка уровня пакета должна использовать %w
        return fmt.Errorf("data validation failed: %w", ErrInvalidInput)
    }
    return nil
}
```

### ❌ Неправильно

```go
package mypackage

import (
    "errors"
    "fmt"
)

var ErrNotFound = errors.New("not found")

func FindUser(id string) (*User, error) {
    // ❌ Экспортная ошибка уровня пакета должна использовать %w, а не %v
    return nil, fmt.Errorf("user lookup failed: %v", ErrNotFound)
}
```

## Неэкспортные/локальные ошибки

### ✅ Правильно

```go
package mypackage

import (
    "fmt"
    "os"
)

func ReadConfig(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        // ✅ Локальная переменная ошибки должна использовать %v
        return fmt.Errorf("failed to open config: %v", err)
    }
    defer file.Close()
    
    return nil
}

func ProcessFile(path string) error {
    localErr := someOperation(path)
    if localErr != nil {
        // ✅ Локальная переменная должна использовать %v
        return fmt.Errorf("processing failed: %v", localErr)
    }
    
    return nil
}
```

### ❌ Неправильно

```go
package mypackage

import (
    "fmt"
    "os"
)

func ReadConfig(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        // ❌ Локальная переменная должна использовать %v, а не %w
        return fmt.Errorf("failed to open config: %w", err)
    }
    defer file.Close()
    
    return nil
}
```

## Смешанные случаи

### ✅ Правильно

```go
package mypackage

import (
    "errors"
    "fmt"
    "net/http"
)

var ErrTimeout = errors.New("operation timeout")

func ComplexOperation() error {
    // Локальная ошибка
    resp, err := http.Get("https://example.com")
    if err != nil {
        // ✅ Локальная ошибка - используем %v
        return fmt.Errorf("HTTP request failed: %v", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode == 408 {
        // ✅ Экспортная ошибка уровня пакета - используем %w
        return fmt.Errorf("request timed out: %w", ErrTimeout)
    }
    
    return nil
}
```

## Граничные случаи

### Случаи, которые НЕ вызывают ошибок линтера

```go
package mypackage

import "fmt"

func EdgeCases() {
    // ✅ Не fmt.Errorf - игнорируется
    fmt.Sprintf("not an error: %v", "something")
    
    // ✅ Без переменных ошибок - игнорируется
    fmt.Errorf("simple error message")
    
    // ✅ С другими типами аргументов - игнорируется
    fmt.Errorf("error with number: %d", 42)
    
    // ✅ Недостаточно аргументов - игнорируется
    fmt.Errorf("no arguments")
}
```

## Результат работы линтера

При запуске на неправильном коде линтер выдаст сообщения:

```
exported_errors.go:19:9: exported package-level error 'ErrNotFound' should use %w instead of %v for error wrapping
unexported_errors.go:20:9: non-exported error 'err' should use %v instead of %w for error formatting
```

## Автоматическое исправление

Линтер предоставляет предложения по исправлению, которые могут быть применены автоматически в поддерживающих редакторах (VS Code, GoLand и др.).