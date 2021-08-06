/*************************************************************************************************
 * Example for basic usage of the hash database
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
)

func main() {
	// Prepares the database.
	dbm := tkrzw.NewDBM()
	dbm.Open("casket.tkt", true, "truncate=true,num_buckets=100")
	defer dbm.Close()

	// Prepares the asynchronous adapter with 4 worker threads.
	async := tkrzw.NewAsyncDBM(dbm, 4)
	defer async.Destruct()

	// Executes the Set method asynchronously.
	future := async.Set("hello", "world", true)
	// Does something in the foreground.
	fmt.Println("Setting a record")
	// Checks the result after awaiting the set operation.
	// The Get method releases the resource of the future.
	status := future.Get()
	if !status.IsOK() {
		fmt.Println("ERROR: " + status.String())
	}

	// Executes the get method asynchronously.
	future = async.Get("hello")
	// Does something in the foreground.
	fmt.Println("Getting a record")
	// Checks the result after awaiting the get operation.
	// The GetStr method releases the resource of the future.
	value, status := future.GetStr()
	if status.IsOK() {
		fmt.Println("VALUE: " + value)
	}
}

// END OF FILE
