package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWKSet struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
	Alg string `json:"alg"`
}

type JWTMiddleware struct {
	idpURL     string
	publicKeys map[string]*rsa.PublicKey
	mu         sync.RWMutex
	lastUpdate time.Time
}

func NewJWTMiddleware(idpURL string) *JWTMiddleware {
	return &JWTMiddleware{
		idpURL:     idpURL,
		publicKeys: make(map[string]*rsa.PublicKey),
	}
}

func (jm *JWTMiddleware) fetchJWKS() error {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	if time.Since(jm.lastUpdate) < 1*time.Hour && len(jm.publicKeys) > 0 {
		return nil
	}

	resp, err := http.Get(jm.idpURL + "/.well-known/jwks.json")
	if err != nil {
		log.Printf("Failed to fetch JWKS: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to fetch JWKS: %d - %s", resp.StatusCode, string(body))
	}

	var jwkSet JWKSet
	if err := json.NewDecoder(resp.Body).Decode(&jwkSet); err != nil {
		log.Printf("Failed to decode JWKS: %v", err)
		return err
	}

	jm.publicKeys = make(map[string]*rsa.PublicKey)
	for _, jwk := range jwkSet.Keys {
		publicKey, err := jwkToPublicKey(jwk)
		if err != nil {
			log.Printf("Failed to convert JWK to public key: %v", err)
			continue
		}
		jm.publicKeys[jwk.Kid] = publicKey
	}

	jm.lastUpdate = time.Now()
	return nil
}

func jwkToPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, err
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, err
	}

	var e int
	for _, b := range eBytes {
		e = e*256 + int(b)
	}

	n := new(big.Int).SetBytes(nBytes)
	publicKey := &rsa.PublicKey{
		N: n,
		E: e,
	}

	return publicKey, nil
}

func (jm *JWTMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := jm.fetchJWKS(); err != nil && len(jm.publicKeys) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch public keys"})
			c.Abort()

			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()

			return
		}

		const bearerSchema = "Bearer "
		if len(authHeader) <= len(bearerSchema) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			c.Abort()

			return
		}

		tokenString := authHeader[len(bearerSchema):]

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, fmt.Errorf("missing kid in token header")
			}

			jm.mu.RLock()
			publicKey, exists := jm.publicKeys[kid]
			jm.mu.RUnlock()

			if !exists {
				return nil, fmt.Errorf("unknown kid: %s", kid)
			}

			return publicKey, nil
		})

		if err != nil {
			log.Printf("Failed to parse token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()

			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()

			return
		}

		c.Set("claims", claims)
		c.Set("user_id", claims["sub"])
		c.Set("username", claims["preferred_username"])
		c.Set("role", claims["role"])

		c.Next()
	}
}
