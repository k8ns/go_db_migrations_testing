package storage

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Config struct {
	Driver   string
	Username string
	Password []byte
	Host     string
	Port     string
	Database string
	Options  map[string]string
}

func NewDB(cfg *Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.Dsn())
	if err != nil {
		return nil, fmt.Errorf("couldn't open connection: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("couldn't ping: %v", err)
	}

	if err = RunMigrations(db); err != nil {
		return nil, fmt.Errorf("couldn't run migration %v", err)
	}

	return db, nil
}

func (c *Config) dsnArgs() []any {
	a := []any{c.Username, c.Password, c.Host, c.Port, c.Database}

	b := bytes.Buffer{}
	for k, v := range c.Options {
		if b.Len() > 0 {
			b.Write([]byte("&"))
		}
		b.Write([]byte(k))
		b.Write([]byte("="))
		b.Write([]byte(v))
	}

	a = append(a, b.String())
	return a
}

func (c *Config) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", c.dsnArgs()...)
}

func DriverName(conn *sql.DB) string {
	dr := fmt.Sprintf("%#v", conn.Driver())
	return dr[1:6]
}

func RunMigrations(db *sql.DB) error {
	driver, _ := mysql.WithInstance(db, &mysql.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://testdata/migrations/",
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("m.Up %v", err)
		}
	}

	return nil
}
