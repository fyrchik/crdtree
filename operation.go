package crdtree

// Timestamp is an alias for integer timestamp type.
type Timestamp = int64

// Node is used to represent nodes.
type Node = uint64

// Meta represents arbitrary meta information.
type Meta []byte

// Move represents a single move operation.
type Move struct {
	Time   Timestamp
	Parent Node
	Meta   Meta
	Child  Node
}

// LogMove represents log record for a single move operation.
type LogMove struct {
	Move
	HasOld bool
	Old    PMPair
}
