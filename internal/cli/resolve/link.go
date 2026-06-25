package resolve

type Link[T, K any] interface {
	Resolver[T, K]
	Matcher[T]
	Node[T, K]
}

type Resolver[T, K any] interface {
	Resolve(T) (K, error)
}

type Matcher[T any] interface {
	Match(T) bool
}

type Node[T, K any] interface {
	Next() Link[T, K]
	SetNext(Link[T, K])
}
