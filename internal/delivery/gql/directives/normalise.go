package directives

import (
	"context"
	"simplestforum/internal/domain"
	"strings"

	"github.com/99designs/gqlgen/graphql"
)

func Normalise(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	val, err := next(ctx)
	if err != nil {
		return nil, err
	}

	switch typedVal := val.(type) {
	case string:
		return strings.TrimSpace(typedVal), nil
	case *string:
		if typedVal == nil {
			return nil, domain.NewError(domain.ErrCodeValidation, "The value must not be null")
		}

		*typedVal = strings.TrimSpace(*typedVal)

		return typedVal, nil
	default:
		return nil, domain.NewError(domain.ErrCodeValidation, "The value must be a string")
	}
}
