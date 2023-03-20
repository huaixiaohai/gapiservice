package dao

import (
	"context"

	"gorm.io/gorm"
)

const GormTransKey = "gorm_trans_key"

func Transaction(ctx context.Context, fc func(context.Context) error) error {
	tx := ctx.Value(GormTransKey)
	if tx != nil {
		return fc(ctx)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		return fc(context.WithValue(ctx, GormTransKey, tx))
	})
}

func getSession(ctx context.Context) *gorm.DB {
	v := ctx.Value(GormTransKey)
	if v != nil {
		return v.(*gorm.DB)
	}
	return db
}
