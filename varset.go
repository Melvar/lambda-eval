package λeval


/* A VarSet represents a set of variables. */
type VarSet struct {
	l *varlist
}

func Singleton(v Variable) VarSet {
	return VarSet{cons(v, nil)}
}

/* Contains is the ∋ relation. It returns true if the Varset contains the
Variable, and false if not. */
func (s VarSet) Contains(v Variable) bool {
	return s.l.contains(v)
}

/* Union is the ∪ operation. It returns a VarSet with all elements in either
argument VarSet, or both. */
func (s VarSet) Union(t VarSet) VarSet {
	return VarSet{s.l.union(t.l)}
}

/* Without is the ∖{·} operation. It returns the VarSet except for the
specified element. */
func (s VarSet) Without(v Variable) VarSet {
	return VarSet{s.l.without(v)}
}

/* A varlist is the internal representation of a VarSet */
type varlist struct {
	head Variable
	tail *varlist
}

func cons(v Variable, t *varlist) *varlist {
	return &varlist{v, t}
}

func (l *varlist) contains(v Variable) bool {
	if l == nil {
		return false
	}
	return l.head == v || l.tail.contains(v)
}

func (l *varlist) union(m *varlist) *varlist {
	if m == nil {
		return l
	} else if l == nil {
		return m
	}
	if l.contains(m.head) {
		return l.union(m.tail)
	}
	return cons(m.head, l).union(m.tail)
}

func (l *varlist) without(v Variable) *varlist {
	if l.head == v {
		return l.tail
	}
	return cons(l.head, l.tail.without(v))
}
