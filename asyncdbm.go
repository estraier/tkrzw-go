/*************************************************************************************************
 * Asynchronous database manager adapter
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

package tkrzw

import (
	"fmt"
)

// Asynchronous database manager adapter.
//
// This class is a wrapper of DBM for asynchronous operations.  A task queue with a thread pool is used inside.  Every method except for the constructor and the destructor is run by a thread in the thread pool and the result is set in the future oject of the return value.  The caller can ignore the future object if it is not necessary.  The destruct method waits for all tasks to be done.  Therefore, the destructor should be called before the database is closed.
type AsyncDBM struct {
	// Pointer to the internal object.
	async uintptr
}

// Makes a new AsyncDBM object.
//
// @param dbm A database object which has been opened.
// @param num_worker_threads The number of threads in the internal thread pool.
// @return The pointer to the created database object.
func NewAsyncDBM(dbm *DBM, num_worker_threads int) *AsyncDBM {
	if dbm.dbm == 0 {
		return nil
	}
	async := async_dbm_new(dbm.dbm, num_worker_threads)
	return &AsyncDBM{async}
}

// Destructs the object and releases resources.
//
// This method waits for all tasks to be done.
func (self *AsyncDBM) Destruct() {
	if self.async == 0 {
		return
	}
	async_dbm_free(self.async)
	self.async = 0
}

// Makes a string representing the adapter.
//
// @return The string representing the adapter.
func (self *AsyncDBM) String() string {
	if self.async == 0 {
		return fmt.Sprintf("#<tkrzw.AsyncDBM:%p:destructed>", &self)
	}
	return fmt.Sprintf("#<tkrzw.AsyncDBM:%p:0x%x>", &self, self.async)
}

// Gets the value of a record of a key.
//
// @param key The key of the record.
// @return The future for the record value and the result status.  If there's no matching record, StatusNotFoundError is set.  The result should be gotten by the GetBytes or GetStr method of the future.
func (self *AsyncDBM) Get(key interface{}) *Future {
	if self.async == 0 {
		return &Future{0}
	}
	return async_dbm_get(self.async, ToByteArray(key))
}

// Gets the values of multiple records of keys.
//
// @param keys The keys of records to retrieve.
// @return The future for a map of retrieved records and the result status.  Keys which don't match existing records are ignored.  The result should be gotten by the GetMap or GetMapStr method of the future.
func (self *AsyncDBM) GetMulti(keys []string) *Future {
	if self.async == 0 {
		return &Future{0}
	}
	return async_dbm_get_multi(self.async, keys)
}

// Sets a record of a key and a value.
//
// @param key The key of the record.
// @param value The value of the record.
// @param overwrite Whether to overwrite the existing value.
// @return The future for the result status.  If overwriting is abandoned, StatusDuplicationError is set.  The result should be gotten by the Get method of the future.
func (self *AsyncDBM) Set(key interface{}, value interface{}, overwrite bool) *Future {
	if self.async == 0 {
		return &Future{0}
	}
	return async_dbm_set(self.async, ToByteArray(key), ToByteArray(value), overwrite)
}

// Sets multiple records.
//
// @param records Records to store.
// @param overwrite Whether to overwrite the existing value if there's a record with the same key.  If true, the existing value is overwritten by the new value.  If false, the operation is given up and an error status is returned.
// @return The future for the result status.  If there are records avoiding overwriting, StatusDuplicationError is set.  The result should be gotten by the Get method of the future.
func (self *AsyncDBM) SetMulti(records map[string][]byte, overwrite bool) *Future {
	if self.async == 0 {
		return &Future{0}
	}
	return async_dbm_set_multi(self.async, records, overwrite)
}

// Sets multiple records, with string data.
//
// @param records Records to store.
// @param overwrite Whether to overwrite the existing value if there's a record with the same key.  If true, the existing value is overwritten by the new value.  If false, the operation is given up and an error status is returned.
// @return The future for the result status.  If there are records avoiding overwriting, StatusDuplicationError is set.  The result should be gotten by the Get method of the future.
func (self *AsyncDBM) SetMultiStr(records map[string]string, overwrite bool) *Future {
	if self.async == 0 {
		return &Future{0}
	}
	rawRecords := make(map[string][]byte)
	for key, value := range records {
		rawRecords[key] = []byte(value)
	}
	return async_dbm_set_multi(self.async, rawRecords, overwrite)
}

// Removes a record of a key.
//
// @param key The key of the record.
// @return The future for the result status.  If there's no matching record, StatusNotFoundError is set.  The result should be gotten by the Get method of the future.
func (self *AsyncDBM) Remove(key interface{}) *Future {
	if self.async == 0 {
		return &Future{0}
	}
	return async_dbm_remove(self.async, ToByteArray(key))
}

// Removes records of keys.
//
// @param key The keys of the records.
// @return The future for the result status.  If there are missing records, StatusNotFoundError is set.  The result should be gotten by the Get method of the future.
func (self *AsyncDBM) RemoveMulti(keys []string) *Future {
	if self.async == 0 {
		return &Future{0}
	}
	return async_dbm_remove_multi(self.async, keys)
}

// Appends data at the end of a record of a key.
//
// @param key The key of the record.
// @param value The value to append.
// @param delim The delimiter to put after the existing record.
// @return The future for the result status.
//
// If there's no existing record, the value is set without the delimiter.  The result should be gotten by the Get method of the future.
func (self *AsyncDBM) Append(key interface{}, value interface{}, delim interface{}) *Future {
	if self.async == 0 {
		return &Future{0}
	}
	return async_dbm_append(self.async, ToByteArray(key), ToByteArray(value), ToByteArray(delim))
}

// Appends data to multiple records.
//
// @param records Records to append.
// @param delim The delimiter to put after the existing record.
// @return The future for the result status.  The result should be gotten by the Get method of the future.
//
// If there's no existing record, the value is set without the delimiter.
func (self *AsyncDBM) AppendMulti(records map[string][]byte, delim interface{}) *Future {
	if self.async == 0 {
		return &Future{0}
	}
	return async_dbm_append_multi(self.async, records, ToByteArray(delim))
}

// Appends data to multiple records, with string data.
//
// @param records Records to append.
// @param delim The delimiter to put after the existing record.
// @return The future for the result status.  The result should be gotten by the Get method of the future.
//
// If there's no existing record, the value is set without the delimiter.
func (self *AsyncDBM) AppendMultiStr(records map[string]string, delim interface{}) *Future {
	if self.async == 0 {
		return &Future{0}
	}
	rawRecords := make(map[string][]byte)
	for key, value := range records {
		rawRecords[key] = []byte(value)
	}
	return async_dbm_append_multi(self.async, rawRecords, ToByteArray(delim))
}

// Compares the value of a record and exchanges if the condition meets.
//
// @param key The key of the record.
// @param expected The expected value.  If it is nil or NilString, no existing record is expected.  If it is AnyBytes or AnyString, an existing record with any value is expacted.
// @param desired The desired value.  If it is nil or NilString, the record is to be removed.  If it is AnyBytes or AnyString, no update is done.
// @return The future for the result status.  If the condition doesn't meet, StatusInfeasibleError is set.  The result should be gotten by the Get method of the future.
func (self *AsyncDBM) CompareExchange(
	key interface{}, expected interface{}, desired interface{}) *Future {
	if self.async == 0 {
		return &Future{0}
	}
	var rawExpected []byte
	if !IsNilData(expected) {
		if (IsAnyData(expected)) {
			rawExpected = AnyBytes;
		} else {
			rawExpected = ToByteArray(expected)
		}
	}
	var rawDesired []byte
	if !IsNilData(desired) {
		if (IsAnyData(desired)) {
			rawDesired = AnyBytes;
		} else {
			rawDesired = ToByteArray(desired)
		}
	}
	return async_dbm_compare_exchange(self.async, ToByteArray(key), rawExpected, rawDesired)
}

// Increments the numeric value of a record.
//
// @param key The key of the record.
// @param inc The incremental value.  If it is Int64Min, the current value is not changed and a new record is not created.
// @param init The initial value.
// @return The future for the result status.and the current value.  The result should be gotten by the GetInt method of the future.
func (self *AsyncDBM) Increment(key interface{}, inc interface{}, init interface{}) *Future {
	if self.async == 0 {
		return &Future{0}
	}
	return async_dbm_increment(self.async, ToByteArray(key), ToInt(inc), ToInt(init))
}

// Compares the values of records and exchanges if the condition meets.
//
// @param expected A sequence of pairs of the record keys and their expected values.  If the value is nil, no existing record is expected.  If the value is AnyBytes, an existing record with any value is expacted.
// @param desired A sequence of pairs of the record keys and their desired values.  If the value is nil, the record is to be removed.
// @return The future for the result status.  If the condition doesn't meet, StatusInfeasibleError is set.  The result should be gotten by the Get method of the future.
func (self *AsyncDBM) CompareExchangeMulti(
	expected []KeyValuePair, desired []KeyValuePair) *Future {
	if self.async == 0 {
		return &Future{0}
	}
	return async_dbm_compare_exchange_multi(self.async, expected, desired)
}

// Compares the values of records and exchanges if the condition meets, using string data.
//
// @param expected A sequence of pairs of the record keys and their expected values.  If the value is NilString, no existing record is expected.  If the value is AnyString, an existing record with any value is expacted.
// @param desired A sequence of pairs of the record keys and their desired values.  If the value is NilString, the record is to be removed.
// @return The future for The result status.  If the condition doesn't meet, StatusInfeasibleError is set.  The result should be gotten by the Get method of the future.
func (self *AsyncDBM) CompareExchangeMultiStr(
	expected []KeyValueStrPair, desired []KeyValueStrPair) *Future {
	if self.async == 0 {
		return &Future{0}
	}
	rawExpected := make([]KeyValuePair, 0, len(expected))
	for _, record := range expected {
		var value []byte
		if !IsNilString(record.Value) {
			if IsAnyString(record.Value) {
				value = AnyBytes
			} else {
				value = []byte(record.Value)
			}
		}
		rawExpected = append(rawExpected, KeyValuePair{[]byte(record.Key), value})
	}
	rawDesired := make([]KeyValuePair, 0, len(desired))
	for _, record := range desired {
		var value []byte
		if !IsNilString(record.Value) {
			value = []byte(record.Value)
		}
		rawDesired = append(rawDesired, KeyValuePair{[]byte(record.Key), value})
	}
	return async_dbm_compare_exchange_multi(self.async, rawExpected, rawDesired)
}

// Changes the key of a record.
//
// @param old_key The old key of the record.
// @param new_key The new key of the record.
// @param overwrite Whether to overwrite the existing record of the new key.
// @param copying Whether to retain the record of the old key.
// @return The future for the result status.  If there's no matching record to the old key, NOT_FOUND_ERROR is set.  If the overwrite flag is false and there is an existing record of the new key, DUPLICATION ERROR is set.  The result should be gotten by the Get method of the future.
//
// This method is done atomically by ProcessMulti.  The other threads observe that the record has either the old key or the new key.  No intermediate states are observed.
func (self *AsyncDBM) Rekey(old_key interface{}, new_key interface{},
	overwrite bool, copying bool) *Future {
	if self.async == 0 {
		return &Future{0}
	}
	return async_dbm_rekey(self.async, ToByteArray(old_key), ToByteArray(new_key),
		overwrite, copying)
}

// Gets the first record and removes it.
//
// @return A tuple of the result status, the key and the value of the first record.  The result should be gotten by the GetPair or GetPairStr method of the future.
func (self *AsyncDBM) PopFirst() *Future {
	if self.async == 0 {
		return &Future{0}
	}
	return async_dbm_pop_first(self.async)
}

// Adds a record with a key of the current timestamp.
//
// @param value The value of the record.
// @param wtime The current wall time used to generate the key.  If it is None, the system clock is used.
// @return The future for the result status.
//
// The key is generated as an 8-bite big-endian binary string of the timestamp.  If there is an existing record matching the generated key, the key is regenerated and the attempt is repeated until it succeeds.
func (self *AsyncDBM) PushLast(value interface{}, wtime float64) *Future {
	if self.async == 0 {
		return &Future{0}
	}
	return async_dbm_push_last(self.async, ToByteArray(value), wtime)
}

// Removes all records.
//
// @return The future for the result status.  The result should be gotten by the Get method of the future.
func (self *AsyncDBM) Clear() *Future {
	if self.async == 0 {
		return &Future{0}
	}
	return async_dbm_clear(self.async)
}

// Rebuilds the entire database.
//
// @param params Optional parameters.  If it is nil, it is ignored.
// @return The future for the result status.  The result should be gotten by the Get method of the future.
//
// The parameters work in the same way as with DBM::Rebuild.
func (self *AsyncDBM) Rebuild(params map[string]string) *Future {
	if self.async == 0 {
		return &Future{0}
	}
	return async_dbm_rebuild(self.async, params)
}

// Synchronizes the content of the database to the file system.
//
// @param hard True to do physical synchronization with the hardware or false to do only logical synchronization with the file system.
// @param params Optional parameters.  If it is nil, it is ignored.
// @return The future for the result status.  The result should be gotten by the Get method of the future.
//
// The parameters work in the same way as with DBM::Synchronize.
func (self *AsyncDBM) Synchronize(hard bool, params map[string]string) *Future {
	if self.async == 0 {
		return &Future{0}
	}
	return async_dbm_synchronize(self.async, hard, params)
}

// Copies the content of the database file to another file.
//
// @param destPath A path to the destination file.
// @param syncHard True to do physical synchronization with the hardware.
// @return The future for the result status.  The result should be gotten by the Get method of the future.
func (self *AsyncDBM) CopyFileData(destPath string, syncHard bool) *Future {
	if self.async == 0 {
		return &Future{0}
	}
	return async_dbm_copy_file_data(self.async, destPath, syncHard)
}

// Exports all records to another database.
//
// @param destDBM The destination database.  The lefetime of the database object must last until the task finishes.
// @return The future for the result status.  The result should be gotten by the Get method of the future.
func (self *AsyncDBM) Export(destDBM *DBM) *Future {
	if self.async == 0 || destDBM.dbm == 0 {
		return &Future{0}
	}
	return async_dbm_export(self.async, destDBM.dbm)
}

// Exports all records of a database to a flat record file.
//
// @param file: The file object to write records in.  The lefetime of the file object must last until the task finishes.
// @return The future for the result status.  The result should be gotten by the Get method of the future.
//
// A flat record file contains a sequence of binary records without any high level structure so it is useful as a intermediate file for data migration.
func (self *AsyncDBM) ExportToFlatRecords(destFile *File) *Future {
	if self.async == 0 || destFile.file == 0 {
		return &Future{0}
	}
	return async_dbm_export_to_flat_records(self.async, destFile.file)
}

// Imports records to a database from a flat record file.
//
// @param file The file object to read records from.  The lefetime of the file object must last until the task finishes.
// @return The future for the result status.  The result should be gotten by the Get method of the future.
func (self *AsyncDBM) ImportFromFlatRecords(srcFile *File) *Future {
	if self.async == 0 || srcFile.file == 0 {
		return &Future{0}
	}
	return async_dbm_import_from_flat_records(self.async, srcFile.file)
}

// Searches the database and get keys which match a pattern.
//
// @param mode The search mode.  "contain" extracts keys containing the pattern.  "begin" extracts keys beginning with the pattern.  "end" extracts keys ending with the pattern.  "regex" extracts keys partially matches the pattern of a regular expression.  "edit" extracts keys whose edit distance to the UTF-8 pattern is the least.  "editbin" extracts keys whose edit distance to the binary pattern is the least.
// @param pattern The pattern for matching.
// @param capacity The maximum records to obtain.  0 means unlimited.
// @return The future for a list of keys matching the condition and the result status.  The result should be gotten by the GetArray or GetArrayStr method of the future.
func (self *AsyncDBM) Search(mode string, pattern string, capacity int) *Future {
	if self.async == 0 {
		return &Future{0}
	}
	return async_dbm_search(self.async, mode, pattern, capacity)
}

// END OF FILE
