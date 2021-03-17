package mysql

import (
	"blogfa/auth/config"
	"blogfa/auth/model/permission"
	"blogfa/auth/model/provider"
	"blogfa/auth/model/role"
	"blogfa/auth/model/user"
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
	AutoMigrate() error
	GetDatabase() *gorm.DB
}

// mysql struct
type msql struct {
	db *gorm.DB
}

// Connect method job is connect to mysql database and check migration
func (m *msql) Connect(config config.GlobalConfig) error {
	logger := zapLogger.GetZapLogger(false)
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

// AutoMigrate method for migrate to database
func (m *msql) AutoMigrate() error {
	sql := m.db.AutoMigrate(
		user.User{},
		role.Role{},
		permission.Permission{},
		provider.Provider{},
	)

	return sql
}

// GetDatabase instance
func (m *msql) GetDatabase() *gorm.DB {
	return m.db
}
