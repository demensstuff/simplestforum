package middleware

import (
	"context"
	"encoding/json"
	"simplestforum/internal/delivery"
	"simplestforum/internal/domain/entity"

	"github.com/99designs/gqlgen/graphql"
)

const introspectionQuery = "IntrospectionQuery"

// WrapResponse is a middleware function which wraps the response into the CustomResponse.
func WrapResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	res := next(ctx)
	sess := entity.GetSession(ctx)

	op := graphql.GetOperationContext(ctx)
	if op.OperationName == introspectionQuery {
		return res
	}

	if len(res.Errors) != 0 {
		return buildErrorResponse(sess, res.Errors[0])
	}

	var dataMap map[string]interface{}

	err := json.Unmarshal(res.Data, &dataMap)
	if err != nil {
		panic("GraphQL response is not a map!")
	}

	if len(dataMap) == 0 {
		panic("GraphQL response doesn't contain entries!")
	}

	for _, v := range dataMap {
		return buildSuccessResponse(v)
	}

	return res
}

// buildSuccessResponse creates a GraphQL response which indicates success.
func buildSuccessResponse(payload interface{}) *graphql.Response {
	return &graphql.Response{
		Data: delivery.BuildSuccessResponse(payload, false),
	}
}

// buildErrorResponse creates an erroneous GraphQL response.
func buildErrorResponse(sess entity.Session, err error) *graphql.Response {
	return &graphql.Response{
		Data: delivery.BuildErrorResponse(sess, err, false),
	}
}
