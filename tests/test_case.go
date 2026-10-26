package tests

import (
	"github.com/goravel/framework/testing"

	"leave/bootstrap"
)

func init() {
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
