package crypto

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "io"

    "golang.org/x/crypto/scrypt"
)

func Encrypt(data, password []byte) ([]byte, error) {
    key, salt, err := deriveKey(password, nil)
    if err != nil {
        return nil, err
    }

    blockCipher, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(blockCipher)
    if err != nil {
        return nil, err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }

    ciphertext := gcm.Seal(nonce, nonce, data, nil)
    ciphertext = append(salt, ciphertext...)

    return ciphertext, nil
}

func Decrypt(data, password []byte) ([]byte, error) {
    salt, data := data[:32], data[32:]

    key, _, err := deriveKey(password, salt)
    if err != nil {
        return nil, err
    }

    blockCipher, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(blockCipher)
    if err != nil {
        return nil, err
    }

    nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, err
    }

    return plaintext, nil
}

func deriveKey(password, salt []byte) ([]byte, []byte, error) {
    if salt == nil {
        salt = make([]byte, 32)
        if _, err := rand.Read(salt); err != nil {
            return nil, nil, err
        }
    }

    key, err := scrypt.Key(password, salt, 1<<15, 8, 1, 32)
    if err != nil {
        return nil, nil, err
    }

    return key, salt, nil
}