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
} RES_DBM_OPEN;

typedef struct {
  char* value_ptr;
  int32_t value_size;
  RES_STATUS status;
} RES_DBM_GET;

typedef struct {
  int64_t count;
  RES_STATUS status;
} RES_DBM_COUNT;

typedef struct {
  char* path;
  RES_STATUS status;
} RES_DBM_GET_FILE_PATH;

typedef struct {
  bool value;
  RES_STATUS status;
} RES_DBM_BOOL;

typedef struct {
  TkrzwKeyValuePair* records;
  int32_t num_records;
} RES_DBM_RECORDS;

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

RES_DBM_OPEN do_dbm_open(const char* path, bool writable, const char* params) {
  RES_DBM_OPEN res;
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

RES_DBM_GET do_dbm_get(TkrzwDBM* dbm, const char* key_ptr, int32_t key_size) {
  RES_DBM_GET res;
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

RES_STATUS do_dbm_remove(TkrzwDBM* dbm, const char* key_ptr, int32_t key_size) {
  RES_STATUS res;
  tkrzw_dbm_remove(dbm, key_ptr, key_size);
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

RES_DBM_COUNT do_dbm_count(TkrzwDBM* dbm) {
  RES_DBM_COUNT res;
  res.count = tkrzw_dbm_count(dbm);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_DBM_COUNT do_dbm_get_file_size(TkrzwDBM* dbm) {
  RES_DBM_COUNT res;
  res.count = tkrzw_dbm_get_file_size(dbm);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_DBM_GET_FILE_PATH do_dbm_get_file_path(TkrzwDBM* dbm) {
  RES_DBM_GET_FILE_PATH res;
  res.path = tkrzw_dbm_get_file_path(dbm);
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

RES_DBM_BOOL do_dbm_should_be_rebuilt(TkrzwDBM* dbm) {
  RES_DBM_BOOL res;
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

RES_DBM_RECORDS do_dbm_inspect(TkrzwDBM* dbm) {
  RES_DBM_RECORDS res;
  res.records = tkrzw_dbm_inspect(dbm, &res.num_records);
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
	res := C.do_dbm_get(xdbm, xkey_ptr, C.int(len(key)))
	var value []byte = nil
	if res.value_ptr != nil {
		defer C.free(unsafe.Pointer(res.value_ptr))
		value = C.GoBytes(unsafe.Pointer(res.value_ptr), res.value_size)
	}
	status := convert_status(res.status)
	return value, status
}

func dbm_set(dbm uintptr, key []byte, value []byte, overwrite bool) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	xvalue_ptr := (*C.char)(C.CBytes(value))
	defer C.free(unsafe.Pointer(xvalue_ptr))
	res := C.do_dbm_set(xdbm, xkey_ptr, C.int(len(key)),
		xvalue_ptr, C.int(len(value)), C.bool(overwrite))
	status := convert_status(res)
	return status
}

func dbm_remove(dbm uintptr, key []byte) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	res := C.do_dbm_remove(xdbm, xkey_ptr, C.int(len(key)))
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
	res := C.do_dbm_append(xdbm, xkey_ptr, C.int(len(key)),
		xvalue_ptr, C.int(len(value)), xdelim_ptr, C.int(len(delim)))
	status := convert_status(res)
	return status
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
	if res.path != nil {
		defer C.free(unsafe.Pointer(res.path))
		path = C.GoString(res.path)
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
		elem := *(*C.TkrzwKeyValuePair)(unsafe.Pointer(rec_ptr))
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
