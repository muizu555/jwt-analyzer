package example

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Issuer のチェックをしていない例（解析ツールで警告対象）
func ParseJWTWithoutIssuerCheck(tokenString string) {
	claims := &MyClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// シンプルなHMAC秘密鍵による署名検証
		return []byte("secret"), nil
	})

	if err != nil {
		fmt.Println("Error parsing token:", err)
		return
	}

	if token.Valid {
		fmt.Println("Token is valid")
		fmt.Println("Username:", claims.Username)
		// NOTE: Issuer チェックが抜けている
	}
}

// Issuer を明示的にチェックしている安全な例
func ParseJWTWithIssuerCheck(tokenString string) {
	claims := &MyClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		fmt.Println("Error parsing token:", err)
		return
	}

	if !token.Valid {
		fmt.Println("Invalid token")
		return
	}

	// Issuer チェック（セキュリティ的に重要）
	if claims.Issuer != "expected-issuer" {
		fmt.Println("Invalid issuer")
		return
	}

	fmt.Println("Token is valid with correct issuer")
	fmt.Println("Username:", claims.Username)
}
