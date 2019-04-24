package fcoin

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
)

//Hmac By MD5
func Hmac(key, data string) string {
	hmac := hmac.New(md5.New, []byte(key))
	hmac.Write([]byte(data))
	return hex.EncodeToString(hmac.Sum([]byte("")))
}

func sortedURI(u string, params url.Values) string {
	uu, _ := url.Parse(u)
	keys, kvs := make([]string,0), make([]string,0)
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		kvs = append(kvs, k + "=" + params.Get(k))
	}
	uu.RawQuery = strings.Join(kvs, "&")
	return uu.String()
}

func sortedBody(values url.Values) string {
	keys, kvs := make([]string,0), make([]string,0)
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		kvs = append(kvs, k + "=" + values.Get(k))
	}
	return strings.Join(kvs, "&")
}
