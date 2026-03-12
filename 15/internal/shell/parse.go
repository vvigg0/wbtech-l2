package shell

import (
	"errors"
)

type parser struct {
	ts []token
	i  int
}

func parse(ts []token) (expression, error) {
	p := &parser{ts: ts}
	return p.parseExpr()
}

func (p *parser) atEnd() bool {
	return p.i >= len(p.ts)
}

func (p *parser) peek() token {
	return p.ts[p.i]
}

func (p *parser) consume() token {
	t := p.ts[p.i]
	p.i++
	return t
}

func (p *parser) parseExpr() (expression, error) {
	var expr expression
	if p.atEnd() {
		return expr, nil
	}

	job, err := p.parseJob()
	if err != nil {
		return expression{}, err
	}
	expr.jobs = append(expr.jobs, job)

	for !p.atEnd() {
		t := p.peek()
		if t.ttype != AND && t.ttype != OR {
			return expression{}, errors.New("ожидался оператор && или ||")
		}
		p.consume()

		var op op
		if t.ttype == AND {
			op = opAnd
		} else {
			op = opOr
		}
		expr.ops = append(expr.ops, op)
		if p.atEnd() {
			return expression{}, errors.New("оператор в конце строки")
		}
		job, err := p.parseJob()
		if err != nil {
			return expression{}, err
		}
		expr.jobs = append(expr.jobs, job)
	}
	if len(expr.ops) != len(expr.jobs)-1 {
		return expression{}, errors.New("несогласованная цепочка jobs/ops")
	}
	return expr, nil
}

func (p *parser) parseJob() (job, error) {
	pl, err := p.parsePipeline()
	if err != nil {
		return job{}, err
	}
	if len(pl) == 0 {
		return job{}, errors.New("пустой job")
	}
	return job{pipeline: pl}, nil
}

func (p *parser) parsePipeline() ([]command, error) {
	cmd, err := p.parseCommand()
	if err != nil {
		return nil, err
	}
	pipeline := []command{cmd}

	for !p.atEnd() {
		t := p.peek()
		if t.ttype == PIPE {
			p.consume()
			if p.atEnd() {
				return nil, errors.New("pipe в конце строки")
			}
			cmd, err := p.parseCommand()
			if err != nil {
				return nil, err
			}
			pipeline = append(pipeline, cmd)
			continue
		}
		if t.ttype == AND || t.ttype == OR {
			break
		}
		return nil, errors.New("неожиданный токен в pipeline")
	}
	return pipeline, nil
}

func (p *parser) parseCommand() (command, error) {
	var c command

	for !p.atEnd() {
		t := p.peek()
		if t.ttype != WORD {
			break
		}
		p.consume()
		c.argv = append(c.argv, t.value)
	}
	if len(c.argv) == 0 {
		return command{}, errors.New("пустая команда")
	}

	for !p.atEnd() {
		t := p.peek()
		if t.ttype == PIPE || t.ttype == AND || t.ttype == OR {
			break
		}
		if t.ttype != IN && t.ttype != OUT {
			return command{}, errors.New("неожиданный токен в команде")
		}
		p.consume()
		if p.atEnd() {
			return command{}, errors.New("редирект без пути")
		}
		n := p.peek()
		if n.ttype != WORD {
			return command{}, errors.New("после редиректа ожидается путь")
		}
		p.consume()

		if t.ttype == IN {
			if c.in != "" {
				return command{}, errors.New("повторный <")
			}
			c.in = n.value
		}
		if t.ttype == OUT {
			if c.out != "" {
				return command{}, errors.New("повторный >")
			}
			c.out = n.value
		}
	}
	return c, nil
}
