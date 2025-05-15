package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

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
	err = writer.Write([]string{"Name", "Email", "Age"})
	if err != nil {
		return
	}

	for i := 0; i < 1000000; i++ {
		name := fmt.Sprintf("User%d", i)
		email := fmt.Sprintf("user+%d@example.com", i)
		age := strconv.Itoa(20 + i%50)
		err := writer.Write([]string{name, email, age})
		if err != nil {
			return
		}
	}
}
