/*************************************************************************************************
 * Future interface
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

// Future containing a status object and extra data.
//
// Future objects are made by methods of AsyncDBM.  Every future object should be destroyed by the "Destruct" method or the "Get" method family to free resources.

type Future struct {
	// Pointer to the internal object.
	future uintptr
}

// Destructs the object and releases resources.
func (self *Future) Destruct() {
	if self.future == 0 {
		return
	}
	future_free(self.future)
	self.future = 0
}

// Makes a string representing the future.
//
// @return The string representing the future.
func (self *Future) String() string {
	if self.future == 0 {
		return fmt.Sprintf("#<tkrzw.Future:%p:destructed>", &self)
	}
	return fmt.Sprintf("#<tkrzw.Future:%p:0x%x>", &self, self.future)
}

// Waits for the operation to be done.
//
// @param timeout The waiting time in seconds.  If it is negative, no timeout is set.
// @return True if the operation has done.  False if timeout occurs.
func (self *Future) Wait(timeout float32) bool {
	if self.future == 0 {
		return false
	}
	return future_wait(self.future, timeout)
}

// Gets the status of the operation.
//
// @return The result status.
//
// The internal resource is released by this method.  "Wait" and "Get" faminly cannot be called after calling this method.
func (self *Future) Get() *Status {
	if self.future == 0 {
		return NewStatus2(StatusPreconditionError, "destructed object")
	}
	status := future_get(self.future)
	self.future = 0
	return status
}

// Gets the extra byte array data and the status of the operation.
//
// @return The bytes value of the matching record and the result status.
//
// The internal resource is released by this method.  "Wait" and "Get" faminly cannot be called after calling this method.
func (self *Future) GetBytes() ([]byte, *Status) {
	if self.future == 0 {
		return nil, NewStatus2(StatusPreconditionError, "destructed object")
	}
	value, status := future_get_bytes(self.future)
	self.future = 0
	return value, status
}

// Gets the extra string data and the status of the operation.
//
// @return The string value and the result status.
//
// The internal resource is released by this method.  "Wait" and "Get" faminly cannot be called after calling this method.
func (self *Future) GetStr() (string, *Status) {
	if self.future == 0 {
		return "", NewStatus2(StatusPreconditionError, "destructed object")
	}
	value, status := future_get_str(self.future)
	self.future = 0
	return value, status
}

// Gets the extra byte array map and the status of the operation.
//
// @return A byte array map and the result status.
//
// The internal resource is released by this method.  "Wait" and "Get" faminly cannot be called after calling this method.
func (self *Future) GetMap() (map[string][]byte, *Status) {
	if self.future == 0 {
		return nil, NewStatus2(StatusPreconditionError, "destructed object")
	}
	value, status := future_get_map(self.future)
	self.future = 0
	return value, status
}

// Gets the extra string map and the status of the operation.
//
// @return A string map and the result status.
//
// The internal resource is released by this method.  "Wait" and "Get" faminly cannot be called after calling this method.
func (self *Future) GetMapStr() (map[string]string, *Status) {
	if self.future == 0 {
		return nil, NewStatus2(StatusPreconditionError, "destructed object")
	}
	value, status := future_get_map_str(self.future)
	self.future = 0
	return value, status
}

// END OF FILE
