package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"reflect"
	"strings"
)

type GraphqlValidator struct {
	validator *validator.Validate
}

func NewGraphqlValidatorProvider() Validator {
	var validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &GraphqlValidator{validator: validate}
}

func (v *GraphqlValidator) Validate(path ast.Path, model any) error {
	if err := v.validator.Struct(model); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		errs := err.(validator.ValidationErrors)
		extensions := make(map[string]interface{}, len(errs))
		for _, ve := range errs {
			extensions[ve.Field()] = translate(ve)
		}

		return &gqlerror.Error{
			Path:    path,                // graphql.GetPath(ctx),
			Message: "Validation failed", // The given input was invalid
			Extensions: map[string]any{
				"code":       "GRAPHQL_VALIDATION_FAILED",
				"validation": extensions,
			},
		}
	}
	return nil
}

func translate(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "field is required"
	case "min":
		return "must be at least " + fe.Param()
	case "max":
		return "may not be greater than " + fe.Param()
	case "oneof":
		return "field does not exist in: " + fe.Param()
	case "fqdn":
		return "field must be a valid hostname"
	}
	return fe.Error()
}
