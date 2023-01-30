package entity

import (
	"context"
	"errors"
	"simplestforum/internal/domain"
)

// AbstractTransaction is any transaction object (for example, *dbr.Tx).
type AbstractTransaction interface {
	RollbackUnlessCommitted()
	Commit() error
}

// Session contains all the information about the API query being processed.
// The context is intentionally put here because the structure is short-lived.
type Session struct {
	Ctx context.Context //nolint:containedctx

	ID          string
	UserID      int64
	Level       UserLevel
	Restriction UserRestriction

	Transaction AbstractTransaction

	RequestedFields RequestFields
}

// ctxKey represents keys in the Context.
type ctxKey string

const ctxSession ctxKey = "session"

// IsAuthorized returns true if the current user is logged in.
func (sess Session) IsAuthorized() bool {
	return sess.UserID != 0
}

// SessionInfoToError sets specific Session-related information needed to track the Error.
func SessionInfoToError(sess Session, err error) *domain.Error {
	var domainErr *domain.Error

	if !errors.As(err, &domainErr) {
		domainErr = domain.NewErrorWrap(err, domain.ErrCodeInternal, err.Error())
	}

	domainErr.UUID = sess.ID
	domainErr.UserID = sess.UserID

	return domainErr
}

// GetSession retrieves the Session out of the Context.
func GetSession(ctx context.Context) Session {
	value := ctx.Value(ctxSession)

	sess, ok := value.(Session)
	if !ok {
		return Session{}
	}

	return sess
}

// PutSession places the Session inside the Context.
func PutSession(ctx context.Context, sess Session) context.Context {
	return context.WithValue(ctx, ctxSession, sess)
}

// RequestFields contains the fields which should be returned in a GraphQL query.
type RequestFields map[string]RequestFields

// ContainsAny checks if the given field should be retrieved.
func (r RequestFields) ContainsAny(keys ...string) bool {
	for _, key := range keys {
		_, ok := r[key]

		if ok {
			return true
		}
	}

	return false
}

// Put registers a field to be retrieved.
func (r RequestFields) Put(key string) {
	r[key] = nil
}

// Delete removes a field to be retrieved.
func (r RequestFields) Delete(key string) {
	delete(r, key)
}

// MergeRequestFields puts requested fields from different GraphQL queries into a single structure.
func MergeRequestFields(r ...RequestFields) RequestFields {
	res := make(RequestFields)

	for _, item := range r {
		for k, v := range item {
			res[k] = v
		}
	}

	return res
}
