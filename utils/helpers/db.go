package helpers

import (
	"context"

	"gorm.io/gorm"
)

type key struct {
	name string
}

// WithTenantDB is a helper function that wraps a context with a tenant database.
func WithTenantDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, key{"db"}, db)
}

// GetTenantDB is a helper function that returns the tenant database from a context.
// If a transaction is present in the context, it will return the transaction instead.
func GetTenantDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(key{"txdb"}).(*gorm.DB); ok {
		return tx
	}
	return ctx.Value(key{"db"}).(*gorm.DB)
}

// WithTransaction is a helper function that projects a transaction into a context.
func WithTransaction(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, key{"txdb"}, tx)
}

// Transaction is a helper function that wraps a function with a transaction.
// The transaction is passed to the function as a context.
// Use GetTenantDB to get the transaction from the context.
func Transaction(ctx context.Context, f func(txCtx context.Context) error) error {
	db := GetTenantDB(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		txCtx := WithTransaction(ctx, tx)
		return f(txCtx)
	})
}
