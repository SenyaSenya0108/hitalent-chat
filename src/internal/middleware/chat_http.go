package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func HttpQueryParameter(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limitQuery := r.URL.Query().Get("limit")
		if limitQuery == "" {
			limitQuery = "20"
		}

		const maxLimit = 100
		limit, err := strconv.ParseInt(limitQuery, 10, 32)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", "Internal Server Error"), http.StatusInternalServerError)
			log.Println("parse page:", err)
			return
		}

		if limit > maxLimit {
			limit = maxLimit
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "limit", int(limit))
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}
