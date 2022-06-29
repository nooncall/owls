package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"github.com/nooncall/owls/go/global"
)

//加密
func AesCrypto(source []byte) ([]byte, error) {
	var block cipher.Block
	var err error
	if block, err = aes.NewCipher([]byte(global.GVA_CONFIG.DBFilter.AesKey)); err != nil {
		return nil, fmt.Errorf("crypto err: %s", err.Error())
	}
	stream := cipher.NewCTR(block, []byte(global.GVA_CONFIG.DBFilter.AesIv))
	stream.XORKeyStream(source, source)
	return source, nil
}

//解密
func AesDeCrypto(cryptoData []byte) ([]byte, error) {
	return AesCrypto(cryptoData)
}
