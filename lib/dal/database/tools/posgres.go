package tools

import (
	"database/sql"
	"flaver/globals"
	"flaver/globals/tools"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormSetting struct {
}

func (this GormSetting) GetDSN(config tools.Postgresql) string {
	output := make([]string, 0)
	if config.GetHost() != "" {
		output = append(output, fmt.Sprintf("host=%s", config.GetHost()))
	}
	if config.GetUsername() != "" {
		output = append(output, fmt.Sprintf("user=%s", config.GetUsername()))
	}
	if config.GetPassword() != "" {
		output = append(output, fmt.Sprintf("password=%s", config.GetPassword()))
	}
	if config.GetDBname() != "" {
		output = append(output, fmt.Sprintf("dbname=%s", config.GetDBname()))
	}
	if config.GetPort() != "" {
		output = append(output, fmt.Sprintf("port=%s", config.GetPort()))
	}
	output = append(output, config.GetConfig())
	return strings.Join(output, " ")
}

func (this GormSetting) GetLogConfig() logger.Interface {
	return logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: 300 * time.Millisecond, Colorful: true,
			IgnoreRecordNotFoundError: true, LogLevel: logger.Info,
		})
}

func (this GormSetting) GetPostgresConfig(config tools.Postgresql) postgres.Config {
	return postgres.Config{
		DSN:                  this.GetDSN(config),
		PreferSimpleProtocol: config.GetPreferSimpleProtocol(),
	}
}

func (this GormSetting) GetGormConfig(config tools.Postgresql) *gorm.Config {
	result := gorm.Config{
		Logger: this.GetLogConfig(),
		NowFunc: func() time.Time {
			location, _ := time.LoadLocation(globals.GetViper().GetString("server.time_location"))
			return time.Now().In(location)
		},
		CreateBatchSize: 200,
	}
	if !config.GetLogger() {
		result.Logger = result.Logger.LogMode(logger.Silent)
	}
	return &result
}

func (this GormSetting) SetupConns(config tools.Postgresql, conn *sql.DB) {
	conn.SetMaxIdleConns(config.MaxIdleConns)
	conn.SetMaxOpenConns(config.MaxOpenConns)

	idleTime := 30 * time.Second
	if config.MaxIdleTime > 0 {
		idleTime = time.Duration(config.MaxIdleTime) * time.Second
	}
	conn.SetConnMaxIdleTime(idleTime)

	lifeTime := 5 * time.Minute
	if config.MaxLifeTime > 0 {
		lifeTime = time.Duration(config.MaxLifeTime) * time.Second
	}
	conn.SetConnMaxLifetime(lifeTime)
}

func (this GormSetting) GetDatabase(config tools.Postgresql) *gorm.DB {
	postgresConfig, gormConfig := this.GetPostgresConfig(config), this.GetGormConfig(config)
	var client *gorm.DB
	var err error
	var conn *sql.DB

	if client, err = gorm.Open(postgres.New(postgresConfig), gormConfig); err != nil {
		globals.GetLogger().Fatalf("PostgreSql startup error")
		return nil
	} else if conn, err = client.DB(); err != nil {
		globals.GetLogger().Fatalf("PostgreSql startup error")
		return nil
	}

	this.SetupConns(config, conn)

	client.Exec("Select 1")
	globals.GetLogger().Infof("[gorm] response time(ms) : %v ", time.Since(time.Now()).Milliseconds())
	return client

}

var (
	clientDB    *gorm.DB
	gormMux     sync.Mutex
)

func getGormSetting(conn *gorm.DB, config tools.Postgresql) *gorm.DB {
	if conn != nil {
		return conn
	}
	gormMux.Lock()
	defer gormMux.Unlock()
	if conn != nil {
		return conn
	}
	return GormSetting{}.GetDatabase(config)
}

func GetClientDB() *gorm.DB {
	if clientDB == nil {
		clientDB = getGormSetting(clientDB, globals.GetConfig().GetPostgresql())
	}
	return clientDB
}
