package gpsdec

type objectqueue struct {
	items      []object
	curr       *object
	len        uint32
	iter       uint32
	desc       string
	descalphaX float64
	descalphaY float64
}

func newObjectQueue(d string, daX, daY float64) objectqueue {
	return objectqueue{
		items:      []object{},
		curr:       nil,
		len:        0,
		iter:       0,
		desc:       d,
		descalphaX: daX,
		descalphaY: daY,
	}
}

func (q *objectqueue) push(o object) {
	o.desc = q.desc
	o.descalphaX = q.descalphaX
	o.descalphaY = q.descalphaY
	q.items = append(q.items, o)
	q.len++
	q.curr = &o
}

func (q *objectqueue) next() *object {
	return q.curr
}

func (q *objectqueue) prev() {
	if q.len != 0 {
		q.iter = (q.iter - 1) % q.len
		q.curr = &q.items[q.iter]
	}
}

func (q *objectqueue) roll() {
	if q.len != 0 {
		q.iter = (q.iter + 1) % q.len
		q.curr = &q.items[q.iter]
	}
}
