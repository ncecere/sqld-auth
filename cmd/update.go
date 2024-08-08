package cmd

import (
    "os"

    "github.com/spf13/cobra"
)

var updateCertCmd = &cobra.Command{
    Use:   "updatecert [name] [namespaces...]",
    Short: "Update certificate with new key pair and token",
    Args:  cobra.MinimumNArgs(2),
    Run:   updateCert,
}

func init() {
    rootCmd.AddCommand(updateCertCmd)
}

func updateCert(cmd *cobra.Command, args []string) {
    name := args[0]
    namespaces := args[1:]

    if err := generateAndSaveCert(cmd, name, namespaces); err != nil {
        cmd.PrintErrf("Error updating certificate: %v\n", err)
        os.Exit(1)
    }

    cmd.Printf("Certificate updated successfully for %s\n", name)
}