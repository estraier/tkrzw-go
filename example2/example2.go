/*************************************************************************************************
 * Example for basic usage of the tree database
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
	// The method OrDie causes panic if the status is not success.
	// You should write your own error handling in large scale programs.
	dbm := tkrzw.NewDBM()
	dbm.Open("casket.tkt", true,
		tkrzw.ParseParams("truncate=true,num_buckets=100")).OrDie()

	// Closes the database for sure and checks the error too.
	defer func() { dbm.Close().OrDie() }()

	// Two bank accounts for Bob and Alice.
	// Numeric values are converted into strings implicitly.
	dbm.Set("Bob", 1000, false).OrDie()
	dbm.Set("Alice", 3000, false).OrDie()

	// Function to do a money transfer atomically.
	transfer := func(src_key string, dest_key string, amount int64) *tkrzw.Status {
		// Gets the old values as numbers.
		old_src_value := tkrzw.ToInt(dbm.GetStrSimple(src_key, "0"))
		old_dest_value := tkrzw.ToInt(dbm.GetStrSimple(dest_key, "0"))

		// Calculates the new values.
		new_src_value := old_src_value - amount
		new_dest_value := old_dest_value + amount
		if new_src_value < 0 {
			return tkrzw.NewStatus(tkrzw.StatusApplicationError, "insufficient value")
		}

		// Prepares the pre-condition and the post-condition of the transaction.
		old_records := []tkrzw.KeyValueStrPair{
			{src_key, tkrzw.ToString(old_src_value)},
			{dest_key, tkrzw.ToString(old_dest_value)},
		}
		new_records := []tkrzw.KeyValueStrPair{
			{src_key, tkrzw.ToString(new_src_value)},
			{dest_key, tkrzw.ToString(new_dest_value)},
		}

		// Performs the transaction atomically.
		// This fails safely if other concurrent transactions break the pre-condition.
		return dbm.CompareExchangeMultiStr(old_records, new_records)
	}

	// Tries a transaction until it succeeds
	var status *tkrzw.Status
	for num_tries := 0; num_tries < 100; num_tries++ {
		status = transfer("Alice", "Bob", 500)
		if !status.Equals(tkrzw.StatusInfeasibleError) {
			break
		}
	}
	status.OrDie()

	// Traverses records in a primitive way.
	iter := dbm.MakeIterator()
	defer iter.Destruct()
	iter.First()
	for {
		key, value, status := iter.GetStr()
		if !status.IsOK() {
			break
		}
		fmt.Println(key, value)
		iter.Next()
	}
}

// END OF FILE
