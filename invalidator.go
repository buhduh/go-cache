package cache

type Invalidator interface {
	IsValid(*Metadata) bool
}
