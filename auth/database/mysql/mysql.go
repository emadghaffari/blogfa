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
)

type store interface {
	Connect(config config.GlobalConfig) error
	AutoMigrate() error
	GetDatabase() *gorm.DB
}

type msql struct {
	db   *gorm.DB
	once sync.Once
}

func (m *msql) Connect(config config.GlobalConfig) error {
	logger := zapLogger.GetZapLogger(false)
	var err error
	m.once.Do(func() {
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

		if config.MYSQL.Automigrate {
			if err := m.AutoMigrate(); err != nil {
				zapLogger.Prepare(logger).
					Add("err", "database automigrate").
					Development().
					Level(zap.ErrorLevel).
					Commit(err.Error())

				return
			}
		}

	})

	return err
}

func (m *msql) AutoMigrate() error {
	sql := m.db.AutoMigrate()

	return sql
}

func (m *msql) GetDatabase() *gorm.DB {
	return m.db
}
