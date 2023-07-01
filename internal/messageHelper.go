package internal

import "fmt"

type MessageHelper struct{}

func (MessageHelper) IsRestricted(msg *Message) (bool, error) {
	d, ok := msg.Delegate.(*ContextMessageDelegate)
	if !ok {
		return false, fmt.Errorf("invalid operation")
	}
	return d.isRestricted(), nil
}

func (MessageHelper) Restrict(msg *Message) error {
	d, ok := msg.Delegate.(*ContextMessageDelegate)
	if !ok {
		return fmt.Errorf("invalid operation")
	}
	d.restrict()
	return nil
}

func (MessageHelper) Unrestrict(msg *Message) error {
	d, ok := msg.Delegate.(*ContextMessageDelegate)
	if !ok {
		return fmt.Errorf("invalid operation")
	}
	d.unrestrict()
	return nil
}
