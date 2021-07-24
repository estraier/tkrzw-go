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
	dbm.Open("casket.tkh", true, "truncate=true,num_buckets=100")

	// Sets records.
	// Keys and values are implicitly converted into bytes.
	dbm.Set("first", "hop", true)
	dbm.Set("second", "step", true)
	dbm.Set("third", "jump", true)

	// Retrieves record values as strings.
	fmt.Println(dbm.GetStrSimple("first", "*"))
	fmt.Println(dbm.GetStrSimple("second", "*"))
	fmt.Println(dbm.GetStrSimple("third", "*"))

	// Traverses records.
	iter := dbm.MakeIterator()
	iter.First()
	for {
		key, value, status := iter.GetStr()
		if !status.IsOK() {
			break
		}
		fmt.Println(key, value)
		iter.Next()
	}
	iter.Destruct()

	// Closes the database.
	dbm.Close()
}

// END OF FILE
