package cmd

import (
    "encoding/base64"
    "encoding/pem"
    "fmt"
    "os"
    "strings"

    "github.com/spf13/cobra"
    "github.com/ncecere/sqld-auth/internal/database"
)

var showAllBase64 bool

var showBase64Cmd = &cobra.Command{
    Use:   "showbase64 [name]",
    Short: "Show base64-encoded public certificate for a name",
    Args:  cobra.MaximumNArgs(1),
    Run:   showBase64,
}

func init() {
    rootCmd.AddCommand(showBase64Cmd)
    showBase64Cmd.Flags().BoolVarP(&showAllBase64, "all", "a", false, "Show all base64-encoded public certificates")
}

func showBase64(cmd *cobra.Command, args []string) {
    if showAllBase64 {
        showAllBase64Certs(cmd)
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
    pemData := pem.EncodeToMemory(pemBlock)
    base64Cert := base64.StdEncoding.EncodeToString(pemData)
    fmt.Println(base64Cert)
}

func showAllBase64Certs(cmd *cobra.Command) {
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

    allPemData := []byte(strings.Join(pemBlocks, ""))
    base64AllCerts := base64.StdEncoding.EncodeToString(allPemData)
    fmt.Println(base64AllCerts)
}