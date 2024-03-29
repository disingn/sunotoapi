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

func getLastUserContent(data map[string]interface{}) string {
	var lastUserContent string
	messages, ok := data["messages"].([]interface{})
	if !ok {
		return lastUserContent
	}

	for i := len(messages) - 1; i >= 0; i-- {
		message, ok := messages[i].(map[string]interface{})
		if !ok {
			continue
		}

		role, ok := message["role"].(string)
		if !ok {
			continue
		}

		if role == "user" {
			lastUserContent, _ = message["content"].(string)
			break
		}
	}

	return lastUserContent
}

func IsJWTExpired() (string, *ErrorResponse) {
	if Jwt == "" {
		var err *ErrorResponse
		Jwt, err = GetJwtToken()
		if err != nil {
			return "", err
		}
	}
	parts := strings.Split(Jwt, ".")
	if len(parts) != 3 {
		return "", NewErrorResponse(ErrCodeResponseInvalid, "invalid JWT format. Expected format: header.payload.signature")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		log.Printf("IsJWTExpired failed, error decoding JWT payload: %v", err)
		return "", NewErrorResponseWithError(ErrCodeResponseInvalid, err)
	}
	var claims Claims
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		log.Printf("IsJWTExpired failed, error unmarshalling JWT claims: %v", err)
		return "", NewErrorResponseWithError(ErrCodeJsonFailed, err)
	}
	expTime := time.Unix(claims.Exp, 0)
	if time.Now().After(expTime) {
		var err *ErrorResponse
		Jwt, err = GetJwtToken()
		if err != nil {
			return "", err
		}
	}
	return Jwt, nil
}
