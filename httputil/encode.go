package httputil

import (
	"net/url"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func utf8ToGBK(utf8str string) string {
	result, _, _ := transform.String(simplifiedchinese.GBK.NewEncoder(), utf8str)
	return result
}

func gbkToUtf8(gbKstr []byte) []byte {
	result, _, _ := transform.Bytes(simplifiedchinese.GBK.NewDecoder(), gbKstr)
	return result
}

func EscapeKeywordByGBK(keyword string) string {
	return url.QueryEscape(utf8ToGBK(keyword))
}
