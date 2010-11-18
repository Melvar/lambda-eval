package λeval

import "testing"
import regex "regexp"

var yexp = Abstraction{Variable("g"), Application{Abstraction{Variable("x"), Application{Variable("g"), Application{Variable("x"), Variable("x")}}}, Abstraction{Variable("x"), Application{Variable("g"), Application{Variable("x"), Variable("x")}}}}}

var ystr = "(λg.((λx.(g (x x))) (λx.(g (x x)))))"

func TestPrint(t *testing.T) {
	var ypr = yexp.String()
	if ypr == ystr {
		t.Logf("Got %s\n from %#v\n", ypr, yexp)
	} else {
		t.Errorf("Got %q\n expected %q\n from %v\n", ypr, ystr, yexp)
	}
}

var redex = Application{Application{Abstraction{Variable("a"),Abstraction{Variable("b"),Application{Application{Variable("a"),Variable("a")},Variable("b")}}},Abstraction{Variable("x"),Abstraction{Variable("y"),Variable("y")}}},Abstraction{Variable("x"),Abstraction{Variable("y"),Variable("x")}}}

var reductregexstr = `\(λ(.)\.\(λ(.)\.(.)\)\)`
var reductexamplestr = "(λx.(λy.x))"

func TestReduce(t *testing.T) {
	var reduct = redex.Evaluate()

	var reductregex, err = regex.Compile(reductregexstr)
	if err != nil {
		t.Fatalf("Regex %s failed to compile: %v\n", reductregexstr, err)
	}
	var reductstr = reduct.String()
	var submatches = reductregex.FindStringSubmatch(reductstr)
	if submatches != nil && submatches[0] == reductstr && submatches[1] == submatches[3] && submatches[1] != submatches[2] {
		t.Logf("Got %s\n from %v\n", reduct, redex)
	} else {
		t.Errorf("Got %s\n expected %s (possibly with variables replaced)\n from %v\n", reductstr, reductexamplestr, redex)
	}
}
