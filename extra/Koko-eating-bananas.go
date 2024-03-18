
package main
import (
	"fmt"
	"math"
	"strconv"
	"bufio"
	"os"
	"strings"
)
func cal(piles []int, h int, k int) bool {
    answer := 0
    for i := 0; i < len(piles); i++{
        answer += int(math.Ceil(float64(piles[i]) / float64(k)))
        if answer > h {
            return false
        } 
    }
	return true
}
func max(x, y int) int {
    if x > y {
        return x
    }
    return y
}
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
func minEatingSpeed(piles []int, h int) int {
	pilesMax := 0
	for _, pile := range piles {
		pilesMax = max(pilesMax, pile)
	}

    r := pilesMax
    l := 1
    answer := r
    for l < r {
        k := (l + r) / 2
        if cal(piles, h, k){
            answer = min(answer, k)
            r = k - 1
        } else {
            l = k + 1
        }
    }
	return answer
}
func main(){
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	p := scanner.Text()
	ppp := strings.Split(p, ",")
	
	var piles []int
	for _, pile := range ppp {
		int, _ := strconv.Atoi(pile)
		piles = append(piles, int)
	}

	var h int
	fmt.Scanf("%d", &h)

	fmt.Println(minEatingSpeed(piles, h))

}