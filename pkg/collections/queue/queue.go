package queue

type Queue[T any] []T

func (q *Queue[T]) Push(v T) {
	*q = append(*q, v)
}

func (q *Queue[T]) Pop() (v T, returned bool) {
	if q.IsEmpty() {
		return
	}

	r := (*q)[len(*q)-1]
	*q = (*q)[:len(*q)-1]

	v = r
	returned = true

	return
}

func (q Queue[T]) IsEmpty() bool {
	return q == nil || len(q) == 0

}

func (q Queue[T]) Len() int {
	if q.IsEmpty() {
		return 0
	}

	return len(q)
}
