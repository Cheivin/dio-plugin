package dao

import (
	"database/sql"
	"github.com/cheivin/dio-plugin/gorm/mapper"
	"github.com/cheivin/dio-plugin/gorm/wrapper"
	"gorm.io/gorm"
)

type Mapper[T any] struct {
	db     *gorm.DB
	dst    T
	tabled bool
}

func NewMapper[T any](db *gorm.DB) *Mapper[T] {
	return &Mapper[T]{
		db: db,
	}
}

func (dao *Mapper[T]) DB() *gorm.DB {
	return dao.withModel()
}

func (dao *Mapper[T]) withModel() *gorm.DB {
	if dao.tabled {
		return dao.db
	}
	return dao.db.Model(dao.dst)
}

func (dao *Mapper[T]) Transaction(fn func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return dao.db.Transaction(fn, opts...)
}

func (dao *Mapper[T]) AutoMigrate(settings ...map[string]interface{}) error {
	//db := dao.withModel()
	//if len(settings) > 0 {
	//	for i := range settings {
	//		for k, v := range settings[i] {
	//			db = db.Set(k, v)
	//		}
	//	}
	//}
	//return db.AutoMigrate(dao.dst)
	return mapper.AutoMigrate[T](dao.db, settings...)
}

func (dao *Mapper[T]) Table(name string, args ...interface{}) *Mapper[T] {
	return &Mapper[T]{db: dao.db.Table(name, args...)}
}

func (dao *Mapper[T]) Select(query interface{}, args ...interface{}) *Mapper[T] {
	return NewMapper[T](dao.withModel().Select(query, args...))
}

func (dao *Mapper[T]) Distinct(args ...interface{}) *Mapper[T] {
	return NewMapper[T](dao.withModel().Distinct(args...))
}

func (dao *Mapper[T]) Where(wrapper *wrapper.Query) *gorm.DB {
	//return wrapper.Scope(dao.withModel())
	return mapper.Where(dao.withModel(), wrapper)
}

func (dao *Mapper[T]) GetOne(cause *wrapper.Query) (record *T, err error) {
	//if cause == nil {
	//	cause = wrapper.Q()
	//}
	//record = new(T)
	//err = dao.Where(cause).Limit(1).Find(record).Error
	//return
	return mapper.GetOne[T](dao.withModel(), cause)
}

func (dao *Mapper[T]) GetAll(cause *wrapper.Query) (records *[]T, err error) {
	//records = new([]T)
	//err = dao.Where(cause).Find(records).Error
	//return
	return mapper.GetAll[T](dao.withModel(), cause)
}

func (dao *Mapper[T]) List(cause *wrapper.Query, limit ...int) (records *[]T, err error) {
	//db := dao.Where(cause)
	//switch len(limit) {
	//case 0:
	//	db = db.Offset(0)
	//case 1:
	//	db = db.Offset(0).Limit(limit[0])
	//default:
	//	db = db.Offset(limit[0]).Limit(limit[1])
	//}
	//err = db.Find(records).Error
	//return
	return mapper.List[T](dao.withModel(), cause, limit...)
}

func (dao *Mapper[T]) Count(cause *wrapper.Query) (total int64, err error) {
	//err = dao.Where(cause).Count(&total).Error
	//return
	return mapper.Count(dao.withModel(), cause)
}

func (dao *Mapper[T]) Page(cause *wrapper.Query, page, size int) (records *[]T, total int64, err error) {
	//total, err = dao.Count(cause)
	//if err != nil {
	//	return
	//}
	//records, err = dao.List(cause, page*size, size)
	//return
	return mapper.Page[T](dao.withModel(), cause, page, size)
}

func (dao *Mapper[T]) Delete(cause *wrapper.Query) (int64, error) {
	//db := dao.Where(cause).Delete(dao.dst)
	//return db.RowsAffected, db.Error
	return mapper.Delete[T](dao.withModel(), cause)
}

func (dao *Mapper[T]) Exist(cause *wrapper.Query) (exist bool, err error) {
	//err = dao.Select("1").Where(cause).Find(cause, &exist).Error
	return mapper.Exist(dao.withModel(), cause)
}

func (dao *Mapper[T]) Sum(field string, cause *wrapper.Query) (sum int64, err error) {
	//err = dao.Select("COALESCE(SUM(" + field + "), 0)").
	//	Where(cause).
	//	Scan(&sum).
	//	Error
	return mapper.Sum[int64](dao.withModel(), field, cause)
}

func (dao *Mapper[T]) SumDecimal(field string, cause *wrapper.Query) (sum float64, err error) {
	//	err = dao.Select("COALESCE(SUM(" + field + "), 0)").
	//		Where(cause).
	//		Scan(&sum).
	//		Error
	return mapper.Sum[float64](dao.withModel(), field, cause)
}

func (dao *Mapper[T]) Insert(value interface{}) error {
	//return dao.db.Create(value).Error
	return mapper.Insert(dao.withModel(), value)
}

func (dao *Mapper[T]) InsertIgnore(value interface{}) error {
	//return dao.db.Clauses(clause.Insert{Modifier: "IGNORE"}).Create(value).Error
	return mapper.InsertIgnore(dao.withModel(), value)
}

func (dao *Mapper[T]) Replace(value interface{}) error {
	//return dao.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(value).Error
	return mapper.Replace(dao.withModel(), value)
}

func (dao *Mapper[T]) Upsert(value interface{}, update *wrapper.Update) error {
	//return dao.db.Clauses(clause.OnConflict{DoUpdates: clause.Assignments(update.Data())}).Create(value).Error
	return mapper.Upsert(dao.withModel(), value, update)
}

func (dao *Mapper[T]) Update(update *wrapper.Update) (int64, error) {
	//if update == nil {
	//	update = wrapper.U()
	//}
	//db := dao.Where(update.Query()).Updates(update.Data())
	//return db.RowsAffected, db.Error
	return mapper.Update(dao.withModel(), update)
}
