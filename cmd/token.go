package cmd

import (
    "os"
    "time"

    "github.com/spf13/cobra"
    "github.com/ncecere/sqld-auth/internal/database"
    "github.com/ncecere/sqld-auth/internal/token"
)

var (
    accessLevel string
    expiration  time.Duration
)

var tokenCmd = &cobra.Command{
    Use:   "token [name]",
    Short: "Generate JWT for a specific certificate",
    Args:  cobra.ExactArgs(1),
    Run:   generateToken,
}

func init() {
    rootCmd.AddCommand(tokenCmd)
    tokenCmd.Flags().StringVarP(&accessLevel, "access", "a", "full", "Access level: full, read, write, or attach-read")
    tokenCmd.Flags().DurationVarP(&expiration, "expiration", "e", 24*time.Hour, "Token expiration time (e.g., 24h, 30m, 1h30m)")
}

func generateToken(cmd *cobra.Command, args []string) {
    name := args[0]

    _, privKey, err := database.GetKeys(name)
    if err != nil {
        cmd.PrintErrf("Error retrieving keys: %v\n", err)
        os.Exit(1)
    }

    tokenStr, err := token.CreateTokenWithNamespaces(privKey, accessLevel, expiration)
    if err != nil {
        cmd.PrintErrf("Error creating token: %v\n", err)
        os.Exit(1)
    }

    cmd.Printf("Token (%s access, expires in %v): %s\n", accessLevel, expiration, tokenStr)
}