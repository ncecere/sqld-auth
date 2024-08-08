package cmd

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/spf13/cobra"
    "github.com/ncecere/sqld-auth/internal/database"
)

var (
    dbPassword string
    dbFile     string
)

var rootCmd = &cobra.Command{
    Use:   "sqld-auth",
    Short: "A tool to manage authentication for sqld",
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    cobra.OnInitialize(initConfig)

    rootCmd.PersistentFlags().StringVar(&dbPassword, "dbpass", "", "Database encryption password")
    rootCmd.PersistentFlags().StringVar(&dbFile, "dbfile", "sqld_auth.db", "Path to the SQLite database file")

    // Add commands
    rootCmd.AddCommand(createCertCmd)
    rootCmd.AddCommand(updateCertCmd)
    rootCmd.AddCommand(tokenCmd)
    rootCmd.AddCommand(showCertCmd)
    rootCmd.AddCommand(showBase64Cmd)  // Add this line for the new showbase64 command
}

func initConfig() {
    if dbFile == "sqld_auth.db" {
        if envFile := os.Getenv("SQLD_AUTH_DBFILE"); envFile != "" {
            dbFile = envFile
        }
    }

    if dbPassword == "" {
        dbPassword = os.Getenv("SQLD_AUTH_DBPASS")
    }

    if dbPassword == "" {
        homeDir, err := os.UserHomeDir()
        if err == nil {
            configFile := filepath.Join(homeDir, ".sqld-auth-token")
            content, err := os.ReadFile(configFile)
            if err == nil {
                dbPassword = strings.TrimSpace(string(content))
            }
        }
    }

    if dbPassword == "" {
        fmt.Println("Error: Database password not provided. Use --dbpass flag, SQLD_AUTH_DBPASS environment variable, or ~/.sqld-auth-token file.")
        os.Exit(1)
    }

    database.SetDBInfo(dbFile, dbPassword)
}