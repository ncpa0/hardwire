package utils

type Option[T any] struct {
	value *T
}

func NewOption[T any](value *T) *Option[T] {
	return &Option[T]{
		value: value,
	}
}

func Empty[T any]() *Option[T] {
	return &Option[T]{
		value: nil,
	}
}

func (o *Option[T]) IsNil() bool {
	return o.value == nil
}

func (o *Option[T]) Get() *T {
	return o.value
}

func (o *Option[T]) GetCopy() T {
	return *o.value
}

func (o *Option[T]) Set(value *T) {
	o.value = value
}
