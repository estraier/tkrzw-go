/*************************************************************************************************
 * Performance tests
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
	"flag"
	"fmt"
	"github.com/estraier/tkrzw-go"
	"math/rand"
	"time"
)

var flagPath = flag.String("path", "", "the file path of the database")
var flagParams = flag.String("params", "", "the parameters for the database")
var flagNumIterations = flag.Int("iter", 10000, "the number of iterations")
var flagNumThreads = flag.Int("threads", 1, "the number of threads")
var flagIsRandom = flag.Bool("random", false, "whether to use random keys")

func main() {
	flag.Parse()
	path := *flagPath
	openParams := *flagParams + ",truncate=true"
	numIterations := *flagNumIterations
	numThreads := *flagNumThreads
	isRandom := *flagIsRandom
	fmt.Printf("path: %s\n", path)
	fmt.Printf("params: %s\n", openParams)
	fmt.Printf("num_iterations: %d\n", numIterations)
	fmt.Printf("num_threads: %d\n", numThreads)
	fmt.Printf("is_random: %t\n", isRandom)
	fmt.Println()
	startMemUsage := tkrzw.GetMemoryUsage()
	dbm := tkrzw.NewDBM()
	dbm.Open(path, true, openParams).OrDie()
	fmt.Println("Setting:")
	startTime := time.Now()
	setter := func(thid int, done chan<- bool) {
		random := rand.New(rand.NewSource(int64(thid)))
		for i := 0; i < numIterations; i++ {
			var keyNum int
			if isRandom {
				keyNum = random.Intn(numIterations * numThreads)
			} else {
				keyNum = thid*numIterations + i
			}
			key := fmt.Sprintf("%08d", keyNum)
			dbm.Set(key, key, true).OrDie()
			seq := i + 1
			if thid == 0 && seq%(numIterations/500) == 0 {
				fmt.Print(".")
				if seq%(numIterations/10) == 0 {
					fmt.Printf(" (%08d)\n", seq)
				}
			}
		}
		done <- true
	}
	dones := make([]chan bool, 0)
	for i := 0; i < numThreads; i++ {
		done := make(chan bool)
		go setter(i, done)
		dones = append(dones, done)
	}
	for _, done := range dones {
		<-done
	}
	dbm.Synchronize(false, "").OrDie()
	endTime := time.Now()
	elapsed := endTime.Sub(startTime).Seconds()
	memUsage := tkrzw.GetMemoryUsage() - startMemUsage
	fmt.Printf("Setting done: num_records=%d file_size=%d time=%.3f qps=%.0f mem=%d\n",
		dbm.CountSimple(), dbm.GetFileSizeSimple(),
		elapsed, float64(numIterations*numThreads)/elapsed, memUsage)
	fmt.Println()
	fmt.Println("Getting:")
	startTime = time.Now()
	getter := func(thid int, done chan<- bool) {
		random := rand.New(rand.NewSource(int64(thid)))
		for i := 0; i < numIterations; i++ {
			var keyNum int
			if isRandom {
				keyNum = random.Intn(numIterations * numThreads)
			} else {
				keyNum = thid*numIterations + i
			}
			key := fmt.Sprintf("%08d", keyNum)
			_, status := dbm.Get(key)
			if !status.Equals(tkrzw.StatusSuccess) && !status.Equals(tkrzw.StatusNotFoundError) {
				panic(status.String())
			}
			seq := i + 1
			if thid == 0 && seq%(numIterations/500) == 0 {
				fmt.Print(".")
				if seq%(numIterations/10) == 0 {
					fmt.Printf(" (%08d)\n", seq)
				}
			}
		}
		done <- true
	}
	dones = make([]chan bool, 0)
	for i := 0; i < numThreads; i++ {
		done := make(chan bool)
		go getter(i, done)
		dones = append(dones, done)
	}
	for _, done := range dones {
		<-done
	}
	endTime = time.Now()
	elapsed = endTime.Sub(startTime).Seconds()
	memUsage = tkrzw.GetMemoryUsage() - startMemUsage
	fmt.Printf("Getting done: num_records=%d file_size=%d time=%.3f qps=%.0f mem=%d\n",
		dbm.CountSimple(), dbm.GetFileSizeSimple(),
		elapsed, float64(numIterations*numThreads)/elapsed, memUsage)
	fmt.Println()
	fmt.Println("Removing:")
	startTime = time.Now()
	remover := func(thid int, done chan<- bool) {
		random := rand.New(rand.NewSource(int64(thid)))
		for i := 0; i < numIterations; i++ {
			var keyNum int
			if isRandom {
				keyNum = random.Intn(numIterations * numThreads)
			} else {
				keyNum = thid*numIterations + i
			}
			key := fmt.Sprintf("%08d", keyNum)
			status := dbm.Remove(key)
			if !status.Equals(tkrzw.StatusSuccess) && !status.Equals(tkrzw.StatusNotFoundError) {
				panic(status.String())
			}
			seq := i + 1
			if thid == 0 && seq%(numIterations/500) == 0 {
				fmt.Print(".")
				if seq%(numIterations/10) == 0 {
					fmt.Printf(" (%08d)\n", seq)
				}
			}
		}
		done <- true
	}
	dones = make([]chan bool, 0)
	for i := 0; i < numThreads; i++ {
		done := make(chan bool)
		go remover(i, done)
		dones = append(dones, done)
	}
	for _, done := range dones {
		<-done
	}
	endTime = time.Now()
	elapsed = endTime.Sub(startTime).Seconds()
	memUsage = tkrzw.GetMemoryUsage() - startMemUsage
	fmt.Printf("Removing done: num_records=%d file_size=%d time=%.3f qps=%.0f mem=%d\n",
		dbm.CountSimple(), dbm.GetFileSizeSimple(),
		elapsed, float64(numIterations*numThreads)/elapsed, memUsage)
	fmt.Println()
	dbm.Close().OrDie()
}

// END OF FILE
