
# sqld-auth

`sqld-auth` is a command-line tool for managing authentication and authorization for sqld. It provides functionality to create and manage Ed25519 key pairs, and generate JSON Web Tokens (JWTs) with various access levels.

## Table of Contents

- [sqld-auth](#sqld-auth)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
  - [Configuration](#configuration)
    - [Database Password](#database-password)
    - [Database File](#database-file)
  - [Usage](#usage)
    - [Creating a Certificate](#creating-a-certificate)
    - [Updating a Certificate](#updating-a-certificate)
    - [Generating Tokens](#generating-tokens)
  - [Command-Line Arguments](#command-line-arguments)
  - [Environment Variables](#environment-variables)
  - [Examples](#examples)

## Installation

To install `sqld-auth`, you need to have Go installed on your system. Then, you can build the tool using:

```bash
go build -o sqld-auth
```

This will create an executable named `sqld-auth` in your current directory.

## Configuration

`sqld-auth` uses an SQLite database to store encrypted private keys. You can configure the tool using a combination of command-line flags, environment variables, and a configuration file.

### Database Password

The database password can be set in three ways (in order of precedence):

1. Command-line flag: `--dbpass`
2. Environment variable: `SQLD_AUTH_DBPASS`
3. Configuration file: `~/.sqld-auth-token`

To use a configuration file, create a file named `.sqld-auth-token` in your home directory and add the password as plain text.

### Database File

The path to the SQLite database file can be set in two ways (in order of precedence):

1. Command-line flag: `--dbfile`
2. Environment variable: `SQLD_AUTH_DBFILE`

If not specified, it defaults to `sqld_auth.db` in the current directory.

## Usage

`sqld-auth` provides three main commands:

- `createcert`: Create a new certificate (key pair)
- `updatecert`: Update an existing certificate
- `token`: Generate a JWT for a specific certificate

### Creating a Certificate

To create a new certificate:

```bash
sqld-auth createcert [name] [namespaces...]
```

- `[name]`: A unique identifier for the certificate
- `[namespaces...]`: One or more namespaces to associate with the certificate

### Updating a Certificate

To update an existing certificate:

```bash
sqld-auth updatecert [name] [namespaces...]
```

- `[name]`: The identifier of the existing certificate
- `[namespaces...]`: One or more namespaces to associate with the updated certificate

### Generating Tokens

To generate a JWT for a specific certificate:

```bash
sqld-auth token [name] --access [access_level]
```

- `[name]`: The identifier of the existing certificate
- `--access`: (Optional) The access level for the token. Can be "full", "read", "write", or "attach-read". Defaults to "full".

## Command-Line Arguments

Global flags (applicable to all commands):

- `--dbpass`: Database encryption password
- `--dbfile`: Path to the SQLite database file

Flags specific to the `token` command:

- `--access` or `-a`: Access level for the token (full, read, write, or attach-read)

## Environment Variables

- `SQLD_AUTH_DBPASS`: Database encryption password
- `SQLD_AUTH_DBFILE`: Path to the SQLite database file

## Examples

1. Create a new certificate:

   ```bash
   sqld-auth createcert myapp namespace1 namespace2 --dbpass mysecretpassword
   ```

2. Update an existing certificate:

   ```bash
   sqld-auth updatecert myapp namespace1 namespace2 namespace3 --dbpass mysecretpassword
   ```

3. Generate a full access token:

   ```bash
   sqld-auth token myapp --dbpass mysecretpassword
   ```

4. Generate a read-only token:

   ```bash
   sqld-auth token myapp --access read --dbpass mysecretpassword
   ```

5. Generate a write access token:

   ```bash
   sqld-auth token myapp --access write --dbpass mysecretpassword
   ```

6. Generate an attach-read token:

   ```bash
   sqld-auth token myapp --access attach-read --dbpass mysecretpassword
   ```

7. Use a custom database file:

   ```bash
   sqld-auth createcert myapp namespace1 --dbfile /path/to/custom.db --dbpass mysecretpassword
   ```

8. Use environment variables:

   ```bash
   export SQLD_AUTH_DBPASS=mysecretpassword
   export SQLD_AUTH_DBFILE=/path/to/custom.db
   sqld-auth token myapp
   ```

9. Use a configuration file for the database password:

   ```bash
   echo "mysecretpassword" > ~/.sqld-auth-token
   sqld-auth createcert myapp namespace1
   ```

Remember to keep your database file and password secure, as they contain sensitive key material.