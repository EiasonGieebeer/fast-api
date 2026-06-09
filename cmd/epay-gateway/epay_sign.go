package main

import (
	"crypto/md5"
	"fmt"
	"sort"
	"strings"
)

// ============================================================================
// Epay Protocol Signing — server-side implementation
// Compatible with github.com/Calcium-Ion/go-epay client
// ============================================================================

// paramsFilter removes sign, sign_type, and empty values before signing
func paramsFilter(params map[string]string) map[string]string {
	filtered := make(map[string]string)
	for k, v := range params {
		if k == "sign" || k == "sign_type" || v == "" {
			continue
		}
		filtered[k] = v
	}
	return filtered
}

// paramsSort returns sorted keys and corresponding values
func paramsSort(params map[string]string) ([]string, []string) {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	values := make([]string, len(keys))
	for i, k := range keys {
		values[i] = params[k]
	}
	return keys, values
}

// createUrlString builds "key1=val1&key2=val2&..."
func createUrlString(keys, values []string) string {
	var sb strings.Builder
	for i, key := range keys {
		sb.WriteString(key)
		sb.WriteByte('=')
		sb.WriteString(values[i])
		sb.WriteByte('&')
	}
	return strings.TrimSuffix(sb.String(), "&")
}

// md5String computes MD5(urlString + key)
func md5String(urlString, key string) string {
	digest := md5.Sum([]byte(urlString + key))
	return fmt.Sprintf("%x", digest)
}

// generateParams signs the params and adds sign + sign_type
func generateParams(params map[string]string, key string) map[string]string {
	filtered := paramsFilter(params)
	ks, vs := paramsSort(filtered)
	sign := md5String(createUrlString(ks, vs), key)
	params["sign"] = sign
	params["sign_type"] = "MD5"
	return params
}

// verifySign checks whether the signature in params matches
func verifySign(params map[string]string, key string) bool {
	sign, ok := params["sign"]
	if !ok {
		return false
	}
	expected := generateParams(params, key)["sign"]
	return sign == expected
}
