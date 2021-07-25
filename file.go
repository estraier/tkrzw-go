/*************************************************************************************************
 * Generic file implementation
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

// Generic file implementation.
//
// All operations except for "Open" and "Close" are thread-safe; Multiple threads can access the same file concurrently.  You can specify a concrete class when you call the "Open" method.  Every opened file must be closed explicitly by the "Close" method to avoid data corruption.
type File struct {
	// Pointer to the internal object.
	file uintptr
}

// Makes a new file object.
//
// @return The pointer to the created file object.
func NewFile() *File {
	return &File{0}
}

// Makes a string representing the file.
//
// @return The string representing the file.
func (self *File) String() string {
	if self.file == 0 {
		return fmt.Sprintf("#<tkrzw.File:%p:unopened>", &self)
	}
	path, _ := file_get_path(self.file)
	size, _ := file_get_size(self.file)
	return fmt.Sprintf("#<tkrzw.File:%s:%d>", path, size)
}

// Opens a file.
//
// @param path A path of the file.
// @param writable If true, the file is writable.  If false, it is read-only.
// @param params Optional parameters.
// @return The result status.
//
// The optional parameters can include options for the file opening operation.
//
// - truncate (bool): True to truncate the file.
// - no_create (bool): True to omit file creation.
// - no_wait (bool): True to fail if the file is locked by another process.
// - no_lock (bool): True to omit file locking.
//
// The optional parameter "file" specifies the internal file implementation class.  The default file class is "MemoryMapAtomicFile".  The other supported classes are "StdFile", "MemoryMapAtomicFile", "PositionalParallelFile", and "PositionalAtomicFile".
//
// For the file "PositionalParallelFile" and "PositionalAtomicFile", these optional parameters are supported.
//
// - block_size (int): The block size to which all blocks should be aligned.
// - access_options (str): Values separated by colon.  "direct" for direct I/O.  "sync" for synchrnizing I/O, "padding" for file size alignment by padding, "pagecache" for the mini page cache in the process.
func (self *File) Open(path string, writable bool, params string) *Status {
	if self.file != 0 {
		return NewStatus2(StatusPreconditionError, "opened file")
	}
	file, status := file_open(path, writable, params)
	if status.code == StatusSuccess {
		self.file = file
	}
	return status
}

// Closes the file.
//
// @return The result status.
func (self *File) Close() *Status {
	if self.file == 0 {
		return NewStatus2(StatusPreconditionError, "not opened file")
	}
	status := file_close(self.file)
	self.file = 0
	return status
}

// Reads data.
//
// @param off The offset of a source region.
// @param size The size to be read.
// @return The bytes value of the read data and the result status.
func (self *File) Read(off int64, size int64) ([]byte, *Status) {
	if self.file == 0 {
		return nil, NewStatus2(StatusPreconditionError, "not opened file")
	}
	return file_read(self.file, off, size)
}

// Reads data as a string.
//
// @param off The offset of a source region.
// @param size The size to be read.
// @return The string value of the read data and the result status.
func (self *File) ReadStr(off int64, size int64) (string, *Status) {
	if self.file == 0 {
		return "", NewStatus2(StatusPreconditionError, "not opened file")
	}
	data, status := file_read(self.file, off, size)
	if status.code == StatusSuccess {
		return string(data), status
	}
	return "", status
}

// Writes data.
//
// @param off The offset of the destination region.
// @param data The data to write.
// @return The result status.
func (self *File) Write(off int64, data interface{}) *Status {
	if self.file == 0 {
		return NewStatus2(StatusPreconditionError, "not opened file")
	}
	return file_write(self.file, off, ToByteArray(data))
}

// Appends data at the end of the file.
//
// @param data The data to write.
// @return The offset at which the data has been put, and the result status.
func (self *File) Append(data interface{}) (int64, *Status) {
	if self.file == 0 {
		return 0, NewStatus2(StatusPreconditionError, "not opened file")
	}
	return file_append(self.file, ToByteArray(data))
}

// Truncates the file.
//
// @param size The new size of the file.
// @return The result status.
//
// If the file is shrunk, data after the new file end is discarded.  If the file is expanded, null codes are filled after the old file end.
func (self *File) Truncate(size int64) *Status {
	if self.file == 0 {
		return NewStatus2(StatusPreconditionError, "not opened file")
	}
	return file_truncate(self.file, size)
}

// Synchronizes the content of the file to the file system.
//
// @param hard True to do physical synchronization with the hardware or false to do only logical synchronization with the file system.
// @param off The offset of the region to be synchronized.
// @param size The size of the region to be synchronized.  If it is zero, the length to the end of file is specified.
// @return The result status.
//
// The pysical file size can be larger than the logical size in order to improve performance by reducing frequency of allocation.  Thus, you should call this function before accessing the file with external tools.
func (self *File) Synchronize(hard bool, off int64, size int64) *Status {
	if self.file == 0 {
		return NewStatus2(StatusPreconditionError, "not opened file")
	}
	return file_synchronize(self.file, hard, off, size)
}

// Gets the size of the file.
//
// @return The size of the file and the result status.
func (self *File) GetSize() (int64, *Status) {
	if self.file == 0 {
		return 0, NewStatus2(StatusPreconditionError, "not opened file")
	}
	return file_get_size(self.file)
}

// Gets the path of the file.
//
// @return The path of the file and the result status.
func (self *File) GetPath() (string, *Status) {
	if self.file == 0 {
		return "", NewStatus2(StatusPreconditionError, "not opened file")
	}
	return file_get_path(self.file)
}

// Searches the file and get lines which match a pattern.
//
// @param mode The search mode.  "contain" extracts lines containing the pattern.  "begin" extracts lines beginning with the pattern.  "end" extracts lines ending with the pattern.  "regex" extracts lines partially matches the pattern of a regular expression.  "edit" extracts lines whose edit distance to the UTF-8 pattern is the least.  "editbin" extracts lines whose edit distance to the binary pattern is the least.
// @param pattern The pattern for matching.
// @param capacity The maximum records to obtain.  0 means unlimited.
// @return A list of lines matching the condition.
func (self *File) Search(mode string, pattern string, capacity int) []string {
	if self.file == 0 {
		return nil
	}
	return file_search(self.file, mode, pattern, capacity)
}
