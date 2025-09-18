package linters

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("errformat", New)
}

type Settings struct {
	// Пока без настроек, можно добавить позже
}

type ErrFormatLinter struct {
	settings Settings
}

// ErrorInfo содержит информацию об ошибке для анализа
type ErrorInfo struct {
	Var           *ast.Ident // Переменная ошибки
	IsExported    bool       // Экспортируемая ли ошибка
	IsPackageLevel bool      // Объявлена на уровне пакета
	Position      token.Pos  // Позиция в коде
}

// New создает новый экземпляр линтера
func New(settings any) (register.LinterPlugin, error) {
	// Настройки пока не используются
	return &ErrFormatLinter{
		settings: Settings{},
	}, nil
}

// BuildAnalyzers создает анализаторы для линтера
func (l *ErrFormatLinter) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		{
			Name: "errformat",
			Doc:  "checks proper error formatting in fmt.Errorf calls",
			Run:  l.run,
		},
	}, nil
}

// GetLoadMode возвращает режим загрузки анализатора
func (l *ErrFormatLinter) GetLoadMode() string {
	return register.LoadModeTypesInfo
}

// run содержит основную логику анализатора
func (l *ErrFormatLinter) run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			return l.inspectNode(n, pass)
		})
	}
	return nil, nil
}

// inspectNode анализирует узлы AST
func (l *ErrFormatLinter) inspectNode(n ast.Node, pass *analysis.Pass) bool {
	call, ok := n.(*ast.CallExpr)
	if !ok {
		return true
	}
	
	// Проверить, что это вызов fmt.Errorf
	if !l.isFmtErrorf(call) {
		return true
	}
	
	// Анализировать аргументы
	l.analyzeErrorfCall(call, pass)
	return true
}

// isFmtErrorf проверяет, является ли вызов fmt.Errorf
func (l *ErrFormatLinter) isFmtErrorf(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	
	return ident.Name == "fmt" && sel.Sel.Name == "Errorf"
}

// analyzeErrorfCall анализирует аргументы вызова fmt.Errorf
func (l *ErrFormatLinter) analyzeErrorfCall(call *ast.CallExpr, pass *analysis.Pass) {
	if len(call.Args) < 2 {
		return // Недостаточно аргументов
	}

	// Получить format string
	formatString := l.extractFormatString(call.Args[0])
	if formatString == "" {
		return
	}

	// Найти переменные ошибок в аргументах
	errorVars := l.findErrorVariables(call.Args[1:], pass)

	// Проверить соответствие форматов
	l.checkFormatCompliance(formatString, errorVars, call, pass)
}

// extractFormatString извлекает строку формата из аргумента
func (l *ErrFormatLinter) extractFormatString(arg ast.Expr) string {
	lit, ok := arg.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return ""
	}

	// Убрать кавычки
	value := lit.Value
	if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
		return value[1 : len(value)-1]
	}
	return ""
}

// findErrorVariables находит переменные ошибок в аргументах
func (l *ErrFormatLinter) findErrorVariables(args []ast.Expr, pass *analysis.Pass) []ErrorInfo {
	var errorVars []ErrorInfo

	for _, arg := range args {
		if ident, ok := arg.(*ast.Ident); ok {
			// Проверить, что это переменная типа error
			if l.isErrorType(ident, pass) {
				errorInfo := l.analyzeErrorVariable(ident, pass)
				errorVars = append(errorVars, errorInfo)
			}
		}
	}

	return errorVars
}

// analyzeErrorVariable анализирует переменную ошибки
func (l *ErrFormatLinter) analyzeErrorVariable(ident *ast.Ident, pass *analysis.Pass) ErrorInfo {
	info := ErrorInfo{
		Var:      ident,
		Position: ident.Pos(),
	}

	// Найти объявление переменной
	obj := ident.Obj
	if obj == nil {
		return info
	}

	// Проверить уровень объявления (пакетный vs локальный)
	info.IsPackageLevel = l.isPackageLevel(obj)

	// Проверить экспортируемость
	info.IsExported = l.isExported(ident.Name)

	return info
}

// isErrorType проверяет, является ли переменная типом error
func (l *ErrFormatLinter) isErrorType(ident *ast.Ident, pass *analysis.Pass) bool {
	// Получить информацию о типе через TypesInfo
	if pass.TypesInfo == nil || pass.TypesInfo.Types == nil {
		return false
	}
	
	// Получить тип выражения
	typeInfo, ok := pass.TypesInfo.Types[ident]
	if !ok {
		return false
	}
	
	// Проверить, что это интерфейс error
	// Тип error - это предопределенный интерфейс
	if typeInfo.Type == nil {
		return false
	}
	
	// Проверить, что тип реализует интерфейс error
	// Для этого сравним строковое представление типа
	typeStr := typeInfo.Type.String()
	
	// Тип error может быть представлен как "error" или более сложными типами, которые его реализуют
	// Но для простоты проверим основные случаи
	return typeStr == "error" ||
		   typeStr == "interface{}" && typeInfo.Type.Underlying().String() == "interface{error() string}"
}

// isPackageLevel проверяет, объявлена ли переменная на уровне пакета
func (l *ErrFormatLinter) isPackageLevel(obj *ast.Object) bool {
	if obj == nil || obj.Kind != ast.Var {
		return false
	}
	
	// Проверить тип объявления
	switch decl := obj.Decl.(type) {
	case *ast.GenDecl:
		// GenDecl на уровне пакета (var declarations)
		return decl.Tok == token.VAR
	case *ast.ValueSpec:
		// ValueSpec может быть частью GenDecl
		return true
	case *ast.AssignStmt:
		// AssignStmt - это локальные присваивания внутри функций
		return false
	default:
		// Неизвестный тип - считаем локальным
		return false
	}
}

// isExported проверяет, является ли имя экспортируемым
func (l *ErrFormatLinter) isExported(name string) bool {
	return len(name) > 0 && name[0] >= 'A' && name[0] <= 'Z'
}

// checkFormatCompliance проверяет соответствие форматов
func (l *ErrFormatLinter) checkFormatCompliance(formatString string, errorVars []ErrorInfo, call *ast.CallExpr, pass *analysis.Pass) {
	// Найти все %w и %v в format string
	formats := l.parseFormatSpecifiers(formatString)
	
	// Сопоставить с переменными ошибок
	if len(formats) != len(errorVars) {
		return // Количество не совпадает, возможно есть другие аргументы
	}
	
	for i, errorVar := range errorVars {
		if i >= len(formats) {
			break
		}
		
		expectedFormat := l.getExpectedFormat(errorVar)
		actualFormat := formats[i]
		
		if expectedFormat != actualFormat {
			l.reportFormatError(errorVar, expectedFormat, actualFormat, call, pass)
		}
	}
}

// parseFormatSpecifiers парсит спецификаторы формата в строке
func (l *ErrFormatLinter) parseFormatSpecifiers(formatString string) []string {
	var formats []string
	
	for i := 0; i < len(formatString); i++ {
		if formatString[i] == '%' && i+1 < len(formatString) {
			next := formatString[i+1]
			if next == 'w' || next == 'v' {
				formats = append(formats, "%"+string(next))
				i++ // Пропустить следующий символ
			}
		}
	}
	
	return formats
}

// getExpectedFormat определяет ожидаемый формат для переменной ошибки
func (l *ErrFormatLinter) getExpectedFormat(errorVar ErrorInfo) string {
	// Для экспортных ошибок уровня пакета используем %w
	if errorVar.IsExported && errorVar.IsPackageLevel {
		return "%w"
	}
	
	// Для всех остальных используем %v
	return "%v"
}

// reportFormatError создает отчет об ошибке форматирования
func (l *ErrFormatLinter) reportFormatError(errorVar ErrorInfo, expected, actual string, call *ast.CallExpr, pass *analysis.Pass) {
	var message string
	
	if errorVar.IsExported && errorVar.IsPackageLevel {
		message = fmt.Sprintf("exported package-level error '%s' should use %%w instead of %s for error wrapping",
			errorVar.Var.Name, actual)
	} else {
		message = fmt.Sprintf("non-exported error '%s' should use %%v instead of %s for error formatting",
			errorVar.Var.Name, actual)
	}
	
	pass.Report(analysis.Diagnostic{
		Pos:      call.Pos(),
		End:      call.End(),
		Category: "errformat",
		Message:  message,
		SuggestedFixes: l.createSuggestedFix(call, actual, expected),
	})
}

// createSuggestedFix создает предложение исправления
func (l *ErrFormatLinter) createSuggestedFix(call *ast.CallExpr, wrong, correct string) []analysis.SuggestedFix {
	if len(call.Args) == 0 {
		return nil
	}
	
	// Получить позицию format string
	formatArg := call.Args[0]
	
	return []analysis.SuggestedFix{
		{
			Message: fmt.Sprintf("Replace %s with %s", wrong, correct),
			TextEdits: []analysis.TextEdit{
				{
					Pos:     formatArg.Pos(),
					End:     formatArg.End(),
					NewText: []byte(l.replaceFormatSpecifier(formatArg, wrong, correct)),
				},
			},
		},
	}
}

// replaceFormatSpecifier заменяет format specifier в строке
func (l *ErrFormatLinter) replaceFormatSpecifier(formatArg ast.Expr, wrong, correct string) string {
	// Извлечь оригинальную строку
	if lit, ok := formatArg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		originalString := lit.Value
		// Заменить первое вхождение wrong на correct
		newString := strings.Replace(originalString, wrong, correct, 1)
		return newString
	}
	return ""
}