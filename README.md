# Highly-available CRDT for trees

Go implementation of a paper [1]. Only basic algorithm is done, without optimizations
(log retention, fast creation, separate list CRDT for metadata).
Also, the `apply` operation is implemented iteratively instead of recursively.

[1] https://martin.kleppmann.com/papers/move-op.pdf