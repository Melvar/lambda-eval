/*	This package contains types and methods for representing and evaluating
	expressions of λ-calculus. A recursive descent parser has been begun.
*/
package λeval

import "fmt"

/* An Expression is a node in the syntax tree. */
type Expression interface {
	/* It can evaluate itself to some kind of minimal form */
	Evaluate() Expression
	/* It can substitute a free variable within itself with an expression */
	Substitute(Variable, Expression) Expression
	/* It can check whether a given Variable is free in it */
	ContainsFree(Variable) bool
	/* Finally, it can make a string representation of itself */
	fmt.Stringer
}

/* A Variable represents a variable in a λ-expression. It’s basically just a
symbol. It implements Expression. */
type Variable string

/* Evaluate yields the Variable itself, since there is no way to reduce it. */
func (v Variable) Evaluate() Expression {
	return v
}

/* Substitute yields the Expression if the Variable is the one to replace, the
same Variable otherwise. */
func (t Variable) Substitute(v Variable, e Expression) Expression {
	if t == v {
		return e
	}
	return t
}

/* Containsfree returns true if this is the Variable being asked about. */
func (t Variable) ContainsFree(v Variable) bool {
	return t == v
}

/* String returns the Variable’s identifying name. */
func (v Variable) String() string {
	return string(v)
}

/* An Abstraction represents a λ-abstraction. For all intents and purposes, it
is a function, consisting of an argument and a body. */
type Abstraction struct {
	Argument Variable
	Body     Expression
}

/* Evaluate yields the Abstraction with its body in simplest form. */
func (a Abstraction) Evaluate() Expression {
	return Abstraction{a.Argument, a.Body.Evaluate()}
}

/* Substitute yields the Abstraction with its Body Substituted. If necessary,
the Argument is changed (the Abstraction α-converted) to prevent capture of a
Variable free in the Expression. */
func (a Abstraction) Substitute(v Variable, e Expression) Expression {
	if v == a.Argument {
		return a
	}
	var n = a.Argument
	for e.ContainsFree(n) {
		n += "′"
	}
	return Abstraction{n, a.Body.Substitute(a.Argument, n).Substitute(v, e)}
}

/* Containsfree yields true if the Variable asked about is free in the Body,
but is not the one bound by the Abstraction */
func (a Abstraction) ContainsFree(v Variable) bool {
	return a.Argument != v && a.Body.ContainsFree(v)
}

/* String returns the Abstraction formatted as in λ-calculus: “(λa.B)”, where
‘a’ is the Argument and ‘B’ is the Body. */
func (a Abstraction) String() string {
	return fmt.Sprintf("(λ%v.%v)", a.Argument, a.Body)
}

/* An Application represents the application of a function to an argument. */
type Application struct {
	Function Expression
	Argument Expression
}

/* Evaluate yields the Evaluation of the result of applying the Function to the
Argument if the Evaluation of the Function is an Abstraction, and itself with
the Function Evaluated otherwise. */
func (a Application) Evaluate() Expression {
	var f = a.Function.Evaluate()
	if l, ok := f.(Abstraction); ok {
		return l.Body.Substitute(l.Argument, a.Argument).Evaluate()
	}
	return Application{f, a.Argument}
}

/* Substitute returns the Application with its Function and Argument
Substituted. */
func (a Application) Substitute(v Variable, e Expression) Expression {
	return Application{a.Function.Substitute(v, e), a.Argument.Substitute(v, e)}
}

/* ContainsFree returns true if the Variable is free in either subexpression */
func (a Application) ContainsFree(v Variable) bool {
	return a.Function.ContainsFree(v) || a.Argument.ContainsFree(v)
}

/* String returns the Application formatted as in λ-calculus: “(F A)”, where
‘F’ is the Function and ‘A’ is the Argument. */
func (a Application) String() string {
	return fmt.Sprintf("(%v %v)", a.Function, a.Argument)
}
