package queue

type Queue[T any] []T

func (q *Queue[T]) Push(v T) {
	*q = append(*q, v)
}

func (q *Queue[T]) Pop() (v T, returned bool) {
	if q.IsEmpty() {
		return
	}

	r := (*q)[0]
	*q = (*q)[1:len(*q)]

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
