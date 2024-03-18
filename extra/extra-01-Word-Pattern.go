package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
)

func main() {

	var p string
	fmt.Scanln(&p)
	pattern := strings.Split(p, "")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	s :=scanner.Text()
	stringSlice := strings.Split(s, " ")

	mapEx   := make(map[string]string)

	for i := 0; i < len(pattern); i++{
		if _, exist := mapEx[pattern[i]]; !exist {
			mapEx[pattern[i]] = stringSlice[i]
		}
	}

	var answer string
	
	for k := 0; k < len(pattern); k++{
		answer += mapEx[pattern[k]] + " "
	}
	answer = strings.TrimSpace(answer)

	if answer == s{
		fmt.Println("sはpに従う")
	} else {
		fmt.Println("sはpに従わない")
	}
}
