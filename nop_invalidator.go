package cache

type NopInvalidator struct{}

func (n *NopInvalidator) IsValid(*Metadata) bool {
	return true
}
