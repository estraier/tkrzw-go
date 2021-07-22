package tkrzw

import (
  "fmt"
)

// Polymorphic database manager.
//
// All operations except for open and close are thread-safe; Multiple threads can access the same database concurrently.  You can specify a data structure when you call the "Open" method.  Every opened database must be closed explicitly by the "Close" method to avoid data corruption.
type DBM struct {
	dbm uintptr
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
	count, _ := dbm_count(self.dbm)
	path, _ := dbm_get_file_path(self.dbm)
	return fmt.Sprintf("#<tkrzw.DBM:%s:%d>", path, count)
}

// Opens a database file.
//
// @param path A path of the file.
// @param writable If true, the file is writable.  If false, it is read-only.
// @param params Optional parameters.  If it is null, it is ignored.
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
//
// The optional parameter "dbm" supercedes the decision of the database type by the extension.  The value is the type name: "HashDBM", "TreeDBM", "SkipDBM", "TinyDBM", "BabyDBM", "CacheDBM", "StdHashDBM", "StdTreeDBM".
//
// The optional parameter "file" specifies the internal file implementation class. The default file class is "MemoryMapAtomicFile".  The other supported classes are "StdFile", "MemoryMapAtomicFile", "PositionalParallelFile", and "PositionalAtomicFile".
//
// For HashDBM, these optional parameters are supported.
//
// - update_mode (string): How to update the database file: "UPDATE_IN_PLACE" for the in-palce or "UPDATE_APPENDING" for the appending mode.
// - record_crc_mode (string): How to add the CRC data to the record: "RECORD_CRC_NONE" to add no CRC to each record, "RECORD_CRC_8" to add CRC-8 to each record, "RECORD_CRC_16" to add CRC-16 to each record, or "RECORD_CRC_32" to add CRC-32 to each record.
// - record_comp_mode (string): How to compress the record data: "RECORD_COMP_NONE" to do no compression, "RECORD_COMP_ZLIB" to compress with ZLib, "RECORD_COMP_ZSTD" to compress with ZStd, "RECORD_COMP_LZ4" to compress with LZ4, "RECORD_COMP_LZMA" to compress with LZMA.
// - offset_width (int): The width to represent the offset of records.
// - align_pow (int): The power to align records.
// - num_buckets (int): The number of buckets for hashing.
// - restore_mode (string): How to restore the database file: "RESTORE_SYNC" to restore to the last synchronized state, "RESTORE_READ_ONLY" to make the database read-only, or "RESTORE_NOOP" to do nothing.  By default, as many records as possible are restored.
// - fbp_capacity (int): The capacity of the free block pool.
// - min_read_size (int): The minimum reading size to read a record.
// - lock_mem_buckets (int): Positive to lock the memory for the hash buckets.
// - cache_buckets (int): Positive to cache the hash buckets on memory.
//
// For TreeDBM, all optional parameters for HashDBM are available.  In addition, these optional parameters are supported.
//
// - max_page_size (int): The maximum size of a page.
// - max_branches (int): The maximum number of branches each inner node can have.
// - max_cached_pages (int): The maximum number of cached pages.
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
// For the file "PositionalParallelFile" and "PositionalAtomicFile", these optional parameters are supported.
//
// - block_size (int): The block size to which all blocks should be aligned.
// - access_options (str): Values separated by colon.  "direct" for direct I/O.  "sync" for synchrnizing I/O, "padding" for file size alignment by padding, "pagecache" for the mini page cache in the process.
//
// If the optional parameter "num_shards" is set, the database is sharded into multiple shard files.  Each file has a suffix like "-00003-of-00015".  If the value is 0, the number of shards is set by patterns of the existing files, or 1 if they doesn't exist.
func (self *DBM) Open(path string, writable bool, params string) *Status {
	if self.dbm != 0 {
		return NewStatus2(STATUS_PRECONDITION_ERROR, "opened database")
	}
	dbm, status := dbm_open(path, writable, params)
	if status.IsOK() {
		self.dbm = dbm
	}
	return status
}

// Closes the database file.
//
// @return The result status.
func (self *DBM) Close() *Status {
	if self.dbm == 0 {
		return NewStatus2(STATUS_PRECONDITION_ERROR, "not opened database")
	}
	return dbm_close(self.dbm)
}

// Gets the value of a record of a key.
//
// @param key The key of the record.
// @return The bytes value of the matching record and the result status.
func (self *DBM) Get(key interface{}) ([]byte, *Status) {
	return dbm_get(self.dbm, ToByteArray(key))
}

// Gets the value of a record of a key, as a string.
//
// @param key The key of the record.
// @return The string value of the matching record and the result status.
func (self *DBM) GetStr(key interface{}) (string, *Status) {
	if self.dbm == 0 {
		return "", NewStatus2(STATUS_PRECONDITION_ERROR, "not opened database")
	}
	value, status := dbm_get(self.dbm, ToByteArray(key))
	if status.IsOK() {
		return string(value), status
	}
	return "", status
}

// Gets the value of a record of a key, in a simple way.
//
// @param key The key of the record.
// @param default_value The value to be returned on failure.
// @return The value of the matching record on success, or the default value on failure.
func (self *DBM) GetSimple(key interface{}, default_value interface{}) []byte {
	if self.dbm == 0 {
		return ToByteArray(default_value)
	}
	value, status := dbm_get(self.dbm, ToByteArray(key))
	if status.IsOK() {
		return value
	}
	return ToByteArray(default_value)
}

// Gets the value of a record of a key, in a simple way, as a string.
//
// @param key The key of the record.
// @param default_value The value to be returned on failure.
// @return The value of the matching record on success, or the default value on failure.
func (self *DBM) GetStrSimple(key interface{}, default_value interface{}) string {
	if self.dbm == 0 {
		return ToString(default_value)
	}
	value, status := dbm_get(self.dbm, ToByteArray(key))
	if status.IsOK() {
		return string(value)
	}
	return ToString(default_value)
}

// Sets a record of a key and a value.
//
// @param key The key of the record.
// @param value The value of the record.
// @param overwrite Whether to overwrite the existing value.  It can be omitted and then false is set.
// @return The result status.  If overwriting is abandoned, STATUS_DUPLICATION_ERROR is returned.
func (self *DBM) Set(key interface{}, value interface{}, overwrite bool) (*Status) {
	if self.dbm == 0 {
		return NewStatus2(STATUS_PRECONDITION_ERROR, "not opened database")
	}
	return dbm_set(self.dbm, ToByteArray(key), ToByteArray(value), overwrite)
}

// Removes a record of a key.
//
// @param key The key of the record.
// @return The result status.  If there's no matching record, STATUS_NOT_FOUND_ERROR is returned.
func (self *DBM) Remove(key interface{}) (*Status) {
	if self.dbm == 0 {
		return NewStatus2(STATUS_PRECONDITION_ERROR, "not opened database")
	}
	return dbm_remove(self.dbm, ToByteArray(key))
}

// Appends data at the end of a record of a key.
//
// @param key The key of the record.
// @param value The value to append.
// @param delim The delimiter to put after the existing record.
// @return The result status.
//
// If there's no existing record, the value is set without the delimiter.
func (self *DBM) Append(key interface{}, value interface{}, delim interface{}) (*Status) {
	if self.dbm == 0 {
		return NewStatus2(STATUS_PRECONDITION_ERROR, "not opened database")
	}
	return dbm_append(self.dbm, ToByteArray(key), ToByteArray(value), ToByteArray(delim))
}

// Gets the number of records.
//
// @return The number of records and the result status.
func (self *DBM) Count() (int64, *Status) {
	if self.dbm == 0 {
		return -1, NewStatus2(STATUS_PRECONDITION_ERROR, "not opened database")
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
	if status.IsOK() {
		return count
	}
	return -1;
}

// Gets the current file size of the database.
//
// @return The current file size of the database and the result status.
func (self *DBM) GetFileSize() (int64, *Status) {
	if self.dbm == 0 {
		return -1, NewStatus2(STATUS_PRECONDITION_ERROR, "not opened database")
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
	if status.IsOK() {
		return file_size
	}
	return -1;
}




/*

func dbm_append(dbm uintptr, key []byte, value []byte, delim []byte) *Status {
	if dbm == 0 || key == nil || value == nil || delim == nil {
		return NewStatus1(STATUS_INVALID_ARGUMENT_ERROR)
	}
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	xvalue_ptr := (*C.char)(C.CBytes(value))
	defer C.free(unsafe.Pointer(xvalue_ptr))
	xdelim_ptr := (*C.char)(C.CBytes(delim))
	defer C.free(unsafe.Pointer(xdelim_ptr))
	res := C.do_dbm_append(xdbm, xkey_ptr, C.int(len(key)),
		xvalue_ptr, C.int(len(value)), xdelim_ptr, C.int(len(delim)))
	status := convert_status(res)
	return status
}

func dbm_count(dbm uintptr) (int64, *Status) {
	if dbm == 0 {
		return -1, NewStatus1(STATUS_INVALID_ARGUMENT_ERROR)
	}
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	res := C.do_dbm_count(xdbm)
	status := convert_status(res.status)
	return int64(res.count), status
}

func dbm_get_file_size(dbm uintptr) (int64, *Status) {
	if dbm == 0 {
		return -1, NewStatus1(STATUS_INVALID_ARGUMENT_ERROR)
	}
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	res := C.do_dbm_get_file_size(xdbm)
	status := convert_status(res.status)
	return int64(res.count), status
}

func dbm_get_file_path(dbm uintptr) (string, *Status) {
	if dbm == 0 {
		return "", NewStatus1(STATUS_INVALID_ARGUMENT_ERROR)
	}
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	res := C.do_dbm_get_file_path(xdbm)
	var path string
	if res.path != nil {
		defer C.free(unsafe.Pointer(res.path))
		path = C.GoString(res.path)
	}
	status := convert_status(res.status)
	return path, status
}

func dbm_clear(dbm uintptr) *Status {
	if dbm == 0 {
		return NewStatus1(STATUS_INVALID_ARGUMENT_ERROR)
	}
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	res := C.do_dbm_clear(xdbm)
	status := convert_status(res)
	return status
}

func dbm_rebuild(dbm uintptr, params string) *Status {
	if dbm == 0 {
		return NewStatus1(STATUS_INVALID_ARGUMENT_ERROR)
	}
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xparams := C.CString(params)
	defer C.free(unsafe.Pointer(xparams))
	res := C.do_dbm_rebuild(xdbm, xparams)
	status := convert_status(res)
	return status
}

func dbm_should_be_rebuilt(dbm uintptr) (bool, *Status) {
	if dbm == 0 {
		return false, NewStatus1(STATUS_INVALID_ARGUMENT_ERROR)
	}
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	res := C.do_dbm_should_be_rebuilt(xdbm)
	status := convert_status(res.status)
	return bool(res.value), status
}

func dbm_synchronize(dbm uintptr, hard bool, params string) *Status {
	if dbm == 0 {
		return NewStatus1(STATUS_INVALID_ARGUMENT_ERROR)
	}
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xparams := C.CString(params)
	defer C.free(unsafe.Pointer(xparams))
	res := C.do_dbm_synchronize(xdbm, C.bool(hard), xparams)
	status := convert_status(res)
	return status
}

*/
