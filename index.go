/*************************************************************************************************
 * Secondary index interface
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

// Secondary index.
//
// All operations except for "Open" and "Close" are thread-safe; Multiple threads can access the same index concurrently.  You can specify a data structure when you call the "Open" method.  Every opened index must be closed explicitly by the "Close" method to avoid data corruption.
type Index struct {
	// Pointer to the internal object.
	index uintptr
}

// Makes a new Index object.
//
// @return The pointer to the created index object.
func NewIndex() *Index {
	return &Index{0}
}

// Releases the resource explicitly.
//
// The index is closed implicitly if it has not been closed.  As long as you close the index explicitly, you don't have to call this method.
func (self *Index) Destruct() {
	if self.index == 0 {
		return
	}
	index_close(self.index)
	self.index = 0
}

// Makes a string representing the index.
//
// @return The string representing the index.
func (self *Index) String() string {
	if self.index == 0 {
		return fmt.Sprintf("#<tkrzw.Index:%p:unopened>", &self)
	}
	path := index_get_file_path(self.index)
	count := index_count(self.index)
	return fmt.Sprintf("#<tkrzw.Index:%s:%d>", path, count)
}

// Opens an index file.
//
// @param path A path of the file.
// @param writable If true, the file is writable.  If false, it is read-only.
// @param params Optional parameters.  If it is nil, it is ignored.
// @return The result status.
//
// If the path is empty, BabyDBM is used internally, which is equivalent to using the MemIndex class.  If the path ends with ".tkt", TreeDBM is used internally, which is equivalent to using the FileIndex class.  If the key comparator of the tuning parameter is not set, PairLexicalKeyComparator is set implicitly.  Other compatible key comparators are PairLexicalCaseKeyComparator, PairDecimalKeyComparator, PairHexadecimalKeyComparator, PairRealNumberKeyComparator, PairSignedBigEndianKeyComparator, and PairFloatBigEndianKeyComparator.  Other options can be specified as with DBM::Open.
func (self *Index) Open(path string, writable bool, params map[string]string) *Status {
	if self.index != 0 {
		return NewStatus2(StatusPreconditionError, "opened index")
	}
	index, status := index_open(path, writable, params)
	if status.code == StatusSuccess {
		self.index = index
	}
	return status
}

// Closes the index file.
//
// @return The result status.
func (self *Index) Close() *Status {
	if self.index == 0 {
		return NewStatus2(StatusPreconditionError, "not opened index")
	}
	status := index_close(self.index)
	self.index = 0
	return status
}

// Checks whether a record exists in the index.
//
// @param key The key of the record.
// @param value The value of the record.
// @return True if the record exists, or false if not.
func (self *Index) Check(key interface{}, value interface{}) bool {
	if self.index == 0 {
		return false
	}
	return index_check(self.index, ToByteArray(key), ToByteArray(value))
}

// Gets all values of records of a key.
//
// @param key The key of the record.
// @param max The maximum number of values to get.  0 means unlimited.
func (self *Index) GetValues(key interface{}, max int) [][]byte {
	if self.index == 0 {
		return make([][]byte, 0)
	}
	return index_get_values(self.index, ToByteArray(key), max)
}

// Gets all values of records of a key, as strings.
//
// @param key The key of the record.
// @param max The maximum number of values to get.  0 means unlimited.
func (self *Index) GetValuesStr(key interface{}, max int) []string {
	if self.index == 0 {
		return make([]string, 0)
	}
	return index_get_values_str(self.index, ToByteArray(key), max)
}

// Adds a record.
//
// @param key The key of the record.  This can be an arbitrary expression to search the index.
// @param value The value of the record.  This should be a primary value of another database.
// @return The result status.
func (self *Index) Add(key interface{}, value interface{}) *Status {
	if self.index == 0 {
		return NewStatus2(StatusPreconditionError, "not opened index")
	}
	return index_add(self.index, ToByteArray(key), ToByteArray(value))
}

// Removes a record.
//
// @param key The key of the record.  This can be an arbitrary expression to search the index.
// @param value The value of the record.  This should be a primary value of another database.
// @return The result status.
func (self *Index) Remove(key interface{}, value interface{}) *Status {
	if self.index == 0 {
		return NewStatus2(StatusPreconditionError, "not opened index")
	}
	return index_remove(self.index, ToByteArray(key), ToByteArray(value))
}

// Gets the number of records.
//
// @return The number of records, or -1 on failure.
func (self *Index) Count() int64 {
	if self.index == 0 {
		return 0
	}
	return index_count(self.index)
}

// Gets the path of the index file.
//
// @return The path of the index file, or an empty string on failure.
func (self *Index) GetFilePath() string {
	if self.index == 0 {
		return ""
	}
	return index_get_file_path(self.index)
}

// Removes all records.
//
// @return The result status.
func (self *Index) Clear() *Status {
	if self.index == 0 {
		return NewStatus2(StatusPreconditionError, "not opened index")
	}
	return index_clear(self.index)
}

// Rebuilds the entire index.
//
// @return The result status.
func (self *Index) Rebuild() *Status {
	if self.index == 0 {
		return NewStatus2(StatusPreconditionError, "not opened index")
	}
	return index_rebuild(self.index)
}

// Synchronizes the content of the index to the file system.
//
// @param hard True to do physical synchronization with the hardware or false to do only logical synchronization with the file system.
// @return The result status.
func (self *Index) Synchronize(hard bool) *Status {
	if self.index == 0 {
		return NewStatus2(StatusPreconditionError, "not opened index")
	}
	return index_synchronize(self.index, hard)
}

// Checks whether the index is open.
//
// @return True if the index is open, or false if not.
func (self *Index) IsOpen() bool {
	if self.index == 0 {
		return false
	}
	return true
}

// Checks whether the index is writable.
//
// @return True if the index is writable, or false if not.
func (self *Index) IsWritable() bool {
	if self.index == 0 {
		return false
	}
	return index_is_writable(self.index)
}

// Makes an iterator for each record.
//
// @return The iterator for each record.
//
// Every iterator should be destructed explicitly by the "Destruct" method.
func (self *Index) MakeIterator() *IndexIterator {
	if self.index == 0 {
		return &IndexIterator{0}
	}
	iter := index_make_iterator(self.index)
	return &IndexIterator{iter}
}

// END OF FILE
