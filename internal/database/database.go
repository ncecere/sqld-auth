package database

import (
    "crypto/ed25519"
    "database/sql"
    "fmt"
    "os"
    "path/filepath"

    _ "github.com/mattn/go-sqlite3"
    "github.com/ncecere/sqld-auth/internal/crypto"
)

var dbFile string
var dbPassword string

func SetDBInfo(file, password string) {
    dbFile = file
    dbPassword = password
}

func SaveKeys(name string, pubKey ed25519.PublicKey, privKey ed25519.PrivateKey) error {
    dir := filepath.Dir(dbFile)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("error creating database directory: %w", err)
    }

    db, err := sql.Open("sqlite3", dbFile)
    if err != nil {
        return fmt.Errorf("error opening database: %w", err)
    }
    defer db.Close()

    if err := initDB(db); err != nil {
        return fmt.Errorf("error initializing database: %w", err)
    }

    encryptedPrivKey, err := crypto.Encrypt(privKey, []byte(dbPassword))
    if err != nil {
        return fmt.Errorf("error encrypting private key: %w", err)
    }

    _, err = db.Exec("INSERT OR REPLACE INTO keys (name, public_key, private_key) VALUES (?, ?, ?)",
        name, pubKey, encryptedPrivKey)
    if err != nil {
        return fmt.Errorf("error saving keys to database: %w", err)
    }

    return nil
}

func GetKeys(name string) (ed25519.PublicKey, ed25519.PrivateKey, error) {
    db, err := sql.Open("sqlite3", dbFile)
    if err != nil {
        return nil, nil, fmt.Errorf("error opening database: %w", err)
    }
    defer db.Close()

    var pubKey, encryptedPrivKey []byte
    err = db.QueryRow("SELECT public_key, private_key FROM keys WHERE name = ?", name).Scan(&pubKey, &encryptedPrivKey)
    if err != nil {
        return nil, nil, fmt.Errorf("error retrieving keys: %w", err)
    }

    privKeyBytes, err := crypto.Decrypt(encryptedPrivKey, []byte(dbPassword))
    if err != nil {
        return nil, nil, fmt.Errorf("error decrypting private key: %w", err)
    }

    return ed25519.PublicKey(pubKey), ed25519.PrivateKey(privKeyBytes), nil
}

func GetAllPublicKeys() (map[string][]byte, error) {
    db, err := sql.Open("sqlite3", dbFile)
    if err != nil {
        return nil, fmt.Errorf("error opening database: %w", err)
    }
    defer db.Close()

    rows, err := db.Query("SELECT name, public_key FROM keys")
    if err != nil {
        return nil, fmt.Errorf("error querying keys: %w", err)
    }
    defer rows.Close()

    certs := make(map[string][]byte)
    for rows.Next() {
        var name string
        var pubKey []byte
        if err := rows.Scan(&name, &pubKey); err != nil {
            return nil, fmt.Errorf("error scanning row: %w", err)
        }
        certs[name] = pubKey
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating rows: %w", err)
    }

    return certs, nil
}

func initDB(db *sql.DB) error {
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS keys (
            name TEXT PRIMARY KEY,
            public_key BLOB,
            private_key BLOB
        )
    `)
    return err
}