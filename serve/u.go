package serve

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"
	"time"
)

type Claims struct {
	Exp int64 `json:"exp"`
}

var Jwt string

func IsJWTExpired(jwt string) string {
	if Jwt == "" {
		Jwt = GetJwtToken()
		return Jwt
	}
	parts := strings.Split(Jwt, ".")
	if len(parts) != 3 {
		return ""
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		log.Print(err)
		return ""
	}
	var claims Claims
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		log.Print(err)
		return ""
	}
	expTime := time.Unix(claims.Exp, 0)
	if time.Now().After(expTime) {
		Jwt = GetJwtToken()
		return Jwt
	}
	return Jwt
}
