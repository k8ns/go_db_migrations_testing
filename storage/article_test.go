package storage

import (
	"testing"
)

func TestArticleString(t *testing.T) {
	a := Article{
		Id:      1,
		Title:   "Article",
		Created: "2024",
	}

	exp := "1:Article-2024"

	if a.String() != exp {
		t.Errorf("expected %s got %s", exp, a.String())
	}
}
