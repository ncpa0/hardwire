package resourceprovider

type Resource[T interface{}] interface {
	Get(c *DynamicRequestContext) (T, error)
}
