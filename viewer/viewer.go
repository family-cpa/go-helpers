package viewer

import (
	"context"
)

type Viewer interface {
	Can(scopes ...string) bool
	GetID() string
}

type ctxKey struct{}

func FromContext(ctx context.Context) Viewer {
	v, _ := ctx.Value(ctxKey{}).(Viewer)
	return v
}

func newContext(parent context.Context, v Viewer) context.Context {
	return context.WithValue(parent, ctxKey{}, v)
}
