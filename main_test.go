package main

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"testing"
	"time"
)

// GORM_REPO: https://github.com/go-gorm/gorm.git
// GORM_BRANCH: master
// TEST_DRIVERS: sqlite, mysql, postgres, sqlserver

func TestGORM(t *testing.T) {
	ctx := context.Background()
	m := User{
		Model: gorm.Model{
			ID: 1,
		},
	}
	var err error
	_ = GetDB(ctx, nil).Take(&m).Error

	err = GetDB(ctx, nil).Transaction(func(tx *gorm.DB) error {
		am := User{
			Model: gorm.Model{
				ID: 1,
			},
		}
		sql := fmt.Sprintf(`UPDATE %s SET age = 10, updated_at = ? WHERE id = 1`,
			GetTableName(ctx, tx, am))
		_err := GetDB(ctx, tx).Exec(sql, time.Now().Format("2006-01-02 15:04:05.999")).Error
		if _err != nil {
			return errors.New("更新失败：" + _err.Error())
		}

		var acc Toy
		_err = GetDB(ctx, tx).Model(Toy{}).Where("owner_id = ?", 10).Take(&acc).Error
		if _err != nil && _err != gorm.ErrRecordNotFound {
			return errors.New("失败：" + _err.Error())
		}

		return nil
	})
	t.Log(err)
}

func GetTableName(ctx context.Context, tx *gorm.DB, dest any) string {
	db := GetDB(ctx, tx)
	err := db.Statement.Parse(dest)
	if err == nil {
		return db.Statement.Schema.Table
	}
	return ""
}

func GetDB(ctx context.Context, tx *gorm.DB, args ...interface{}) (db *gorm.DB) {
	if tx != nil {
		db = tx
	} else {
		db = DB.WithContext(ctx) //TODO create session； because not is NewDB, so clone = 2; getInstance() does not satisfy clone ==1; Table != ""
	}
	return
}
func (p User) UpdateById(ctx context.Context, tx *gorm.DB, fields any) (int64, error) {
	res := GetDB(ctx, tx).Select(fields).Updates(p)
	return res.RowsAffected, res.Error
}
