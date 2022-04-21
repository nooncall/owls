package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/qingfeng777/owls/server/config"
)

//加密
func AesCrypto(source []byte) ([]byte, error) {
	var block cipher.Block
	var err error
	if block, err = aes.NewCipher([]byte(config.Conf.Server.AesKey)); err != nil {
		return nil, fmt.Errorf("crypto err: %s", err.Error())
	}
	stream := cipher.NewCTR(block, []byte(config.Conf.Server.AesIv))
	stream.XORKeyStream(source, source)
	return source, nil
}

//解密
func AesDeCrypto(cryptoData []byte) ([]byte, error) {
	return AesCrypto(cryptoData)
}
