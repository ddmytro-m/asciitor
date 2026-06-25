package resolve

import "errors"

type ChainNode[T, K any] struct {
	next     Link[T, K]
	resolver Resolver[T, K]
	matcher  Matcher[T]
}

func NewNode[T, K any](r Resolver[T, K], m Matcher[T]) *ChainNode[T, K] {
	return &ChainNode[T, K]{
		resolver: r,
		matcher:  m,
	}
}

func (cn *ChainNode[T, K]) Next() Link[T, K] {
	return cn.next
}

func (cn *ChainNode[T, K]) SetNext(l Link[T, K]) {
	cn.next = l
}

func (cn *ChainNode[T, K]) Resolve(val T) (k K, _ error) {
	if cn.resolver == nil {
		return k, errors.New("resolver is not present")
	}

	return cn.resolver.Resolve(val)
}

func (cn *ChainNode[T, K]) Match(val T) bool {
	if cn.matcher == nil {
		return false
	}

	return cn.matcher.Match(val)
}
