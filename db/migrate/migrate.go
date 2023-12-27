package migrate

import (
	"database/sql"

	"github.com/NBCUniversal/gvs-ms-common/logx"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/pkg/errors"

	// register migration from source file
	_ "github.com/golang-migrate/migrate/source/file"
)

const (
	driver = "mysql"
)

//go:generate mockery --output migratesmocks --outpkg migratesmocks --name Migrator
type Migrator interface {
	Migrate() error
}

type Options struct {
	DSN       string
	DBName    string
	SourceURL string
}

type migratorImpl struct {
	DNS        string
	DBName     string
	SourceURL  string
	SQLMigrate *migrate.Migrate
	Logger     migrate.Logger
}

func New(config Options) (Migrator, error) {
	mi := &migratorImpl{
		DNS:       config.DSN,
		DBName:    config.DBName,
		SourceURL: config.SourceURL,
	}
	mi.Logger = mi
	return mi, nil
}

func (mi *migratorImpl) Migrate() error {
	if mi == nil {
		return errors.New("migrate: migration is null")
	}
	session, err := sql.Open(driver, mi.DNS)
	if err != nil {
		return errors.Wrap(err, "db: failed to open connection")
	}

	driver, err := mysql.WithInstance(session, &mysql.Config{})
	if err != nil {
		return err
	}
	defer func() {
		if err := driver.Close(); err != nil {
			return
		}
	}()

	m, err := migrate.NewWithDatabaseInstance(mi.SourceURL, mi.DBName, driver)
	if err != nil {
		return err
	}

	m.Log = mi.Logger
	mi.SQLMigrate = m

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return errors.Wrap(err, "migrate : failed to get migration version")
	}

	if dirty {
		// force to prior version to reattempt migration
		err := m.Force(int(version) - 1)
		if err != nil {
			return errors.Wrap(err, "migrate : failed to force to previous previous schema version")
		}
		m.Log.Printf("migrate : forced to previous version: %v to reattempt migration", int(version)-1)
	}

	m.Log.Printf("migrate : run current schema version: %v dirty: %v", version, dirty)

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "migrate: failed to invoke up migrate")
	}

	if err != nil && err == migrate.ErrLocked {
		return errors.Wrap(err, "migrate : failed to database locked")
	}

	if _, ok := err.(migrate.ErrDirty); ok {
		// race condition if two or more front-ends reach db to run migration process.
		version, _, err := m.Version()
		if err != nil {
			return errors.Wrap(err, "migrate : failed to reattempt migration migration version")
		}

		mi.Logger.Printf("migrate : attempt to set recorded schema level to %d and rerun subsequent migration(s)", version)
		err = m.Force(int(version) - 1)
		if err != nil {
			return errors.Wrap(err, "migrate : failed to force to previous schema version")
		}

		if err := m.Up(); err != nil {
			return errors.Wrap(err, "migrate : re-migration attempt after force schema version")
		}
	}

	m.Log.Printf("migrate: completed db migration")

	_, err = m.Close()
	return errors.Wrap(err, "migrate:  failed to close migrations connection")
}

func (mi *migratorImpl) Printf(format string, v ...interface{}) {
	logx.Infof(format, v...)
}

func (mi *migratorImpl) Verbose() bool {
	return true
}
