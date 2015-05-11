package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidUrl(t *testing.T) {
	opts := Options{}
	examples := []string{
		"postgresql://foobar",
		"foobar",
	}

	for _, val := range examples {
		opts.URL = val
		str, err := buildConnectionString(opts)

		assert.Equal(t, "", str)
		assert.Error(t, err)
		assert.Equal(t, "Invalid URL. Valid format: postgres://user:password@host:port/db?sslmode=mode", err.Error())
	}
}

func TestValidUrl(t *testing.T) {
	url := "postgres://myhost/database"
	str, err := buildConnectionString(Options{URL: url})

	assert.Equal(t, nil, err)
	assert.Equal(t, url, str)
}

func TestUrlAndSslFlag(t *testing.T) {
	str, err := buildConnectionString(Options{
		URL: "postgres://myhost/database",
		Ssl: "disable",
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://myhost/database?sslmode=disable", str)
}

func TestLocalhostUrlAndNoSslFlag(t *testing.T) {
	str, err := buildConnectionString(Options{
		URL: "postgres://localhost/database",
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://localhost/database?sslmode=disable", str)

	str, err = buildConnectionString(Options{
		URL: "postgres://127.0.0.1/database",
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://127.0.0.1/database?sslmode=disable", str)
}

func TestLocalhostUrlAndSslFlag(t *testing.T) {
	str, err := buildConnectionString(Options{
		URL: "postgres://localhost/database",
		Ssl: "require",
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://localhost/database?sslmode=require", str)

	str, err = buildConnectionString(Options{
		URL: "postgres://127.0.0.1/database",
		Ssl: "require",
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://127.0.0.1/database?sslmode=require", str)
}

func TestLocalhostUrlAndSslArg(t *testing.T) {
	str, err := buildConnectionString(Options{
		URL: "postgres://localhost/database?sslmode=require",
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://localhost/database?sslmode=require", str)

	str, err = buildConnectionString(Options{
		URL: "postgres://127.0.0.1/database?sslmode=require",
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://127.0.0.1/database?sslmode=require", str)
}

func TestFlagArgs(t *testing.T) {
	str, err := buildConnectionString(Options{
		Host:   "host",
		Port:   5432,
		User:   "user",
		Pass:   "password",
		DbName: "db",
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://user:password@host:5432/db", str)
}

func TestLocalhost(t *testing.T) {
	opts := Options{
		Host:   "localhost",
		Port:   5432,
		User:   "user",
		Pass:   "password",
		DbName: "db",
	}

	str, err := buildConnectionString(opts)
	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://user:password@localhost:5432/db?sslmode=disable", str)

	opts.Host = "127.0.0.1"
	str, err = buildConnectionString(opts)
	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://user:password@127.0.0.1:5432/db?sslmode=disable", str)
}

func TestLocalhostAndSsl(t *testing.T) {
	opts := Options{
		Host:   "localhost",
		Port:   5432,
		User:   "user",
		Pass:   "password",
		DbName: "db",
		Ssl:    "require",
	}

	str, err := buildConnectionString(opts)
	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://user:password@localhost:5432/db?sslmode=require", str)
}

func TestPort(t *testing.T) {
	opts := Options{Host: "host", User: "user", Port: 5000, DbName: "db"}
	str, err := buildConnectionString(opts)

	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres://user@host:5000/db", str)
}

func TestBlank(t *testing.T) {
	assert.Equal(t, true, connectionSettingsBlank(Options{}))
	assert.Equal(t, false, connectionSettingsBlank(Options{Host: "host", User: "user"}))
	assert.Equal(t, false, connectionSettingsBlank(Options{Host: "host", User: "user", DbName: "db"}))
	assert.Equal(t, false, connectionSettingsBlank(Options{URL: "url"}))
}
