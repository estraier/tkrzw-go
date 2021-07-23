/*************************************************************************************************
 * Iterator interface
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

// Iterator for each record.
//
// An iterator is made by the "MakeIerator" method of DBM.  Every unused iterator object should be destructed explicitly by the "Destruct" method to free resources.
type Iterator struct {
	// Pointer to the internal object.
	iter uintptr
}

// Releases the resource explicitly.
func (self *Iterator) Destruct() {
	if self.iter == 0 {
		return
	}
	dbm_iter_free(self.iter)
	self.iter = 0
}

// Makes a string representing the iterator.
//
// @return The string representing the iterator.
func (self *Iterator) String() string {
	if self.iter == 0 {
		return fmt.Sprintf("#<tkrzw.Iter:%p:destructed>", &self)
	}
	key, status := dbm_iter_get_key_esc(self.iter)
	if status.code == StatusSuccess {
		return fmt.Sprintf("#<tkrzw.Iter:%s>", key)
	}
	return fmt.Sprintf("#<tkrzw.Iter:%p:unlocated>", &self)
}

// Initializes the iterator to indicate the first record.
//
// @return The result status.
//
// Even if there's no record, the operation doesn't fail.
func (self *Iterator) First() *Status {
	if self.iter == 0 {
		return NewStatus2(StatusPreconditionError, "destructed Iterator")
	}
	return dbm_iter_first(self.iter)
}

// Initializes the iterator to indicate the last record.
//
// @return The result status.
//
// Even if there's no record, the operation doesn't fail.  This method is suppoerted only by ordered databases.
func (self *Iterator) Last() *Status {
	if self.iter == 0 {
		return NewStatus2(StatusPreconditionError, "destructed Iterator")
	}
	return dbm_iter_last(self.iter)
}

// Initializes the iterator to indicate a specific record.
//
// @param key The key of the record to look for.
// @return The result status.
//
// Ordered databases can support "lower bound" jump; If there's no record with the same key, the iterator refers to the first record whose key is greater than the given key.  The operation fails with unordered databases if there's no record with the same key.
func (self *Iterator) Jump(key interface{}) *Status {
	if self.iter == 0 {
		return NewStatus2(StatusPreconditionError, "destructed Iterator")
	}
	return dbm_iter_jump(self.iter, ToByteArray(key))
}

// Initializes the iterator to indicate the last record whose key is lower than a given key.
//
// @param key The key to compare with.
// @param inclusive If true, the considtion is inclusive: equal to or lower than the key.
// @return The result status.
//
// Even if there's no matching record, the operation doesn't fail.  This method is suppoerted only by ordered databases.
func (self *Iterator) JumpLower(key interface{}, inclusive bool) *Status {
	if self.iter == 0 {
		return NewStatus2(StatusPreconditionError, "destructed Iterator")
	}
	return dbm_iter_jump_lower(self.iter, ToByteArray(key), inclusive)
}

// Initializes the iterator to indicate the first record whose key is upper than a given key.
//
// @param key The key to compare with.
// @param inclusive If true, the considtion is inclusive: equal to or upper than the key.
// @return The result status.
//
// Even if there's no matching record, the operation doesn't fail.  This method is suppoerted only by ordered databases.
func (self *Iterator) JumpUpper(key interface{}, inclusive bool) *Status {
	if self.iter == 0 {
		return NewStatus2(StatusPreconditionError, "destructed Iterator")
	}
	return dbm_iter_jump_upper(self.iter, ToByteArray(key), inclusive)
}

// Moves the iterator to the next record.
//
// @return The result status.
//
// If the current record is missing, the operation fails.  Even if there's no next record, the operation doesn't fail.
func (self *Iterator) Next() *Status {
	if self.iter == 0 {
		return NewStatus2(StatusPreconditionError, "destructed Iterator")
	}
	return dbm_iter_next(self.iter)
}

// Moves the iterator to the previous record.
//
// @return The result status.
//
// If the current record is missing, the operation fails.  Even if there's no previous record, the operation doesn't fail.  This method is suppoerted only by ordered databases.
func (self *Iterator) Previous() *Status {
	if self.iter == 0 {
		return NewStatus2(StatusPreconditionError, "destructed Iterator")
	}
	return dbm_iter_previous(self.iter)
}

// Gets the key and the value of the current record of the iterator.
//
// @return The key and the value of the current record, and the result status.
func (self *Iterator) Get() ([]byte, []byte, *Status) {
	if self.iter == 0 {
		return nil, nil, NewStatus2(StatusPreconditionError, "destructed Iterator")
	}
	return dbm_iter_get(self.iter)
}

// Gets the key and the value of the current record of the iterator, as strings.
//
// @return The key and the value of the current record, and the result status.
func (self *Iterator) GetStr() (string, string, *Status) {
	if self.iter == 0 {
		return "", "", NewStatus2(StatusPreconditionError, "destructed Iterator")
	}
	key, value, status := dbm_iter_get(self.iter)
	if status.code == StatusSuccess {
		return string(key), string(value), status
	}
	return "", "", status
}

// Gets the key of the current record.
//
// @return The key of the current record and the result status.
func (self *Iterator) GetKey() ([]byte, *Status) {
	if self.iter == 0 {
		return nil, NewStatus2(StatusPreconditionError, "destructed Iterator")
	}
	return dbm_iter_get_key(self.iter)
}

// Gets the key of the current record, as a string.
//
// @return The key of the current record and the result status.
func (self *Iterator) GetKeyStr() (string, *Status) {
	if self.iter == 0 {
		return "", NewStatus2(StatusPreconditionError, "destructed Iterator")
	}
	key, status := dbm_iter_get_key(self.iter)
	if status.code == StatusSuccess {
		return string(key), status
	}
	return "", status
}

// Gets the value of the current record.
//
// @return The value of the current record and the result status.
func (self *Iterator) GetValue() ([]byte, *Status) {
	if self.iter == 0 {
		return nil, NewStatus2(StatusPreconditionError, "destructed Iterator")
	}
	return dbm_iter_get_value(self.iter)
}

// Gets the value of the current record, as a string.
//
// @return The value of the current record and the result status.
func (self *Iterator) GetValueStr() (string, *Status) {
	if self.iter == 0 {
		return "", NewStatus2(StatusPreconditionError, "destructed Iterator")
	}
	value, status := dbm_iter_get_value(self.iter)
	if status.code == StatusSuccess {
		return string(value), status
	}
	return "", status
}

// END OF FILE
