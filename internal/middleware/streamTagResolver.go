package middleware

import "github.com/Bofry/structproto"

func StreamTagResolve(fieldname, token string) (*structproto.Tag, error) {
	var tag *structproto.Tag
	if token != "" && token != "-" {
		tag = &structproto.Tag{
			Name: token,
		}
	}
	return tag, nil
}
