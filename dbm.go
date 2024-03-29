/*************************************************************************************************
 * Database manager interface
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

// Polymorphic database manager.
//
// All operations except for "Open" and "Close" are thread-safe; Multiple threads can access the same database concurrently.  You can specify a data structure when you call the "Open" method.  Every opened database must be closed explicitly by the "Close" method to avoid data corruption.
type DBM struct {
	// Pointer to the internal object.
	dbm uintptr
}

// A function to process a record.
type RecordProcessor func(key []byte, value []byte) interface{}

// A pair of the key and the value of a record.
type KeyValuePair struct {
	// The key.
	Key []byte
	// The value
	Value []byte
}

// A string pair of the key and the value of a record.
type KeyValueStrPair struct {
	// The key.
	Key string
	// The value
	Value string
}

// A pair of the key and the record processor function.
type KeyProcPair struct {
	// The key.
	Key interface{}
	// The processor function.
	Proc RecordProcessor
}

// A pair of the key bytes and the record processor function.
type KeyBytesProcPair struct {
	// The key.
	Key []byte
	// The processor function.
	Proc RecordProcessor
}

// Makes a new DBM object.
//
// @return The pointer to the created database object.
func NewDBM() *DBM {
	return &DBM{0}
}

// Makes a string representing the database.
//
// @return The string representing the database.
func (self *DBM) String() string {
	if self.dbm == 0 {
		return fmt.Sprintf("#<tkrzw.DBM:%p:unopened>", &self)
	}
	path, _ := dbm_get_file_path(self.dbm)
	count, _ := dbm_count(self.dbm)
	return fmt.Sprintf("#<tkrzw.DBM:%s:%d>", path, count)
}

// Opens a database file.
//
// @param path A path of the file.
// @param writable If true, the file is writable.  If false, it is read-only.
// @param params Optional parameters.  If it is nil, it is ignored.
// @return The result status.
//
// The extension of the path indicates the type of the database.
//
// - .tkh : File hash database (HashDBM)
// - .tkt : File tree database (TreeDBM)
// - .tks : File skip database (SkipDBM)
// - .tkmt : On-memory hash database (TinyDBM)
// - .tkmb : On-memory tree database (BabyDBM)
// - .tkmc : On-memory cache database (CacheDBM)
// - .tksh : On-memory STL hash database (StdHashDBM)
// - .tkst : On-memory STL tree database (StdTreeDBM)
//
// The optional parameters can include options for the file opening operation.
//
// - truncate (bool): True to truncate the file.
// - no_create (bool): True to omit file creation.
// - no_wait (bool): True to fail if the file is locked by another process.
// - no_lock (bool): True to omit file locking.
// - sync_hard (bool): True to do physical synchronization when closing.
//
// The optional parameter "dbm" supercedes the decision of the database type by the extension.  The value is the type name: "HashDBM", "TreeDBM", "SkipDBM", "TinyDBM", "BabyDBM", "CacheDBM", "StdHashDBM", "StdTreeDBM".
//
// The optional parameter "file" specifies the internal file implementation class. The default file class is "MemoryMapAtomicFile".  The other supported classes are "StdFile", "MemoryMapAtomicFile", "PositionalParallelFile", and "PositionalAtomicFile".
//
// For HashDBM, these optional parameters are supported.
//
// - update_mode (string): How to update the database file: "UPDATE_IN_PLACE" for the in-palce or "UPDATE_APPENDING" for the appending mode.
// - record_crc_mode (string): How to add the CRC data to the record: "RECORD_CRC_NONE" to add no CRC to each record, "RECORD_CRC_8" to add CRC-8 to each record, "RECORD_CRC_16" to add CRC-16 to each record, or "RECORD_CRC_32" to add CRC-32 to each record.
// - record_comp_mode (string): How to compress the record data: "RECORD_COMP_NONE" to do no compression, "RECORD_COMP_ZLIB" to compress with ZLib, "RECORD_COMP_ZSTD" to compress with ZStd, "RECORD_COMP_LZ4" to compress with LZ4, "RECORD_COMP_LZMA" to compress with LZMA, "RECORD_COMP_RC4" to cipher with RC4, "RECORD_COMP_AES" to cipher with AES.
// - offset_width (int): The width to represent the offset of records.
// - align_pow (int): The power to align records.
// - num_buckets (int): The number of buckets for hashing.
// - restore_mode (string): How to restore the database file: "RESTORE_SYNC" to restore to the last synchronized state, "RESTORE_READ_ONLY" to make the database read-only, or "RESTORE_NOOP" to do nothing.  By default, as many records as possible are restored.
// - fbp_capacity (int): The capacity of the free block pool.
// - min_read_size (int): The minimum reading size to read a record.
// - cache_buckets (bool): True to cache the hash buckets on memory.
// - cipher_key (string): The encryption key for cipher compressors.
//
// For TreeDBM, all optional parameters for HashDBM are available.  In addition, these optional parameters are supported.
//
// - max_page_size (int): The maximum size of a page.
// - max_branches (int): The maximum number of branches each inner node can have.
// - max_cached_pages (int): The maximum number of cached pages.
// - page_update_mode (string): What to do when each page is updated: "PAGE_UPDATE_NONE" is to do no operation or "PAGE_UPDATE_WRITE" is to write immediately.
// - key_comparator (string): The comparator of record keys: "LexicalKeyComparator" for the lexical order, "LexicalCaseKeyComparator" for the lexical order ignoring case, "DecimalKeyComparator" for the order of the decimal integer numeric expressions, "HexadecimalKeyComparator" for the order of the hexadecimal integer numeric expressions, "RealNumberKeyComparator" for the order of the decimal real number expressions.
//
// For SkipDBM, these optional parameters are supported.
//
// - offset_width (int): The width to represent the offset of records.
// - step_unit (int): The step unit of the skip list.
// - max_level (int): The maximum level of the skip list.
// - restore_mode (string): How to restore the database file: "RESTORE_SYNC" to restore to the last synchronized state or "RESTORE_NOOP" to do nothing make the database read-only.  By default, as many records as possible are restored.
// - sort_mem_size (int): The memory size used for sorting to build the database in the at-random mode.
// - insert_in_order (bool): If true, records are assumed to be inserted in ascending order of the key.
// - max_cached_records (int): The maximum number of cached records.
//
// For TinyDBM, these optional parameters are supported.
//
// - num_buckets (int): The number of buckets for hashing.
//
// For BabyDBM, these optional parameters are supported.
//
// - key_comparator (string): The comparator of record keys. The same ones as TreeDBM.
//
// For CacheDBM, these optional parameters are supported.
//
// - cap_rec_num (int): The maximum number of records.
// - cap_mem_size (int): The total memory size to use.
//
// All databases support taking update logs into files.  It is enabled by setting the prefix of update log files.
//
// - ulog_prefix (str): The prefix of the update log files.
// - ulog_max_file_size (num): The maximum file size of each update log file.  By default, it is 1GiB.
// - ulog_server_id (num): The server ID attached to each log.  By default, it is 0.
// - ulog_dbm_index (num): The DBM index attached to each log.  By default, it is 0.
//
// For the file "PositionalParallelFile" and "PositionalAtomicFile", these optional parameters are supported.
//
// - block_size (int): The block size to which all blocks should be aligned.
// - access_options (str): Values separated by colon.  "direct" for direct I/O.  "sync" for synchrnizing I/O, "padding" for file size alignment by padding, "pagecache" for the mini page cache in the process.
//
// If the optional parameter "num_shards" is set, the database is sharded into multiple shard files.  Each file has a suffix like "-00003-of-00015".  If the value is 0, the number of shards is set by patterns of the existing files, or 1 if they doesn't exist.
func (self *DBM) Open(path string, writable bool, params map[string]string) *Status {
	if self.dbm != 0 {
		return NewStatus2(StatusPreconditionError, "opened database")
	}
	dbm, status := dbm_open(path, writable, params)
	if status.code == StatusSuccess {
		self.dbm = dbm
	}
	return status
}

// Closes the database file.
//
// @return The result status.
func (self *DBM) Close() *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	status := dbm_close(self.dbm)
	self.dbm = 0
	return status
}

// Processes a record with an arbitrary function.
//
// @param key The key of the record.
// @param proc The function to process a record.  The first parameter is the key bytes of the record.  The second parameter is the value bytes of the existing record, or nil if it the record doesn't exist.  The return value is bytes or a string to update the record value.  If the return value is nil or NilString, the record is not modified.  If the return value is RemoveBytes or RemoveString, the record is removed.
// @param writable True if the processor can edit the record.
// @return The result status.
func (self *DBM) Process(key interface{}, proc RecordProcessor, writable bool) *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_process(self.dbm, ToByteArray(key), proc, writable)
}

// Checks if a record exists or not.
//
// @param key The key of the record.
// @return True if the record exists, or false if not.
func (self *DBM) Check(key interface{}) bool {
	if self.dbm == 0 {
		return false
	}
	return dbm_check(self.dbm, ToByteArray(key))
}

// Gets the value of a record of a key.
//
// @param key The key of the record.
// @return The bytes value of the matching record and the result status.  If there's no matching record, the status is StatusNotFoundError.
func (self *DBM) Get(key interface{}) ([]byte, *Status) {
	if self.dbm == 0 {
		return nil, NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_get(self.dbm, ToByteArray(key))
}

// Gets the value of a record of a key, as a string.
//
// @param key The key of the record.
// @return The string value of the matching record and the result status.  If there's no matching record, the status is StatusNotFoundError.
func (self *DBM) GetStr(key interface{}) (string, *Status) {
	if self.dbm == 0 {
		return "", NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_get_str(self.dbm, ToByteArray(key))
}

// Gets the value of a record of a key, in a simple way.
//
// @param key The key of the record.
// @param defaultValue The value to be returned on failure.
// @return The value of the matching record on success, or the default value on failure.
func (self *DBM) GetSimple(key interface{}, defaultValue interface{}) []byte {
	if self.dbm == 0 {
		return ToByteArray(defaultValue)
	}
	value, status := dbm_get(self.dbm, ToByteArray(key))
	if status.code == StatusSuccess {
		return value
	}
	return ToByteArray(defaultValue)
}

// Gets the value of a record of a key, in a simple way, as a string.
//
// @param key The key of the record.
// @param defaultValue The value to be returned on failure.
// @return The value of the matching record on success, or the default value on failure.
func (self *DBM) GetStrSimple(key interface{}, defaultValue interface{}) string {
	if self.dbm == 0 {
		return ToString(defaultValue)
	}
	value, status := dbm_get(self.dbm, ToByteArray(key))
	if status.code == StatusSuccess {
		return string(value)
	}
	return ToString(defaultValue)
}

// Gets the values of multiple records of keys.
//
// @param keys The keys of records to retrieve.
// @return A map of retrieved records.  Keys which don't match existing records are ignored.
func (self *DBM) GetMulti(keys []string) map[string][]byte {
	if self.dbm == 0 {
		return make(map[string][]byte)
	}
	return dbm_get_multi(self.dbm, keys)
}

// Gets the values of multiple records of keys, as strings.
//
// @param keys The keys of records to retrieve.
// @eturn A map of retrieved records.  Keys which don't match existing records are ignored.
func (self *DBM) GetMultiStr(keys []string) map[string]string {
	if self.dbm == 0 {
		return make(map[string]string)
	}
	return dbm_get_multi_str(self.dbm, keys)
}

// Sets a record of a key and a value.
//
// @param key The key of the record.
// @param value The value of the record.
// @param overwrite Whether to overwrite the existing value.
// @return The result status.  If overwriting is abandoned, StatusDuplicationError is returned.
func (self *DBM) Set(key interface{}, value interface{}, overwrite bool) *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_set(self.dbm, ToByteArray(key), ToByteArray(value), overwrite)
}

// Sets a record and get the old value.
//
// @param key: The key of the record.
// @param value The value of the record.
// @param overwrite Whether to overwrite the existing value.
// @return The old value and the result status.
func (self *DBM) SetAndGet(key interface{}, value interface{}, overwrite bool) ([]byte, *Status) {
	if self.dbm == 0 {
		return nil, NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_set_and_get(self.dbm, ToByteArray(key), ToByteArray(value), overwrite)
}

// Sets a record and get the old value, as a string.
//
// @param key: The key of the record.
// @param value The value of the record.
// @param overwrite Whether to overwrite the existing value.
// @return The old value and the result status.
func (self *DBM) SetAndGetStr(key interface{}, value interface{},
	overwrite bool) (*string, *Status) {
	if self.dbm == 0 {
		return nil, NewStatus2(StatusPreconditionError, "not opened database")
	}
	old_value, status := dbm_set_and_get(self.dbm, ToByteArray(key), ToByteArray(value), overwrite)
	if old_value != nil {
		old_value_str := string(old_value)
		return &old_value_str, status
	}
	return nil, status
}

// Sets multiple records.
//
// @param records Records to store.
// @param overwrite Whether to overwrite the existing value if there's a record with the same key.  If true, the existing value is overwritten by the new value.  If false, the operation is given up and an error status is returned.
// @return The result status.  If there are records avoiding overwriting, StatusDuplicationError is returned.
func (self *DBM) SetMulti(records map[string][]byte, overwrite bool) *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_set_multi(self.dbm, records, overwrite)
}

// Sets multiple records, with string data.
//
// @param records Records to store.
// @param overwrite Whether to overwrite the existing value if there's a record with the same key.  If true, the existing value is overwritten by the new value.  If false, the operation is given up and an error status is returned.
// @return The result status.  If there are records avoiding overwriting, StatusDuplicationError is set.
func (self *DBM) SetMultiStr(records map[string]string, overwrite bool) *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	rawRecords := make(map[string][]byte)
	for key, value := range records {
		rawRecords[key] = []byte(value)
	}
	return dbm_set_multi(self.dbm, rawRecords, overwrite)
}

// Removes a record of a key.
//
// @param key The key of the record.
// @return The result status.  If there's no matching record, StatusNotFoundError is returned.
func (self *DBM) Remove(key interface{}) *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_remove(self.dbm, ToByteArray(key))
}

// Removes a record and get the value.
//
// @param key The key of the record.
// @return The old value and the result status.
func (self *DBM) RemoveAndGet(key interface{}) ([]byte, *Status) {
	if self.dbm == 0 {
		return nil, NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_remove_and_get(self.dbm, ToByteArray(key))
}

// Removes a record and get the value, as a string.
//
// @param key The key of the record.
// @return The old value and the result status.
func (self *DBM) RemoveAndGetStr(key interface{}) (*string, *Status) {
	if self.dbm == 0 {
		return nil, NewStatus2(StatusPreconditionError, "not opened database")
	}
	old_value, status := dbm_remove_and_get(self.dbm, ToByteArray(key))
	if old_value != nil {
		old_value_str := string(old_value)
		return &old_value_str, status
	}
	return nil, status
}

// Removes records of keys.
//
// @param key The keys of the records.
// @return The result status.  If there are missing records, StatusNotFoundError is returned.
func (self *DBM) RemoveMulti(keys []string) *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_remove_multi(self.dbm, keys)
}

// Appends data at the end of a record of a key.
//
// @param key The key of the record.
// @param value The value to append.
// @param delim The delimiter to put after the existing record.
// @return The result status.
//
// If there's no existing record, the value is set without the delimiter.
func (self *DBM) Append(key interface{}, value interface{}, delim interface{}) *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_append(self.dbm, ToByteArray(key), ToByteArray(value), ToByteArray(delim))
}

// Appends data to multiple records.
//
// @param records Records to append.
// @param delim The delimiter to put after the existing record.
// @return The result status.
//
// If there's no existing record, the value is set without the delimiter.
func (self *DBM) AppendMulti(records map[string][]byte, delim interface{}) *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_append_multi(self.dbm, records, ToByteArray(delim))
}

// Appends data to multiple records, with string data.
//
// @param records Records to append.
// @param delim The delimiter to put after the existing record.
// @return The result status.
//
// If there's no existing record, the value is set without the delimiter.
func (self *DBM) AppendMultiStr(records map[string]string, delim interface{}) *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	rawRecords := make(map[string][]byte)
	for key, value := range records {
		rawRecords[key] = []byte(value)
	}
	return dbm_append_multi(self.dbm, rawRecords, ToByteArray(delim))
}

// Compares the value of a record and exchanges if the condition meets.
//
// @param key The key of the record.
// @param expected The expected value.  If it is nil or NilString, no existing record is expected.  If it is AnyBytes or AnyString, an existing record with any value is expacted.
// @param desired The desired value.  If it is nil or NilString, the record is to be removed.  If it is AnyBytes or AnyString, no update is done.
// @return The result status.  If the condition doesn't meet, StatusInfeasibleError is returned.
func (self *DBM) CompareExchange(
	key interface{}, expected interface{}, desired interface{}) *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	var rawExpected []byte
	if !IsNilData(expected) {
		if IsAnyData(expected) {
			rawExpected = AnyBytes
		} else {
			rawExpected = ToByteArray(expected)
		}
	}
	var rawDesired []byte
	if !IsNilData(desired) {
		if IsAnyData(desired) {
			rawDesired = AnyBytes
		} else {
			rawDesired = ToByteArray(desired)
		}
	}
	return dbm_compare_exchange(self.dbm, ToByteArray(key), rawExpected, rawDesired)
}

// Does compare-and-exchange and/or gets the old value of the record.
//
// @param key The key of the record.
// @param expected The expected value.  If it is nil or NilString, no existing record is expected.  If it is AnyBytes or AnyString, an existing record with any value is expacted.
// @param desired The desired value.  If it is nil or NilString, the record is to be removed.  If it is AnyBytes or AnyString, no update is done.
// @return The old value and the result status.  If the condition doesn't meet, the state is INFEASIBLE_ERROR.  If there's no existing record, the value is nil.
func (self *DBM) CompareExchangeAndGet(
	key interface{}, expected interface{}, desired interface{}) ([]byte, *Status) {
	if self.dbm == 0 {
		return nil, NewStatus2(StatusPreconditionError, "not opened database")
	}
	var rawExpected []byte
	if !IsNilData(expected) {
		if IsAnyData(expected) {
			rawExpected = AnyBytes
		} else {
			rawExpected = ToByteArray(expected)
		}
	}
	var rawDesired []byte
	if !IsNilData(desired) {
		if IsAnyData(desired) {
			rawDesired = AnyBytes
		} else {
			rawDesired = ToByteArray(desired)
		}
	}
	return dbm_compare_exchange_and_get(self.dbm, ToByteArray(key), rawExpected, rawDesired)
}

// Does compare-and-exchange and/or gets the old value of the record, as a string.
//
// @param key The key of the record.
// @param expected The expected value.  If it is nil or NilString, no existing record is expected.  If it is AnyBytes or AnyString, an existing record with any value is expacted.
// @param desired The desired value.  If it is nil or NilString, the record is to be removed.  If it is AnyBytes or AnyString, no update is done.
// @return The old value and the result status.  If the condition doesn't meet, the state is INFEASIBLE_ERROR.  If there's no existing record, the value is NilString.
func (self *DBM) CompareExchangeAndGetStr(
	key interface{}, expected interface{}, desired interface{}) (string, *Status) {
	rawActual, status := self.CompareExchangeAndGet(key, expected, desired)
	actual := NilString
	if rawActual != nil {
		actual = string(rawActual)
	}
	return actual, status
}

// Increments the numeric value of a record.
//
// @param key The key of the record.
// @param inc The incremental value.  If it is Int64Min, the current value is not changed and a new record is not created.
// @param init The initial value.
// @return The current value and the result status.
func (self *DBM) Increment(key interface{}, inc interface{}, init interface{}) (int64, *Status) {
	if self.dbm == 0 {
		return 0, NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_increment(self.dbm, ToByteArray(key), ToInt(inc), ToInt(init))
}

// Processes multiple records with arbitrary functions.
//
// @param keyProcPairs A list of pairs of keys and their functions.  The first parameter is the key bytes of the record.  The second parameter is the value bytes of the existing record, or nil if it the record doesn't exist.  The return value is bytes or a string to update the record value.  If the return value is nil or NilString, the record is not modified.  If the return value is RemoveBytes or RemoveString, the record is removed.
// @param writable True if the processor can edit the record.
// @return The result status.
func (self *DBM) ProcessMulti(keyProcPairs []KeyProcPair, writable bool) *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	rawPairs := make([]KeyBytesProcPair, 0, len(keyProcPairs))
	for _, pair := range keyProcPairs {
		rawPair := KeyBytesProcPair{ToByteArray(pair.Key), pair.Proc}
		rawPairs = append(rawPairs, rawPair)
	}
	return dbm_process_multi(self.dbm, rawPairs, writable)
}

// Compares the values of records and exchanges if the condition meets.
//
// @param expected A sequence of pairs of the record keys and their expected values.  If the value is nil, no existing record is expected.  If the value is AnyBytes, an existing record with any value is expacted.
// @param desired A sequence of pairs of the record keys and their desired values.  If the value is nil, the record is to be removed.
// @return The result status.  If the condition doesn't meet, StatusInfeasibleError is returned.
func (self *DBM) CompareExchangeMulti(expected []KeyValuePair, desired []KeyValuePair) *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_compare_exchange_multi(self.dbm, expected, desired)
}

// Compares the values of records and exchanges if the condition meets, using string data.
//
// @param expected A sequence of pairs of the record keys and their expected values.  If the value is NilString, no existing record is expected.  If the value is AnyString, an existing record with any value is expacted.
// @param desired A sequence of pairs of the record keys and their desired values.  If the value is NilString, the record is to be removed.
// @return The result status.  If the condition doesn't meet, StatusInfeasibleError is returned.
func (self *DBM) CompareExchangeMultiStr(
	expected []KeyValueStrPair, desired []KeyValueStrPair) *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
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
	return dbm_compare_exchange_multi(self.dbm, rawExpected, rawDesired)
}

// Changes the key of a record.
//
// @param old_key The old key of the record.
// @param new_key The new key of the record.
// @param overwrite Whether to overwrite the existing record of the new key.
// @param copying Whether to retain the record of the old key.
// @return The result status.  If there's no matching record to the old key, NOT_FOUND_ERROR is returned.  If the overwrite flag is false and there is an existing record of the new key, DUPLICATION ERROR is returned.
//
// This method is done atomically by ProcessMulti.  The other threads observe that the record has either the old key or the new key.  No intermediate states are observed.
func (self *DBM) Rekey(oldKey interface{}, newKey interface{},
	overwrite bool, copying bool) *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_rekey(self.dbm, ToByteArray(oldKey), ToByteArray(newKey), overwrite, copying)
}

// Gets the first record and removes it.
//
// @return The key and the value of the first record, and the result status.
func (self *DBM) PopFirst() ([]byte, []byte, *Status) {
	if self.dbm == 0 {
		return nil, nil, NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_pop_first(self.dbm)
}

// Gets the first record as strings and removes it.
//
// @return The key and the value of the first record, and the result status.
func (self *DBM) PopFirstStr() (string, string, *Status) {
	if self.dbm == 0 {
		return "", "", NewStatus2(StatusPreconditionError, "not opened database")
	}
	key, value, status := dbm_pop_first(self.dbm)
	if status.code == StatusSuccess {
		return string(key), string(value), status
	}
	return "", "", status
}

// Adds a record with a key of the current timestamp.
//
// @param value The value of the record.
// @param wtime The current wall time used to generate the key.  If it is None, the system clock is used.
// @return The result status.
//
// The key is generated as an 8-bite big-endian binary string of the timestamp.  If there is an existing record matching the generated key, the key is regenerated and the attempt is repeated until it succeeds.
func (self *DBM) PushLast(value interface{}, wtime float64) *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_push_last(self.dbm, ToByteArray(value), wtime)
}

// Processes each and every record in the database with an arbitrary function.
//
// @param proc The function to process a record.  The first parameter is the key bytes of the record.  The second parameter is the value bytes of the existing record, or nil if it the record doesn't exist.  The return value is bytes or a string to update the record value.  If the return value is nil or NilString, the record is not modified.  If the return value is RemoveBytes or RemoveString, the record is removed.
// @param writable True if the processor can edit the record.
// @return The result status.
//
// The given function is called repeatedly for each record.  It is also called once before the iteration and once after the iteration with both the key and the value being nil.
func (self *DBM) ProcessEach(proc RecordProcessor, writable bool) *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_process_each(self.dbm, proc, writable)
}

// Gets the number of records.
//
// @return The number of records and the result status.
func (self *DBM) Count() (int64, *Status) {
	if self.dbm == 0 {
		return -1, NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_count(self.dbm)
}

// Gets the number of records, in a simple way.
//
// @return The number of records or -1 on failure.
func (self *DBM) CountSimple() int64 {
	if self.dbm == 0 {
		return -1
	}
	count, status := dbm_count(self.dbm)
	if status.code == StatusSuccess {
		return count
	}
	return -1
}

// Gets the current file size of the database.
//
// @return The current file size of the database and the result status.
func (self *DBM) GetFileSize() (int64, *Status) {
	if self.dbm == 0 {
		return -1, NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_get_file_size(self.dbm)
}

// Gets the current file size of the database, in a simple way.
//
// @return The current file size of the database, or -1 on failure.
func (self *DBM) GetFileSizeSimple() int64 {
	if self.dbm == 0 {
		return -1
	}
	file_size, status := dbm_get_file_size(self.dbm)
	if status.code == StatusSuccess {
		return file_size
	}
	return -1
}

// Gets the path of the database file.
//
// @return The path of the database file and the result status.
func (self *DBM) GetFilePath() (string, *Status) {
	if self.dbm == 0 {
		return "", NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_get_file_path(self.dbm)
}

// Gets the path of the database file, in a simple way.
//
// @return The path of the database file, or an empty string on failure.
func (self *DBM) GetFilePathSimple() string {
	if self.dbm == 0 {
		return ""
	}
	path, status := dbm_get_file_path(self.dbm)
	if status.code == StatusSuccess {
		return path
	}
	return ""
}

// Gets the timestamp in seconds of the last modified time.
//
// @return The timestamp in seconds of the last modified time and the result status.
func (self *DBM) GetTimestamp() (float64, *Status) {
	if self.dbm == 0 {
		return -1, NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_get_timestamp(self.dbm)
}

// Gets the timestamp in seconds of the last modified time, in a simple way.
//
// @return The timestamp in seconds of the last modified, or -1 on failure.
func (self *DBM) GetTimestampSimple() float64 {
	if self.dbm == 0 {
		return -1
	}
	timestamp, status := dbm_get_timestamp(self.dbm)
	if status.code == StatusSuccess {
		return timestamp
	}
	return -1
}

// Removes all records.
//
// @return The result status.
func (self *DBM) Clear() *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_clear(self.dbm)
}

// Rebuilds the entire database.
//
// @param params Optional parameters.  If it is nil, it is ignored.
// @return The result status.
//
// The optional parameters are the same as the Open method.  Omitted tuning parameters are kept the same or implicitly optimized.
//
// In addition, HashDBM, TreeDBM, and SkipDBM supports the following parameters.
//
// - skip_broken_records (bool): If true, the operation continues even if there are broken records which can be skipped.
// - sync_hard (bool): If true, physical synchronization with the hardware is done before finishing the rebuilt file.
func (self *DBM) Rebuild(params map[string]string) *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_rebuild(self.dbm, params)
}

// Checks whether the database should be rebuilt.
//
// @return The result decision and the result status.  The decision is true to be optimized or false with no necessity.
func (self *DBM) ShouldBeRebuilt() (bool, *Status) {
	if self.dbm == 0 {
		return false, NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_should_be_rebuilt(self.dbm)
}

// Checks whether the database should be rebuilt, in a simple way.
//
// @return True to be optimized or false with no necessity.
func (self *DBM) ShouldBeRebuiltSimple() bool {
	if self.dbm == 0 {
		return false
	}
	tobe, status := dbm_should_be_rebuilt(self.dbm)
	if status.code == StatusSuccess {
		return tobe
	}
	return false
}

// Synchronizes the content of the database to the file system.
//
// @param hard True to do physical synchronization with the hardware or false to do only logical synchronization with the file system.
// @param params Optional parameters.  If it is nil, it is ignored.
// @return The result status.
//
// Only SkipDBM uses the optional parameters.  The "merge" parameter specifies paths of databases to merge, separated by colon.  The "reducer" parameter specifies the reducer to apply to records of the same key.  "ReduceToFirst", "ReduceToSecond", "ReduceToLast", etc are supported.
func (self *DBM) Synchronize(hard bool, params map[string]string) *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_synchronize(self.dbm, hard, params)
}

// Copies the content of the database file to another file.
//
// @param destPath A path to the destination file.
// @param syncHard True to do physical synchronization with the hardware.
// @return The result status.
func (self *DBM) CopyFileData(destPath string, syncHard bool) *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_copy_file_data(self.dbm, destPath, syncHard)
}

// Exports all records to another database.
//
// @param destDBM The destination database.
// @return The result status.
func (self *DBM) Export(destDBM *DBM) *Status {
	if self.dbm == 0 || destDBM.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	return dbm_export(self.dbm, destDBM.dbm)
}

// Exports all records of a database to a flat record file.
//
// @param file: The file object to write records in.
// @return The result status.
//
// A flat record file contains a sequence of binary records without any high level structure so it is useful as a intermediate file for data migration.
func (self *DBM) ExportToFlatRecords(destFile *File) *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	if destFile.file == 0 {
		return NewStatus2(StatusPreconditionError, "not opened file")
	}
	return dbm_export_to_flat_records(self.dbm, destFile.file)
}

// Imports records to a database from a flat record file.
//
// @param file The file object to read records from.
// @return The result status.
func (self *DBM) ImportFromFlatRecords(srcFile *File) *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	if srcFile.file == 0 {
		return NewStatus2(StatusPreconditionError, "not opened file")
	}
	return dbm_import_from_flat_records(self.dbm, srcFile.file)
}

// Exports the keys of all records as lines to a text file.
//
// @param file The file object to write keys in.
// @return The result status.
//
// As the exported text file is smaller than the database file, scanning the text file by the search method is often faster than scanning the whole database.
func (self *DBM) ExportKeysAsLines(destFile *File) *Status {
	if self.dbm == 0 {
		return NewStatus2(StatusPreconditionError, "not opened database")
	}
	if destFile.file == 0 {
		return NewStatus2(StatusPreconditionError, "not opened file")
	}
	return dbm_export_keys_as_lines(self.dbm, destFile.file)
}

// Inspects the database.
//
// return A map of property names and their values.
func (self *DBM) Inspect() map[string]string {
	if self.dbm == 0 {
		return nil
	}
	records := dbm_inspect(self.dbm)
	return records
}

// Checks whether the database is open.
//
// @return True if the database is open, or false if not.
func (self *DBM) IsOpen() bool {
	if self.dbm == 0 {
		return false
	}
	return true
}

// Checks whether the database is writable.
//
// @return True if the database is writable, or false if not.
func (self *DBM) IsWritable() bool {
	if self.dbm == 0 {
		return false
	}
	return dbm_is_writable(self.dbm)
}

// Checks whether the database condition is healthy.
//
// @return True if the database condition is healthy, or false if not.
func (self *DBM) IsHealthy() bool {
	if self.dbm == 0 {
		return false
	}
	return dbm_is_healthy(self.dbm)
}

// Checks whether ordered operations are supported.
//
// @return True if ordered operations are supported, or false if not.
func (self *DBM) IsOrdered() bool {
	if self.dbm == 0 {
		return false
	}
	return dbm_is_ordered(self.dbm)
}

// Searches the database and get keys which match a pattern.
//
// @param mode The search mode.  "contain" extracts keys containing the pattern.  "begin" extracts keys beginning with the pattern.  "end" extracts keys ending with the pattern.  "regex" extracts keys partially matches the pattern of a regular expression.  "edit" extracts keys whose edit distance to the UTF-8 pattern is the least.  "editbin" extracts keys whose edit distance to the binary pattern is the least.  "containcase", "containword", and "containcaseword" extract keys considering case and word boundary.  Ordered databases support "upper" and "lower" which extract keys whose positions are upper/lower than the pattern. "upperinc" and "lowerinc" are their inclusive versions.
// @param pattern The pattern for matching.
// @param capacity The maximum records to obtain.  0 means unlimited.
// @return A list of keys matching the condition.
func (self *DBM) Search(mode string, pattern string, capacity int) []string {
	if self.dbm == 0 {
		return make([]string, 0)
	}
	return dbm_search(self.dbm, mode, pattern, capacity)
}

// Makes an iterator for each record.
//
// @return The iterator for each record.
//
// Every iterator should be destructed explicitly by the "Destruct" method.
func (self *DBM) MakeIterator() *Iterator {
	if self.dbm == 0 {
		return &Iterator{0}
	}
	iter := dbm_make_iterator(self.dbm)
	return &Iterator{iter}
}

// Makes a channel to read each records.
//
// @return the channel to read each records.  All values should be read from the channel to avoid resource leak.
func (self *DBM) Each() <-chan KeyValuePair {
	chan_record := make(chan KeyValuePair)
	reader := func(chan_send chan<- KeyValuePair) {
		defer close(chan_record)
		iter := self.MakeIterator()
		defer iter.Destruct()
		if !iter.First().IsOK() {
			return
		}
		for {
			key, value, status := iter.Get()
			if !status.IsOK() {
				break
			}
			chan_send <- KeyValuePair{key, value}
			if !iter.Next().IsOK() {
				return
			}
		}
	}
	go reader(chan_record)
	return chan_record
}

// Makes a channel to read each records, as strings.
//
// @return the channel to read each records.  All values should be read from the channel to avoid resource leak.
func (self *DBM) EachStr() <-chan KeyValueStrPair {
	chan_record := make(chan KeyValueStrPair)
	reader := func(chan_send chan<- KeyValueStrPair) {
		defer close(chan_record)
		iter := self.MakeIterator()
		defer iter.Destruct()
		if !iter.First().IsOK() {
			return
		}
		for {
			key, value, status := iter.GetStr()
			if !status.IsOK() {
				break
			}
			chan_send <- KeyValueStrPair{key, value}
			if !iter.Next().IsOK() {
				return
			}
		}
	}
	go reader(chan_record)
	return chan_record
}

// Restores a broken database as a new healthy database.
//
// @param old_file_path The path of the broken database.
// @param new_file_path The path of the new database to be created.
// @param class_name The name of the database class.  If it is nil or empty, the class is guessed from the file extension.
// @param end_offset The exclusive end offset of records to read.  Negative means unlimited.  0 means the size when the database is synched or closed properly.  Using a positive value is not meaningful if the number of shards is more than one.
// @param cipherKey The encryption key for cipher compressors.
// @return The result status.
func RestoreDatabase(
	oldFilePath string, newFilePath string, className string,
	endOffset int64, cipherKey string) *Status {
	return dbm_restore_database(oldFilePath, newFilePath, className, endOffset, cipherKey)
}

// END OF FILE
