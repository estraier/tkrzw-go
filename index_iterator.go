/*************************************************************************************************
 * Index iterator interface
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

// Iterator for each record of the secondary index.
//
// An iterator is made by the "MakeIerator" method of Index.  Every unused iterator object should be destructed explicitly by the "Destruct" method to free resources.
type IndexIterator struct {
	// Pointer to the internal object.
	iter uintptr
}

// Releases the resource explicitly.
func (self *IndexIterator) Destruct() {
	if self.iter == 0 {
		return
	}
	index_iter_free(self.iter)
	self.iter = 0
}

// Makes a string representing the iterator.
//
// @return The string representing the iterator.
func (self *IndexIterator) String() string {
	if self.iter == 0 {
		return fmt.Sprintf("#<tkrzw.Iterator:%p:destructed>", &self)
	}
	return fmt.Sprintf("#<tkrzw.Iterator:%p>", &self)
}

// Initializes the iterator to indicate the first record.
func (self *IndexIterator) First() {
	if self.iter == 0 {
		return
	}
	index_iter_first(self.iter)
}

// Initializes the iterator to indicate the last record.
func (self *IndexIterator) Last() {
	if self.iter == 0 {
		return
	}
	index_iter_last(self.iter)
}

// Initializes the iterator to indicate a specific range.
//
// @param key The key of the lower bound.
// @param value The value of the lower bound.
func (self *IndexIterator) Jump(key interface{}, value interface{}) {
	if self.iter == 0 {
		return
	}
	index_iter_jump(self.iter, ToByteArray(key), ToByteArray(value))
}

// Gets the key and the value of the current record of the iterator.
//
// @return The key and the value of the current record, and a boolean status.
func (self *IndexIterator) Get() ([]byte, []byte, bool) {
	if self.iter == 0 {
		return nil, nil, false
	}
	return index_iter_get(self.iter)
}

// Gets the key and the value of the current record of the iterator, as strings.
//
// @return The key and the value of the current record, and a boolean status.
func (self *IndexIterator) GetStr() (string, string, bool) {
	if self.iter == 0 {
		return "", "", false
	}
	key, value, status := index_iter_get(self.iter)
	if status {
		return string(key), string(value), true
	}
	return "", "", false
}

// Moves the iterator to the next record.
func (self *IndexIterator) Next() {
	if self.iter == 0 {
		return
	}
	index_iter_next(self.iter)
}

// Moves the iterator to the previous record.
func (self *IndexIterator) Previous() {
	if self.iter == 0 {
		return
	}
	index_iter_previous(self.iter)
}

// END OF FILE
