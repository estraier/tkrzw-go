/*************************************************************************************************
 * Example for key comparators of the tree database
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
	// Opens a new database with the default key comparator (LexicalKeyComparator).
	dbm := tkrzw.NewDBM()
	dbm.Open("casket.tkt", true, tkrzw.ParseParams("truncate=true")).OrDie()

	// Sets records with the key being a big-endian binary of an integer.
	// e.g: "\x00\x00\x00\x00\x00\x00\x00\x31" -> "hop"
	dbm.Set(tkrzw.SerializeInt(1), "hop", true).OrDie()
	dbm.Set(tkrzw.SerializeInt(256), "step", true).OrDie()
	dbm.Set(tkrzw.SerializeInt(32), "jump", true).OrDie()

	// Gets records with the key being a decimal string of an integer.
	fmt.Println(dbm.GetStrSimple(tkrzw.SerializeInt(1), ""))
	fmt.Println(dbm.GetStrSimple(tkrzw.SerializeInt(256), ""))
	fmt.Println(dbm.GetStrSimple(tkrzw.SerializeInt(32), ""))

	// Lists up all records, restoring keys into integers.
	iter := dbm.MakeIterator()
	iter.First()
	for {
		key, value, status := iter.Get()
		if !status.IsOK() {
			break
		}
		fmt.Printf("%d: %s\n", tkrzw.DeserializeInt(key), value)
		iter.Next()
	}
	iter.Destruct()

	// Closes the database.
	dbm.Close().OrDie()

	// Opens a new database with the decimal integer comparator.
	dbm = tkrzw.NewDBM()
	dbm.Open("casket.tkt", true, tkrzw.ParseParams(
		"truncate=true,key_comparator=Decimal")).OrDie()

	// Sets records with the key being a decimal string of an integer.
	// e.g: "1" -> "hop"
	dbm.Set("1", "hop", true).OrDie()
	dbm.Set("256", "step", true).OrDie()
	dbm.Set("32", "jump", true).OrDie()

	// Gets records with the key being a decimal string of an integer.
	fmt.Println(dbm.GetStrSimple("1", ""))
	fmt.Println(dbm.GetStrSimple("256", ""))
	fmt.Println(dbm.GetStrSimple("32", ""))

	// Lists up all records, restoring keys into integers.
	iter = dbm.MakeIterator()
	iter.First()
	for {
		key, value, status := iter.Get()
		if !status.IsOK() {
			break
		}
		fmt.Printf("%d: %s\n", tkrzw.ToInt(key), value)
		iter.Next()
	}
	iter.Destruct()

	// Closes the database.
	dbm.Close().OrDie()

	// Opens a new database with the decimal real number comparator.
	dbm = tkrzw.NewDBM()
	dbm.Open("casket.tkt", true, tkrzw.ParseParams(
		"truncate=true,key_comparator=RealNumber")).OrDie()

	// Sets records with the key being a decimal string of a real number.
	// e.g: "1.5" -> "hop"
	dbm.Set("1.5", "hop", true).OrDie()
	dbm.Set("256.5", "step", true).OrDie()
	dbm.Set("32.5", "jump", true).OrDie()

	// Gets records with the key being a decimal string of a real number.
	fmt.Println(dbm.GetStrSimple("1.5", ""))
	fmt.Println(dbm.GetStrSimple("256.5", ""))
	fmt.Println(dbm.GetStrSimple("32.5", ""))

	// Lists up all records, restoring keys into floating-point numbers.
	iter = dbm.MakeIterator()
	iter.First()
	for {
		key, value, status := iter.Get()
		if !status.IsOK() {
			break
		}
		fmt.Printf("%.3f: %s\n", tkrzw.ToFloat(key), value)
		iter.Next()
	}
	iter.Destruct()

	// Closes the database.
	dbm.Close().OrDie()

	// Opens a new database with the big-endian floating-point numbers comparator.
	dbm = tkrzw.NewDBM()
	dbm.Open("casket.tkt", true, tkrzw.ParseParams(
		"truncate=true,key_comparator=FloatBigEndian")).OrDie()

	// Sets records with the key being a big-endian binary of a floating-point number.
	// e.g: "\x3F\xF8\x00\x00\x00\x00\x00\x00" -> "hop"
	dbm.Set(tkrzw.SerializeFloat(1.5), "hop", true).OrDie()
	dbm.Set(tkrzw.SerializeFloat(256.5), "step", true).OrDie()
	dbm.Set(tkrzw.SerializeFloat(32.5), "jump", true).OrDie()

	// Gets records with the key being a big-endian binary of a floating-point number.
	fmt.Println(dbm.GetStrSimple(tkrzw.SerializeFloat(1.5), ""))
	fmt.Println(dbm.GetStrSimple(tkrzw.SerializeFloat(256.5), ""))
	fmt.Println(dbm.GetStrSimple(tkrzw.SerializeFloat(32.5), ""))

	// Lists up all records, restoring keys into floating-point numbers.
	iter = dbm.MakeIterator()
	iter.First()
	for {
		key, value, status := iter.Get()
		if !status.IsOK() {
			break
		}
		fmt.Printf("%.3f: %s\n", tkrzw.DeserializeFloat(key), value)
		iter.Next()
	}
	iter.Destruct()

	// Closes the database.
	dbm.Close().OrDie()
}

// END OF FILE
