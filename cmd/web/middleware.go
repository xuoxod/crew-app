package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/xuoxod/crew-app/internal/helpers"
)

func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteAddr := r.RemoteAddr
		host := r.Host
		path := r.URL.Path
		method := r.Method
		protocol := r.Proto
		protocolMajor := r.ProtoMajor
		protocolMinor := r.ProtoMinor

		fmt.Printf("\nPage Hit\n\tHost: %v\n\taddress: %v\n\tPath: %v\n\tMethod: %v\n\tProtocol: %v\n\t\tMajor: %v\n\t\tMinor: %v\n", host, remoteAddr, path, method, protocol, protocolMajor, protocolMinor)
		next.ServeHTTP(w, r)
	})
}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "Sign in first")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Unauth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "warning", "Resource not found")
			http.Redirect(w, r, "/user/dashboard", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
