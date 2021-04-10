package datahandler

import (
	"fmt"
)

type DataHandler interface {
	Put(string, interface{}) error
	//Must throw a ValueNotPresentError and be nil
	Get(string) (interface{}, error)
	Clear() error
	//Must throw a ValueNotPresentError
	Remove(string) error
	Range(func(string, interface{}) bool)
}

type ValueNotPresentError struct {
	Key string
}

func (v ValueNotPresentError) Error() string {
	return fmt.Sprintf("no value found for key '%s'", v.Key)
}

func IsValueNotPresentError(err error) bool {
	_, ok := err.(ValueNotPresentError)
	return ok
}
