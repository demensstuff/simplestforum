package middleware

import (
	"context"
	"simplestforum/internal/domain/entity"

	"github.com/99designs/gqlgen/graphql"
)

// CollectFields puts the fields requested in a GraphQL query into the Session.
func CollectFields(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	op := graphql.GetOperationContext(ctx)
	if op.OperationName == "IntrospectionQuery" {
		return next(ctx)
	}

	sess := entity.GetSession(ctx)
	sess.RequestedFields = extractFields(op, graphql.CollectFieldsCtx(ctx, nil))
	ctx = entity.PutSession(ctx, sess)

	return next(ctx)
}

// extractFields recursively gets requested fields from the GraphQL provided structure.
func extractFields(op *graphql.OperationContext, fields []graphql.CollectedField) entity.RequestFields {
	if fields == nil {
		return nil
	}

	res := make(entity.RequestFields)
	for _, f := range fields {
		res[f.Name] = extractFields(op, graphql.CollectFields(op, f.Selections, nil))
	}

	return res
}
