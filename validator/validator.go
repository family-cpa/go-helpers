package validator

import (
	"github.com/vektah/gqlparser/v2/ast"
)

type Validator interface {
	Validate(path ast.Path, model any) error
}
