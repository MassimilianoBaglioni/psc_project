
package main
import "testing"

func TestRandomWorkerNumber(t *testing.T){
  got := randomWorkerNumber(100, 10)
  if got > 10 || got <= 0{
    t.Error()
  }

  got = randomWorkerNumber(10, 100)
  if got != 10{
    t.Error()
  }

  got = randomWorkerNumber(10, -11)
  if got != 10{
    t.Error()
  }

  got = randomWorkerNumber(10, -9)
  if got > 10 || got <= 0{
    t.Error()
  }

  got = randomWorkerNumber(0, 0)
  if got != 1 {
    t.Error()
  }

  got = randomWorkerNumber(10, 0)
  if got != 1 {
    t.Error()
  }
} 
