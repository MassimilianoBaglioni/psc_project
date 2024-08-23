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

/*
*Input: Takes file path of the input text.
*output: List of the words in the file.
*/
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

/*
* The function returns the number of words for that each worker should work on.
* It simply divides the number of total words and distributed the remainder equally.
*
* Input: Lenght of the list and the number of workers for the list.
*
* Output: A list, with the number of words for each worker.
*
*/
func createSegments(listLen int, noWorkers int) []int {

	if noWorkers <= 0 {
		panic("Values error in slices function")
	}

	if listLen < noWorkers {
		panic("less words than workers")
	}

	segmentsSize := listLen / noWorkers

  //Extra remainder words.
	extra := listLen % noWorkers

	segmentsList := make([]int, noWorkers)

	for i := 0; i < noWorkers; i++ {
		end := segmentsSize

    //Distriburtes the remainder to the first slices.
		if extra > 0 {
			end++
			extra--
		}
		segmentsList[i] = end
	}

	return segmentsList
}

func worker(inputChan chan string, segmentLen int, doneChan chan bool, channelsList []chan string) {

	segment := make([]string, segmentLen)

  //Collect the words in the shared unbuffered channel.
	for i := 0; i < segmentLen; i++ {
		segment[i] = <-inputChan
	}

	shuffle(segment)

  //Write in a random channel.
	for _, word := range segment {
		insertWord(channelsList, word)
	}

  //Communicte that the routine is done.
	doneChan <- true
}

/*
* Input: A list of channels and a word to insert.
*
* Inserts one word into one of the channels, choosen at random.
*
*/
func insertWord(channelList []chan string, word string) {
	index := rand.Intn(len(channelList))
	channelList[index] <- word
}

/*
* Input: Slice with the text to shuffle and the number of workers.
*
* Output: Slice with the shuffled text.
*
* The function starts the routines and collects the result.
*/
func textShuffle(wordsList []string, workersNumber int) []string {

	initialChan := make(chan string) // Initial text to shuffle.
	doneChan := make(chan bool, workersNumber) // Channel used to wait routines.
  
  segments := createSegments(len(wordsList), workersNumber) //List of segments for each worker.
  
  //List of channels that will be filled by workers.
	workersChannels := make([]chan string, len(segments))

  //Start workers and create channels.
	for index, segment := range segments {
		workersChannels[index] = make(chan string, len(wordsList))
		go worker(initialChan, segment, doneChan, workersChannels)
	}

  //Distribute words between workers concurrently on an unbuffered channel.
	go func() {
		for _, word := range wordsList {
			initialChan <- word
		}
		close(initialChan)
	}()

	//Wait for all workers to finish.
	for i := 0; i < workersNumber; i++ {
		<-doneChan
	}

	var result []string
  
  //Fill the returned list.
	for _, channel := range workersChannels {
		close(channel)
		for word := range channel {
			result = append(result, word)
		}
	}

	return result
}

/*
* Input: File path where the shuffled sequence will be written.
* 
* Utility function used just to debug or check results of the program over long runs.
*/
func appendToFile(filePath string, lines []string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, line := range lines {
		if _, err := writer.WriteString(line + " "); err != nil {
			return fmt.Errorf("error writing to file: %w", err)
		}
	}

	if _, err := writer.WriteString("\n" + "-------------" + "\n"); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}

	return nil
}

/*
* 
* Input: number of words in the input text to shuffle and the max number of workers accepted.
* 
* The function returns a random number (used as number of worker routines to spawn) that is never
* bigger than the number of words and never higher than the given max number of workers decided by
* the second parameter.
*/
func randomWorkerNumber(words int, maxWorkers int) int{
  if maxWorkers < 0{
    maxWorkers = -maxWorkers
  }

  if maxWorkers == 0 {
    return 1
  }

  if words < maxWorkers{
    return words
  }

  return rand.Intn(maxWorkers) + 1 
}

func main() {
  start := time.Now()
  rand.Seed(time.Now().UnixNano())

	wordsList, err := getWords("./text.txt")
  
  if len(wordsList) == 0{
    fmt.Println("No words to shuffle")
    return
  }
	
  check(err)

	workersNumber1 := randomWorkerNumber(len(wordsList), 10)
	workersNumber2 := randomWorkerNumber(len(wordsList), 10)

  //Comment the random generated to hardcode the number of routines.
  //workersNumber1 := 3
	//workersNumber2 := 3

	fmt.Println("Specs:")
	fmt.Println("First pass ", workersNumber1, " workers")
	fmt.Println("Second pass ", workersNumber2, " workers")
	fmt.Println("On ", len(wordsList), " words")

	shuffledText := textShuffle(textShuffle(wordsList, workersNumber1), workersNumber2)
  fmt.Println(time.Since(start))
	//fmt.Println(shuffledText)
	appendToFile("./output.txt", shuffledText)
}
