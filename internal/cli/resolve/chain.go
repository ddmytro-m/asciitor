package resolve

import "errors"

type Chain[T, K any] struct {
	head   Link[T, K]
	tail   Link[T, K]
	length int
}

func NewChain[T, K any]() *Chain[T, K] {
	return &Chain[T, K]{}
}

func (c *Chain[T, K]) AddLink(l Link[T, K]) {
	if c.length == 0 {
		c.head = l
		c.tail = l
		c.length = 1
	} else {
		c.tail.SetNext(l)
		c.tail = l
		c.length++
	}
}

func (c *Chain[T, K]) Resolve(val T) (k K, _ error) {
	for cn := c.head; cn != nil; cn = cn.Next() {
		if !cn.Match(val) {
			continue
		}

		resolved, err := cn.Resolve(val)
		if err != nil {
			continue
		}

		return resolved, nil
	}
	return k, errors.New("value is not resolvable")
}

func (c *Chain[T, K]) Match(val T) bool {
	for cn := c.head; cn != nil; cn = cn.Next() {
		if cn.Match(val) {
			return true
		}
	}
	return false
}
