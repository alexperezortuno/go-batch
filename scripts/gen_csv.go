package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// generateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func main() {
	file, err := os.Create("large_users.csv")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// encabezado
	err = writer.Write([]string{"Username", "Password", "Email", "Name", "Age"})
	if err != nil {
		return
	}

	for i := 0; i < 1000000; i++ {
		username := fmt.Sprintf("user%d", i)
		password, err := generateRandomString(32)
		if err != nil {
			continue
		}
		name := fmt.Sprintf("User%d", i)
		email := fmt.Sprintf("user+%d@example.com", i)
		age := strconv.Itoa(20 + i%50)
		err = writer.Write([]string{username, password, email, name, age})
		if err != nil {
			return
		}
	}
}
