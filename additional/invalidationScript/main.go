package main

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"go.mongodb.org/mongo-driver/bson"
// )

// func main() {
// 	// Define the interval duration
// 	interval := 2 * time.Minute // Change this to your desired interval

// 	// Create a new ticker that ticks at the specified interval
// 	ticker := time.NewTicker(interval)
// 	defer ticker.Stop() // Ensure the ticker is stopped when the main function exits

// 	// Channel to catch OS signals
// 	sigChan := make(chan os.Signal, 1)
// 	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
// 	signal.Notify(sigChan, os.Interrupt)

// 	// Use a goroutine to run the periodic function
// 	go func() {
// 		for {
// 			select {
// 			case <-ticker.C:
// 				filter := bson.D{{Key: "valid", Value: false}}
// 				res, err := db_config.KEYS.DeleteMany(context.TODO(), filter)
// 				if err != nil {
// 					fmt.Println(err)
// 				} else {
// 					fmt.Println("Number of keys deleted: ", res.DeletedCount)
// 				}
// 			}
// 		}
// 	}()

// 	// Wait for an OS signal to exit
// 	<-sigChan
// }
