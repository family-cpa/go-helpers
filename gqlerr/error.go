package gqlerr

import (
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func Default(path ...ast.Path) error {
	msg := "Internal server error"

	if len(path) > 0 {
		return &gqlerror.Error{
			Path:    path[0],
			Message: msg,
			Extensions: map[string]any{
				"code": "INTERNAL_SERVER_ERROR",
			},
		}
	}
	return &gqlerror.Error{Message: msg}
}

func Forbidden(path ...ast.Path) error {
	msg := "You are not authorized to perform this action"

	if len(path) > 0 {
		return &gqlerror.Error{
			Path:    path[0],
			Message: msg,
			Extensions: map[string]any{
				"code": "FORBIDDEN",
			},
		}
	}
	return &gqlerror.Error{Message: msg}
}
