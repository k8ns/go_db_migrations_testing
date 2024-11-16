//go:build integration

//go:generate go run gen.go -dir=../migrations -out=testdata/schema.sql
package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

type TestEnv struct {
	ctx            context.Context
	mysqlContainer *mysql.MySQLContainer
	db             *sql.DB
}

var env = &TestEnv{}

func (e *TestEnv) initContainer() (err error) {
	e.mysqlContainer, err = mysql.Run(e.ctx,
		"mysql:8",
		mysql.WithDatabase("app"),
		mysql.WithUsername("root"),
		mysql.WithPassword("mainpswd"),
		mysql.WithScripts(filepath.Join("testdata", "schema.sql")),
	)

	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}

	state, err := e.mysqlContainer.State(e.ctx)
	if err != nil {
		log.Printf("failed to get container state: %s", err)
	}
	fmt.Println("is running:", state.Running)

	return
}

func (e *TestEnv) initDB() (err error) {
	conn, err := e.mysqlContainer.ConnectionString(e.ctx, "multiStatements=true", "charset=utf8", "parseTime=True", "loc=Local")
	if err != nil {
		return
	}

	e.db, err = sql.Open("mysql", conn)
	if err != nil {
		err = fmt.Errorf("couldn't open connection: %w", err)
		return
	}

	if err = e.db.Ping(); err != nil {
		err = fmt.Errorf("couldn't ping: %w", err)
		return
	}

	return
}

func (e *TestEnv) cleanup() (err error) {
	if err = e.db.Close(); err != nil {
		log.Printf("failed to close DB connection: %s", err)
	}

	if err := testcontainers.TerminateContainer(e.mysqlContainer); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}

	return
}

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	env.ctx = ctx

	if err := env.initContainer(); err != nil {
		panic(err)
	}

	if err := env.initDB(); err != nil {
		panic(err)
	}

	start := time.Now()

	fmt.Println("done testing:", m.Run(), time.Since(start))

	if err := env.cleanup(); err != nil {
		panic(err)
	}
}
