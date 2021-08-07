/*************************************************************************************************
 * Example to compare performance of goroutines and the asynchronous API
 *
 * Copyright 2020 Google LLC
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License.  You may obtain a copy of the License at
 *     https://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software distributed under the
 * License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied.  See the License for the specific language governing permissions
 * and limitations under the License.
 *************************************************************************************************/

package main

import (
	"fmt"
	"github.com/estraier/tkrzw-go"
	"time"
)

// Sets a record in the database.
func set(dbm *tkrzw.DBM, key string, value string, done chan<- *tkrzw.Status) {
	done <- dbm.Set(key, value, true)
}

func main() {
	// Prepares resources.
	numBatches := 100
	numRecords := 100
	dbm := tkrzw.NewDBM()
	path := "casket.tkh"
	params := fmt.Sprintf("truncate=true,num_buckets=%d", numBatches*numRecords*2)

	// Evaluates asynchrnouns operations with goroutines and channels.
	fmt.Println("Setting with goroutines and channels ...")
	dbm.Open(path, true, params).OrDie()
	startTime := time.Now()
	for batchID := 0; batchID < numBatches; batchID++ {
		dones := make([]chan *tkrzw.Status, 0, numRecords)
		for i := 0; i < numRecords; i++ {
			key := fmt.Sprintf("%08d", batchID*numBatches+i)
			done := make(chan *tkrzw.Status)
			go set(dbm, key, key, done)
			dones = append(dones, done)
		}
		for _, done := range dones {
			status := <-done
			status.OrDie()
		}
	}
	endTime := time.Now()
	elapsed := endTime.Sub(startTime).Seconds()
	fmt.Printf("time=%.3f, qps=%.0f\n", elapsed, float64(numBatches*numRecords)/elapsed)
	dbm.Close().OrDie()

	// Evaluates asynchrnouns operations with goroutines and channels.
	fmt.Println("Setting with the asynchrnouns API ...")
	dbm.Open(path, true, params).OrDie()
	async := tkrzw.NewAsyncDBM(dbm, 2)
	startTime = time.Now()
	for batchID := 0; batchID < numBatches; batchID++ {
		futures := make([]*tkrzw.Future, 0, numRecords)
		for i := 0; i < numRecords; i++ {
			key := fmt.Sprintf("%08d", batchID*numBatches+i)
			future := async.Set(key, key, true)
			futures = append(futures, future)
		}
		for _, future := range futures {
			status := future.Get()
			status.OrDie()
		}
	}
	endTime = time.Now()
	elapsed = endTime.Sub(startTime).Seconds()
	fmt.Printf("time=%.3f, qps=%.0f\n", elapsed, float64(numBatches*numRecords)/elapsed)
	async.Destruct()

	// Releases resources.
	dbm.Close().OrDie()
}
