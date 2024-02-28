package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sort"
	"strconv"
	
)

func main(){
	var n int
	fmt.Scanf("%d", &n)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	nums := scanner.Text()
	numsArray := strings.Split(nums, ",")

	sort.Slice(numsArray, func(i, j int) bool {
		num1, _ := strconv.Atoi(numsArray[i])
		num2, _ := strconv.Atoi(numsArray[j])
        return num1 < num2
    })
	
	mapEx := make(map[int]bool)
	for _, num := range numsArray {
		i, _ := strconv.Atoi(num)
		mapEx[i] = true
	}

	var answer []int
	for i := 1; i <= n; i++ {
		if !mapEx[i]{
			answer = append(answer, i)
		}
	}
	fmt.Println(answer)
}