/*************************************************************************************************
 * Example for process methods
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
	"regexp"
	"strings"
)

func main() {
	// Opens the database.
	dbm := tkrzw.NewDBM()
	dbm.Open("casket.tkh", true,
		tkrzw.ParseParams("truncate=true,num_buckets=100"))
	defer dbm.Close()

	// Sets records with lambda functions.
	dbm.Process("doc-1", func(k []byte, v []byte) interface{} {
		return "Tokyo is the capital city of Japan."
	}, true)
	dbm.Process("doc-2", func(k []byte, v []byte) interface{} {
		return "Is she living in Tokyo, Japan?"
	}, true)
	dbm.Process("doc-3", func(k []byte, v []byte) interface{} {
		return "She must leave Tokyo!"
	}, true)

	// Lowers record values.
	lower := func(k []byte, v []byte) interface{} {
		// If no matching record, nil is given as the value.
		if v == nil {
			return nil
		}
		// Sets the new value.
		// Note that the key and the value are a "bytes" object.
		return strings.ToLower(string(v))
	}
	dbm.Process("doc-1", lower, true)
	dbm.Process("doc-2", lower, true)
	dbm.Process("doc-3", lower, true)
	dbm.Process("non-existent", lower, true)

	// Adds multiple records at once.
	ops := []tkrzw.KeyProcPair{
		{"doc-4", func(k []byte, v []byte) interface{} { return "Tokyo Go!" }},
		{"doc-5", func(k []byte, v []byte) interface{} { return "Japan Go!" }},
	}
	dbm.ProcessMulti(ops, true)

	// Modifies multiple records at once.
	dbm.ProcessMulti([]tkrzw.KeyProcPair{{"doc-4", lower}, {"doc-5", lower}}, true)

	// Checks the whole content.
	// This uses an external iterator and is relavively slow.
	for record := range dbm.EachStr() {
		fmt.Println(record.Key, record.Value)
	}

	// Function for word counting.
	wordCounts := make(map[string]int)
	wordSplitter := regexp.MustCompile("\\W")
	wordCounter := func(key []byte, value []byte) interface{} {
		if key == nil {
			return nil
		}
		words := wordSplitter.Split(string(value), -1)
		for _, word := range words {
			if len(word) == 0 {
				continue
			}
			wordCounts[word] += 1
		}
		return nil
	}

	// The second parameter should be false if the value is not updated.
	dbm.ProcessEach(wordCounter, false)
	for word, count := range wordCounts {
		fmt.Printf("%s = %d\n", word, count)
	}

	// Returning RemoveBytes by the callbacks removes the record.
	dbm.Process("doc-1", func(k []byte, v []byte) interface{} {
		return tkrzw.RemoveBytes
	}, true)
	println(dbm.CountSimple())
	dbm.ProcessMulti([]tkrzw.KeyProcPair{
		{"doc-4", func(k []byte, v []byte) interface{} { return tkrzw.RemoveBytes }},
		{"doc-5", func(k []byte, v []byte) interface{} { return tkrzw.RemoveBytes }},
	}, true)
	println(dbm.CountSimple())
	dbm.ProcessEach(func(k []byte, v []byte) interface{} {
		return tkrzw.RemoveBytes
	}, true)
	println(dbm.CountSimple())
}

// END OF FILE
