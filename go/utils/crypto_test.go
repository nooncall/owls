package utils

import (
	"testing"

	"github.com/nooncall/owls/go/config"
	"github.com/nooncall/owls/go/utils/logger"
)

func TestCryPto(t *testing.T) {
	logger.InitLog(".", "test.log", "debug")
	config.InitConfig("../config/config.yaml")

	f := "TestCryPto"
	pwd := "aaaaaa"

	cryptoPwd, err := AesCrypto([]byte(pwd))
	if err != nil {
		t.Errorf("%s crypto err: %s", f, err.Error())
		t.FailNow()
	}
	deCryptoPwd, err := AesDeCrypto(cryptoPwd)
	if err != nil {
		t.Errorf("%s decrypto err: %s", f, err.Error())
		t.FailNow()
	}

	if pwd != string(deCryptoPwd) {
		t.Errorf("%s failed, not equal", f)
		t.FailNow()
	}
}
