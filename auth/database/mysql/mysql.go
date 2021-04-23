package mysql

import (
	"blogfa/auth/config"
	zapLogger "blogfa/auth/pkg/logger"
	"fmt"
	"sync"

	"go.uber.org/zap"
	mql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Storage   store  = &msql{}
	namespace string = ""
	once      sync.Once
)

// store interface is interface for store things into mysql
type store interface {
	Connect(config config.GlobalConfig) error
	GetDatabase() *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Table(name string, args ...interface{}) *gorm.DB
	Preload(query string, args ...interface{}) *gorm.DB
}

// mysql struct
type msql struct {
	db *gorm.DB
}

// Connect method job is connect to mysql database and check migration
func (m *msql) Connect(config config.GlobalConfig) error {
	logger := zapLogger.GetZapLogger(config.Debug())
	var err error
	once.Do(func() {
		if config.MYSQL.Namespace != "" {
			namespace = config.MYSQL.Namespace
		}

		conf := &gorm.Config{}

		datasource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.MYSQL.Username,
			config.MYSQL.Password,
			config.MYSQL.Host,
			config.MYSQL.Schema,
		)

		m.db, err = gorm.Open(mql.Open(datasource), conf)
		if err != nil {
			zapLogger.Prepare(logger).
				Development().
				Level(zap.ErrorLevel).
				Commit(err.Error())

			return
		}

	})

	return err
}

// GetDatabase instance
func (m *msql) GetDatabase() *gorm.DB {
	return m.db
}

func (m *msql) Where(query interface{}, args ...interface{}) *gorm.DB {
	return m.db.Where(query, args...)
}

func (m *msql) First(dest interface{}, conds ...interface{}) *gorm.DB {
	return m.db.First(dest, conds...)
}

func (m *msql) Table(name string, args ...interface{}) *gorm.DB {
	return m.db.Table(name, args...)
}

func (m *msql) Preload(query string, args ...interface{}) *gorm.DB {
	return m.db.Preload(query, args...)
}
