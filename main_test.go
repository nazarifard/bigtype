package bigtype

import (
	"os"
	"testing"

	"github.com/nazarifard/bigtype/log"
)

func TestMain(m *testing.M) {
	log.SetLogger(log.ZeroLog{})
	code := m.Run()
	os.Exit(code)
}
