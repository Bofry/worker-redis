package middleware

import (
	"fmt"
	"reflect"

	"github.com/Bofry/host"
	"github.com/Bofry/structproto"
	"github.com/Bofry/structproto/reflecting"
)

var _ structproto.StructBinder = new(MessageHandlerBinder)

type MessageHandlerBinder struct {
	messageHandlerType string
	components         map[string]reflect.Value
}

func (b *MessageHandlerBinder) Init(context *structproto.StructProtoContext) error {
	return nil
}

func (b *MessageHandlerBinder) Bind(field structproto.FieldInfo, target reflect.Value) error {
	if v, ok := b.components[field.Name()]; ok {
		if !target.IsValid() {
			return fmt.Errorf("specifiec argument 'target' is invalid. cannot bind '%s' to '%s'",
				field.Name(),
				b.messageHandlerType)
		}

		target = reflecting.AssignZero(target)
		if v.Type().ConvertibleTo(target.Type()) {
			target.Set(v.Convert(target.Type()))
		}
	}
	return nil
}

func (b *MessageHandlerBinder) Deinit(context *structproto.StructProtoContext) error {
	return b.preformInitMethod(context)
}

func (b *MessageHandlerBinder) preformInitMethod(context *structproto.StructProtoContext) error {
	rv := context.Target()
	if rv.CanAddr() {
		rv = rv.Addr()
		// call resource.Init()
		fn := rv.MethodByName(host.APP_COMPONENT_INIT_METHOD)
		if fn.IsValid() {
			if fn.Kind() != reflect.Func {
				return fmt.Errorf("fail to Init() resource. cannot find func %s() within type %s\n", host.APP_COMPONENT_INIT_METHOD, rv.Type().String())
			}
			if fn.Type().NumIn() != 0 || fn.Type().NumOut() != 0 {
				return fmt.Errorf("fail to Init() resource. %s.%s() type should be func()\n", rv.Type().String(), host.APP_COMPONENT_INIT_METHOD)
			}
			fn.Call([]reflect.Value(nil))
		}
	}
	return nil
}
