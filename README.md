# ErrFormat Linter

Линтер для golangci-lint, который проверяет правильное форматирование ошибок в вызовах `fmt.Errorf()`.

## Правила

### Экспортные ошибки уровня пакета
Должны использовать `%w` для правильного wrapping:
```go
var ErrNotFound = errors.New("not found")

func process() error {
    return fmt.Errorf("processing failed: %w", ErrNotFound) // ✅ Правильно
    return fmt.Errorf("processing failed: %v", ErrNotFound) // ❌ Ошибка
}
```

### Неэкспортные ошибки
Должны использовать `%v` для форматирования:
```go
func process() error {
    err := someFunc()
    if err != nil {
        return fmt.Errorf("processing failed: %v", err) // ✅ Правильно
        return fmt.Errorf("processing failed: %w", err) // ❌ Ошибка
    }
}
```

## Установка

### Как плагин для golangci-lint

1. Склонируйте репозиторий:
```bash
git clone https://github.com/yourusername/errformat
cd errformat
```

2. Соберите плагин:
```bash
go build -buildmode=plugin -o errformat.so .
```

3. Скопируйте пример конфигурации:
```bash
cp .golangci.example.yml .golangci.yml
```

4. Запустите линтер:
```bash
golangci-lint run --disable-all -E errformat
```

## Разработка

### Требования
- Go 1.23+
- golangci-lint

### Тестирование
```bash
go test -v ./...
```

### Сборка
```bash
go build -buildmode=plugin -o errformat.so .
```

## Архитектура

Линтер построен на основе `golang.org/x/tools/go/analysis` и включает:

- **Анализатор AST** - обнаруживает вызовы `fmt.Errorf`
- **Детектор типов ошибок** - различает экспортные и неэкспортные ошибки
- **Валидатор форматов** - проверяет соответствие `%w` и `%v` правилам
- **Система отчетов** - генерирует диагностику с предложениями исправлений

## Тестовые случаи

Проект включает комплексные тесты в директории `testdata/src/testlintdata/`:

- `exported_errors/` - тесты для экспортных ошибок уровня пакета
- `unexported_errors/` - тесты для локальных и неэкспортных ошибок
- `mixed_cases/` - смешанные сценарии использования
- `edge_cases/` - граничные случаи и исключения

## Техническое задание

Подробные технические требования и спецификации находятся в папке `tasks/`:

- `task_01_analysis.md` - анализ требований
- `task_02_structure.md` - архитектура проекта
- `task_03_implementation.md` - детали реализации
- И другие файлы с техническими заданиями