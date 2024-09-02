package viewer

import (
	"net/http"
	"slices"
	"strings"
)

const (
	UserIDHeader = "x-user-id"
	ScopesHeader = "x-scopes"
)

type UserViewer struct {
	ID     string
	Scopes []string
}

func (v UserViewer) GetID() string {
	return v.ID
}

func (v UserViewer) Can(scopes ...string) bool {
	for _, s := range scopes {
		if slices.Contains(v.Scopes, s) {
			return true
		}
	}
	return false
}

func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get(UserIDHeader)
		scopes := r.Header.Get(ScopesHeader)

		ctx := newContext(r.Context(), UserViewer{
			ID:     userId,
			Scopes: strings.Split(scopes, " "),
		})

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
