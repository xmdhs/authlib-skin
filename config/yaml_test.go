package config

import (
	"os"
	"testing"
)

func TestYamlDeCode(t *testing.T) {
	b, err := os.ReadFile("config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	c, err := YamlDeCode(b)
	if err != nil {
		t.Fatal(err)
	}
	if c.Sql.MysqlDsn != "123" {
		t.FailNow()
	}
}
