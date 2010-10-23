package λeval



/* A VarSet represents a set of variables. */
type VarSet *varlist

/* Contains is the ∋ relation. It returns true if the Varset contains the
Variable, and false if not. */
func (s VarSet) Contains(v Variable) bool {
	return (*varlist)(s).contains(v)
}

/* Union is the ∪ operation. It returns a VarSet with all elements in either
argument VarSet, or both. */
func (s VarSet) Union(t Varset) Varset {
	return (*varlist)(s).union((*varlist)(t))
}

/* A varlist is the internal representation of a VarSet */
type varlist struct{
	head Variable
	tail *varlist
}

func cons(v Variable, t *varlist) *varlist {
	return &varlist{ v, t }
}

func (l varlist) contains(v Variable) bool {
	return l.head == v || tail.contains(v)
}

func (l varlist) union(m varlist) *varlist {
	if m == nil {
		return l
	}
	if l.contains(m.head) {return l.union(m.tail)}
	return cons(m.head, l).union(m.tail)
}

