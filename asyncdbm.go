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

// END OF FILE
