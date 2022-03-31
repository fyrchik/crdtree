package crdtree

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestState_Apply(t *testing.T) {
	rand.Seed(42)

	const (
		nodeCount = 4
		opCount   = 10
		iterCount = 100
	)

	s := NewState()
	ops := make([]Move, nodeCount)
	for i := range ops {
		ops[i] = Move{
			Time:   Timestamp(i),
			Parent: 0,
			Meta:   make([]byte, 5),
			Child:  Node(i + 1),
		}
		rand.Read(ops[i].Meta)
		s.Apply(ops[i])
	}

	ops = make([]Move, opCount)
	for i := range ops {
		ops[i] = Move{
			Time:   Timestamp(i + nodeCount),
			Parent: Node(rand.Uint32() % (nodeCount + 1)),
			Meta:   make([]byte, 5),
			Child:  Node(rand.Uint32()%nodeCount + 1),
		}
		rand.Read(ops[i].Meta)
	}

	expected := s.Clone()
	for i := range ops {
		expected.Apply(ops[i])
	}

	for i := 0; i < iterCount; i++ {
		rand.Shuffle(len(ops), func(i, j int) { ops[i], ops[j] = ops[j], ops[i] })

		actual := s.Clone()
		for i := range ops {
			actual.Apply(ops[i])
		}
		require.Equal(t, expected, actual)
	}
}

const benchNodeCount = 1000

func BenchmarkApplySequential(b *testing.B) {
	benchmarkApply(b, benchNodeCount, func(nodeCount, opCount int) []Move {
		ops := make([]Move, opCount)
		for i := range ops {
			ops[i] = Move{
				Time:   Timestamp(i),
				Parent: Node(rand.Intn(nodeCount)),
				Meta:   []byte{0, 1, 2, 3, 4},
				Child:  Node(rand.Intn(nodeCount)),
			}
		}
		return ops
	})
}

func BenchmarkApplyReorderLast(b *testing.B) {
	// Group operations in a blocks of 10, order blocks in increasing timestamp order,
	// and operations in a single block in reverse.
	const blockSize = 10

	benchmarkApply(b, benchNodeCount, func(nodeCount, opCount int) []Move {
		ops := make([]Move, opCount)
		for i := range ops {
			ops[i] = Move{
				Time:   Timestamp(i),
				Parent: Node(rand.Intn(nodeCount)),
				Meta:   []byte{0, 1, 2, 3, 4},
				Child:  Node(rand.Intn(nodeCount)),
			}
			if i != 0 && i%blockSize == 0 {
				for j := 0; j < blockSize/2; j++ {
					ops[i-j], ops[i+j-blockSize] = ops[i+j-blockSize], ops[i-j]
				}
			}
		}
		return ops
	})
}

func benchmarkApply(b *testing.B, n int, genFunc func(int, int) []Move) {
	rand.Seed(42)

	s := NewState()
	ops := genFunc(n, b.N)

	b.ResetTimer()
	b.ReportAllocs()
	for i := range ops {
		s.Apply(ops[i])
	}
}
