package main

import (
	"fmt"
	"time"

	"github.com/robfig/cron"
)

func GoogleSheetCronJob() {
	fmt.Printf("GSCronJOB called")
}

func main() {
	c := cron.New()
	c.AddFunc("0 30 * * * *", func() { fmt.Println("Every hour on the half hour") })
	c.AddFunc("@hourly", func() { fmt.Println("Every hour") })
	c.AddFunc("@every 1h30m", func() { fmt.Println("Every hour thirty") })
	c.AddFunc("* * * * * *", GoogleSheetCronJob) //every second
	c.Start()

	// Funcs are invoked in their own goroutine, asynchronously.

	// Funcs may also be added to a running Cron
	c.AddFunc("@daily", func() { fmt.Println("Every day") })
	time.Sleep(10 * time.Second)
	// Inspect the cron job entries' next and previous run times.
	fmt.Println(c.Entries())
	c.Stop() // Stop the scheduler (does not stop any jobs already running).
}
