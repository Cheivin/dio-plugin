package mapper

import (
	"errors"
	"github.com/cheivin/dio-plugin/gorm/wrapper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func AutoMigrate[T any](db *gorm.DB, settings ...map[string]any) error {
	return Migrate(db, new(T), settings...)
}

func Migrate(db *gorm.DB, model any, settings ...map[string]any) error {
	db = db.Model(model)
	if len(settings) > 0 {
		for i := range settings {
			for k, v := range settings[i] {
				db = db.Set(k, v)
			}
		}
	}
	return db.AutoMigrate(model)
}

func Insert(db *gorm.DB, value any) error {
	return db.Create(value).Error
}

func InsertIgnore(db *gorm.DB, value any) error {
	return db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(value).Error
}

func Replace(db *gorm.DB, value any) error {
	return db.Clauses(clause.OnConflict{UpdateAll: true}).Create(value).Error
}

func Upsert(db *gorm.DB, value any, update *wrapper.Update) error {
	return db.Clauses(clause.OnConflict{DoUpdates: clause.Assignments(update.Data())}).Create(value).Error
}

func Where(db *gorm.DB, wrapper *wrapper.Query) *gorm.DB {
	return wrapper.Scope(db)
}

func Update(db *gorm.DB, update *wrapper.Update) (int64, error) {
	if update == nil {
		update = wrapper.U()
	}
	db = Where(db, update.Query()).Updates(update.Data())
	return db.RowsAffected, db.Error
}

func GetOne[T any](db *gorm.DB, cause *wrapper.Query) (record *T, err error) {
	if cause == nil {
		cause = wrapper.Q()
	}
	record = new(T)
	err = Where(db, cause).Limit(1).Take(record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return
}

func GetAll[T any](db *gorm.DB, cause *wrapper.Query) (records *[]T, err error) {
	records = new([]T)
	err = Where(db, cause).Find(records).Error
	return
}

func List[T any](db *gorm.DB, cause *wrapper.Query, limit ...int) (records *[]T, err error) {
	db = Where(db, cause)
	switch len(limit) {
	case 0:
		db = db.Offset(0)
	case 1:
		db = db.Offset(0).Limit(limit[0])
	default:
		db = db.Offset(limit[0]).Limit(limit[1])
	}
	records = new([]T)
	err = db.Find(records).Error
	return
}

func Count(db *gorm.DB, cause *wrapper.Query) (total int64, err error) {
	err = Where(db, cause).Count(&total).Error
	return
}

func Page[T any](db *gorm.DB, cause *wrapper.Query, page, size int) (records *[]T, total int64, err error) {
	total, err = Count(db, cause)
	if err != nil {
		return
	}
	records, err = List[T](db, cause, page*size, size)
	return
}

func Delete[T any](db *gorm.DB, cause *wrapper.Query) (int64, error) {
	model := new(T)
	db = Where(db, cause).Delete(model)
	return db.RowsAffected, db.Error
}

func Exist(db *gorm.DB, cause *wrapper.Query) (exist bool, err error) {
	err = Where(db.Select("1"), cause).Find(&exist).Error
	return
}

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

func Sum[V number](db *gorm.DB, field string, cause *wrapper.Query) (sum V, err error) {
	err = Where(db.Select("COALESCE(SUM("+field+"), 0)"), cause).
		Scan(&sum).
		Error
	return
}
