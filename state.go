package crdtree

// PMPair couples parent and metadata.
type PMPair struct {
	Parent Node
	Meta   Meta
}

// Tree is a mapping from the child nodes to their parent and metadata.
type Tree map[Node]PMPair

// State represents state being replicated.
type State struct {
	Operations []LogMove
	Tree       Tree
}

// NewState constructs new empty tree.
func NewState() *State {
	return &State{
		Tree: Tree{},
	}
}

// Clone returns deep-copy of s.
func (s *State) Clone() *State {
	ns := &State{
		Operations: make([]LogMove, len(s.Operations)),
		Tree:       make(Tree, len(s.Tree)),
	}

	copy(ns.Operations, s.Operations)
	for c := range s.Tree {
		ns.Tree[c] = s.Tree[c]
	}

	return ns
}

// undo un-does op and changes s in-place.
func (s *State) undo(op LogMove) {
	delete(s.Tree, op.Child)
	if op.HasOld {
		s.Tree[op.Child] = op.Old
	}
}

// redo applies op as the last operation and changes s in-place.
func (s *State) redo(op LogMove) {
	s.Operations = append(s.Operations, s.do(op.Move))
}

// Apply puts op in log at a proper position, re-applies all subsequent operations
// from log and changes s in-place.
func (s *State) Apply(op Move) {
	var index int
	for index = len(s.Operations); index > 0; index-- {
		if s.Operations[index-1].Time <= op.Time {
			break
		}
	}

	if index == len(s.Operations) {
		s.Operations = append(s.Operations, s.do(op))
		return
	}

	s.Operations = append(s.Operations[:index+1], s.Operations[index:]...)
	for i := len(s.Operations) - 1; i > index; i-- {
		s.undo(s.Operations[i])
	}

	s.Operations[index] = s.do(op)
	
	for i := index + 1; i < len(s.Operations); i++ {
		s.Operations[i] = s.do(s.Operations[i].Move)
	}
}

// do executes a single move operation on a tree.
func (s *State) do(op Move) LogMove {
	p, ok := s.Tree[op.Child]
	lm := LogMove{
		Move: Move{
			Time:   op.Time,
			Parent: op.Parent,
			Meta:   op.Meta,
			Child:  op.Child,
		},
		HasOld: ok,
		Old:    p,
	}

	if !s.Tree.isAncestor(op.Child, op.Parent) {
		p.Meta = op.Meta
		p.Parent = op.Parent
		s.Tree[op.Child] = p
	}

	return lm
}

// isAncestor returns true if parent is an ancestor of a child.
// For convenience, also return true if parent == child.
func (t Tree) isAncestor(parent, child Node) bool {
	for c := child; c != parent; {
		p, ok := t[c]
		if !ok {
			return false
		}
		c = p.Parent
	}
	return true
}
