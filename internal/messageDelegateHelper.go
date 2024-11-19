package internal

import "fmt"

type MessageDelegateHelper struct{}

func (MessageDelegateHelper) IsRestricted(msg *Message) (bool, error) {
	d, ok := msg.Delegate.(*ContextMessageDelegate)
	if !ok {
		return false, fmt.Errorf("invalid operation")
	}
	return d.isRestricted(), nil
}

func (MessageDelegateHelper) Restrict(msg *Message) error {
	d, ok := msg.Delegate.(*ContextMessageDelegate)
	if !ok {
		return fmt.Errorf("invalid operation")
	}
	d.restrict()
	return nil
}

func (MessageDelegateHelper) Unrestrict(msg *Message) error {
	d, ok := msg.Delegate.(*ContextMessageDelegate)
	if !ok {
		return fmt.Errorf("invalid operation")
	}
	d.unrestrict()
	return nil
}
