package cmd

import (
    "encoding/pem"
    "fmt"
    "os"
    "strings"

    "github.com/spf13/cobra"
    "github.com/ncecere/sqld-auth/internal/database"
)

var showAll bool

var showCertCmd = &cobra.Command{
    Use:   "showcert [name]",
    Short: "Show public certificate for a name",
    Args:  cobra.MaximumNArgs(1),
    Run:   showCert,
}

func init() {
    rootCmd.AddCommand(showCertCmd)
    showCertCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all public certificates")
}

func showCert(cmd *cobra.Command, args []string) {
    if showAll {
        showAllCerts(cmd)
        return
    }

    if len(args) == 0 {
        cmd.PrintErr("Error: name is required when not using --all flag\n")
        cmd.Usage()
        os.Exit(1)
    }

    name := args[0]
    pubKey, _, err := database.GetKeys(name)
    if err != nil {
        cmd.PrintErrf("Error retrieving keys: %v\n", err)
        os.Exit(1)
    }

    pemBlock := &pem.Block{
        Type:  "PUBLIC KEY",
        Bytes: pubKey,
    }
    fmt.Print(string(pem.EncodeToMemory(pemBlock)))
}

func showAllCerts(cmd *cobra.Command) {
    certs, err := database.GetAllPublicKeys()
    if err != nil {
        cmd.PrintErrf("Error retrieving all public keys: %v\n", err)
        os.Exit(1)
    }

    var pemBlocks []string
    for _, pubKey := range certs {
        pemBlock := &pem.Block{
            Type:  "PUBLIC KEY",
            Bytes: pubKey,
        }
        pemBlocks = append(pemBlocks, string(pem.EncodeToMemory(pemBlock)))
    }

    fmt.Print(strings.Join(pemBlocks, ""))
}