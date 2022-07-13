package v2

import (
	"testing"

	"myst/test/v2/suite"
)

var setup func(*testing.T)

func TestIntegration(t *testing.T) {
	setup = suite.SetupSuite(t)
}
