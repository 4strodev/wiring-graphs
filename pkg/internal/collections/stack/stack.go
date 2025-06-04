package stack

type Stack[T any] []T

func (q *Stack[T]) Add(v T) {
	*q = append(*q, v)
}

func (q *Stack[T]) Pop() (v T, returned bool) {
	if q.IsEmpty() {
		return
	}

	r := (*q)[len(*q)-1]
	*q = (*q)[:len(*q)-1]

	v = r
	returned = true

	return
}

func (q Stack[T]) IsEmpty() bool {
	return q == nil || len(q) == 0

}

func (q Stack[T]) Len() int {
	if q.IsEmpty() {
		return 0
	}

	return len(q)
}
