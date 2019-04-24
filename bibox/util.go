package bibox

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
)

//Hmac By MD5
func Hmac(key, data string) string {
	hmac := hmac.New(md5.New, []byte(key))
	hmac.Write([]byte(data))
	return hex.EncodeToString(hmac.Sum([]byte("")))
}
