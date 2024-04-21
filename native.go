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
  TkrzwFile* file;
  RES_STATUS status;
} RES_FILE;

typedef struct {
  TkrzwIndex* index;
  RES_STATUS status;
} RES_INDEX;

typedef struct {
  char* value_ptr;
  int32_t value_size;
  RES_STATUS status;
} RES_BYTES;

typedef struct {
  int64_t num;
  RES_STATUS status;
} RES_INT;

typedef struct {
  double num;
  RES_STATUS status;
} RES_FLOAT;

typedef struct {
  char* str;
  RES_STATUS status;
} RES_STR;

typedef struct {
  TkrzwKeyValuePair* str_pair;
  RES_STATUS status;
} RES_STRPAIR;

typedef struct {
  bool value;
  RES_STATUS status;
} RES_BOOL;

typedef struct {
  TkrzwStr* values;
  int32_t num_values;
  RES_STATUS status;
} RES_STRARRAY;

typedef struct {
  TkrzwKeyValuePair* records;
  int32_t num_records;
  RES_STATUS status;
} RES_STRMAP;

typedef struct {
  char* key_ptr;
  int32_t key_size;
  char* value_ptr;
  int32_t value_size;
  RES_STATUS status;
} RES_REC;

typedef struct {
  char* key_ptr;
  int32_t key_size;
  char* value_ptr;
  int32_t value_size;
  bool status;
} RES_REC_BOOL;

typedef struct {
  void* proc_up;
  void* buffer;
} RecordProcessorArg;

struct callRecordProcessor_return {
  void* ptr;
  int32_t size;
};

extern struct callRecordProcessor_return callRecordProcessor(
  void* up, void* keyPtr, int32_t keySize, void* valuePtr, int32_t valueSize);

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

void free_str_pairs(TkrzwKeyValuePair* array, int32_t size) {
  for (int32_t i = 0; i < size; i++) {
    free((char*)(array[i].key_ptr));
    if (array[i].value_ptr != TKRZW_ANY_DATA) {
      free((char*)(array[i].value_ptr));
    }
  }
  free(array);
}

const char* run_record_processor(
  RecordProcessorArg* arg, const char* key_ptr, int32_t key_size,
  const char* value_ptr, int32_t value_size, int32_t* ret_size) {
  if (arg->buffer) {
    free(arg->buffer);
    arg->buffer = NULL;
  }
  struct callRecordProcessor_return rv = callRecordProcessor(
    arg->proc_up, (void*)key_ptr, key_size, (void*)value_ptr, value_size);
  const char* ret_ptr;
  if (rv.ptr == NULL) {
    ret_ptr = TKRZW_REC_PROC_NOOP;
    *ret_size = 0;
  } else if (rv.ptr == (void*)1) {
    ret_ptr = TKRZW_REC_PROC_REMOVE;
    *ret_size = 0;
  } else {
    ret_ptr = rv.ptr;
    *ret_size = rv.size;
    arg->buffer = rv.ptr;
  }
  return ret_ptr;
}

RES_STATUS do_future_get(TkrzwFuture* future) {
  RES_STATUS res;
  tkrzw_future_get(future);
  tkrzw_future_free(future);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_BYTES do_future_get_bytes(TkrzwFuture* future) {
  RES_BYTES res;
  res.value_ptr = tkrzw_future_get_str(future, &res.value_size);
  tkrzw_future_free(future);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_STRPAIR do_future_get_str_pair(TkrzwFuture* future) {
  RES_STRPAIR res;
  res.str_pair = tkrzw_future_get_str_pair(future);
  tkrzw_future_free(future);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_STRARRAY do_future_get_str_array(TkrzwFuture* future) {
  RES_STRARRAY res;
  res.values = tkrzw_future_get_str_array(future, &res.num_values);
  tkrzw_future_free(future);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_STRMAP do_future_get_str_map(TkrzwFuture* future) {
  RES_STRMAP res;
  res.records = tkrzw_future_get_str_map(future, &res.num_records);
  tkrzw_future_free(future);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_INT do_future_get_int(TkrzwFuture* future) {
  RES_INT res;
  res.num = tkrzw_future_get_int(future);
  tkrzw_future_free(future);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
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

RES_STATUS do_dbm_process(TkrzwDBM* dbm, const char* key_ptr, int32_t key_size,
                          void* proc_up, bool writable) {
  RES_STATUS res;
  RecordProcessorArg proc_arg;
  proc_arg.proc_up = proc_up;
  proc_arg.buffer = NULL;
  tkrzw_dbm_process(dbm, key_ptr, key_size, (tkrzw_record_processor)run_record_processor,
    &proc_arg, writable);
  free(proc_arg.buffer);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_BYTES do_dbm_get(TkrzwDBM* dbm, const char* key_ptr, int32_t key_size) {
  RES_BYTES res;
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

RES_BYTES do_dbm_set_and_get(
    TkrzwDBM* dbm, const char* key_ptr, int32_t key_size,
    const char* value_ptr, int32_t value_size, bool overwrite) {
  RES_BYTES res;
  res.value_ptr = tkrzw_dbm_set_and_get(
      dbm, key_ptr, key_size, value_ptr, value_size, overwrite, &res.value_size);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
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

RES_BYTES do_dbm_remove_and_get(TkrzwDBM* dbm, const char* key_ptr, int32_t key_size) {
  RES_BYTES res;
  res.value_ptr = tkrzw_dbm_remove_and_get(dbm, key_ptr, key_size, &res.value_size);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
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

RES_STATUS do_dbm_append_multi(
    TkrzwDBM* dbm, const TkrzwKeyValuePair* records, int32_t num_records,
    const char* delim_ptr, int32_t delim_size) {
  RES_STATUS res;
  tkrzw_dbm_append_multi(dbm, records, num_records, delim_ptr, delim_size);
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

RES_BYTES do_dbm_compare_exchange_and_get(
    TkrzwDBM* dbm, const char* key_ptr, int32_t key_size,
    const char* expected_ptr, int32_t expected_size,
    const char* desired_ptr, int32_t desired_size) {
  RES_BYTES res;
  res.value_ptr = tkrzw_dbm_compare_exchange_and_get(
      dbm, key_ptr, key_size, expected_ptr, expected_size, desired_ptr, desired_size,
      &res.value_size);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_INT do_dbm_increment(
    TkrzwDBM* dbm, const char* key_ptr, int32_t key_size, int64_t inc, int64_t init) {
  RES_INT res;
  res.num = tkrzw_dbm_increment(dbm, key_ptr, key_size, inc, init);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_process_multi(
    TkrzwDBM* dbm, TkrzwKeyProcPair* pairs, int32_t num_pairs, bool writable) {
  RES_STATUS res;
  for (int i = 0; i < num_pairs; i++) {
    TkrzwKeyProcPair* pair = pairs + i;
    void* proc_up = pair->proc_arg;
    pair->proc = (tkrzw_record_processor)run_record_processor;
    RecordProcessorArg* proc_arg = malloc(sizeof(RecordProcessorArg));
    proc_arg->proc_up = proc_up;
    proc_arg->buffer = NULL;
    pair->proc_arg = proc_arg;
  }
  tkrzw_dbm_process_multi(dbm, pairs, num_pairs, writable);
  for (int i = 0; i < num_pairs; i++) {
    TkrzwKeyProcPair* pair = pairs + i;
    RecordProcessorArg* proc_arg = (RecordProcessorArg*)pair->proc_arg;
    free(proc_arg->buffer);
    free(proc_arg);
  }
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_compare_exchange_multi(
    TkrzwDBM* dbm, const TkrzwKeyValuePair* expected, int32_t num_expected,
    const TkrzwKeyValuePair* desired, int32_t num_desired) {
  RES_STATUS res;
  tkrzw_dbm_compare_exchange_multi(dbm, expected, num_expected, desired, num_desired);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_rekey(
    TkrzwDBM* dbm, const char* old_key_ptr, int32_t old_key_size,
    const char* new_key_ptr, int32_t new_key_size, bool overwrite, bool copying) {
  RES_STATUS res;
  tkrzw_dbm_rekey(dbm, old_key_ptr, old_key_size, new_key_ptr, new_key_size, overwrite, copying);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_REC do_dbm_pop_first(TkrzwDBM* dbm) {
  RES_REC res;
  res.key_ptr = NULL;
  res.value_ptr = NULL;
  tkrzw_dbm_pop_first(dbm, &res.key_ptr, &res.key_size, &res.value_ptr, &res.value_size);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_push_last(
    TkrzwDBM* dbm, const char* value_ptr, int32_t value_size, double wtime) {
  RES_STATUS res;
  tkrzw_dbm_push_last(dbm, value_ptr, value_size, wtime);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_process_each(TkrzwDBM* dbm, void* proc_up, bool writable) {
  RES_STATUS res;
  RecordProcessorArg proc_arg;
  proc_arg.proc_up = proc_up;
  proc_arg.buffer = NULL;
  tkrzw_dbm_process_each(dbm, (tkrzw_record_processor)run_record_processor, &proc_arg, writable);
  free(proc_arg.buffer);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_INT do_dbm_count(TkrzwDBM* dbm) {
  RES_INT res;
  res.num = tkrzw_dbm_count(dbm);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_INT do_dbm_get_file_size(TkrzwDBM* dbm) {
  RES_INT res;
  res.num = tkrzw_dbm_get_file_size(dbm);
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

RES_FLOAT do_dbm_get_timestamp(TkrzwDBM* dbm) {
  RES_FLOAT res;
  res.num = tkrzw_dbm_get_timestamp(dbm);
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

RES_STATUS do_dbm_copy_file_data(TkrzwDBM* dbm, const char* dest_path, bool sync_hard) {
  RES_STATUS res;
  tkrzw_dbm_copy_file_data(dbm, dest_path, sync_hard);
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

RES_STATUS do_dbm_export_to_flat_records(TkrzwDBM* dbm, TkrzwFile* dest_file) {
  RES_STATUS res;
  tkrzw_dbm_export_to_flat_records(dbm, dest_file);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_import_from_flat_records(TkrzwDBM* dbm, TkrzwFile* src_file) {
  RES_STATUS res;
  tkrzw_dbm_import_from_flat_records(dbm, src_file);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_export_keys_as_lines(TkrzwDBM* dbm, TkrzwFile* dest_file) {
  RES_STATUS res;
  tkrzw_dbm_export_keys_as_lines(dbm, dest_file);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STRMAP do_dbm_inspect(TkrzwDBM* dbm) {
  RES_STRMAP res;
  res.records = tkrzw_dbm_inspect(dbm, &res.num_records);
  return res;
}

RES_STATUS do_dbm_restore_database(const char* old_file_path, const char* new_file_path,
    const char* class_name, int64_t end_offset, const char* cipher_key) {
  RES_STATUS res;
  bool r = tkrzw_dbm_restore_database(
    old_file_path, new_file_path, class_name, end_offset, cipher_key);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
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

RES_BYTES do_dbm_iter_get_key(TkrzwDBMIter* iter) {
  RES_BYTES res;
  res.value_ptr = tkrzw_dbm_iter_get_key(iter, &res.value_size);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_BYTES do_dbm_iter_get_value(TkrzwDBMIter* iter) {
  RES_BYTES res;
  res.value_ptr = tkrzw_dbm_iter_get_value(iter, &res.value_size);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_iter_set(TkrzwDBMIter* iter, const char* value_ptr, int32_t value_size) {
  RES_STATUS res;
  tkrzw_dbm_iter_set(iter, value_ptr, value_size);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_dbm_iter_remove(TkrzwDBMIter* iter) {
  RES_STATUS res;
  tkrzw_dbm_iter_remove(iter);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_REC do_dbm_iter_step(TkrzwDBMIter* iter) {
  RES_REC res;
  res.key_ptr = NULL;
  res.value_ptr = NULL;
  tkrzw_dbm_iter_step(iter, &res.key_ptr, &res.key_size, &res.value_ptr, &res.value_size);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_FILE do_file_open(const char* path, bool writable, const char* params) {
  RES_FILE res;
  res.file = tkrzw_file_open(path, writable, params);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_file_close(TkrzwFile* file) {
  RES_STATUS res;
  tkrzw_file_close(file);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_file_read(TkrzwFile* file, int64_t off, char* buf, size_t size) {
  RES_STATUS res;
  tkrzw_file_read(file, off, buf, size);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_file_write(TkrzwFile* file, int64_t off, const char* buf, size_t size) {
  RES_STATUS res;
  tkrzw_file_write(file, off, buf, size);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_INT do_file_append(TkrzwFile* file, const char* buf, size_t size) {
  RES_INT res;
  tkrzw_file_append(file, buf, size, &res.num);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_INT do_file_truncate(TkrzwFile* file, int64_t size) {
  RES_INT res;
  tkrzw_file_truncate(file, size);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_INT do_file_synchronize(TkrzwFile* file, bool hard, int64_t off, int64_t size) {
  RES_INT res;
  tkrzw_file_synchronize(file, hard, off, size);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_INT do_file_get_size(TkrzwFile* file) {
  RES_INT res;
  res.num = tkrzw_file_get_size(file);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_STR do_file_get_path(TkrzwFile* file) {
  RES_STR res;
  res.str = tkrzw_file_get_path(file);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_INDEX do_index_open(const char* path, bool writable, const char* params) {
  RES_INDEX res;
  res.index = tkrzw_index_open(path, writable, params);
  TkrzwStatus status = tkrzw_get_last_status();
  res.status.code = status.code;
  res.status.message = copy_status_message(status.message);
  return res;
}

RES_STATUS do_index_close(TkrzwIndex* index) {
  RES_STATUS res;
  tkrzw_index_close(index);
  TkrzwStatus status = tkrzw_get_last_status();
  res.code = status.code;
  res.message = copy_status_message(status.message);
  return res;
}

RES_REC_BOOL do_index_iter_get(TkrzwIndexIter* iter) {
  RES_REC_BOOL res;
  res.key_ptr = NULL;
  res.value_ptr = NULL;
  res.status = tkrzw_index_iter_get(
    iter, &res.key_ptr, &res.key_size, &res.value_ptr, &res.value_size);
  return res;
}

*/
import "C"

import (
	"strings"
	"unsafe"
)

// The package version numbers.
var Version string

// The recognized OS name.
var OSName string

// The recognized OS name.
var PageSize int

// The minimum value of int64. */
var Int64Min int64

// The minimum value of int64. */
var Int64Max int64

func init() {
	Version = C.GoString(C.TKRZW_PACKAGE_VERSION)
	OSName = C.GoString(C.TKRZW_OS_NAME)
	PageSize = int(C.TKRZW_PAGE_SIZE)
	Int64Min = int64(C.TKRZW_INT64MIN)
	Int64Max = int64(C.TKRZW_INT64MAX)
}

func get_memory_capacity() int64 {
	return int64(C.tkrzw_get_memory_capacity())
}

func get_memory_usage() int64 {
	return int64(C.tkrzw_get_memory_usage())
}

func primary_hash(data []byte, num_buckets uint64) uint64 {
	xdata := (*C.char)(C.CBytes(data))
	defer C.free(unsafe.Pointer(xdata))
	return uint64(C.tkrzw_primary_hash(xdata, C.int32_t(len(data)), C.uint64_t(num_buckets)))
}

func secondary_hash(data []byte, num_shards uint64) uint64 {
	xdata := (*C.char)(C.CBytes(data))
	defer C.free(unsafe.Pointer(xdata))
	return uint64(C.tkrzw_secondary_hash(xdata, C.int32_t(len(data)), C.uint64_t(num_shards)))
}

func edit_distance_lev(a string, b string, utf bool) int {
	xa := C.CString(a)
	defer C.free(unsafe.Pointer(xa))
	xb := C.CString(b)
	defer C.free(unsafe.Pointer(xb))
	return int(C.tkrzw_str_edit_distance_lev(xa, xb, C.bool(utf)))
}

func convert_status(res C.RES_STATUS) *Status {
	if res.message == nil {
		return NewStatus1(StatusCode(res.code))
	}
	defer C.free(unsafe.Pointer(res.message))
	return NewStatus2(StatusCode(res.code), C.GoString(res.message))
}

func join_params(params map[string]string) string {
	if params == nil {
		return ""
	}
	fields := make([]string, 0, 4)
	for name, value := range params {
		fields = append(fields, name+"="+value)
	}
	return strings.Join(fields, ",")
}

func future_free(future uintptr) {
	xfuture := (*C.TkrzwFuture)(unsafe.Pointer(future))
	C.tkrzw_future_free(xfuture)
}

func future_wait(future uintptr, timeout float64) bool {
	xfuture := (*C.TkrzwFuture)(unsafe.Pointer(future))
	return bool(C.tkrzw_future_wait(xfuture, C.double(timeout)))
}

func future_get(future uintptr) *Status {
	xfuture := (*C.TkrzwFuture)(unsafe.Pointer(future))
	res := C.do_future_get(xfuture)
	status := convert_status(res)
	return status
}

func future_get_bytes(future uintptr) ([]byte, *Status) {
	xfuture := (*C.TkrzwFuture)(unsafe.Pointer(future))
	res := C.do_future_get_bytes(xfuture)
	defer C.free(unsafe.Pointer(res.value_ptr))
	value := C.GoBytes(unsafe.Pointer(res.value_ptr), res.value_size)
	status := convert_status(res.status)
	return value, status
}

func future_get_str(future uintptr) (string, *Status) {
	xfuture := (*C.TkrzwFuture)(unsafe.Pointer(future))
	res := C.do_future_get_bytes(xfuture)
	defer C.free(unsafe.Pointer(res.value_ptr))
	value := C.GoStringN(res.value_ptr, res.value_size)
	status := convert_status(res.status)
	return value, status
}

func future_get_pair(future uintptr) ([]byte, []byte, *Status) {
	xfuture := (*C.TkrzwFuture)(unsafe.Pointer(future))
	res := C.do_future_get_str_pair(xfuture)
	defer C.free(unsafe.Pointer(res.str_pair))
	key := C.GoBytes(unsafe.Pointer(res.str_pair.key_ptr), res.str_pair.key_size)
	value := C.GoBytes(unsafe.Pointer(res.str_pair.value_ptr), res.str_pair.value_size)
	status := convert_status(res.status)
	return key, value, status
}

func future_get_pair_str(future uintptr) (string, string, *Status) {
	xfuture := (*C.TkrzwFuture)(unsafe.Pointer(future))
	res := C.do_future_get_str_pair(xfuture)
	defer C.free(unsafe.Pointer(res.str_pair))
	key := C.GoStringN(res.str_pair.key_ptr, res.str_pair.key_size)
	value := C.GoStringN(res.str_pair.value_ptr, res.str_pair.value_size)
	status := convert_status(res.status)
	return key, value, status
}

func future_get_array(future uintptr) ([][]byte, *Status) {
	xfuture := (*C.TkrzwFuture)(unsafe.Pointer(future))
	res := C.do_future_get_str_array(xfuture)
	defer C.tkrzw_free_str_array(res.values, res.num_values)
	values := make([][]byte, 0, res.num_values)
	value_ptr := uintptr(unsafe.Pointer(res.values))
	for i := C.int32_t(0); i < res.num_values; i++ {
		xvalue := (*C.TkrzwStr)(unsafe.Pointer(value_ptr))
		value := C.GoBytes(unsafe.Pointer(xvalue.ptr), xvalue.size)
		values = append(values, value)
		value_ptr += unsafe.Sizeof(C.TkrzwStr{})
	}
	status := convert_status(res.status)
	return values, status
}

func future_get_array_str(future uintptr) ([]string, *Status) {
	xfuture := (*C.TkrzwFuture)(unsafe.Pointer(future))
	res := C.do_future_get_str_array(xfuture)
	defer C.tkrzw_free_str_array(res.values, res.num_values)
	values := make([]string, 0, res.num_values)
	value_ptr := uintptr(unsafe.Pointer(res.values))
	for i := C.int32_t(0); i < res.num_values; i++ {
		xvalue := (*C.TkrzwStr)(unsafe.Pointer(value_ptr))
		value := C.GoStringN(xvalue.ptr, xvalue.size)
		values = append(values, value)
		value_ptr += unsafe.Sizeof(C.TkrzwStr{})
	}
	status := convert_status(res.status)
	return values, status
}

func future_get_map(future uintptr) (map[string][]byte, *Status) {
	xfuture := (*C.TkrzwFuture)(unsafe.Pointer(future))
	res := C.do_future_get_str_map(xfuture)
	defer C.tkrzw_free_str_map(res.records, res.num_records)
	records := make(map[string][]byte)
	rec_ptr := uintptr(unsafe.Pointer(res.records))
	for i := C.int32_t(0); i < res.num_records; i++ {
		elem := (*C.TkrzwKeyValuePair)(unsafe.Pointer(rec_ptr))
		key := C.GoStringN(elem.key_ptr, elem.key_size)
		value := C.GoBytes(unsafe.Pointer(elem.value_ptr), elem.value_size)
		records[key] = value
		rec_ptr += unsafe.Sizeof(C.TkrzwKeyValuePair{})
	}
	status := convert_status(res.status)
	return records, status
}

func future_get_map_str(future uintptr) (map[string]string, *Status) {
	xfuture := (*C.TkrzwFuture)(unsafe.Pointer(future))
	res := C.do_future_get_str_map(xfuture)
	defer C.tkrzw_free_str_map(res.records, res.num_records)
	records := make(map[string]string)
	rec_ptr := uintptr(unsafe.Pointer(res.records))
	for i := C.int32_t(0); i < res.num_records; i++ {
		elem := (*C.TkrzwKeyValuePair)(unsafe.Pointer(rec_ptr))
		key := C.GoStringN(elem.key_ptr, elem.key_size)
		value := C.GoStringN(elem.value_ptr, elem.value_size)
		records[key] = value
		rec_ptr += unsafe.Sizeof(C.TkrzwKeyValuePair{})
	}
	status := convert_status(res.status)
	return records, status
}

func future_get_int(future uintptr) (int64, *Status) {
	xfuture := (*C.TkrzwFuture)(unsafe.Pointer(future))
	res := C.do_future_get_int(xfuture)
	status := convert_status(res.status)
	return int64(res.num), status
}

func dbm_open(path string, writable bool, params map[string]string) (uintptr, *Status) {
	xpath := C.CString(path)
	defer C.free(unsafe.Pointer(xpath))
	xparams := C.CString(join_params(params))
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

func dbm_process(dbm uintptr, key []byte, proc RecordProcessor, writable bool) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	proc_up := registerRecordProcessor(proc)
	defer deregisterRecordProcessor(proc_up)
	res := C.do_dbm_process(xdbm, xkey_ptr, C.int32_t(len(key)), proc_up, C.bool(writable))
	status := convert_status(res)
	return status
}

func dbm_check(dbm uintptr, key []byte) bool {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	return (bool)(C.tkrzw_dbm_check(xdbm, xkey_ptr, C.int32_t(len(key))))
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

func dbm_get_str(dbm uintptr, key []byte) (string, *Status) {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	res := C.do_dbm_get(xdbm, xkey_ptr, C.int32_t(len(key)))
	var value string
	if res.value_ptr != nil {
		defer C.free(unsafe.Pointer(res.value_ptr))
		value = C.GoStringN(res.value_ptr, res.value_size)
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

func dbm_get_multi_str(dbm uintptr, keys []string) map[string]string {
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
	records := make(map[string]string)
	rec_ptr := uintptr(unsafe.Pointer(xrecords))
	for i := C.int32_t(0); i < num_records; i++ {
		elem := (*C.TkrzwKeyValuePair)(unsafe.Pointer(rec_ptr))
		key := C.GoStringN(elem.key_ptr, elem.key_size)
		value := C.GoStringN(elem.value_ptr, elem.value_size)
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

func dbm_set_and_get(dbm uintptr, key []byte, value []byte, overwrite bool) ([]byte, *Status) {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	xvalue_ptr := (*C.char)(C.CBytes(value))
	defer C.free(unsafe.Pointer(xvalue_ptr))
	res := C.do_dbm_set_and_get(xdbm, xkey_ptr, C.int32_t(len(key)),
		xvalue_ptr, C.int32_t(len(value)), C.bool(overwrite))
	var old_value []byte = nil
	if res.value_ptr != nil {
		defer C.free(unsafe.Pointer(res.value_ptr))
		old_value = C.GoBytes(unsafe.Pointer(res.value_ptr), res.value_size)
	}
	status := convert_status(res.status)
	return old_value, status
}

func dbm_set_multi(dbm uintptr, records map[string][]byte, overwrite bool) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xrecs_size := len(records) * int(unsafe.Sizeof(C.TkrzwKeyValuePair{}))
	xrecs := (*C.TkrzwKeyValuePair)(unsafe.Pointer(C.malloc(C.size_t(xrecs_size + 1))))
	defer C.free_str_pairs(xrecs, C.int32_t(len(records)))
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

func dbm_remove_and_get(dbm uintptr, key []byte) ([]byte, *Status) {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	res := C.do_dbm_remove_and_get(xdbm, xkey_ptr, C.int32_t(len(key)))
	var old_value []byte = nil
	if res.value_ptr != nil {
		defer C.free(unsafe.Pointer(res.value_ptr))
		old_value = C.GoBytes(unsafe.Pointer(res.value_ptr), res.value_size)
	}
	status := convert_status(res.status)
	return old_value, status
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

func dbm_append_multi(dbm uintptr, records map[string][]byte, delim []byte) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xrecs_size := len(records) * int(unsafe.Sizeof(C.TkrzwKeyValuePair{}))
	xrecs := (*C.TkrzwKeyValuePair)(unsafe.Pointer(C.malloc(C.size_t(xrecs_size + 1))))
	defer C.free_str_pairs(xrecs, C.int32_t(len(records)))
	xrec_ptr := uintptr(unsafe.Pointer(xrecs))
	for key, value := range records {
		xrec := (*C.TkrzwKeyValuePair)(unsafe.Pointer(xrec_ptr))
		xrec.key_ptr = C.CString(key)
		xrec.key_size = C.int32_t(len(key))
		xrec.value_ptr = (*C.char)(C.CBytes(value))
		xrec.value_size = C.int32_t(len(value))
		xrec_ptr += unsafe.Sizeof(C.TkrzwKeyValuePair{})
	}
	xdelim_ptr := (*C.char)(C.CBytes(delim))
	defer C.free(unsafe.Pointer(xdelim_ptr))
	res := C.do_dbm_append_multi(
		xdbm, xrecs, C.int32_t(len(records)), xdelim_ptr, C.int32_t(len(delim)))
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
		if IsAnyBytes(expected) {
			xexpected_ptr = C.TKRZW_ANY_DATA
		} else {
			xexpected_ptr = (*C.char)(C.CBytes(expected))
			defer C.free(unsafe.Pointer(xexpected_ptr))
			xexpected_size = C.int32_t(len(expected))
		}
	}
	var xdesired_ptr *C.char
	var xdesired_size C.int32_t
	if desired != nil {
		if IsAnyBytes(desired) {
			xdesired_ptr = C.TKRZW_ANY_DATA
		} else {
			xdesired_ptr = (*C.char)(C.CBytes(desired))
			defer C.free(unsafe.Pointer(xdesired_ptr))
			xdesired_size = C.int32_t(len(desired))
		}
	}
	res := C.do_dbm_compare_exchange(
		xdbm, xkey_ptr, C.int32_t(len(key)),
		xexpected_ptr, xexpected_size, xdesired_ptr, xdesired_size)
	status := convert_status(res)
	return status
}

func dbm_compare_exchange_and_get(
	dbm uintptr, key []byte, expected []byte, desired []byte) ([]byte, *Status) {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	var xexpected_ptr *C.char
	var xexpected_size C.int32_t
	if expected != nil {
		if IsAnyBytes(expected) {
			xexpected_ptr = C.TKRZW_ANY_DATA
		} else {
			xexpected_ptr = (*C.char)(C.CBytes(expected))
			defer C.free(unsafe.Pointer(xexpected_ptr))
			xexpected_size = C.int32_t(len(expected))
		}
	}
	var xdesired_ptr *C.char
	var xdesired_size C.int32_t
	if desired != nil {
		if IsAnyBytes(desired) {
			xdesired_ptr = C.TKRZW_ANY_DATA
		} else {
			xdesired_ptr = (*C.char)(C.CBytes(desired))
			defer C.free(unsafe.Pointer(xdesired_ptr))
			xdesired_size = C.int32_t(len(desired))
		}
	}
	res := C.do_dbm_compare_exchange_and_get(
		xdbm, xkey_ptr, C.int32_t(len(key)),
		xexpected_ptr, xexpected_size, xdesired_ptr, xdesired_size)
	var actual []byte = nil
	if res.value_ptr != nil {
		defer C.free(unsafe.Pointer(res.value_ptr))
		actual = C.GoBytes(unsafe.Pointer(res.value_ptr), res.value_size)
	}
	status := convert_status(res.status)
	return actual, status
}

func dbm_increment(dbm uintptr, key []byte, inc int64, init int64) (int64, *Status) {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	res := C.do_dbm_increment(
		xdbm, xkey_ptr, C.int32_t(len(key)), C.int64_t(inc), C.int64_t(init))
	status := convert_status(res.status)
	return int64(res.num), status
}

func dbm_process_multi(dbm uintptr, pairs []KeyBytesProcPair, writable bool) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xpairs_size := len(pairs) * int(unsafe.Sizeof(C.TkrzwKeyProcPair{}))
	xpairs := (*C.TkrzwKeyProcPair)(unsafe.Pointer(C.malloc(C.size_t(xpairs_size + 1))))
	defer C.free(unsafe.Pointer(xpairs))
	xpair_ptr := uintptr(unsafe.Pointer(xpairs))
	for _, pair := range pairs {
		xpair := (*C.TkrzwKeyProcPair)(unsafe.Pointer(xpair_ptr))
		xkey_ptr := (*C.char)(C.CBytes(pair.Key))
		defer C.free(unsafe.Pointer(xkey_ptr))
		proc_up := registerRecordProcessor(pair.Proc)
		defer deregisterRecordProcessor(proc_up)
		xpair.key_ptr = xkey_ptr
		xpair.key_size = C.int32_t(len(pair.Key))
		xpair.proc = nil
		xpair.proc_arg = proc_up
		xpair_ptr += unsafe.Sizeof(C.TkrzwKeyProcPair{})
	}
	res := C.do_dbm_process_multi(xdbm, xpairs, C.int32_t(len(pairs)), C.bool(writable))
	status := convert_status(res)
	return status
}

func dbm_compare_exchange_multi(
	dbm uintptr, expected []KeyValuePair, desired []KeyValuePair) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xexpected_size := len(expected) * int(unsafe.Sizeof(C.TkrzwKeyValuePair{}))
	xexpected := (*C.TkrzwKeyValuePair)(unsafe.Pointer(C.malloc(C.size_t(xexpected_size + 1))))
	defer C.free_str_pairs(xexpected, C.int32_t(len(expected)))
	xexp_ptr := uintptr(unsafe.Pointer(xexpected))
	for _, pair := range expected {
		xexp := (*C.TkrzwKeyValuePair)(unsafe.Pointer(xexp_ptr))
		xexp.key_ptr = (*C.char)(C.CBytes(pair.Key))
		xexp.key_size = C.int32_t(len(pair.Key))
		if pair.Value == nil {
			xexp.value_ptr = nil
			xexp.value_size = 0
		} else if IsAnyBytes(pair.Value) {
			xexp.value_ptr = C.TKRZW_ANY_DATA
			xexp.value_size = 0
		} else {
			xexp.value_ptr = (*C.char)(C.CBytes(pair.Value))
			xexp.value_size = C.int32_t(len(pair.Value))
		}
		xexp_ptr += unsafe.Sizeof(C.TkrzwKeyValuePair{})
	}
	xdesired_size := len(desired) * int(unsafe.Sizeof(C.TkrzwKeyValuePair{}))
	xdesired := (*C.TkrzwKeyValuePair)(unsafe.Pointer(C.malloc(C.size_t(xdesired_size + 1))))
	defer C.free_str_pairs(xdesired, C.int32_t(len(desired)))
	xdes_ptr := uintptr(unsafe.Pointer(xdesired))
	for _, pair := range desired {
		xdes := (*C.TkrzwKeyValuePair)(unsafe.Pointer(xdes_ptr))
		xdes.key_ptr = (*C.char)(C.CBytes(pair.Key))
		xdes.key_size = C.int32_t(len(pair.Key))
		if pair.Value == nil {
			xdes.value_ptr = nil
			xdes.value_size = 0
		} else {
			xdes.value_ptr = (*C.char)(C.CBytes(pair.Value))
			xdes.value_size = C.int32_t(len(pair.Value))
		}
		xdes_ptr += unsafe.Sizeof(C.TkrzwKeyValuePair{})
	}
	res := C.do_dbm_compare_exchange_multi(
		xdbm, xexpected, C.int32_t(len(expected)), xdesired, C.int32_t(len(desired)))
	status := convert_status(res)
	return status
}

func dbm_rekey(dbm uintptr, old_key []byte, new_key []byte,
	overwrite bool, copying bool) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xold_key_ptr := (*C.char)(C.CBytes(old_key))
	defer C.free(unsafe.Pointer(xold_key_ptr))
	xnew_key_ptr := (*C.char)(C.CBytes(new_key))
	defer C.free(unsafe.Pointer(xnew_key_ptr))
	res := C.do_dbm_rekey(xdbm, xold_key_ptr, C.int32_t(len(old_key)),
		xnew_key_ptr, C.int32_t(len(new_key)), C.bool(overwrite), C.bool(copying))
	status := convert_status(res)
	return status
}

func dbm_pop_first(dbm uintptr) ([]byte, []byte, *Status) {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	res := C.do_dbm_pop_first(xdbm)
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

func dbm_push_last(dbm uintptr, value []byte, wtime float64) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xvalue_ptr := (*C.char)(C.CBytes(value))
	defer C.free(unsafe.Pointer(xvalue_ptr))
	res := C.do_dbm_push_last(xdbm, xvalue_ptr, C.int32_t(len(value)), C.double(wtime))
	status := convert_status(res)
	return status
}

func dbm_process_each(dbm uintptr, proc RecordProcessor, writable bool) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	proc_up := registerRecordProcessor(proc)
	defer deregisterRecordProcessor(proc_up)
	res := C.do_dbm_process_each(xdbm, proc_up, C.bool(writable))
	status := convert_status(res)
	return status
}

func dbm_count(dbm uintptr) (int64, *Status) {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	res := C.do_dbm_count(xdbm)
	status := convert_status(res.status)
	return int64(res.num), status
}

func dbm_get_file_size(dbm uintptr) (int64, *Status) {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	res := C.do_dbm_get_file_size(xdbm)
	status := convert_status(res.status)
	return int64(res.num), status
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

func dbm_get_timestamp(dbm uintptr) (float64, *Status) {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	res := C.do_dbm_get_timestamp(xdbm)
	status := convert_status(res.status)
	return float64(res.num), status
}

func dbm_clear(dbm uintptr) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	res := C.do_dbm_clear(xdbm)
	status := convert_status(res)
	return status
}

func dbm_rebuild(dbm uintptr, params map[string]string) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xparams := C.CString(join_params(params))
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

func dbm_synchronize(dbm uintptr, hard bool, params map[string]string) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xparams := C.CString(join_params(params))
	defer C.free(unsafe.Pointer(xparams))
	res := C.do_dbm_synchronize(xdbm, C.bool(hard), xparams)
	status := convert_status(res)
	return status
}

func dbm_copy_file_data(dbm uintptr, dest_path string, sync_hard bool) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xdest_path := C.CString(dest_path)
	defer C.free(unsafe.Pointer(xdest_path))
	res := C.do_dbm_copy_file_data(xdbm, xdest_path, C.bool(sync_hard))
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

func dbm_export_to_flat_records(dbm uintptr, dest_file uintptr) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xdest_file := (*C.TkrzwFile)(unsafe.Pointer(dest_file))
	res := C.do_dbm_export_to_flat_records(xdbm, xdest_file)
	status := convert_status(res)
	return status
}

func dbm_import_from_flat_records(dbm uintptr, src_file uintptr) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xsrc_file := (*C.TkrzwFile)(unsafe.Pointer(src_file))
	res := C.do_dbm_import_from_flat_records(xdbm, xsrc_file)
	status := convert_status(res)
	return status
}

func dbm_export_keys_as_lines(dbm uintptr, dest_file uintptr) *Status {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xdest_file := (*C.TkrzwFile)(unsafe.Pointer(dest_file))
	res := C.do_dbm_export_keys_as_lines(xdbm, xdest_file)
	status := convert_status(res)
	return status
}

func dbm_inspect(dbm uintptr) map[string]string {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	res := C.do_dbm_inspect(xdbm)
	defer C.tkrzw_free_str_map(res.records, res.num_records)
	rec_ptr := uintptr(unsafe.Pointer(res.records))
	records := make(map[string]string)
	for i := C.int32_t(0); i < res.num_records; i++ {
		elem := (*C.TkrzwKeyValuePair)(unsafe.Pointer(rec_ptr))
		name := C.GoStringN(elem.key_ptr, elem.key_size)
		value := C.GoStringN(elem.value_ptr, elem.value_size)
		records[name] = value
		rec_ptr += unsafe.Sizeof(C.TkrzwKeyValuePair{})
	}
	return records
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

func dbm_search(dbm uintptr, mode string, pattern string, capacity int) []string {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	xmode := C.CString(mode)
	defer C.free(unsafe.Pointer(xmode))
	xpattern := C.CString(pattern)
	defer C.free(unsafe.Pointer(xpattern))
	var num_matched C.int32_t = 0
	xkeys := C.tkrzw_dbm_search(
		xdbm, xmode, xpattern, C.int32_t(len(pattern)), C.int32_t(capacity), &num_matched)
	keys := make([]string, 0, num_matched)
	key_ptr := uintptr(unsafe.Pointer(xkeys))
	for i := C.int32_t(0); i < num_matched; i++ {
		xkey := (*C.TkrzwStr)(unsafe.Pointer(key_ptr))
		key := C.GoStringN(xkey.ptr, xkey.size)
		keys = append(keys, key)
		key_ptr += unsafe.Sizeof(C.TkrzwStr{})
	}
	return keys
}

func dbm_make_iterator(dbm uintptr) uintptr {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	return uintptr(unsafe.Pointer(C.tkrzw_dbm_make_iterator(xdbm)))
}

func dbm_restore_database(
	old_file_path string, new_file_path string, class_name string,
	end_offset int64, cipher_key string) *Status {
	xold_file_path := C.CString(old_file_path)
	defer C.free(unsafe.Pointer(xold_file_path))
	xnew_file_path := C.CString(new_file_path)
	defer C.free(unsafe.Pointer(xnew_file_path))
	xclass_name := C.CString(class_name)
	defer C.free(unsafe.Pointer(xclass_name))
	xcipher_key := (*C.char)(C.CString(cipher_key))
	defer C.free(unsafe.Pointer(xcipher_key))
	res := C.do_dbm_restore_database(
		xold_file_path, xnew_file_path, xclass_name, C.int64_t(end_offset), xcipher_key)
	status := convert_status(res)
	return status
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

func dbm_iter_set(iter uintptr, value []byte) *Status {
	xiter := (*C.TkrzwDBMIter)(unsafe.Pointer(iter))
	xvalue_ptr := (*C.char)(C.CBytes(value))
	defer C.free(unsafe.Pointer(xvalue_ptr))
	res := C.do_dbm_iter_set(xiter, xvalue_ptr, C.int32_t(len(value)))
	status := convert_status(res)
	return status
}

func dbm_iter_remove(iter uintptr) *Status {
	xiter := (*C.TkrzwDBMIter)(unsafe.Pointer(iter))
	res := C.do_dbm_iter_remove(xiter)
	status := convert_status(res)
	return status
}

func dbm_iter_step(iter uintptr) ([]byte, []byte, *Status) {
	xiter := (*C.TkrzwDBMIter)(unsafe.Pointer(iter))
	res := C.do_dbm_iter_step(xiter)
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

func async_dbm_new(dbm uintptr, num_worker_threads int) uintptr {
	xdbm := (*C.TkrzwDBM)(unsafe.Pointer(dbm))
	return uintptr(unsafe.Pointer(C.tkrzw_async_dbm_new(xdbm, C.int32_t(num_worker_threads))))
}

func async_dbm_free(async uintptr) {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
	C.tkrzw_async_dbm_free(xasync)
}

func async_dbm_get(async uintptr, key []byte) *Future {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	xfuture := C.tkrzw_async_dbm_get(xasync, xkey_ptr, C.int32_t(len(key)))
	return &Future{uintptr(unsafe.Pointer(xfuture))}
}

func async_dbm_get_multi(async uintptr, keys []string) *Future {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
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
	xfuture := C.tkrzw_async_dbm_get_multi(xasync, xkeys, C.int32_t(len(keys)))
	return &Future{uintptr(unsafe.Pointer(xfuture))}
}

func async_dbm_set(async uintptr, key []byte, value []byte, overwrite bool) *Future {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	xvalue_ptr := (*C.char)(C.CBytes(value))
	defer C.free(unsafe.Pointer(xvalue_ptr))
	xfuture := C.tkrzw_async_dbm_set(xasync, xkey_ptr, C.int32_t(len(key)),
		xvalue_ptr, C.int32_t(len(value)), C.bool(overwrite))
	return &Future{uintptr(unsafe.Pointer(xfuture))}
}

func async_dbm_set_multi(async uintptr, records map[string][]byte, overwrite bool) *Future {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
	xrecs_size := len(records) * int(unsafe.Sizeof(C.TkrzwKeyValuePair{}))
	xrecs := (*C.TkrzwKeyValuePair)(unsafe.Pointer(C.malloc(C.size_t(xrecs_size + 1))))
	defer C.free_str_pairs(xrecs, C.int32_t(len(records)))
	xrec_ptr := uintptr(unsafe.Pointer(xrecs))
	for key, value := range records {
		xrec := (*C.TkrzwKeyValuePair)(unsafe.Pointer(xrec_ptr))
		xrec.key_ptr = C.CString(key)
		xrec.key_size = C.int32_t(len(key))
		xrec.value_ptr = (*C.char)(C.CBytes(value))
		xrec.value_size = C.int32_t(len(value))
		xrec_ptr += unsafe.Sizeof(C.TkrzwKeyValuePair{})
	}
	xfuture := C.tkrzw_async_dbm_set_multi(
		xasync, xrecs, C.int32_t(len(records)), C.bool(overwrite))
	return &Future{uintptr(unsafe.Pointer(xfuture))}
}

func async_dbm_remove(async uintptr, key []byte) *Future {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	xfuture := C.tkrzw_async_dbm_remove(xasync, xkey_ptr, C.int32_t(len(key)))
	return &Future{uintptr(unsafe.Pointer(xfuture))}
}

func async_dbm_remove_multi(async uintptr, keys []string) *Future {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
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
	xfuture := C.tkrzw_async_dbm_remove_multi(xasync, xkeys, C.int32_t(len(keys)))
	return &Future{uintptr(unsafe.Pointer(xfuture))}
}

func async_dbm_append(async uintptr, key []byte, value []byte, delim []byte) *Future {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	xvalue_ptr := (*C.char)(C.CBytes(value))
	defer C.free(unsafe.Pointer(xvalue_ptr))
	xdelim_ptr := (*C.char)(C.CBytes(delim))
	defer C.free(unsafe.Pointer(xdelim_ptr))
	xfuture := C.tkrzw_async_dbm_append(xasync, xkey_ptr, C.int32_t(len(key)),
		xvalue_ptr, C.int32_t(len(value)), xdelim_ptr, C.int32_t(len(delim)))
	return &Future{uintptr(unsafe.Pointer(xfuture))}
}

func async_dbm_append_multi(async uintptr, records map[string][]byte, delim []byte) *Future {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
	xrecs_size := len(records) * int(unsafe.Sizeof(C.TkrzwKeyValuePair{}))
	xrecs := (*C.TkrzwKeyValuePair)(unsafe.Pointer(C.malloc(C.size_t(xrecs_size + 1))))
	defer C.free_str_pairs(xrecs, C.int32_t(len(records)))
	xrec_ptr := uintptr(unsafe.Pointer(xrecs))
	for key, value := range records {
		xrec := (*C.TkrzwKeyValuePair)(unsafe.Pointer(xrec_ptr))
		xrec.key_ptr = C.CString(key)
		xrec.key_size = C.int32_t(len(key))
		xrec.value_ptr = (*C.char)(C.CBytes(value))
		xrec.value_size = C.int32_t(len(value))
		xrec_ptr += unsafe.Sizeof(C.TkrzwKeyValuePair{})
	}
	xdelim_ptr := (*C.char)(C.CBytes(delim))
	defer C.free(unsafe.Pointer(xdelim_ptr))
	xfuture := C.tkrzw_async_dbm_append_multi(
		xasync, xrecs, C.int32_t(len(records)), xdelim_ptr, C.int32_t(len(delim)))
	return &Future{uintptr(unsafe.Pointer(xfuture))}
}

func async_dbm_compare_exchange(
	async uintptr, key []byte, expected []byte, desired []byte) *Future {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	var xexpected_ptr *C.char
	var xexpected_size C.int32_t
	if expected != nil {
		if IsAnyBytes(expected) {
			xexpected_ptr = C.TKRZW_ANY_DATA
		} else {
			xexpected_ptr = (*C.char)(C.CBytes(expected))
			defer C.free(unsafe.Pointer(xexpected_ptr))
			xexpected_size = C.int32_t(len(expected))
		}
	}
	var xdesired_ptr *C.char
	var xdesired_size C.int32_t
	if desired != nil {
		if IsAnyBytes(desired) {
			xdesired_ptr = C.TKRZW_ANY_DATA
		} else {
			xdesired_ptr = (*C.char)(C.CBytes(desired))
			defer C.free(unsafe.Pointer(xdesired_ptr))
			xdesired_size = C.int32_t(len(desired))
		}
	}
	xfuture := C.tkrzw_async_dbm_compare_exchange(xasync, xkey_ptr, C.int32_t(len(key)),
		xexpected_ptr, xexpected_size, xdesired_ptr, xdesired_size)
	return &Future{uintptr(unsafe.Pointer(xfuture))}
}

func async_dbm_increment(async uintptr, key []byte, inc int64, init int64) *Future {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	xfuture := C.tkrzw_async_dbm_increment(
		xasync, xkey_ptr, C.int32_t(len(key)), C.int64_t(inc), C.int64_t(init))
	return &Future{uintptr(unsafe.Pointer(xfuture))}
}

func async_dbm_compare_exchange_multi(
	async uintptr, expected []KeyValuePair, desired []KeyValuePair) *Future {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
	xexpected_size := len(expected) * int(unsafe.Sizeof(C.TkrzwKeyValuePair{}))
	xexpected := (*C.TkrzwKeyValuePair)(unsafe.Pointer(C.malloc(C.size_t(xexpected_size + 1))))
	defer C.free_str_pairs(xexpected, C.int32_t(len(expected)))
	xexp_ptr := uintptr(unsafe.Pointer(xexpected))
	for _, pair := range expected {
		xexp := (*C.TkrzwKeyValuePair)(unsafe.Pointer(xexp_ptr))
		xexp.key_ptr = (*C.char)(C.CBytes(pair.Key))
		xexp.key_size = C.int32_t(len(pair.Key))
		if pair.Value == nil {
			xexp.value_ptr = nil
			xexp.value_size = 0
		} else if IsAnyBytes(pair.Value) {
			xexp.value_ptr = C.TKRZW_ANY_DATA
			xexp.value_size = 0
		} else {
			xexp.value_ptr = (*C.char)(C.CBytes(pair.Value))
			xexp.value_size = C.int32_t(len(pair.Value))
		}
		xexp_ptr += unsafe.Sizeof(C.TkrzwKeyValuePair{})
	}
	xdesired_size := len(desired) * int(unsafe.Sizeof(C.TkrzwKeyValuePair{}))
	xdesired := (*C.TkrzwKeyValuePair)(unsafe.Pointer(C.malloc(C.size_t(xdesired_size + 1))))
	defer C.free_str_pairs(xdesired, C.int32_t(len(desired)))
	xdes_ptr := uintptr(unsafe.Pointer(xdesired))
	for _, pair := range desired {
		xdes := (*C.TkrzwKeyValuePair)(unsafe.Pointer(xdes_ptr))
		xdes.key_ptr = (*C.char)(C.CBytes(pair.Key))
		xdes.key_size = C.int32_t(len(pair.Key))
		if pair.Value == nil {
			xdes.value_ptr = nil
			xdes.value_size = 0
		} else {
			xdes.value_ptr = (*C.char)(C.CBytes(pair.Value))
			xdes.value_size = C.int32_t(len(pair.Value))
		}
		xdes_ptr += unsafe.Sizeof(C.TkrzwKeyValuePair{})
	}
	xfuture := C.tkrzw_async_dbm_compare_exchange_multi(
		xasync, xexpected, C.int32_t(len(expected)), xdesired, C.int32_t(len(desired)))
	return &Future{uintptr(unsafe.Pointer(xfuture))}
}

func async_dbm_rekey(async uintptr, old_key []byte, new_key []byte,
	overwrite bool, copying bool) *Future {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
	xold_key_ptr := (*C.char)(C.CBytes(old_key))
	defer C.free(unsafe.Pointer(xold_key_ptr))
	xnew_key_ptr := (*C.char)(C.CBytes(new_key))
	defer C.free(unsafe.Pointer(xnew_key_ptr))
	xfuture := C.tkrzw_async_dbm_rekey(xasync, xold_key_ptr, C.int32_t(len(old_key)),
		xnew_key_ptr, C.int32_t(len(new_key)), C.bool(overwrite), C.bool(copying))
	return &Future{uintptr(unsafe.Pointer(xfuture))}
}

func async_dbm_pop_first(async uintptr) *Future {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
	xfuture := C.tkrzw_async_dbm_pop_first(xasync)
	return &Future{uintptr(unsafe.Pointer(xfuture))}
}

func async_dbm_push_last(async uintptr, value []byte, wtime float64) *Future {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
	xvalue_ptr := (*C.char)(C.CBytes(value))
	defer C.free(unsafe.Pointer(xvalue_ptr))
	xfuture := C.tkrzw_async_dbm_push_last(
		xasync, xvalue_ptr, C.int32_t(len(value)), C.double(wtime))
	return &Future{uintptr(unsafe.Pointer(xfuture))}
}

func async_dbm_clear(async uintptr) *Future {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
	xfuture := C.tkrzw_async_dbm_clear(xasync)
	return &Future{uintptr(unsafe.Pointer(xfuture))}
}

func async_dbm_rebuild(async uintptr, params map[string]string) *Future {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
	xparams := C.CString(join_params(params))
	defer C.free(unsafe.Pointer(xparams))
	xfuture := C.tkrzw_async_dbm_rebuild(xasync, xparams)
	return &Future{uintptr(unsafe.Pointer(xfuture))}
}

func async_dbm_synchronize(async uintptr, hard bool, params map[string]string) *Future {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
	xparams := C.CString(join_params(params))
	defer C.free(unsafe.Pointer(xparams))
	xfuture := C.tkrzw_async_dbm_synchronize(xasync, C.bool(hard), xparams)
	return &Future{uintptr(unsafe.Pointer(xfuture))}
}

func async_dbm_copy_file_data(async uintptr, dest_path string, sync_hard bool) *Future {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
	xdest_path := C.CString(dest_path)
	defer C.free(unsafe.Pointer(xdest_path))
	xfuture := C.tkrzw_async_dbm_copy_file_data(xasync, xdest_path, C.bool(sync_hard))
	return &Future{uintptr(unsafe.Pointer(xfuture))}
}

func async_dbm_export(async uintptr, dest_dbm uintptr) *Future {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
	xdest_dbm := (*C.TkrzwDBM)(unsafe.Pointer(dest_dbm))
	xfuture := C.tkrzw_async_dbm_export(xasync, xdest_dbm)
	return &Future{uintptr(unsafe.Pointer(xfuture))}
}

func async_dbm_export_to_flat_records(async uintptr, dest_file uintptr) *Future {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
	xdest_file := (*C.TkrzwFile)(unsafe.Pointer(dest_file))
	xfuture := C.tkrzw_async_dbm_export_to_flat_records(xasync, xdest_file)
	return &Future{uintptr(unsafe.Pointer(xfuture))}
}

func async_dbm_import_from_flat_records(async uintptr, src_file uintptr) *Future {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
	xsrc_file := (*C.TkrzwFile)(unsafe.Pointer(src_file))
	xfuture := C.tkrzw_async_dbm_import_from_flat_records(xasync, xsrc_file)
	return &Future{uintptr(unsafe.Pointer(xfuture))}
}

func async_dbm_search(async uintptr, mode string, pattern string, capacity int) *Future {
	xasync := (*C.TkrzwAsyncDBM)(unsafe.Pointer(async))
	xmode := C.CString(mode)
	defer C.free(unsafe.Pointer(xmode))
	xpattern := C.CString(pattern)
	defer C.free(unsafe.Pointer(xpattern))
	xfuture := C.tkrzw_async_dbm_search(
		xasync, xmode, xpattern, C.int32_t(len(pattern)), C.int32_t(capacity))
	return &Future{uintptr(unsafe.Pointer(xfuture))}
}

func file_open(path string, writable bool, params map[string]string) (uintptr, *Status) {
	xpath := C.CString(path)
	defer C.free(unsafe.Pointer(xpath))
	xparams := C.CString(join_params(params))
	defer C.free(unsafe.Pointer(xparams))
	res := C.do_file_open(xpath, C.bool(writable), xparams)
	status := convert_status(res.status)
	return uintptr(unsafe.Pointer(res.file)), status
}

func file_close(file uintptr) *Status {
	xfile := (*C.TkrzwFile)(unsafe.Pointer(file))
	res := C.do_file_close(xfile)
	status := convert_status(res)
	return status
}

func file_read(file uintptr, off int64, size int64) ([]byte, *Status) {
	xfile := (*C.TkrzwFile)(unsafe.Pointer(file))
	xdata := (*C.char)(unsafe.Pointer(C.malloc(C.size_t(size + 1))))
	defer C.free(unsafe.Pointer(xdata))
	res := C.do_file_read(xfile, C.int64_t(off), xdata, C.size_t(size))
	var data []byte = nil
	if res.code == C.TKRZW_STATUS_SUCCESS {
		data = C.GoBytes(unsafe.Pointer(xdata), C.int(size))
	}
	status := convert_status(res)
	return data, status
}

func file_write(file uintptr, off int64, data []byte) *Status {
	xfile := (*C.TkrzwFile)(unsafe.Pointer(file))
	xdata := (*C.char)(C.CBytes(data))
	defer C.free(unsafe.Pointer(xdata))
	res := C.do_file_write(xfile, C.int64_t(off), xdata, C.size_t(len(data)))
	status := convert_status(res)
	return status
}

func file_append(file uintptr, data []byte) (int64, *Status) {
	xfile := (*C.TkrzwFile)(unsafe.Pointer(file))
	xdata := (*C.char)(C.CBytes(data))
	defer C.free(unsafe.Pointer(xdata))
	res := C.do_file_append(xfile, xdata, C.size_t(len(data)))
	status := convert_status(res.status)
	return int64(res.num), status
}

func file_truncate(file uintptr, size int64) *Status {
	xfile := (*C.TkrzwFile)(unsafe.Pointer(file))
	res := C.do_file_truncate(xfile, C.int64_t(size))
	status := convert_status(res.status)
	return status
}

func file_synchronize(file uintptr, hard bool, off int64, size int64) *Status {
	xfile := (*C.TkrzwFile)(unsafe.Pointer(file))
	res := C.do_file_synchronize(xfile, C.bool(hard), C.int64_t(off), C.int64_t(size))
	status := convert_status(res.status)
	return status
}

func file_get_size(file uintptr) (int64, *Status) {
	xfile := (*C.TkrzwFile)(unsafe.Pointer(file))
	res := C.do_file_get_size(xfile)
	status := convert_status(res.status)
	return int64(res.num), status
}

func file_get_path(file uintptr) (string, *Status) {
	xfile := (*C.TkrzwFile)(unsafe.Pointer(file))
	res := C.do_file_get_path(xfile)
	var path string
	if res.str != nil {
		defer C.free(unsafe.Pointer(res.str))
		path = C.GoString(res.str)
	}
	status := convert_status(res.status)
	return path, status
}

func file_search(file uintptr, mode string, pattern string, capacity int) []string {
	xfile := (*C.TkrzwFile)(unsafe.Pointer(file))
	xmode := C.CString(mode)
	defer C.free(unsafe.Pointer(xmode))
	xpattern := C.CString(pattern)
	defer C.free(unsafe.Pointer(xpattern))
	var num_matched C.int32_t = 0
	xlines := C.tkrzw_file_search(
		xfile, xmode, xpattern, C.int32_t(len(pattern)), C.int32_t(capacity), &num_matched)
	lines := make([]string, 0, num_matched)
	line_ptr := uintptr(unsafe.Pointer(xlines))
	for i := C.int32_t(0); i < num_matched; i++ {
		xline := (*C.TkrzwStr)(unsafe.Pointer(line_ptr))
		line := C.GoStringN(xline.ptr, xline.size)
		lines = append(lines, line)
		line_ptr += unsafe.Sizeof(C.TkrzwStr{})
	}
	return lines
}

func index_open(path string, writable bool, params map[string]string) (uintptr, *Status) {
	xpath := C.CString(path)
	defer C.free(unsafe.Pointer(xpath))
	xparams := C.CString(join_params(params))
	defer C.free(unsafe.Pointer(xparams))
	res := C.do_index_open(xpath, C.bool(writable), xparams)
	status := convert_status(res.status)
	return uintptr(unsafe.Pointer(res.index)), status
}

func index_close(index uintptr) *Status {
	xindex := (*C.TkrzwIndex)(unsafe.Pointer(index))
	res := C.do_index_close(xindex)
	status := convert_status(res)
	return status
}

func index_check(index uintptr, key []byte, value []byte) bool {
	xindex := (*C.TkrzwIndex)(unsafe.Pointer(index))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	xvalue_ptr := (*C.char)(C.CBytes(value))
	defer C.free(unsafe.Pointer(xvalue_ptr))
	return (bool)(C.tkrzw_index_check(xindex, xkey_ptr, C.int32_t(len(key)),
		xvalue_ptr, C.int32_t(len(value))))
}

func index_get_values(index uintptr, key []byte, max int) []string {
	xindex := (*C.TkrzwIndex)(unsafe.Pointer(index))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	var num_values C.int32_t = 0
	xvalues := C.tkrzw_index_get_values(
		xindex, xkey_ptr, C.int32_t(len(key)), C.int32_t(max), &num_values)
	values := make([]string, 0, num_values)
	value_ptr := uintptr(unsafe.Pointer(xvalues))
	for i := C.int32_t(0); i < num_values; i++ {
		xvalue := (*C.TkrzwStr)(unsafe.Pointer(value_ptr))
		value := C.GoStringN(xvalue.ptr, xvalue.size)
		values = append(values, value)
		value_ptr += unsafe.Sizeof(C.TkrzwStr{})
	}
	return values
}

func index_add(index uintptr, key []byte, value []byte) bool {
	xindex := (*C.TkrzwIndex)(unsafe.Pointer(index))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	xvalue_ptr := (*C.char)(C.CBytes(value))
	defer C.free(unsafe.Pointer(xvalue_ptr))
	return (bool)(C.tkrzw_index_add(xindex, xkey_ptr, C.int32_t(len(key)),
		xvalue_ptr, C.int32_t(len(value))))
}

func index_remove(index uintptr, key []byte, value []byte) bool {
	xindex := (*C.TkrzwIndex)(unsafe.Pointer(index))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	xvalue_ptr := (*C.char)(C.CBytes(value))
	defer C.free(unsafe.Pointer(xvalue_ptr))
	return (bool)(C.tkrzw_index_remove(xindex, xkey_ptr, C.int32_t(len(key)),
		xvalue_ptr, C.int32_t(len(value))))
}

func index_count(index uintptr) int64 {
	xindex := (*C.TkrzwIndex)(unsafe.Pointer(index))
	return (int64)(C.tkrzw_index_count(xindex))
}

func index_clear(index uintptr) bool {
	xindex := (*C.TkrzwIndex)(unsafe.Pointer(index))
	return (bool)(C.tkrzw_index_clear(xindex))
}

func index_rebuild(index uintptr) bool {
	xindex := (*C.TkrzwIndex)(unsafe.Pointer(index))
	return (bool)(C.tkrzw_index_rebuild(xindex))
}

func index_synchronize(index uintptr, hard bool) bool {
	xindex := (*C.TkrzwIndex)(unsafe.Pointer(index))
	return (bool)(C.tkrzw_index_synchronize(xindex, C.bool(hard)))
}

func index_is_writable(index uintptr) bool {
	xindex := (*C.TkrzwIndex)(unsafe.Pointer(index))
	return (bool)(C.tkrzw_index_is_writable(xindex))
}

func index_make_iterator(index uintptr) uintptr {
	xindex := (*C.TkrzwIndex)(unsafe.Pointer(index))
	return uintptr(unsafe.Pointer(C.tkrzw_index_make_iterator(xindex)))
}

func index_iter_free(iter uintptr) {
	xiter := (*C.TkrzwIndexIter)(unsafe.Pointer(iter))
	C.tkrzw_index_iter_free(xiter)
}

func index_iter_first(iter uintptr) {
	xiter := (*C.TkrzwIndexIter)(unsafe.Pointer(iter))
	C.tkrzw_index_iter_first(xiter)
}

func index_iter_last(iter uintptr) {
	xiter := (*C.TkrzwIndexIter)(unsafe.Pointer(iter))
	C.tkrzw_index_iter_last(xiter)
}

func index_iter_jump(iter uintptr, key []byte) {
	xiter := (*C.TkrzwIndexIter)(unsafe.Pointer(iter))
	xkey_ptr := (*C.char)(C.CBytes(key))
	defer C.free(unsafe.Pointer(xkey_ptr))
	C.tkrzw_index_iter_jump(xiter, xkey_ptr, C.int32_t(len(key)))
}

func index_iter_get(iter uintptr) ([]byte, []byte, bool) {
	xiter := (*C.TkrzwIndexIter)(unsafe.Pointer(iter))
	res := C.do_index_iter_get(xiter)
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
	return key, value, bool(res.status)
}

func index_iter_next(iter uintptr) {
	xiter := (*C.TkrzwIndexIter)(unsafe.Pointer(iter))
	C.tkrzw_index_iter_next(xiter)
}

func index_iter_previous(iter uintptr) {
	xiter := (*C.TkrzwIndexIter)(unsafe.Pointer(iter))
	C.tkrzw_index_iter_previous(xiter)
}

// END OF FILE
