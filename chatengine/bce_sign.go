package chatengine

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	BCE_AUTH_VERSION   = "bce-auth-v1"
	ISO8601Format      = "2006-01-02T15:04:05Z"
	AUTHORIZATION      = "Authorization"
	BCE_PREFIX         = "x-bce-"
	BCE_REQUEST_ID     = "x-bce-request-id"
	SIGN_JOINER        = "\n"
	SIGN_HEADER_JOINER = ";"
	expireSeconds      = 300
)

// 签名鉴权算法：https://cloud.baidu.com/doc/Reference/s/Njwvz1wot
func bceSign(
	accessKey, secretKey string,
	uri, method string,
	queryParams map[string]string,
	headers map[string]string,
	headersToSign map[string]struct{},
) string {
	signDate := formatISO8601Date(nowUTCSeconds())
	signKeyInfo := fmt.Sprintf("%s/%s/%s/%d",
		BCE_AUTH_VERSION,
		accessKey,
		signDate,
		expireSeconds)
	signKey := hmacSha256Hex(secretKey, signKeyInfo)
	canonicalUri := getCanonicalURIPath(uri)
	canonicalQueryString := getCanonicalQueryString(queryParams)
	canonicalHeaders, signedHeadersArr := getCanonicalHeaders(headers, headersToSign)

	// Generate signed headers string
	signedHeaders := ""
	if len(signedHeadersArr) > 0 {
		sort.Strings(signedHeadersArr)
		signedHeaders = strings.Join(signedHeadersArr, SIGN_HEADER_JOINER)
	}

	// Generate signature
	canonicalParts := []string{method, canonicalUri, canonicalQueryString, canonicalHeaders}
	canonicalReq := strings.Join(canonicalParts, SIGN_JOINER)
	signature := hmacSha256Hex(signKey, canonicalReq)

	// Generate auth string and add to the reqeust header
	authStr := signKeyInfo + "/" + signedHeaders + "/" + signature
	return authStr
}

func getCanonicalURIPath(path string) string {
	if len(path) == 0 {
		return "/"
	}
	canonical_path := path
	if strings.HasPrefix(path, "/") {
		canonical_path = path[1:]
	}
	canonical_path = uriEncode(canonical_path, false)
	return "/" + canonical_path
}

func getCanonicalQueryString(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	result := make([]string, 0, len(params))
	for k, v := range params {
		if strings.EqualFold(k, AUTHORIZATION) {
			continue
		}
		item := ""
		if len(v) == 0 {
			item = fmt.Sprintf("%s=", uriEncode(k, true))
		} else {
			item = fmt.Sprintf("%s=%s", uriEncode(k, true), uriEncode(v, true))
		}
		result = append(result, item)
	}
	sort.Strings(result)
	return strings.Join(result, "&")
}

func getCanonicalHeaders(headers map[string]string,
	headersToSign map[string]struct{}) (string, []string) {
	canonicalHeaders := make([]string, 0, len(headers))
	signHeaders := make([]string, 0, len(headersToSign))
	for k, v := range headers {
		headKey := strings.ToLower(k)
		if headKey == strings.ToLower(AUTHORIZATION) {
			continue
		}
		_, headExists := headersToSign[headKey]
		if headExists ||
			(strings.HasPrefix(headKey, BCE_PREFIX) &&
				(headKey != BCE_REQUEST_ID)) {

			headVal := strings.TrimSpace(v)
			encoded := uriEncode(headKey, true) + ":" + uriEncode(headVal, true)
			canonicalHeaders = append(canonicalHeaders, encoded)
			signHeaders = append(signHeaders, headKey)
		}
	}
	sort.Strings(canonicalHeaders)
	sort.Strings(signHeaders)
	return strings.Join(canonicalHeaders, SIGN_JOINER), signHeaders
}

func formatISO8601Date(timestamp_second int64) string {
	tm := time.Unix(timestamp_second, 0).UTC()
	return tm.Format(ISO8601Format)
}

func nowUTCSeconds() int64 { return time.Now().UTC().Unix() }

func hmacSha256Hex(key, data string) string {
	hasher := hmac.New(sha256.New, []byte(key))
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}

func uriEncode(uri string, encodeSlash bool) string {
	var byte_buf bytes.Buffer
	for _, b := range []byte(uri) {
		if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') ||
			b == '-' || b == '_' || b == '.' || b == '~' || (b == '/' && !encodeSlash) {
			byte_buf.WriteByte(b)
		} else {
			byte_buf.WriteString(fmt.Sprintf("%%%02X", b))
		}
	}
	return byte_buf.String()
}
