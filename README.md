# cache
--
    import "github.com/buhduh/go-cache"

Package cache provides a modifiable thread safe cache system that separates
concerns along the way data is stored in the backend with the DataHandler
interface and when/how data is validated with the Invalidator interface.

## Examples

```
func ExampleNewTimedInvalidator() {
	lifetime, _ := time.ParseDuration(".5s")
	myCache := NewCache(nil, NewTimedInvalidator(lifetime))
	myCache.Put("foo", "a string")
	oneSec, _ := time.ParseDuration("1s")
	// IsValid is run every half second, make sure enough time has passed.
	time.Sleep(2 * oneSec)
	item, err := myCache.Get("foo")
	if item != nil || !IsValueNotPresentError(err) {
		fmt.Println("won't see this")
	} else {
		fmt.Println("expired")
	}
	// Output: expired
}

func ExampleNewInMemoryDataHandler() {
	// InMemoryDataHandler is the default.
	// could also do: NewCache(NewInMemoryDataHandler(), nil)
	myCache := NewCache(nil, nil)
	newItem, _ := myCache.Get("foo", 42)
	found, _ := myCache.Get("foo")
	if newItem.(int) != found.(int) {
		fmt.Println("not equal")
	} else {
		fmt.Println("equal")
	}
	// Output: equal
}
```

## Usage

#### func  IsValueNotPresentError

```go
func IsValueNotPresentError(err error) bool
```
IsValueNotPresentError is a simple test to determine if an error is of type
'ValueNotPresentError'.

#### type Cacher

```go
type Cacher interface {
	// Clear remove all elements from the cache.
	Clear()
	// Get a single element from the cache, if a second parameter is
	// provided, will set the cache to that value if nothing is present
	// returns a ValueNotPresentError if no value was found at key.
	Get(string, ...interface{}) (interface{}, error)
	// Put a value at key
	Put(string, interface{}) (interface{}, error)
	// Remove a single item, returning the item or a ValueNotPresentError
	// if no item is present.
	Remove(string) (interface{}, error)
	// Destroy the cache releasing resources.
	Destroy()
}
```

Cacher primary interface for this package.

#### func  NewCache

```go
func NewCache(

	dataHandler DataHandler,

	inv Invalidator,
) Cacher
```
NewCache returns a Cacher Interface whose behavior is determined by datahandler
and inv.

#### type DataHandler

```go
type DataHandler interface {
	// Put a single item in the cache
	Put(string, interface{}) error
	// Get a sinle item from the cache, must return a ValueNotPresentError
	// if there is nothing in the cache there.
	Get(string) (interface{}, error)
	// Clear removes all elements from the cache.
	Clear() error
	// Remove an item from the cache, returns a ValueNotPresentError
	// if there was nothing in the cache at key.
	Remove(string) error
	// Range iterates through all items in the cache and calls
	// the passed in function.  If the function returns false, iteration halts.
	Range(func(string, interface{}) bool)
}
```

DataHandler is the interface that the Cacher interface uses to actually store
data to the cache.

#### func  NewInMemoryDataHandler

```go
func NewInMemoryDataHandler() DataHandler
```
NewInMemoryDataHandler returns a Datahandler that is backed with a sync.Map.

#### type Invalidator

```go
type Invalidator interface {
	// IsValid determines whether or not a cache item is valid.
	// This package constantly loops through the cache in a background go routine
	// calling IsValid on all items.
	IsValid(*Metadata) bool

	// AccessExtra is called whenever Cacher.Get() is called and an item
	// is present.
	AccessExtra(*Metadata)
	// CreateExtra is called when Cacher.Put or in the case of an insertion when Cacher.Get
	// is called.
	CreateExtra(*Metadata)
	// UpdateExtra is called whenever Cacher.Put overwrites an existing item.
	UpdateExtra(*Metadata)
}
```

Invalidator is the interface that Cacher uses to determine if an item is valid.
If IsValid returns false, the item will be removed from the cache.

#### func  NewTimedInvalidator

```go
func NewTimedInvalidator(lifetime time.Duration) Invalidator
```
NewTimedInvalidator returns an Invaidator that validates cache based on
lifetime. Takes the most recent value of Metadata.Accessed, Metadata.Created, or
Metadata.Updated and compares to lifefime.

#### type Metadata

```go
type Metadata struct {
	// KeyCount is a thread safe pointer to the total count of the cache
	// KeyCount is managed in a background go routine.
	KeyCount *int64
	// Created is a Unix time stamp when an item was originally inserted into the cache.
	Created int64
	// Accessed is a Unix time stamp of the last time an item was retrieved with Cacher.Get
	Accessed int64
	// Modified is a Unix time stamp of the last time an item was modfied with Cacher.Put
	Modified int64
	// Extra provides a means for an outside implementation of Invalidator to determine
	// if an item is valid.
	Extra interface{}
}
```

Metadata is primarily used by Invalidator to determine if a cache item is valid.
Invalidator.AccessExtra, Invalidator.CreateExtra, and Invalidator.UpdateExtra
are intended to modify the Extra field in Metadata.

#### type ValueNotPresentError

```go
type ValueNotPresentError struct {
	Key string // The item key.
}
```

ValueNotPresentError is returned when an item isn't found in the cache.

#### func (ValueNotPresentError) Error

```go
func (v ValueNotPresentError) Error() string
```
Error satisfies the Error interface.
