package middleware

import (
	"context"
	"log"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

// Recover intercepts panics and returns a correct GraphQL error.
func Recover(_ context.Context, err interface{}) error {
	log.Println("PANIC RECOVERED: ", err)

	return gqlerror.Errorf("Internal server error: %v", err)
}
