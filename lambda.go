/*	This package contains types and methods for representing and evaluating
	expressions of λ-calculus. It may, in the future, include a parser, or a
	parser may be written as a separate package.
*/
package λeval

import "fmt"

/* An Expression is a node in the syntax tree. */
type Expression interface {
	/* It can evaluate itself to some kind of minimal form */
	Evaluate() Expression
	/* It can cause all instances of a variable within it to be renamed */
	AlphaConvert(Variable, Variable) Expression
	/* It can substitute a free variable within itself with an expression */
	Substitute(Variable, Expression) Expression
	/* It can generate the set of its free variables */
	FreeVariables() VarSet
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

/* AlphaConvert yields the new name if the Variable is the one to replace,
otherwise it remains the same. */
func (v Variable) AlphaConvert(x, y Variable) Expression {
	if string(v) == string(x) {
		return y
	}
	return v
}

/* Substituting a Variable is analogous to α-converting it. */
func (t Variable) Substitute(v Variable, e Expression) Expression {
	if string(t) == string(v) {
		return e
	}
	return t
}

/* The set of free variables of a variable contains just that variable */
func (v Variable) FreeVariables() VarSet {
	return Singleton(v)
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
	return Abstraction{ a.Argument, a.Body.Evaluate() }
}

/* AlphaConvert yields the Abstraction with its Argument and Body
AlphaConverted. For now, it can still capture variables. */
func (a Abstraction) AlphaConvert(x, y Variable) Expression {
	return Abstraction{a.Argument.AlphaConvert(x, y), a.Body.AlphaConvert(x, y)} //TODO: Make non-capturing
}

/* Substitute yields the Abstraction with its Body Substituted. Its Argument is
not, since that cannot be any Expression but only a Variable. For now, it
can still capture variables. */
func (a Abstraction) Substitute(v Variable, e Expression) Expression {
	return Abstraction{a.Argument, a.Body.Substitute(v, e)} //TODO: Make non-capturing
}

/* The set of free variables of an abstraction is the set of free variables of
its body, minus its argument, which it binds. */
func (a Abstraction) FreeVariables() VarSet {
	return a.Body.FreeVariables().Without(a.Argument)
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
	return Application{ f, a.Argument }
}

/* AlphaConvert returns the Application with its Function and Argument
AlphaConverted. */
func (a Application) AlphaConvert(x, y Variable) Expression {
	return Application{a.Function.AlphaConvert(x, y), a.Argument.AlphaConvert(x, y)}
}

/* Substitute returns the Application with its Function and Argument
Substituted. */
func (a Application) Substitute(v Variable, e Expression) Expression {
	return Application{a.Function.Substitute(v, e), a.Argument.Substitute(v, e)}
}

/* The set of free variables of an Application is the union of the sets of free
variables of each of its subterms. */
func (a Application) FreeVariables() VarSet {
	return a.Function.FreeVariables().Union(a.Argument.FreeVariables())
}

/* String returns the Application formatted as in λ-calculus: “(F A)”, where
‘F’ is the Function and ‘A’ is the Argument. */
func (a Application) String() string {
	return fmt.Sprintf("(%v %v)", a.Function, a.Argument)
}
