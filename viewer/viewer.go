package viewer

import (
	"context"
	"errors"
	"slices"
)

type Viewer interface {
	Can(scopes ...string) bool
	ID() string
	IDX() (string, error)
	Sess() string
}

type ctxKey struct{}

func FromContext(ctx context.Context) Viewer {
	v, _ := ctx.Value(ctxKey{}).(Viewer)
	return v
}

func newContext(parent context.Context, v Viewer) context.Context {
	return context.WithValue(parent, ctxKey{}, v)
}

type DefaultViewer struct {
	sub    string
	jti    string
	scopes []string
}

func (v DefaultViewer) ID() string {
	return v.sub
}

func (v DefaultViewer) IDX() (string, error) {
	if len(v.sub) > 0 && v.sub != "" {
		return v.sub, nil
	}
	return "", errors.New("unauthenticated viewer")
}

func (v DefaultViewer) Sess() string {
	return v.jti
}

func (v DefaultViewer) Can(scopes ...string) bool {
	for _, s := range scopes {
		if slices.Contains(v.scopes, s) {
			return true
		}
	}
	return false
}
