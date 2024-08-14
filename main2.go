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

func worker(inputChan chan string, sliceLen int, doneChan chan bool, channelList []chan string, id int) {

  channelsListCopy := make([]chan string, len(channelList))
  copy(channelsListCopy, channelList)

	slice := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		slice[i] = <-inputChan
	}

	shuffle(slice)

	for _, word := range slice {
    insertWord(channelsListCopy, word)
	}
	doneChan <- true
}

func insertWord(channelList []chan string, word string){
  index := rand.Intn(len(channelList))
  inserted := false
  
  for !inserted{
    select{
    case channelList[index] <- word:
      inserted = true
    default:
      removeElement(channelList, index)
      index = rand.Intn(len(channelList))
    }
  }
}

func removeElement(slice []chan string, index int)([]chan string, error){
  if index < 0 || index >= len(slice){
    return nil, fmt.Errorf("Index out of range")
  }
  
  return append(slice[:index], slice[index+1:]...),nil
}

func textShuffle(wordsList []string, workersNumber int) []string{

	sortedChan := make(chan string)
	doneChan := make(chan bool, workersNumber)
  slices := slices(len(wordsList), workersNumber)

  workersChannels := make([]chan string, len(slices)) //list of channels that workers will fill
  
  for index, slice := range slices {
		workersChannels[index] = make(chan string, slice)
    go worker(sortedChan, slice, doneChan, workersChannels, index)
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
	}

  var result []string
  
  for _, channel := range(workersChannels) {
    close(channel)
    for word := range(channel){
      result = append(result, word)
    }
  }

  return result
}

func appendToFile(filePath string, lines []string) error {
	// Open the file in append mode, create if it does not exist, and set permissions.
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Create a buffered writer for efficient writing.
	writer := bufio.NewWriter(file)

	// Write each line to the file with a newline.
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}
	
  if _, err := writer.WriteString("-------------" + "\n"); err != nil {
			return fmt.Errorf("error writing to file: %w", err)
	}

	// Flush buffered writer to ensure all data is written.
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}

	return nil
}

func main(){
	rand.Seed(time.Now().UnixNano())
  
  wordsList, err := getWords("./text.txt")
	check(err)

  workersNumber1 := rand.Intn(len(wordsList))
  workersNumber2 := rand.Intn(len(wordsList))

  fmt.Println("Specs:")
  fmt.Println("First pass ", workersNumber1, " workers")
  fmt.Println("Second pass ", workersNumber2, " workers")
  fmt.Println("On ", len(wordsList), " words")
  
  shuffledText := textShuffle(textShuffle(wordsList,workersNumber1), workersNumber2)
  fmt.Println(shuffledText)
  appendToFile("./output.txt", shuffledText)
}





























