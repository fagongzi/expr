package expr

import (
	"fmt"
)

const (
	tokenUnknown    = 10000
	tokenLeftParen  = 1  // (
	tokenRightParen = 2  // )
	tokenLogicAnd   = 3  // &&
	tokenLogicOr    = 4  // ||
	tokenEqual      = 5  // ==
	tokenNotEqual   = 6  // !=
	tokenGT         = 7  // >
	tokenGE         = 8  // >=
	tokenLT         = 9  // <
	tokenLE         = 10 // <=
	tokenIn         = 11 // in
	tokenNotIn      = 12 // !in
	tokenMatch      = 13 // ~
	tokenNotMatch   = 14 // !~
	tokenVarStart   = 15 // {
	tokenVarEnd     = 16 // }
	tokenNumberVal  = 17 // num:
	tokenStringVal  = 18 // str:
)

type parser struct {
	expr      *node
	stack     stack
	prevToken int
	lexer     Lexer
	factory   VarExprFactory
}

func newParser(input []byte, factory VarExprFactory) *parser {
	p := &parser{
		expr:      &node{},
		prevToken: tokenUnknown,
		lexer:     NewScanner(input),
		factory:   factory,
	}

	p.lexer.AddSymbol([]byte("("), tokenLeftParen)
	p.lexer.AddSymbol([]byte(")"), tokenRightParen)
	p.lexer.AddSymbol([]byte("&&"), tokenLogicAnd)
	p.lexer.AddSymbol([]byte("||"), tokenLogicOr)
	p.lexer.AddSymbol([]byte("=="), tokenEqual)
	p.lexer.AddSymbol([]byte("!="), tokenNotEqual)
	p.lexer.AddSymbol([]byte(">"), tokenGT)
	p.lexer.AddSymbol([]byte(">="), tokenGE)
	p.lexer.AddSymbol([]byte("<"), tokenLT)
	p.lexer.AddSymbol([]byte("<"), tokenLE)
	p.lexer.AddSymbol([]byte("in"), tokenIn)
	p.lexer.AddSymbol([]byte("!in"), tokenNotIn)
	p.lexer.AddSymbol([]byte("~"), tokenMatch)
	p.lexer.AddSymbol([]byte("!~"), tokenNotMatch)
	p.lexer.AddSymbol([]byte("{"), tokenVarStart)
	p.lexer.AddSymbol([]byte("}"), tokenVarEnd)
	p.lexer.AddSymbol([]byte("num:"), tokenNumberVal)
	p.lexer.AddSymbol([]byte("str:"), tokenStringVal)

	return p
}

// (
//   ( { origin.a } > 1 || { origin.g } > 0 ) &&
//   { origin.h } > 0
// ) &&
// ( { origin.a } > 1 || { origin.b } > 2 ) &&
// (
//	{ origin.c } < 2 &&
//  ( { num: origin.d } > 2 || { str: origin.e } < 3 )
// )

// Parse return a parsed Expr
func Parse(input []byte, factory VarExprFactory) (Expr, error) {
	p := newParser(input, factory)
	return p.parser()
}

func (p *parser) parser() (Expr, error) {
	p.stack.push(p.expr)
	for {
		p.lexer.NextToken()
		token := p.lexer.Token()
		var err error
		switch token {
		case tokenLeftParen:
			err = p.doLeftParen()
		case tokenRightParen:
			err = p.doRightParen()
		case tokenLogicAnd, tokenLogicOr:
			err = p.doLogic()
		case tokenEqual, tokenNotEqual, tokenGT, tokenGE, tokenLT, tokenLE,
			tokenIn, tokenNotIn, tokenMatch, tokenNotMatch:
			err = p.doCMP()
		case tokenVarStart:
			err = p.doVarStart()
		case tokenNumberVal, tokenStringVal:
			err = p.doVarType()
		case tokenVarEnd:
			err = p.doVarEnd()
		case TokenEOI:
			err = p.doEOI()
			if err != nil {
				return nil, err
			}

			return p.stack.pop(), nil
		}

		if err != nil {
			return nil, err
		}

		p.prevToken = token
	}
}

func (p *parser) doLeftParen() error {
	switch p.prevToken {
	case tokenUnknown:
		p.stack.push(&node{})
	case tokenLeftParen:
		p.stack.append(&node{})
	case tokenLogicAnd:
		err := p.stack.appendWithLogic(&node{}, and)
		if err != nil {
			return err
		}
	case tokenLogicOr:
		err := p.stack.appendWithLogic(&node{}, or)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("expect <(,and,or> before %d",
			p.lexer.TokenIndex())
	}

	p.lexer.SkipString()
	return nil
}

func (p *parser) doRightParen() error {
	var err error
	switch p.prevToken {
	case tokenRightParen: // (a > 1 && (b > 2 && (c > 1)))
		p.stack.pop()
		p.lexer.SkipString()
	case tokenEqual: // (a > 1 && (b == 1))
		err = p.addCmp(equal)
		p.stack.pop()
	case tokenNotEqual:
		err = p.addCmp(notEqual)
		p.stack.pop()
	case tokenGT:
		err = p.addCmp(gt)
		p.stack.pop()
	case tokenGE:
		err = p.addCmp(ge)
		p.stack.pop()
	case tokenLT:
		err = p.addCmp(lt)
		p.stack.pop()
	case tokenLE:
		err = p.addCmp(le)
		p.stack.pop()
	case tokenIn:
		err = p.addCmp(in)
		p.stack.pop()
	case tokenNotIn:
		err = p.addCmp(notIn)
		p.stack.pop()
	case tokenMatch:
		err = p.addCmp(match)
		p.stack.pop()
	case tokenNotMatch:
		err = p.addCmp(notMatch)
		p.stack.pop()
	default:
		return fmt.Errorf("expect <),cmp operator value> before %d",
			p.lexer.TokenIndex())
	}

	return err
}

func (p *parser) doLogic() error {
	var err error
	switch p.prevToken {
	case tokenRightParen: // ((a > 1) && b > 1)
		p.lexer.SkipString()
	case tokenEqual: // (a > 1 && (b == 1))
		err = p.addCmp(equal)
	case tokenNotEqual:
		err = p.addCmp(notEqual)
	case tokenGT:
		err = p.addCmp(gt)
	case tokenGE:
		err = p.addCmp(ge)
	case tokenLT:
		err = p.addCmp(lt)
	case tokenLE:
		err = p.addCmp(le)
	case tokenIn:
		err = p.addCmp(in)
	case tokenNotIn:
		err = p.addCmp(notIn)
	case tokenMatch:
		err = p.addCmp(match)
	case tokenNotMatch:
		err = p.addCmp(notMatch)
	default:
		return fmt.Errorf("expect <),cmp operator value> before %d",
			p.lexer.TokenIndex())
	}

	return err
}

func (p *parser) doCMP() error {
	switch p.prevToken {
	case tokenVarEnd: // { a } > 0
	default:
		return fmt.Errorf("expect <}> before %d",
			p.lexer.TokenIndex())
	}

	p.lexer.SkipString()
	return nil
}

func (p *parser) doVarType() error {
	switch p.prevToken {
	case tokenVarStart:
	default:
		return fmt.Errorf("expect <{> before %d",
			p.lexer.TokenIndex())
	}

	p.lexer.SkipString()
	return nil
}

func (p *parser) doVarStart() error {
	switch p.prevToken {
	case tokenUnknown: // { a } > 0
	case tokenLeftParen: // ( { a } > 1 )
	case tokenLogicAnd, tokenLogicOr: // { a } > 1 && { a } > 1

	default:
		return fmt.Errorf("expect <}> before %d",
			p.lexer.TokenIndex())
	}

	p.lexer.SkipString()
	return nil
}

func (p *parser) doVarEnd() error {
	switch p.prevToken {
	case tokenVarStart, tokenStringVal:
		value, err := p.factory(p.lexer.ScanString())
		if err != nil {
			return err
		}

		p.stack.current().(*node).add(&stringNode{
			varExpr: value,
		})
	case tokenNumberVal:
		value, err := p.factory(p.lexer.ScanString())
		if err != nil {
			return err
		}

		p.stack.current().(*node).add(&numberNode{
			varExpr: value,
		})
	default:
		return fmt.Errorf("expect <{> before %d",
			p.lexer.TokenIndex())
	}

	return nil
}

func (p *parser) doEOI() error {
	var err error
	switch p.prevToken {
	case tokenRightParen:
		p.lexer.SkipString()
	case tokenEqual: // b == 1
		err = p.addCmp(equal)
	case tokenNotEqual:
		err = p.addCmp(notEqual)
	case tokenGT:
		err = p.addCmp(gt)
	case tokenGE:
		err = p.addCmp(ge)
	case tokenLT:
		err = p.addCmp(lt)
	case tokenLE:
		err = p.addCmp(le)
	case tokenIn:
		err = p.addCmp(in)
	case tokenNotIn:
		err = p.addCmp(notIn)
	case tokenMatch:
		err = p.addCmp(match)
	case tokenNotMatch:
		err = p.addCmp(notMatch)
	default:
		return fmt.Errorf("expect <), cmp operator> before %d",
			p.lexer.TokenIndex())
	}

	return err
}

func (p *parser) addCmp(cmp cmp) error {
	value := p.lexer.ScanString()
	if len(value) == 0 {
		return fmt.Errorf("missing cmp value before %d",
			p.lexer.TokenIndex())
	}

	return p.stack.addCmp(cmp, value)
}

type stack struct {
	nodes []Expr
}

func (s *stack) push(v *node) {
	s.nodes = append(s.nodes, v)
}

func (s *stack) append(v *node) {
	s.current().(*node).add(v)
	s.push(v)
}

func (s *stack) addCmp(cmp cmp, value []byte) error {
	return s.current().(*node).lastExpr().(calcNode).AddCmp(cmp, value)
}

func (s *stack) appendWithLogic(v *node, logic logic) error {
	err := s.current().(*node).append(logic, v)
	if err != nil {
		return nil
	}
	s.push(v)
	return nil
}

func (s *stack) current() Expr {
	return s.nodes[len(s.nodes)-1]
}

func (s *stack) pop() Expr {
	n := len(s.nodes) - 1
	v := s.nodes[n]
	s.nodes[n] = nil
	s.nodes = s.nodes[:n]
	return v
}
