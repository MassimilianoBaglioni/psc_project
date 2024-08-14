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

func worker(inputChan chan string, sliceLen int, doneChan chan bool, channelsList []chan string) {

	slice := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		slice[i] = <-inputChan
	}

	shuffle(slice)

	for _, word := range slice {
		insertWord(channelsList, word)
	}
	doneChan <- true
}

func insertWord(channelList []chan string, word string) {
	index := rand.Intn(len(channelList))
	channelList[index] <- word
}

func removeElement(slice []chan string, index int) ([]chan string, error) {
	if index < 0 || index >= len(slice) {
		return nil, fmt.Errorf("Index out of range")
	}

	return append(slice[:index], slice[index+1:]...), nil
}

func textShuffle(wordsList []string, workersNumber int) []string {

	//Channel that will store initial text to shuffle.
	initialChan := make(chan string)

	//Channel used to wait all workers to stop.
	doneChan := make(chan bool, workersNumber)

	//Each worker will have a portion of the original list of words.
	slices := slices(len(wordsList), workersNumber)

	//List of channels that will be filled by workers.
	workersChannels := make([]chan string, len(slices))

	for index, slice := range slices {
		workersChannels[index] = make(chan string, len(wordsList))
		go worker(initialChan, slice, doneChan, workersChannels)
	}

	go func() {
		for _, word := range wordsList {
			initialChan <- word
		}
		close(initialChan)
	}()

	// Wait for all workers to finish
	for i := 0; i < workersNumber; i++ {
		<-doneChan
	}

	var result []string

	for _, channel := range workersChannels {
		close(channel)
		for word := range channel {
			result = append(result, word)
		}
	}

	return result
}

// Appends the passed list of strings.
func appendToFile(filePath string, lines []string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}

	if _, err := writer.WriteString("-------------" + "\n"); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}

	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	wordsList, err := getWords("./text.txt")
	check(err)

	workersNumber1 := rand.Intn(len(wordsList)-1) + 1
	workersNumber2 := rand.Intn(len(wordsList)-1) + 1

	fmt.Println("Specs:")
	fmt.Println("First pass ", workersNumber1, " workers")
	fmt.Println("Second pass ", workersNumber2, " workers")
	fmt.Println("On ", len(wordsList), " words")

	shuffledText := textShuffle(textShuffle(wordsList, workersNumber1), workersNumber2)
	fmt.Println(shuffledText)
	appendToFile("./output.txt", shuffledText)
}
