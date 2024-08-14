package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getWords(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	check(err)
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

func shuffle(words []string) {
	for i := len(words) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		words[i], words[j] = words[j], words[i]
	}
}

func slices(listLen int, noWorkers int) []int {

	if noWorkers <= 0 {
		panic("Values error in slices function")
	}

	if listLen < noWorkers {
		panic("less words than workers")
	}

	sliceSize := listLen / noWorkers
	extra := listLen % noWorkers

	slicesList := make([]int, noWorkers)

	for i := 0; i < noWorkers; i++ {
		end := sliceSize
		if extra > 0 {
			end++
			extra--
		}
		slicesList[i] = end
	}

	return slicesList
}

func worker(outputChan chan string, inputChan chan string, sliceLen int, doneChan chan bool) {
	fmt.Println("Worker started")
	slice := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		slice[i] = <-inputChan
	}

	fmt.Println("slice: ", slice)
	//time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
	shuffle(slice)
	for _, word := range slice {
		outputChan <- word
	}
	doneChan <- true
}

func main(){
	rand.Seed(time.Now().UnixNano())

	wordsList, err := getWords("./text.txt")
	check(err)

	workersNumber := 3

	sortedChan := make(chan string)
	doneChan := make(chan bool, workersNumber)
  slices := slices(len(wordsList), workersNumber)
 	shuffledChan := make(chan string, len(wordsList)) // Buffer the shuffledChan to prevent blocking
 
  for _, slice := range slices {
		go worker(shuffledChan, sortedChan, slice, doneChan)
	}
  
  go func() {
		for _, word := range wordsList {
			sortedChan <- word
		}
		close(sortedChan)
	}()

	// Wait for all workers to finish
	for i := 0; i < workersNumber; i++ {
    <-doneChan
		fmt.Println("Done")
	}
  close(shuffledChan)

  for word := range(shuffledChan){
    fmt.Println(word)
  }

}





























