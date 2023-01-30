package directives

import (
	"context"
	"simplestforum/internal/domain"
	"unicode/utf8"

	"github.com/99designs/gqlgen/graphql"
)

func Range(ctx context.Context, _ interface{}, next graphql.Resolver, min int64, max int64) (interface{}, error) {
	val, err := next(ctx)
	if err != nil {
		return nil, err
	}

	var l int64

	switch typedVal := val.(type) {
	case string:
		l = int64(utf8.RuneCountInString(typedVal))
	case []int64:
		l = int64(len(typedVal))
	case int64:
		l = typedVal
	case *string:
		if typedVal == nil {
			return nil, domain.NewError(domain.ErrCodeValidation, "The value must not be null")
		}

		l = int64(utf8.RuneCountInString(*typedVal))
	case *int64:
		if typedVal == nil {
			return nil, domain.NewError(domain.ErrCodeValidation, "The value must not be null")
		}

		l = *typedVal
	default:
		return nil, domain.NewError(domain.ErrCodeValidation, "The value must be a string or an int")
	}

	if l < min || l > max {
		return nil, domain.NewError(domain.ErrCodeValidation, "The value must be between %d and %d (or %d and %d characters long)", min,
			max, min, max)
	}

	return val, nil
}
