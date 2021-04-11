package cache

type Reaper struct {
	*metadataHelper
	Invalidator
}

type ExtraCallback func(*Metadata)

func NewReaper(inv Invalidator, accessCB, createCB, updateCB ExtraCallback) *Reaper {
	return &Reaper{
		newMetadataHelper(accessCB, createCB, updateCB),
		inv,
	}
}
