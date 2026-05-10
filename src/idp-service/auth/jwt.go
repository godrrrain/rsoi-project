package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"lab2/src/idp-service/models"

	"github.com/golang-jwt/jwt/v5"
)

const kid = "kid-1"
const rsaPrivateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA6gnYIy4Lv3Y/SQ4D6NpCIf3S8pYpY7QQhOkiZEeQPKBK9Xeb
PZYhFBNDF1AQMty8c9/cKfouuha0XFmRhhfmsvH0cd4NsxG3Qnwk0bf6v2sMpVV7
0sihyPzAw2IoJPk5pBN83qu9DjiC5WMcdtFIlasU/ezlG6h3ZpEydvE/wU945HoV
SUZpXCDUlPBHzcBQKoST20dQTzM7DXlLz86uPkRwySMyzNAgGu+JELWqgu3NDtn+
bs3fjghZqRj81ubEO6EMB/A3rpJ8agkfVo8VlJvr2MFnjZqmwU0GGi80BHWpKW2j
iyg9rYdXtkxVeiXtVzcHxJg5IqzUZ8AE68dS0QIDAQABAoIBAASm6lEdM0I0lNKU
PO1ICNpdfe/g3pze7Xjh3DNVdZNR7ZmCXY9oZHsOXLYZ+x9ysiuEuK7sLFl6Ne6z
wFPztcwK6YypoPi1LoELoAz74OflBibja4sbnjUlOoTrtp/xBQV7Fm9j5yWxJ39t
JFnDhFOGx3uaZc2yYCC1nETfW+DBFJp8QT57dGHdorqOdCNEXdvVhRQe4QMsK8Pi
ZlOlVYcuSPeQB+1wKGvb/S0x5xLGd3zfvMJ28o41uu5KqebmsNfOWppyTu5hAN5W
6dFj0jm7TiTtNUMRuqHPNXKcVjHd6LagSDn0NP1KzNUxvxrOG5zFEh4jmG4M4hzi
UjfTECECgYEA+xW9zrKaruwi8SLUe+bDpCMEc9zBYOZCVU04Z2Nxb7PfgC/EQBzF
P+VS26d3W5FKMHroRil4wtXViGPuhEiKvatSFuSYJOtFSLt6upCtCV07pwyk/lkC
RgbpuL3AMeqTuLctixpdKJxvbnFqKtHlg6b84eVW7C1NUwRjeGejtbECgYEA7p6t
lxO29/hiF+1WgSuVMYXlnlbXXYK47U5W7zK9rOHCGLsS7CjYwTrmFOcJEsM0Qfpl
UbLQI9lj+v0mAuEASlBBHvYv5blHoFAWN5zen3I26f9xdBHmRPyxyOH/NCxC1LGU
LBhTHrI94DSxLNh6skBvOsWtpDuBP+uJPIR6FyECgYEAgsz/tVcr5+ZSCaaoZOeB
kddAMY+WGgG6GrAAqzON27ArxZ6csP2L8E5qDM3ACy60JG9S44IlS/KTq9rLXZRg
2pAOUqjBbbI2xL4OIHTP/+nW8p5OscXyvkJJrZkEL7zROdALZMTWNRrRngptUWNJ
Gn16jb+ouZ6cApxtqULscPECgYEA5v/N5Lc9JYjazXcBi0J5x9trkoFXNDtccr6o
AiAI5tgWYoKXqu9QBp/SJOIUMomuiUCx3QlR3aKR22Q97AONmGNg52xEqgtXf6aI
G4ZNLeYPqy+S0V6SoK5QHbxKpmNCv0y5uIZD0S+UHvxjmJppDS67fxXnJ1pDoXGP
BXrqBoECgYEAr64rTOCngLOUFVCX4uuNWBcOCmZ7pMXFQ8f5MjxSg4R+JLGVsLo2
T+rzjrYbYdEkBBNofMjI174FyqOO1LKzJ5jr6f6r708CjCQBpQtUtlXQjHcdIkSR
69eO3ac+oZmLd9VeaYUnPaXrqcjDAS0ct5wimhSLMYfshDy2JLBbNPE=
-----END RSA PRIVATE KEY-----
`

type JWTManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	issuer     string
}

func NewJWTManager(issuer string) (*JWTManager, error) {
	block, _ := pem.Decode([]byte(rsaPrivateKey))
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key %w", err)
	}

	return &JWTManager{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
		issuer:     issuer,
	}, nil
}

func (jm *JWTManager) GenerateAccessToken(user *models.User, clientId string, scopes string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(24 * 30 * time.Hour)

	claims := jwt.MapClaims{
		"sub":                user.UserUid.String(),
		"preferred_username": user.Username,
		"email":              user.Email,
		"name":               user.FullName,
		"role":               user.Role,
		"scope":              scopes,
		"iat":                now.Unix(),
		"exp":                expiresAt.Unix(),
		"iss":                jm.issuer,
		"aud":                clientId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid

	tokenString, err := token.SignedString(jm.privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (jm *JWTManager) GenerateIdToken(user *models.User, clientId string, scopes string, nonce string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(24 * 30 * time.Hour)

	claims := jwt.MapClaims{
		"iss": jm.issuer,
		"sub": user.UserUid.String(),
		"aud": clientId,
		"iat": now.Unix(),
		"exp": expiresAt.Unix(),
	}

	if Contains(scopes, "profile") {
		claims["preferred_username"] = user.Username
		claims["name"] = user.FullName
	}

	if Contains(scopes, "email") {
		claims["email"] = user.Email
		claims["email_verified"] = true
	}

	if nonce != "" {
		claims["nonce"] = nonce
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(jm.privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateAuthorizationCode() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func GenerateRefreshToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (jm *JWTManager) GetPublicKeyJWK() (models.JWK, error) {
	publicKey := jm.publicKey

	nBytes := publicKey.N.Bytes()
	eBytes := make([]byte, 4)
	eBytes[0] = byte(publicKey.E >> 24)
	eBytes[1] = byte(publicKey.E >> 16)
	eBytes[2] = byte(publicKey.E >> 8)
	eBytes[3] = byte(publicKey.E)

	for i := 0; i < len(eBytes); i++ {
		if eBytes[i] != 0 {
			eBytes = eBytes[i:]
			break
		}
	}

	return models.JWK{
		Kty: "RSA",
		Use: "sig",
		Kid: kid,
		N:   base64.RawURLEncoding.EncodeToString(nBytes),
		E:   base64.RawURLEncoding.EncodeToString(eBytes),
		Alg: "RS256",
	}, nil
}

func (jm *JWTManager) GetJWKS() models.JWKS {
	jwk, _ := jm.GetPublicKeyJWK()
	return models.JWKS{
		Keys: []models.JWK{jwk},
	}
}

func (jm *JWTManager) VerifyAccessToken(tokenString string) (*jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jm.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &claims, nil
}

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func VerifyPassword(password, hash string) bool {
	return HashPassword(password) == hash
}

func (jm *JWTManager) ExportJWKS() []byte {
	jwks := jm.GetJWKS()
	data, _ := json.MarshalIndent(jwks, "", "  ")
	return data
}
