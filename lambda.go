package λeval

import "fmt"

type Expression interface{
	Evaluate() Expression
	AlphaConvert(Variable, Variable) Expression
	Substitute(Variable, Expression) Expression
	fmt.Stringer
}

type Variable string

func (v Variable) Evaluate() Expression {
	return v
}

func (v Variable) AlphaConvert(x, y Variable) Expression {
	if string(v) == string(x) { return y }
	return v
}

func (t Variable) Substitute(v Variable, e Expression) Expression {
	if string(t) == string(v) { return e }
	return t
}

func (v Variable) String() string {
	return string(v)
}

type Abstraction struct{
	Argument Variable
	Body Expression
}

func (a Abstraction) Evaluate() Expression {
	return a
}

func (a Abstraction) AlphaConvert(x, y Variable) Expression {
	return Abstraction{ a.Argument.AlphaConvert(x, y), a.Body.AlphaConvert(x, y) } //TODO: Make non-capturing
}

func (a Abstraction) Substitute(v Variable, e Expression) Expression {
	return Abstraction{ a.Argument, a.Body.Substitute(v, e) } //TODO: Make non-capturing
}

func (a Abstraction) String() string {
	return fmt.Sprintf("(λ%v.%v)", a.Argument, a.Body)
}

type Application struct{
	Function Expression
	Argument Expression
}

func (a Application) Evaluate() Expression {
	return a //TODO: stub
}

func (a Application) AlphaConvert(x, y Variable) Expression {
	return Application{ a.Function.AlphaConvert(x, y), a.Argument.AlphaConvert(x, y) }
}

func (a Application) Substitute(v Variable, e Expression) Expression {
	return Application{ a.Function.Substitute(v, e), a.Argument.Substitute(v, e) }
}

func (a Application) String() string {
	return fmt.Sprintf("(%v %v)", a.Function, a.Argument)
}
