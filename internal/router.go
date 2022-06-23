package internal

type Router map[string]MessageHandler

func NewRouter() *Router {
	var router Router = make(map[string]MessageHandler)
	return &router
}

func (r *Router) Add(stream string, handler MessageHandler) {
	route := *r
	route[stream] = handler
}

func (r *Router) Remove(stream string) {
	route := *r
	delete(route, stream)
}

func (r *Router) Get(stream string) MessageHandler {
	route := *r
	if v, ok := route[stream]; ok {
		return v
	}
	return nil
}
