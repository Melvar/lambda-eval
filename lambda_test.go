package λeval

import "testing"

var yexp = Abstraction{Variable("g"), Application{Abstraction{Variable("x"), Application{Variable("g"), Application{Variable("x"), Variable("x")}}}, Abstraction{Variable("x"), Application{Variable("g"), Application{Variable("x"), Variable("x")}}}}}

var ystr = "(λg.((λx.(g (x x))) (λx.(g (x x)))))"

func TestPrint(t *testing.T) {
	var ypr = yexp.String()
	if ypr == ystr {
		t.Logf("Got %s\n from %#v\n", ypr, yexp)
	} else {
		t.Errorf("Got %q\n expected %q\n from %#v\n", ypr, ystr, yexp)
	}
}
