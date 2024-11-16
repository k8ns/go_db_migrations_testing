//go:build integration

package storage


import (
	"testing"
)

func TestStorageGet(t *testing.T) {
	_, err := env.db.Exec("INSERT INTO articles (title, created) VALUES('Atricle First', '2024-10-24 10:10:10')")
	if err != nil {
		t.Fatalf("failed to insert: %s", err)
	}

	s := &Storage{
		db: env.db,
	}

	_, err = s.Get(1)
	if err != nil {
		t.Fatalf("didn't expect to get error %s", err)
	}
}
