package beauter

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatHCL(t *testing.T) {
	input := `job "example" {
  datacenters = ["dc1"]

  group "cache" {
    task "redis" {
      driver = "docker"

      config {
        image       = "redis:3.2"

        port_map {
          db = 6379
        }
      }

      resources {
        cpu    = 500
        memory = 256

        network {
          mbits = 10
          port  "db"  {}
        }
      }
    }
  }
}`
	expected := `job "example" {
  datacenters = ["dc1"]

  group "cache" {
    task "redis" {
      driver = "docker"

      config {
        image = "redis:3.2"

        port_map {
          db = 6379
        }
      }

      resources {
        cpu    = 500
        memory = 256

        network {
          mbits = 10
          port "db" {}
        }
      }
    }
  }
}`
	result, err := FormatHCL(input)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestFormatJSON(t *testing.T) {
	input := `{"a":"b"}`
	expected := `{
	"a": "b"
}`
	result, err := FormatJSON(input)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}
