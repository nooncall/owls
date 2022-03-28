package utils

import (
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/logger"
)

func TestCryPto(t *testing.T) {
	logger.InitLog(".", "test.log", "debug")
	config.InitConfig("../config/config.yml")

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
