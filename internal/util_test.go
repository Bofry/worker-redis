package internal

import (
	"reflect"
	"testing"
)

var _ MessageObserverAffinity = new(MockMessageObserverAffair)

type MockMessageObserverAffair struct{}

// MessageObserverTypes implements MessageObserverAffair.
func (*MockMessageObserverAffair) MessageObserverTypes() []reflect.Type {
	return nil
}

func TestIsMessageObserverAffair(t *testing.T) {
	var (
		affair = new(MockMessageObserverAffair)
	)

	rv := reflect.ValueOf(affair)

	ok := isMessageObserverAffinity(rv)
	var expectedOK bool = true
	if expectedOK != ok {
		t.Errorf("assert OK:: expected '%v', got '%v'", expectedOK, ok)
	}
}

func TestAsMessageObserverAffair(t *testing.T) {
	var (
		affair = new(MockMessageObserverAffair)
	)

	rv := reflect.ValueOf(affair)

	v := asMessageObserverAffinity(rv)
	if v == nil {
		t.Error("asMessageObserverAffair() should return non-nil")
	}
	_, ok := v.(MessageObserverAffinity)
	var expectedOK bool = true
	if expectedOK != ok {
		t.Errorf("assert OK:: expected '%v', got '%v'", expectedOK, ok)
	}
}
