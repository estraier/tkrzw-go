/*************************************************************************************************
 * Bridging code to C native functions
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

/*
#cgo pkg-config: tkrzw

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "tkrzw_langc.h"

typedef struct {
  int32_t code;
  char* message;
} RES_STATUS;

typedef struct {
  TkrzwDBM* dbm;
  RES_STATUS status;
} RES_DBM;

typedef struct {
  char* value_ptr;
  int32_t value_size;
  RES_STATUS status;
} RES_VALUE;

typedef struct {
  int64_t count;
  RES_STATUS status;
} RES_INT;

typedef struct {
  char* str;
  RES_STATUS status;
} RES_STR;

typedef struct {
  bool value;
  RES_STATUS status;
} RES_BOOL;

typedef struct {
  TkrzwKeyValuePair* records;
  int32_t num_records;
} RES_MAP;

typedef struct {
  char* key_ptr;
  int32_t key_size;
  char* value_ptr;
  int32_t value_size;
  RES_STATUS status;
} RES_REC;

char* copy_status_message(const char* message) {
  if (*message == '\0') {
    return NULL;
  }
  size_t len = strlen(message);
  char* new_message = malloc(len + 1);
  if (new_message) {
    memcpy(new_message, message, len + 1);
  }
  return new_message;
}

RES_DBM do_dbm_open(const char* path, bool writable, const char* params) {
  RES_DBM res;
  res.dbm = tkrzw_dbm_open(path, writable, params);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_close(TkrzwDBM* dbm) {
  RES_STATUS res;
  tkrzw_dbm_close(dbm);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_VALUE do_dbm_get(TkrzwDBM* dbm, const char* key_ptr, int32_t key_size) {
  RES_VALUE res;
  res.value_ptr = tkrzw_dbm_get(dbm, key_ptr, key_size, &res.value_size);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_set(
    TkrzwDBM* dbm, const char* key_ptr, int32_t key_size,
    const char* value_ptr, int32_t value_size, bool overwrite) {
  RES_STATUS res;
  tkrzw_dbm_set(dbm, key_ptr, key_size, value_ptr, value_size, overwrite);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_set_multi(
    TkrzwDBM* dbm, const TkrzwKeyValuePair* records, int32_t num_records, bool overwrite) {
  RES_STATUS res;
  tkrzw_dbm_set_multi(dbm, records, num_records, overwrite);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_remove(TkrzwDBM* dbm, const char* key_ptr, int32_t key_size) {
  RES_STATUS res;
  tkrzw_dbm_remove(dbm, key_ptr, key_size);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_remove_multi(TkrzwDBM* dbm, const TkrzwStr* keys, int32_t num_keys) {
  RES_STATUS res;
  tkrzw_dbm_remove_multi(dbm, keys, num_keys);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_append(
    TkrzwDBM* dbm, const char* key_ptr, int32_t key_size,
    const char* value_ptr, int32_t value_size, const char* delim_ptr, int32_t delim_size) {
  RES_STATUS res;
  tkrzw_dbm_append(dbm, key_ptr, key_size, value_ptr, value_size, delim_ptr, delim_size);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_compare_exchange(
    TkrzwDBM* dbm, const char* key_ptr, int32_t key_size,
    const char* expected_ptr, int32_t expected_size,
    const char* desired_ptr, int32_t desired_size) {
  RES_STATUS res;
  tkrzw_dbm_compare_exchange(
      dbm, key_ptr, key_size, expected_ptr, expected_size, desired_ptr, desired_size);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_INT do_dbm_increment(
    TkrzwDBM* dbm, const char* key_ptr, int32_t key_size, int64_t inc, int64_t init) {
  RES_INT res;
  res.count = tkrzw_dbm_increment(dbm, key_ptr, key_size, inc, init);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_INT do_dbm_count(TkrzwDBM* dbm) {
  RES_INT res;
  res.count = tkrzw_dbm_count(dbm);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_INT do_dbm_get_file_size(TkrzwDBM* dbm) {
  RES_INT res;
  res.count = tkrzw_dbm_get_file_size(dbm);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_STR do_dbm_get_file_path(TkrzwDBM* dbm) {
  RES_STR res;
  res.str = tkrzw_dbm_get_file_path(dbm);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_clear(TkrzwDBM* dbm) {
  RES_STATUS res;
  tkrzw_dbm_clear(dbm);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_rebuild(TkrzwDBM* dbm, const char* params) {
  RES_STATUS res;
  tkrzw_dbm_rebuild(dbm, params);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_BOOL do_dbm_should_be_rebuilt(TkrzwDBM* dbm) {
  RES_BOOL res;
  res.value = tkrzw_dbm_should_be_rebuilt(dbm);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_synchronize(TkrzwDBM* dbm, bool hard, const char* params) {
  RES_STATUS res;
  tkrzw_dbm_synchronize(dbm, hard, NULL, NULL, params);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_copy_file_data(TkrzwDBM* dbm, const char* dest_path) {
  RES_STATUS res;
  tkrzw_dbm_copy_file_data(dbm, dest_path);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_export(TkrzwDBM* dbm, TkrzwDBM* dest_dbm) {
  RES_STATUS res;
  tkrzw_dbm_export(dbm, dest_dbm);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_MAP do_dbm_inspect(TkrzwDBM* dbm) {
  RES_MAP res;
  res.records = tkrzw_dbm_inspect(dbm, &res.num_records);
  return res;
}

RES_STATUS do_dbm_iter_first(TkrzwDBMIter* iter) {
  RES_STATUS res;
  tkrzw_dbm_iter_first(iter);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_iter_last(TkrzwDBMIter* iter) {
  RES_STATUS res;
  tkrzw_dbm_iter_last(iter);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_iter_jump(TkrzwDBMIter* iter, const char* key_ptr, int32_t key_size) {
  RES_STATUS res;
  tkrzw_dbm_iter_jump(iter, key_ptr, key_size);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_iter_jump_lower(
    TkrzwDBMIter* iter, const char* key_ptr, int32_t key_size, bool inclusive) {
  RES_STATUS res;
  tkrzw_dbm_iter_jump_lower(iter, key_ptr, key_size, inclusive);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_iter_jump_upper(
    TkrzwDBMIter* iter, const char* key_ptr, int32_t key_size, bool inclusive) {
  RES_STATUS res;
  tkrzw_dbm_iter_jump_upper(iter, key_ptr, key_size, inclusive);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_iter_next(TkrzwDBMIter* iter) {
  RES_STATUS res;
  tkrzw_dbm_iter_next(iter);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_iter_previous(TkrzwDBMIter* iter) {
  RES_STATUS res;
  tkrzw_dbm_iter_previous(iter);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STR do_dbm_iter_get_key_esc(TkrzwDBMIter* iter) {
  int32_t key_size = 0;
  char* key_ptr = tkrzw_dbm_iter_get_key(iter, &key_size);
  RES_STR res;
  if (key_ptr == NULL) {
    res.str = NULL;
  } else {
    res.str = tkrzw_str_escape_c(key_ptr, key_size, true, NULL);
  }
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_REC do_dbm_iter_get(TkrzwDBMIter* iter) {
  RES_REC res;
  res.key_ptr = NULL;
  res.value_ptr = NULL;
  tkrzw_dbm_iter_get(iter, &res.key_ptr, &res.key_size, &res.value_ptr, &res.value_size);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_VALUE do_dbm_iter_get_key(TkrzwDBMIter* iter) {
  RES_VALUE res;
  res.value_ptr = tkrzw_dbm_iter_get_key(iter, &res.value_size);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_VALUE do_dbm_iter_get_value(TkrzwDBMIter* iter) {
  RES_VALUE res;
  res.value_ptr = tkrzw_dbm_iter_get_value(iter, &res.value_size);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

*/
import "C"

import (
	"unsafe"
)

var VERSION string

func init() {
	VERSION = C.GoString(C.TKRZW_PACKAGE_VERSION)
}

func convert_status(res C.RES_STATUS) *Status {
	if res.message == nil {
		return NewStatus1(StatusCode(res.code))
	}
	defer C.free(unsafe.Pointer(res.message))
	return NewStatus2(StatusCode(res.code), C.GoString(res.message))
}

func dbm_open(path string, writable bool, params string) (uintptr, *Status) {
	xpath := C.CString(path)
	defer C.free(unsafe.Pointer(xpath))
	xparams := C.CString(params)
	defer C.free(unsafe.Pointer(xparams))
	res := C.do_dbm_open(xpath, C.bool(writable), xparams)
	status := convert_status(res.status)
	return uintptr(unsafe.Pointer(res.dbm)), status
}

func dbm_close(dbm uintptr) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	res := C.do_dbm_close(xdbm)
	status := convert_status(res)
	return status
}

func dbm_get(dbm uintptr, key []byte) ([]byte, *Status) {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	res := C.do_dbm_get(xdbm, xkey_ptr, C.int32_t(len(key)))
	var value []byte = nil
	if res.value_ptr != nil {
		defer C.free(unsafe.Pointer(res.value_ptr))
		value = C.GoBytes(unsafe.Pointer(res.value_ptr), res.value_size)
	}
	status := convert_status(res.status)
	return value, status
}

func dbm_get_multi(dbm uintptr, keys []string) map[string][]byte {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xkeys_size := len(keys) * int(unsafe.Sizeof(C.TkrzwStr{}))
	xkeys := (*C.TkrzwStr)(unsafe.Pointer(C.malloc(C.size_t(xkeys_size + 1))))
	defer C.tkrzw_free_str_array(xkeys, C.int32_t(len(keys)))
	xkey_ptr := uintptr(unsafe.Pointer(xkeys))
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		xkey := (*C.TkrzwStr)(unsafe.Pointer(xkey_ptr))
		xkey.ptr = C.CString(key)
		xkey.size = C.int32_t(len(key))
		xkey_ptr += unsafe.Sizeof(C.TkrzwStr{})
	}
	var num_records C.int32_t
	xrecords := C.tkrzw_dbm_get_multi(xdbm, xkeys, C.int32_t(len(keys)), &num_records)
	defer C.tkrzw_free_str_map(xrecords, num_records)
	records := make(map[string][]byte)
	rec_ptr := uintptr(unsafe.Pointer(xrecords))
	for i := C.int32_t(0); i < num_records; i++ {
		elem := (*C.TkrzwKeyValuePair)(unsafe.Pointer(rec_ptr))
		key := C.GoStringN(elem.key_ptr, elem.key_size)
		value := C.GoBytes(unsafe.Pointer(elem.value_ptr), elem.value_size)
		records[key] = value
		rec_ptr += unsafe.Sizeof(C.TkrzwKeyValuePair{})
	}
	return records
}

func dbm_set(dbm uintptr, key []byte, value []byte, overwrite bool) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	xvalue_ptr := (*C.char)(C.CBytes(value))
	defer C.free(unsafe.Pointer(xvalue_ptr))
	res := C.do_dbm_set(xdbm, xkey_ptr, C.int32_t(len(key)),
		xvalue_ptr, C.int32_t(len(value)), C.bool(overwrite))
	status := convert_status(res)
	return status
}

func dbm_set_multi(dbm uintptr, records map[string][]byte, overwrite bool) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xrecs_size := len(records) * int(unsafe.Sizeof(C.TkrzwKeyValuePair{}))
	xrecs := (*C.TkrzwKeyValuePair)(unsafe.Pointer(C.malloc(C.size_t(xrecs_size + 1))))
	defer C.tkrzw_free_str_map(xrecs, C.int32_t(len(records)))
	xrec_ptr := uintptr(unsafe.Pointer(xrecs))
	for key, value := range records {
		xrec := (*C.TkrzwKeyValuePair)(unsafe.Pointer(xrec_ptr))
		xrec.key_ptr = C.CString(key)
		xrec.key_size = C.int32_t(len(key))
		xrec.value_ptr = (*C.char)(C.CBytes(value))
		xrec.value_size = C.int32_t(len(value))
		xrec_ptr += unsafe.Sizeof(C.TkrzwKeyValuePair{})
	}
	res := C.do_dbm_set_multi(xdbm, xrecs, C.int32_t(len(records)), C.bool(overwrite))
	status := convert_status(res)
	return status
}

func dbm_remove(dbm uintptr, key []byte) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	res := C.do_dbm_remove(xdbm, xkey_ptr, C.int32_t(len(key)))
	status := convert_status(res)
	return status
}

func dbm_remove_multi(dbm uintptr, keys []string) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xkeys_size := len(keys) * int(unsafe.Sizeof(C.TkrzwStr{}))
	xkeys := (*C.TkrzwStr)(unsafe.Pointer(C.malloc(C.size_t(xkeys_size + 1))))
	defer C.tkrzw_free_str_array(xkeys, C.int32_t(len(keys)))
	xkey_ptr := uintptr(unsafe.Pointer(xkeys))
	for i := 0; i < len(keys); i++ {
		key := keys[i]
		xkey := (*C.TkrzwStr)(unsafe.Pointer(xkey_ptr))
		xkey.ptr = C.CString(key)
		xkey.size = C.int32_t(len(key))
		xkey_ptr += unsafe.Sizeof(C.TkrzwStr{})
	}
	res := C.do_dbm_remove_multi(xdbm, xkeys, C.int32_t(len(keys)))
	status := convert_status(res)
	return status
}

func dbm_append(dbm uintptr, key []byte, value []byte, delim []byte) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	xvalue_ptr := (*C.char)(C.CBytes(value))
	defer C.free(unsafe.Pointer(xvalue_ptr))
	xdelim_ptr := (*C.char)(C.CBytes(delim))
	defer C.free(unsafe.Pointer(xdelim_ptr))
	res := C.do_dbm_append(xdbm, xkey_ptr, C.int32_t(len(key)),
		xvalue_ptr, C.int32_t(len(value)), xdelim_ptr, C.int32_t(len(delim)))
	status := convert_status(res)
	return status
}

func dbm_compare_exchange(dbm uintptr, key []byte, expected []byte, desired []byte) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	var xexpected_ptr *C.char
	var xexpected_size C.int32_t
	if expected != nil {
		xexpected_ptr = (*C.char)(C.CBytes(expected))
		defer C.free(unsafe.Pointer(xexpected_ptr))
		xexpected_size = C.int32_t(len(expected))
	}
	var xdesired_ptr *C.char
	var xdesired_size C.int32_t
	if desired != nil {
		xdesired_ptr = (*C.char)(C.CBytes(desired))
		defer C.free(unsafe.Pointer(xdesired_ptr))
		xdesired_size = C.int32_t(len(desired))
	}
	res := C.do_dbm_compare_exchange(
		xdbm, xkey_ptr, C.int32_t(len(key)),
		xexpected_ptr, xexpected_size, xdesired_ptr, xdesired_size)
	status := convert_status(res)
	return status
}

func dbm_increment(dbm uintptr, key []byte, inc int64, init int64) (int64, *Status) {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	res := C.do_dbm_increment(
		xdbm, xkey_ptr, C.int32_t(len(key)), C.int64_t(inc), C.int64_t(init))
	status := convert_status(res.status)
	return int64(res.count), status
}

func dbm_count(dbm uintptr) (int64, *Status) {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	res := C.do_dbm_count(xdbm)
	status := convert_status(res.status)
	return int64(res.count), status
}

func dbm_get_file_size(dbm uintptr) (int64, *Status) {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	res := C.do_dbm_get_file_size(xdbm)
	status := convert_status(res.status)
	return int64(res.count), status
}

func dbm_get_file_path(dbm uintptr) (string, *Status) {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	res := C.do_dbm_get_file_path(xdbm)
	var path string
	if res.str != nil {
		defer C.free(unsafe.Pointer(res.str))
		path = C.GoString(res.str)
	}
	status := convert_status(res.status)
	return path, status
}

func dbm_clear(dbm uintptr) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	res := C.do_dbm_clear(xdbm)
	status := convert_status(res)
	return status
}

func dbm_rebuild(dbm uintptr, params string) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xparams := C.CString(params)
	defer C.free(unsafe.Pointer(xparams))
	res := C.do_dbm_rebuild(xdbm, xparams)
	status := convert_status(res)
	return status
}

func dbm_should_be_rebuilt(dbm uintptr) (bool, *Status) {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	res := C.do_dbm_should_be_rebuilt(xdbm)
	status := convert_status(res.status)
	return bool(res.value), status
}

func dbm_synchronize(dbm uintptr, hard bool, params string) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xparams := C.CString(params)
	defer C.free(unsafe.Pointer(xparams))
	res := C.do_dbm_synchronize(xdbm, C.bool(hard), xparams)
	status := convert_status(res)
	return status
}

func dbm_copy_file_data(dbm uintptr, dest_path string) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xdest_path := C.CString(dest_path)
	defer C.free(unsafe.Pointer(xdest_path))
	res := C.do_dbm_copy_file_data(xdbm, xdest_path)
	status := convert_status(res)
	return status
}

func dbm_export(dbm uintptr, dest_dbm uintptr) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xdest_dbm := (*C.TkrzwDBM)(unsafe.Pointer(dest_dbm))
	res := C.do_dbm_export(xdbm, xdest_dbm)
	status := convert_status(res)
	return status
}

func dbm_inspect(dbm uintptr, records map[string]string) {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	res := C.do_dbm_inspect(xdbm)
	defer C.tkrzw_free_str_map(res.records, res.num_records)
	rec_ptr := uintptr(unsafe.Pointer(res.records))
	for i := C.int32_t(0); i < res.num_records; i++ {
		elem := (*C.TkrzwKeyValuePair)(unsafe.Pointer(rec_ptr))
		name := C.GoStringN(elem.key_ptr, elem.key_size)
		value := C.GoStringN(elem.value_ptr, elem.value_size)
		records[name] = value
		rec_ptr += unsafe.Sizeof(C.TkrzwKeyValuePair{})
	}
}

func dbm_is_writable(dbm uintptr) bool {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	return bool(C.tkrzw_dbm_is_writable(xdbm))
}

func dbm_is_healthy(dbm uintptr) bool {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	return bool(C.tkrzw_dbm_is_healthy(xdbm))
}

func dbm_is_ordered(dbm uintptr) bool {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	return bool(C.tkrzw_dbm_is_ordered(xdbm))
}

func dbm_make_iterator(dbm uintptr) uintptr {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	return uintptr(unsafe.Pointer(C.tkrzw_dbm_make_iterator(xdbm)))
}

func dbm_iter_free(iter uintptr) {
	xiter := (*C.TkrzwDBMIter)(unsafe.Pointer(iter))
	C.tkrzw_dbm_iter_free(xiter)
}

func dbm_iter_first(iter uintptr) *Status {
	xiter := (*C.TkrzwDBMIter)(unsafe.Pointer(iter))
	res := C.do_dbm_iter_first(xiter)
	status := convert_status(res)
	return status
}

func dbm_iter_last(iter uintptr) *Status {
	xiter := (*C.TkrzwDBMIter)(unsafe.Pointer(iter))
	res := C.do_dbm_iter_last(xiter)
	status := convert_status(res)
	return status
}

func dbm_iter_jump(iter uintptr, key []byte) *Status {
	xiter := (*C.TkrzwDBMIter)(unsafe.Pointer(iter))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	res := C.do_dbm_iter_jump(xiter, xkey_ptr, C.int32_t(len(key)))
	status := convert_status(res)
	return status
}

func dbm_iter_jump_lower(iter uintptr, key []byte, inclusive bool) *Status {
	xiter := (*C.TkrzwDBMIter)(unsafe.Pointer(iter))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	res := C.do_dbm_iter_jump_lower(xiter, xkey_ptr, C.int32_t(len(key)), C.bool(inclusive))
	status := convert_status(res)
	return status
}

func dbm_iter_jump_upper(iter uintptr, key []byte, inclusive bool) *Status {
	xiter := (*C.TkrzwDBMIter)(unsafe.Pointer(iter))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	res := C.do_dbm_iter_jump_upper(xiter, xkey_ptr, C.int32_t(len(key)), C.bool(inclusive))
	status := convert_status(res)
	return status
}

func dbm_iter_next(iter uintptr) *Status {
	xiter := (*C.TkrzwDBMIter)(unsafe.Pointer(iter))
	res := C.do_dbm_iter_next(xiter)
	status := convert_status(res)
	return status
}

func dbm_iter_previous(iter uintptr) *Status {
	xiter := (*C.TkrzwDBMIter)(unsafe.Pointer(iter))
	res := C.do_dbm_iter_previous(xiter)
	status := convert_status(res)
	return status
}

func dbm_iter_get_key_esc(iter uintptr) (string, *Status) {
	xiter := (*C.TkrzwDBMIter)(unsafe.Pointer(iter))
	res := C.do_dbm_iter_get_key_esc(xiter)
	var key string
	if res.str != nil {
		defer C.free(unsafe.Pointer(res.str))
		key = C.GoString(res.str)
	}
	status := convert_status(res.status)
	return key, status
}

func dbm_iter_get(iter uintptr) ([]byte, []byte, *Status) {
	xiter := (*C.TkrzwDBMIter)(unsafe.Pointer(iter))
	res := C.do_dbm_iter_get(xiter)
	var key []byte = nil
	if res.key_ptr != nil {
		defer C.free(unsafe.Pointer(res.key_ptr))
		key = C.GoBytes(unsafe.Pointer(res.key_ptr), res.key_size)
	}
	var value []byte = nil
	if res.value_ptr != nil {
		defer C.free(unsafe.Pointer(res.value_ptr))
		value = C.GoBytes(unsafe.Pointer(res.value_ptr), res.value_size)
	}
	status := convert_status(res.status)
	return key, value, status
}

func dbm_iter_get_key(iter uintptr) ([]byte, *Status) {
	xiter := (*C.TkrzwDBMIter)(unsafe.Pointer(iter))
	res := C.do_dbm_iter_get_key(xiter)
	var key []byte = nil
	if res.value_ptr != nil {
		defer C.free(unsafe.Pointer(res.value_ptr))
		key = C.GoBytes(unsafe.Pointer(res.value_ptr), res.value_size)
	}
	status := convert_status(res.status)
	return key, status
}

func dbm_iter_get_value(iter uintptr) ([]byte, *Status) {
	xiter := (*C.TkrzwDBMIter)(unsafe.Pointer(iter))
	res := C.do_dbm_iter_get_value(xiter)
	var value []byte = nil
	if res.value_ptr != nil {
		defer C.free(unsafe.Pointer(res.value_ptr))
		value = C.GoBytes(unsafe.Pointer(res.value_ptr), res.value_size)
	}
	status := convert_status(res.status)
	return value, status
}

// END OF FILE
