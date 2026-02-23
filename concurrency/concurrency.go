package main

import (
	"fmt"
	"sync"
	"time"
)

// === Goroutines ===

// Lightweight threads managed by the Go runtime

func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func goroutineExample() {
	// starts a new goroutine that executes the function,
	// and the evaluation of that function happens in the current goroutine (main goroutine)
	go say("world")
	say("hello")

	// In above example, both functions called at around same time in non blocking order
	// the above two functions execute on different goroutines simultaneously
	// the goroutines share the same address space
	// the order of thread execution may be different depending on os thread schedular
	// 		(ex. "hello" then "world" for one program exec, but "world then "hello" on next)
	// bc of this it is not thread safe to modify the same data asynchronously using different goroutines
	// (without using a syncronization method to access shared memory like channels or Go's sync package)
}

// === Channels ===

// A typed conduit through which you can send / recieve values
// Used to safely transfer data and sync executation between conc. goroutines
// Use the channel operator, <- (data flows in the direction of the arrow)
// Channels are bidirectional by default (can set as unidirectional for send only or receive only)
func channelExample() {

	// Channels must be created before use, like maps and slices
	ch := make(chan int)

	s := []int{7, 2, -5, 4, 1}

	go sum(s[:len(s)/2], ch)
	go sum(s[len(s)/2:], ch)

	// sends and receives BLOCK until the other side is ready
	halfSum1 := <-ch
	halfSum2 := <-ch
	// Note - there is no guarantee which goroutine sum goes into which variable
	// (same reason as mentioned in goroutineExample above)

	sum := halfSum1 + halfSum2

	fmt.Printf("Value from channel, blocking: %d\n", sum)
}

func sum(s []int, ch chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	ch <- sum // send sum to the channel
}

// --- Buffered channels ---

func bufferedChannelsExample() {
	// Channels can be buffered to store a certain number of elements
	// Ex. creating a buffered channel that can store 10 ints
	// 		Provide a buffer length as the second arg
	ch := make(chan int, 10)

	// -- Buffered channels store elements in a FIFO order using a queue --

	for i := range cap(ch) {
		fmt.Printf("Sending element %d to channel. Value is %d\n", i+1, i)
		ch <- i
	}

	fmt.Println()

	for i := range cap(ch) {
		// value is removed from the channel when recieving
		fmt.Printf("Recieving element %d from channel. Value is: %d\n", i+1, <-ch)
	}
}

// In the below two examples a deadlock occurs
// bc the goroutine sleeps after it finishes executing synchronous code
// and then waits forever for the channel to be ready

// Sends to a buffered channel block when the buffer is full
// (Also sends to an unbuffered channel blocks IF there is no reciever directly after)
func deadlockExampleOverfilledBufferBlock() {
	ch := make(chan int, 2)
	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Println(<-ch)
}

// Recieves to a buffered channel block when the buffer is empty
func deadlockExampleEmptyBufferBlock() {
	ch := make(chan int, 1)
	fmt.Println(<-ch)
}

// --- Unbuffered channels ---

// See notes in channelExample()

// Unbuffered channels require immediate synchronization, unlike buffered
// Bc the capacity of an unbuffered channel is 0 (has no storage queue)
// The sender (one goroutine) waits / blocks until a reciever (another goroutine) is ready

// If no reciever is present then a deadlock occurs
func deadlockExampleUnbufferedChNoReciever() {
	ch := make(chan int)

	// Blocks here bc no reciever goroutine present
	ch <- 1

	fmt.Println(<-ch)
}

// To fix this, a separate goroutine needs to be started
// So at least one goroutine is active
// (main goroutine sleeps while waiting for other goroutine to finish sending)
func noDeadlockUnbufferedChannel() {
	ch := make(chan int)

	// start a separate goroutine that handles sending
	go func() {
		ch <- 1
	}()

	// current goroutine (main) sleeps until the sender goroutine finishes
	fmt.Println(<-ch)
}

// --- Channel Summary and Use cases: ---

// Both: for synchronized communication between goroutines

// Unbuffered channel
// - For synchronous communcation
// - Syncs sender with reciever, as both must be ready at same time for data transfer to occur.
// - Reciever waits until sender is ready, and vice versa

// Buffered channel
// - For semi-synchronous communication - asynchronous but blocks if producer empty or if consumer full
// - Allows a sender to continously send messages with set "rate" (based on num elements)
// 	 without waiting for a receiver, reducing blocking

// === Range and Close ===

// A sender can close a channel to indicate that no more values will be sent

// Recievers can test where a channel has been closed,
// by assigning a second paramter to the recieve expression
func addToChannelThenClose(ch chan int) {
	ch <- 1
	close(ch)
}

func testClosedChannelEx() {
	ch := make(chan int, 1)
	addToChannelThenClose(ch)

	v, ok := <-ch
	fmt.Println(v, ok)

	// if there are no more values to receive and the channel is closed:
	// ok is false, v is it's zero value
	v, ok = <-ch
	fmt.Println(v, ok)
}

// Only the sender should close a channel, never the receiver.
// Sending on a closed channel will cause a panic.
func sendingOnClosedChannelPanicEx() {
	ch := make(chan int, 1)
	close(ch)
	ch <- 5
	val := <-ch
	fmt.Println(val)
}

// Don't usually need to close channels.
// Closing is only necessary when the receiver must be told there are no more values coming,
// such as to terminate a range loop.
const MAX_CH_BUFFER_CAP = 10

func fillChannelToIndex(ch chan int, endAtIndex int) {
	defer close(ch)
	endAtIndex = min(endAtIndex, MAX_CH_BUFFER_CAP)
	for i := range endAtIndex {
		ch <- i
	}
}

// Using range with channel receives values from it repeatedly until the sender closes it
func loopThroughValsUntilChannelClosedEx() {
	ch := make(chan int, MAX_CH_BUFFER_CAP)
	endAtIndex := 6

	go fillChannelToIndex(ch, endAtIndex)
	for i := range ch {
		fmt.Println(i)
	}
}

// === Select ===

// The select statement lets a goroutine wait on multiple communication operations.
// A select blocks until one of its cases can run, then it executes that case.
// It chooses one at random if multiple are ready.
func fibonacci(c, quit chan int) {
	x, y := 0, 1
	defaultExecuted := false
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return

		// The default case in a select is run if no other case is ready
		default:
			if !defaultExecuted {
				fmt.Println("no channel ready. Executing default case")
				defaultExecuted = true
			}
		}
	}
}

func selectEx() {
	c := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		time.Sleep(1 * time.Nanosecond)
		quit <- 0
	}()
	fibonacci(c, quit)
}

// === sync.Mutex ===

// Go provides mutual exclusion for when communication among goroutines are not needed
// std library provides the libraries sync.Mutex, with Lock and Unlock methods
type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

func Inc(counterValue *int) {
	(*counterValue)++
}

func (counter *SafeCounter) SafeInc(key string) {
	// Lock so only one goroutine at a time can access the critical section map v
	counter.mu.Lock()
	val := counter.v[key]
	Inc(&val)
	counter.v[key] = val
	counter.mu.Unlock()
}

func (counter *SafeCounter) GetValue(key string) int {
	counter.mu.Lock()
	// can unlock once function finishes executing using defer
	defer counter.mu.Unlock()
	return counter.v[key]
}

const COUNTER_KEY = "key1"

func safeIncrementMutuxExample() {
	counter := SafeCounter{v: make(map[string]int)}
	for range 100 {
		go counter.SafeInc(COUNTER_KEY)
	}
	time.Sleep(time.Second)

	// is 100 every time
	fmt.Printf("Counter value after safe increment using mutex: %v", counter.GetValue(COUNTER_KEY))
}

func unsafeIncrementExample() {
	val := 0
	for range 100 {
		go Inc(&val)
	}
	time.Sleep(time.Second)

	// some executions is 99, others it is 100, depending on thread scheduling
	fmt.Printf("Counter value after unsafe increment with race conditions: %v", val)
}

func main() {
	// goroutineExample()
	// channelExample()
	// bufferedChannelsExample()
	// deadlockExampleOverfilledBufferBlock()
	// deadlockExampleEmptyBufferBlock()
	// deadlockExampleUnbufferedChNoReciever()
	// noDeadlockUnbufferedChannel()
	// testClosedChannelEx()
	// sendingOnClosedChannelPanicEx()
	// loopThroughValsUntilChannelClosedEx()
	// selectEx()
	// safeIncrementMutuxExample()
	// unsafeIncrementExample()
}
