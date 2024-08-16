package util

import (
	"context"
	"fmt"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/luc-wallace/idolvault/internal/models"
)

// JSONContentType sets the Content-Type of outbound requests to application/json before they are dispatched.
func JSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func HTMXOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("hx-request") != "true" {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, "forbidden")
			return
		}
		next.ServeHTTP(w, r)
	})
}

type session struct {
	Registering bool
	LoggedIn    bool
	User        *models.User
}

type PageParams struct {
	Session session
}

type contextKey int

const (
	KeyDefaultParams contextKey = iota
)

func UseDefaultParams(sess *scs.SessionManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := sess.Get(r.Context(), "user").(models.User)
			var registering bool
			if !ok {
				_, registering = sess.Get(r.Context(), "newUser").(models.NewUser)
			}

			defaultParams := PageParams{
				Session: session{
					LoggedIn:    ok,
					Registering: registering,
					User:        &user,
				},
			}
			req := r.WithContext(context.WithValue(r.Context(), KeyDefaultParams, defaultParams))

			next.ServeHTTP(w, req)
		})
	}
}

func GetParams(r *http.Request) PageParams {
	params, _ := r.Context().Value(KeyDefaultParams).(PageParams)
	return params
}

func RegisterRedirect(sess *scs.SessionManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/register" && sess.Exists(r.Context(), "newUser") {
				http.Redirect(w, r, "/register", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func CheckAuth(sess *scs.SessionManager, authState bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			exists := sess.Exists(r.Context(), "user")

			if !authState && exists {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			} else if authState && !exists {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func CheckAuthAPI(sess *scs.SessionManager, authState bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			exists := sess.Exists(r.Context(), "user")

			if !authState && exists {
				w.WriteHeader(http.StatusForbidden)
				fmt.Fprint(w, `forbidden`)
				return
			} else if authState && !exists {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, `you must be logged in to access this endpoint`)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
