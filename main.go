package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func genData(n int) {
	rand.Seed(time.Now().UnixNano())
	var filename = "data" + strconv.Itoa(n) + ".txt"

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()

	var my_slice []int

	for i := 0; i < n; i++ {
		my_slice = append(my_slice, rand.Intn(n))
		fmt.Fprint(w, my_slice[i], " ")
	}
}

func readData(filename string) []int {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print(err)
	}
	str := strings.Split(string(b), " ")
	var result []int
	for _, item := range str {
		val, _ := strconv.Atoi(item)
		result = append(result, val)
	}
	return result
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sequential(my_slice []int) {
	var (
		val1       = my_slice[0]
		val2       = my_slice[1]
		tmpAbsDiff = abs(my_slice[0] - my_slice[1])
	)
	t0 := time.Now()
	for i := 0; i < len(my_slice)-1; i++ {
		currentDiff := my_slice[i] - my_slice[i+1]
		currentAbsDiff := abs(currentDiff)
		if currentAbsDiff > tmpAbsDiff {
			val1 = my_slice[i]
			val2 = my_slice[i+1]
			tmpAbsDiff = currentAbsDiff
		}
	}
	fmt.Printf("Elapsed time sequential: %v\n", time.Since(t0))
	fmt.Printf("Val 1: %v\n", val1)
	fmt.Printf("Val 2: %v\n", val2)
}

func splitSlice(my_slice []int) [][]int {
	size := len(my_slice) / processors
	var j int
	var result [][]int
	for i := 0; i < len(my_slice)-1; i += size {
		j += size
		if j > len(my_slice) {
			j = len(my_slice)
		}
		// do what do you want to with the sub-slice, here just printing the sub-slices
		result = append(result, my_slice[i:j])
	}
	return result
}

func parallelGoroutine(slice []int, t int, wg *sync.WaitGroup, result [][3]int) {
	defer wg.Done()
	var (
		val1       = slice[0]
		val2       = slice[1]
		tmpAbsDiff = abs(slice[0] - slice[1])
	)
	for i := 0; i < len(slice)-1; i++ {
		currentDiff := slice[i] - slice[i+1]
		currentAbsDiff := abs(currentDiff)
		if currentAbsDiff > tmpAbsDiff {
			val1 = slice[i]
			val2 = slice[i+1]
			tmpAbsDiff = currentAbsDiff
		}
	}
	result[t] = [3]int{val1, val2, tmpAbsDiff}
	//fmt.Print(t, [3]int{val1, val2, tmpAbsDiff}, slice, "\n")
}

func parallel(my_slice []int) {
	t0 := time.Now()
	splitted := splitSlice(my_slice)
	var wg sync.WaitGroup
	var result = [][3]int{}
	for i := 0; i < processors; i++ {
		result = append(result, [3]int{})
	}
	for i, part := range splitted {
		wg.Add(1)
		go func(partD []int, iD int, wgD *sync.WaitGroup, resultD [][3]int) {
			parallelGoroutine(partD, iD, wgD, resultD)
		}(part, i, &wg, result)
	}
	wg.Wait()
	var (
		val1 = result[0][0]
		val2 = result[0][1]
		min  = result[0][2]
	)
	for i, _ := range result {
		if result[i][2] > min {
			min = result[i][2]
			val1 = result[i][0]
			val2 = result[i][1]
		}
	}
	for i := 0; i < len(splitted)-2; i++ {
		currentDiff := splitted[i][len(splitted[i])-1] - splitted[i+1][0]
		currentAbsDiff := abs(currentDiff)
		if currentAbsDiff > min {
			min = result[i][2]
			val1 = result[i][0]
			val2 = result[i][1]
		}
	}
	fmt.Printf("Elapsed time parallel by goroutines: %v\n", time.Since(t0))
	fmt.Printf("Val 1: %v\n", val1)
	fmt.Printf("Val 2: %v\n", val2)
}

var processors = 16

func main() {
	// genData(100000000)
	runtime.GOMAXPROCS(processors)
	var my_slice = readData("data1000000000.txt")
	parallel(my_slice)   // n + 6Tn / p + 4Tp + 5Tn
	sequential(my_slice) // 6Tn

	//var p = 1
	//fmt.Print(p)

}
