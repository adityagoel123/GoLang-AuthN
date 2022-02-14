package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if err := checkHealth(); err != nil {
		http.Error(w, "health check failed", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "OK..Health Check Passed\n")
}

func messagesHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*User)
	if !ok {
		http.Error(w, "no user", http.StatusInternalServerError)
		return
	}
	log.Printf("User thus received from Context is : %s", user)

	// FIXME:
	fmt.Fprint(w, "[]\n")
}

func authToken(r *http.Request) string {
	hdr := r.Header.Get("Authorization")
	return strings.TrimPrefix(hdr, "Bearer ")
}

func requireAuth(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token := authToken(r)
		fmt.Println("The token thus received in the request is :", token)

		user := userFromToken(token)
		fmt.Println("The User Object thus received from Token :", user)

		if user == nil {
			http.Error(w, "bad authentication", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func main() {
	http.HandleFunc("/health", healthHandler)
	h := requireAuth(http.HandlerFunc(messagesHandler))
	http.Handle("/messages", h)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
