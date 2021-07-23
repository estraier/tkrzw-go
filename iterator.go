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
	fmt.Println("Bye")
}

// END OF FILE
