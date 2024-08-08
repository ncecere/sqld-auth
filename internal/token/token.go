package token

import (
    "crypto/ed25519"
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v4"
)

func CreateTokenWithNamespaces(privKey ed25519.PrivateKey, access string, expiration time.Duration) (string, error) {
    exp := time.Now().Add(expiration)
    claims := jwt.MapClaims{
        "exp": exp.Unix(),
    }

    switch access {
    case "read":
        claims["a"] = "ro"
    case "write":
        claims["a"] = "rw"
    case "attach-read":
        claims["a"] = "roa"
    case "full":
        // No claims["a"] set for full access
    default:
        return "", fmt.Errorf("invalid access level: %s", access)
    }

    token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
    return token.SignedString(privKey)
}