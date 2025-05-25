package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadTPSFile(fileName string) []int {

	transaction := make([]int, 0)

	// Path from where you are executing the command
	file, err := os.Open(fmt.Sprintf("config/%v", fileName))

	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		tpsDetails := strings.Split(line, " ")
		tps, err := strconv.Atoi(tpsDetails[1])
		if err != nil {
			continue
		}
		transaction = append(transaction, tps)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return transaction

}
