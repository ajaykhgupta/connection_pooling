package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
	_ "github.com/lib/pq" // PostgreSQL driver
)

const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "user"
	dbPassword = "password"
	dbName     = "dbname"
	poolSize = 3
)

func performDBTaskNoConnectionPool(wg *sync.WaitGroup, id int) {
	fmt.Println("hello from worker id", id)
	defer wg.Done()

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("[Worker %d] Error opening DB: %v", id, err)
		return
	}
	defer db.Close()

	_, err = db.Exec("SELECT pg_sleep(1);")
	if err != nil {
		log.Printf("[Worker %d] Query error: %v", id, err)
		return
	}
}

func performDBTaskwConnectionPool(wg *sync.WaitGroup, id int, queue chan *sql.DB){
	fmt.Println("hello from worker id", id)

	defer wg.Done()
	db := <-queue
	fmt.Printf("[Worker %d] Got DB connection\n", id)

	// Perform the query
	_, err := db.Exec("SELECT pg_sleep(1);")
	if err != nil {
		log.Printf("[Worker %d] Query error: %v", id, err)
	} else {
		fmt.Printf("[Worker %d] Query completed\n", id)
	}

	// Put the DB connection back into the queue
	queue <- db

}

func getConnectionPoolQueue() chan *sql.DB {
	queue := make(chan *sql.DB, poolSize)

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	for i := 0; i < poolSize; i++ {
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatalf("Failed to open DB connection %d: %v", i+1, err)
		}
		queue <- db
	}

	return queue
}


func main() {
	// without connection pool main
	/*
	start := time.Now()

	
	var wg sync.WaitGroup
	numRequests := 150

	wg.Add(numRequests)

	for i := 1; i <= numRequests; i++ {
		go performDBTaskNoConnectionPool(&wg, i)
	}

	wg.Wait()
	endTime := time.Since(start).Seconds()
	fmt.Printf("Total time taken to run the program is %.2f seconds\n", endTime)
	*/

	// with connection pool main

	queue := getConnectionPoolQueue()
	start := time.Now()

	
	var wg sync.WaitGroup
	numRequests := 100

	wg.Add(numRequests)

	for i := 1; i <= numRequests; i++ {
		go performDBTaskwConnectionPool(&wg, i, queue)
	}

	wg.Wait()
	endTime := time.Since(start).Seconds()
	fmt.Printf("Total time taken to run the program is %.2f seconds\n", endTime)
}