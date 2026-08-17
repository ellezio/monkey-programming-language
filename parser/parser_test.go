package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`

	l := lexer.New(input)
	program := prepareProgram(t, l, 3)

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral() not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%q", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%q", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322;
	`

	l := lexer.New(input)
	program := prepareProgram(t, l, 3)

	for _, stmt := range program.Statements {
		if stmt.TokenLiteral() != "return" {
			t.Fatalf("stmt.TokenLiteral() not 'return'. got=%q", stmt.TokenLiteral())
		}

		_, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt is not *ast.ReturnStatement. got=%T", stmt)
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`
	l := lexer.New(input)
	program := prepareProgram(t, l, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	testLiteralExpression(t, stmt.Expression, "foobar")
}

func TestIntegerLiteral(t *testing.T) {
	input := `5;`
	l := lexer.New(input)
	program := prepareProgram(t, l, 1)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	if !testLiteralExpression(t, stmt.Expression, 5) {
		return
	}
}

func TestPrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		operand  int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		program := prepareProgram(t, l, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements is not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not *ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Oparator is not %q. got=%q", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Operand, tt.operand) {
			return
		}
	}
}

func TestInfixExpression(t *testing.T) {
	tests := []struct {
		input    string
		left     int64
		operator string
		right    int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		program := prepareProgram(t, l, 1)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statement[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.left, tt.operator, tt.right) {
			return
		}
	}
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) bool {
	infix, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("stmt.Expression is not *ast.InfixExpression. got=%T", exp)
		return false
	}

	if !testLiteralExpression(t, infix.Left, left) {
		return false
	}

	if infix.Operator != operator {
		t.Errorf("infix.Operator is not '%s'. got=%T", operator, infix.Operator)
		return false
	}

	if !testLiteralExpression(t, infix.Right, right) {
		return false
	}

	return true
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Fatalf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}

	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, expectedValue int64) bool {
	il, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("expression is not *ast.IntegerLiteral. got=%T", exp)
		return false
	}

	if il.Value != expectedValue {
		t.Errorf("lit.Value is not %d. got=%d", expectedValue, il.Value)
		return false
	}

	if il.TokenLiteral() != fmt.Sprintf("%d", expectedValue) {
		t.Errorf("ident.TokenLiteral() is not %d. got=%s", expectedValue, il.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("expression not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %q. got=%q", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() not %q. got=%q", value, ident.TokenLiteral())
		return false
	}

	return true
}

func prepareProgram(t *testing.T, l *lexer.Lexer, expectedStatementsLen int) *ast.Program {
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != expectedStatementsLen {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", expectedStatementsLen, len(program.Statements))
	}

	return program
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}
