package cmd

import (
    "crypto/ed25519"
    "fmt"
    "os"
    "strings"
    "time"

    "github.com/spf13/cobra"
    "github.com/ncecere/sqld-auth/internal/database"
    "github.com/ncecere/sqld-auth/internal/token"
)

var createCertCmd = &cobra.Command{
    Use:   "createcert [name] [namespaces...]",
    Short: "Create a new certificate with key pair and token",
    Args:  cobra.MinimumNArgs(2),
    Run:   createCert,
}

func init() {
    rootCmd.AddCommand(createCertCmd)
}

func createCert(cmd *cobra.Command, args []string) {
    name := args[0]
    namespaces := args[1:]

    if err := generateAndSaveCert(cmd, name, namespaces); err != nil {
        cmd.PrintErrf("Error creating certificate: %v\n", err)
        os.Exit(1)
    }

    cmd.Printf("Certificate created successfully for %s\n", name)
}

func generateAndSaveCert(cmd *cobra.Command, name string, namespaces []string) error {
    pubKey, privKey, err := ed25519.GenerateKey(nil)
    if err != nil {
        return fmt.Errorf("error generating key pair: %w", err)
    }

    if err := database.SaveKeys(name, pubKey, privKey); err != nil {
        return fmt.Errorf("error saving keys: %w", err)
    }

    // Use a default expiration of 24 hours for the initial token
    tokenStr, err := token.CreateTokenWithNamespaces(privKey, "full", 24*time.Hour)
    if err != nil {
        return fmt.Errorf("error creating token: %w", err)
    }

    cmd.Printf("Name: %s\n", name)
    cmd.Printf("Namespaces: %s\n", strings.Join(namespaces, ", "))
    cmd.Printf("Token (full access, expires in 24h): %s\n", tokenStr)

    return nil
}