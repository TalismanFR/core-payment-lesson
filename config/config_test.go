package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	p, _ := filepath.Abs("cfg1.yaml")
	p = filepath.Dir(p)
	f, err := os.CreateTemp(p, "cfg1.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name()) // clean up

	y := `http:
  port: 8000
  maxHeaderBytes: 1
  readTimeout: 10s
  writeTimeout: 10s

cache:
  ttl: 60s

auth:
  accessTokenTTL: 2h
  refreshTokenTTL: 720h #30 days
  verificationCodeLength: 8

limiter:
  rps: 10
  burst: 20
  ttl: 10m

vault:
  mountPath: "terminals"

postgres:
  dbName: "payservice-db"
  sslMode: "disable"
`

	if _, err := f.Write([]byte(y)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	cfg := Parse(f.Name())

	fmt.Printf("%#v\n", cfg)
}
