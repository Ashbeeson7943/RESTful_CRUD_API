To run a function every X minutes in Go, you can use a combination of goroutines and the `time` package, particularly the `time.Ticker` type, which ticks at regular intervals. Here's an example of how you can achieve this:

```go
package main

import (
	"fmt"
	"time"
)

// Function to be run periodically
func periodicFunction() {
	fmt.Println("Function is running at", time.Now())
}

func main() {
	// Define the interval duration
	interval := 2 * time.Minute // Change this to your desired interval

	// Create a new ticker that ticks at the specified interval
	ticker := time.NewTicker(interval)
	defer ticker.Stop() // Ensure the ticker is stopped when the main function exits

	// Use a goroutine to run the periodic function
	go func() {
		for {
			select {
			case <-ticker.C:
				periodicFunction()
			}
		}
	}()

	// Keep the main function running indefinitely
	select {}
}
```

### Explanation:

1. **Importing Packages**:
   - `fmt`: For printing messages.
   - `time`: For handling time-related functions.

2. **Defining the Periodic Function**:
   - `periodicFunction()`: A simple function that prints the current time. Replace this with the function you want to run periodically.

3. **Main Function**:
   - **Interval Definition**: `interval := 2 * time.Minute` sets the interval to 2 minutes. Adjust this to the desired interval.
   - **Creating a Ticker**: `ticker := time.NewTicker(interval)` creates a new ticker that ticks at the specified interval.
   - **Defer Ticker Stop**: `defer ticker.Stop()` ensures the ticker is stopped when the main function exits to prevent resource leakage.
   - **Goroutine**: The goroutine runs an infinite loop that waits for a tick from the ticker. When it receives a tick (`case <-ticker.C`), it calls `periodicFunction()`.

4. **Keeping the Main Function Running**:
   - `select {}` is used to block the main function indefinitely, ensuring that the program keeps running and the ticker keeps ticking. 

### Important Considerations:

- **Error Handling**: Add error handling inside the `periodicFunction` if needed.
- **Program Exit**: Ensure that your program handles graceful shutdowns if necessary, such as cleaning up resources or stopping the ticker explicitly.
- **Concurrency**: The example runs `periodicFunction` in a goroutine. If your function performs I/O operations or is CPU-intensive, ensure it is thread-safe and consider using synchronization primitives if needed.

### Example with Graceful Shutdown

To handle graceful shutdowns, you can use channels to catch signals like `os.Interrupt`:

```go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Function to be run periodically
func periodicFunction() {
	fmt.Println("Function is running at", time.Now())
}

func main() {
	// Define the interval duration
	interval := 2 * time.Minute // Change this to your desired interval

	// Create a new ticker that ticks at the specified interval
	ticker := time.NewTicker(interval)
	defer ticker.Stop() // Ensure the ticker is stopped when the main function exits

	// Channel to catch OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Use a goroutine to run the periodic function
	go func() {
		for {
			select {
			case <-ticker.C:
				periodicFunction()
			}
		}
	}()

	// Wait for an OS signal to exit
	<-sigChan
	fmt.Println("Shutting down gracefully...")
}
```

### Explanation of Graceful Shutdown:

1. **Signal Handling**:
   - `os/signal`: Package to handle incoming signals.
   - `signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)`: Notify the `sigChan` channel on receiving `os.Interrupt` (Ctrl+C) or `SIGTERM`.

2. **Waiting for Signal**:
   - `<-sigChan`: Blocks until a signal is received, then the program prints "Shutting down gracefully..." and exits.

This ensures that your program can be terminated gracefully, allowing you to clean up resources if necessary.