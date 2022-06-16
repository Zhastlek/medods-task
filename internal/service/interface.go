package service

type Creater interface {
}

type Validator interface {
}

type Token struct {
	Creater
	Validator
}

func NewToken(c Creater, v Validator) *Token {
	return &Token{
		Creater:   c,
		Validator: v,
	}
}
