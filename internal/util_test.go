package internal

import (
	"reflect"
	"testing"
)

var _ MessageObserverAffair = new(MockMessageObserverAffair)

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

	ok := isMessageObserverAffair(rv)
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

	v := asMessageObserverAffair(rv)
	if v == nil {
		t.Error("asMessageObserverAffair() should return non-nil")
	}
	_, ok := v.(MessageObserverAffair)
	var expectedOK bool = true
	if expectedOK != ok {
		t.Errorf("assert OK:: expected '%v', got '%v'", expectedOK, ok)
	}
}
