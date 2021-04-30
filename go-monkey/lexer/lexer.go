package lexer

import (
	"github.com/htamakos/go-monkey/token"
)

type Lexer struct {
	input        string // 入力における現在位置
	position     int    // 入力における現在位置
	readPosition int    // これから読み込む位置（現在の文字の次）
	ch           byte   //　現在検査中の文字
}

// Lexer 型の構造体を初期化するファクトリ関数
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

const NULL = 0
const INCREMENTAL_INPUT = 1

// ヘルパーメソッド
// input文字列の1文字を読み、現在位置を更新する
// この字句解析は ASCII のみに対応している
// UNICODE と UTF-8 をサポートするためには、byte から rune に変更し、次の文字を読む処理を変更する必要がある
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = NULL
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += INCREMENTAL_INPUT
}

// 次のトークンを返す
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case NULL:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// 識別子を読み取る
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

// 数値を読み取る
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

// 空白文字をスキップする
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		l.readChar()
	}
}

// 次の文字を覗き見する
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return NULL
	} else {
		return l.input[l.readPosition]
	}
}

// トークンを構造体を返す
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// アルファベット文字かどうかを判定する
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// 数値文字かどうかを判定する
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
