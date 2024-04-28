/*************************************************************************************************
 * Example for secondary index
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
	// Opens the index.
	index := tkrzw.NewIndex()
	index.Open("casket.tkt", true,
		tkrzw.ParseParams("truncate=true,num_buckets=100")).OrDie()
	defer func() { index.Close().OrDie() }()

	// Adds records to the index.
	// The key is a division name and the value is person name.
	index.Add("general", "anne").OrDie()
	index.Add("general", "matthew").OrDie()
	index.Add("general", "marilla").OrDie()
	index.Add("sales", "gilbert").OrDie()

	// Anne moves to the sales division.
	index.Remove("general", "anne").OrDie()
	index.Add("sales", "anne").OrDie()

	// Prints all members for each division.
	divisions := [] string{"general", "sales"}
	for _, division := range divisions {
		fmt.Printf("%s\n", division)
		members := index.GetValuesStr(division, 0)
		for _, member := range members {
			fmt.Printf(" -- %s\n", member)
		}
	}

	// Prints every record by iterator.
	iter := index.MakeIterator()
	defer iter.Destruct()
	iter.First()
	for {
		key, value, ok := iter.GetStr()
		if !ok {
			break
		}
		fmt.Printf("%s: %s\n", key, value)
		iter.Next()
	}
}

// END OF FILE
