package dao

import (
	"database/sql"
	"github.com/cheivin/gorm-ext/mapper"
	"github.com/cheivin/gorm-ext/wrapper"
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

func ToMapper[T any](dao *Dao) *Mapper[T] {
	return NewMapper[T](dao.DB())
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
	return mapper.Where(dao.withModel(), wrapper)
}

func (dao *Mapper[T]) GetOne(cause *wrapper.Query) (record *T, err error) {
	return mapper.GetOne[T](dao.withModel(), cause)
}

func (dao *Mapper[T]) GetAll(cause *wrapper.Query) (records *[]T, err error) {
	return mapper.GetAll[T](dao.withModel(), cause)
}

func (dao *Mapper[T]) List(cause *wrapper.Query, limit ...int) (records *[]T, err error) {
	return mapper.List[T](dao.withModel(), cause, limit...)
}

func (dao *Mapper[T]) Count(cause *wrapper.Query) (total int64, err error) {
	return mapper.Count(dao.withModel(), cause)
}

func (dao *Mapper[T]) Page(cause *wrapper.Query, page, size int) (records *[]T, total int64, err error) {
	return mapper.Page[T](dao.withModel(), cause, page, size)
}

func (dao *Mapper[T]) Delete(cause *wrapper.Query) (int64, error) {
	return mapper.Delete[T](dao.withModel(), cause)
}

func (dao *Mapper[T]) Exist(cause *wrapper.Query) (exist bool, err error) {
	return mapper.Exist(dao.withModel(), cause)
}

func (dao *Mapper[T]) Sum(field string, cause *wrapper.Query) (sum int64, err error) {
	return mapper.Sum[int64](dao.withModel(), field, cause)
}

func (dao *Mapper[T]) SumDecimal(field string, cause *wrapper.Query) (sum float64, err error) {
	return mapper.Sum[float64](dao.withModel(), field, cause)
}

func (dao *Mapper[T]) Insert(value interface{}) error {
	return mapper.Insert(dao.withModel(), value)
}

func (dao *Mapper[T]) InsertIgnore(value interface{}) error {
	return mapper.InsertIgnore(dao.withModel(), value)
}

func (dao *Mapper[T]) Replace(value interface{}) error {
	return mapper.Replace(dao.withModel(), value)
}

func (dao *Mapper[T]) Upsert(value interface{}, update *wrapper.Update) error {
	return mapper.Upsert(dao.withModel(), value, update)
}

func (dao *Mapper[T]) Update(update *wrapper.Update) (int64, error) {
	return mapper.Update(dao.withModel(), update)
}
