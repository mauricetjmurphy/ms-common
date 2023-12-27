package db

import (
	"context"
	"time"

	"github.com/mauricetjmurphy/ms-common/clients/aws/secrets"
	"github.com/mauricetjmurphy/ms-common/db/migrate"
	"github.com/mauricetjmurphy/ms-common/db/query"
	"github.com/mauricetjmurphy/ms-common/logx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

//go:generate mockery --output dbmocks --outpkg dbmocks --name DB
type DB interface {
	Save(ctx context.Context, value interface{}) error
	Update(ctx context.Context, value interface{}) error
	Delete(ctx context.Context, value interface{}) error
	Query() *query.Query
	DBInstance() *gorm.DB
	Close() error
}

// DB is a wrapper for gorm.DB.
type dbImpl struct {
	*gorm.DB
}

// New creates new DB instance on given context and configuration options.
func New(ctx context.Context, opts ...Option) (DB, error) {
	dbConfig := NewConfigs(opts...)

	dsnDBConfig, err := loadDSNConfig(ctx, dbConfig)
	if err != nil {
		return nil, err
	}

	if dbConfig.RequiredMigration() {
		if err := runMigration(dsnDBConfig, dbConfig.MigrationSourceURL); err != nil {
			return nil, err
		}
	}

	tx, err := gorm.Open(mysql.Open(dsnDBConfig.toDSNConnString()), &gorm.Config{
		Logger: &dbLogger{},
	})

	if err != nil {
		return nil, errors.Wrap(err, "db : failed initialize db session")
	}

	_, err = tx.DB()
	if err != nil {
		return nil, errors.Wrap(err, "db : failed open connection")
	}

	return &dbImpl{tx}, nil
}

// Wrap creates one DB instance from existing connection.
func Wrap(tx *gorm.DB) DB {
	return &dbImpl{tx}
}

func (db *dbImpl) Save(ctx context.Context, value interface{}) error {
	if ctx != nil {
		db.WithContext(ctx)
	}
	results := db.Session(&gorm.Session{FullSaveAssociations: true}).Save(value)
	if results.Error != nil {
		return results.Error
	}
	return nil
}

func (db *dbImpl) Update(ctx context.Context, value interface{}) error {
	if ctx != nil {
		db.WithContext(ctx)
	}
	if result := db.Session(&gorm.Session{FullSaveAssociations: true}).Select("*").Updates(value); result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *dbImpl) Delete(ctx context.Context, value interface{}) error {
	if ctx != nil {
		db.WithContext(ctx)
	}
	if result := db.Unscoped().Select(clause.Associations).Delete(value); result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *dbImpl) Query() *query.Query {
	return query.New(db.DB)
}

func (db *dbImpl) DBInstance() *gorm.DB {
	return db.DB
}

func (db *dbImpl) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return errors.Wrap(err, "db : failed to get db")
	}
	return sqlDB.Close()
}

type dbLogger struct {
}

func (l *dbLogger) LogMode(logger.LogLevel) logger.Interface {
	return l
}

func (l *dbLogger) Info(ctx context.Context, s string, args ...interface{}) {
	logx.WithContext(ctx).Infof(s, args...)
}

func (l *dbLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	logx.WithContext(ctx).Warnf(s, args...)
}

func (l *dbLogger) Error(ctx context.Context, s string, args ...interface{}) {
	logx.WithContext(ctx).Errorf(s, args...)
}

func (l *dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := logx.Fields{}
	logx.WithContext(ctx).WithFields(logrus.Fields(fields)).Debugf("%s [%s]", sql, elapsed)
}

func loadDSNConfig(ctx context.Context, cfg *Config) (*dsnConf, error) {
	if cfg.IsLocal() {
		return cfg.toDSN(), nil
	}

	sm, err := secrets.New(ctx, cfg.AWSRegion)
	if err != nil {
		return nil, errors.Wrapf(err, "db : failed to create secrets managers instance on region %v", cfg.AWSRegion)
	}

	secret, err := sm.GetSecret(ctx, cfg.AWSSecretID)
	if err != nil {
		return nil, errors.Wrapf(err, "db : failed to load GetSecret(%v) on region %v", cfg.AWSSecretID, cfg.AWSRegion)
	}

	return &dsnConf{
		Host:     secret.Host,
		Port:     secret.Port,
		Name:     secret.Name,
		User:     secret.Username,
		Password: secret.Password,
	}, nil
}

func runMigration(dns *dsnConf, sourceURL string) error {
	mi, err := migrate.New(migrate.Options{
		DSN:       dns.toDSNConnString(),
		DBName:    dns.Name,
		SourceURL: sourceURL,
	})

	if err != nil {
		return errors.Wrap(err, "db : failed to create migrate instance")
	}

	if err = mi.Migrate(); err != nil {
		return errors.Wrap(err, "db : failed to run db migration")
	}

	return nil
}
