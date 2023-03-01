/*************************************************************************************************
 * Utility interface
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

// #include <stdlib.h>
import "C"

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

// The special bytes value to remove a record.
var RemoveBytes = []byte("\x00[REMOVE]\x00")

// The special string value to remove a record.
var RemoveString = string([]byte("\x00[REMOVE]\x00"))

// The special bytes value for no-operation or any data.
var AnyBytes = []byte("\x00[ANY]\x00")

// The special string value for no-operation or any data.
var AnyString = string([]byte("\x00[ANY]\x00"))

// The special string value for non-existing data.
var NilString = string([]byte("\x00[NIL]\x00"))

// Converts any object into a string.
//
// @param x The object to convert.
// @return The result string.
func ToString(x interface{}) string {
	switch x := x.(type) {
	case []byte:
		return string(x)
	case string:
		return x
	case bool:
		if x {
			return "true"
		}
		return "false"
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		return fmt.Sprintf("%d", x)
	case uintptr:
		return fmt.Sprintf("0x%016x", x)
	case float32, float64, complex64, complex128:
		return fmt.Sprintf("%.6f", x)
	case error:
		return x.Error()
	}
	if x == nil {
		return ""
	}
	if str, ok := x.(fmt.Stringer); ok {
		return str.String()
	}
	return fmt.Sprintf("#<%T>", x)
}

// Converts any object into a byte array.
//
// @param x The object to convert.
// @return The result byte array.
func ToByteArray(x interface{}) []byte {
	switch x := x.(type) {
	case []byte:
		return x
	case string:
		return []byte(x)
	case bool:
		if x {
			return []byte("true")
		}
		return []byte("false")
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		return []byte(fmt.Sprintf("%d", x))
	case uintptr:
		return []byte(fmt.Sprintf("0x%016x", x))
	case float32, float64, complex64, complex128:
		return []byte(fmt.Sprintf("%.6f", x))
	case error:
		return []byte(x.Error())
	}
	if x == nil {
		return make([]byte, 0)
	}
	if str, ok := x.(fmt.Stringer); ok {
		return []byte(str.String())
	}
	return []byte(fmt.Sprintf("#<%T>", x))
}

// Converts any object into an integer.
//
// @param x The object to convert.
// @return The result integer.
func ToInt(value interface{}) int64 {
	switch value := value.(type) {
	case int:
		return int64(value)
	case uint:
		return int64(value)
	case int8:
		return int64(value)
	case uint8:
		return int64(value)
	case int16:
		return int64(value)
	case uint16:
		return int64(value)
	case int32:
		return int64(value)
	case uint32:
		return int64(value)
	case int64:
		return int64(value)
	case uint64:
		return int64(value)
	case float32:
		return int64(value)
	case float64:
		return int64(value)
	case string:
		int_value, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			return int64(int_value)
		}
		float_value, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return int64(float_value)
		}
	case []byte:
		int_value, err := strconv.ParseInt(string(value), 10, 64)
		if err == nil {
			return int64(int_value)
		}
		float_value, err := strconv.ParseFloat(string(value), 64)
		if err == nil {
			return int64(float_value)
		}
	case bool:
		if value {
			return 1
		}
	}
	return 0
}

// Converts any object into a real number.
//
// @param x The object to convert.
// @return The result real number.
func ToFloat(value interface{}) float64 {
	switch value := value.(type) {
	case int:
		return float64(value)
	case uint:
		return float64(value)
	case int8:
		return float64(value)
	case uint8:
		return float64(value)
	case int16:
		return float64(value)
	case uint16:
		return float64(value)
	case int32:
		return float64(value)
	case uint32:
		return float64(value)
	case int64:
		return float64(value)
	case uint64:
		return float64(value)
	case float32:
		return float64(value)
	case float64:
		return float64(value)
	case string:
		int_value, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return float64(int_value)
		}
	case []byte:
		int_value, err := strconv.ParseFloat(string(value), 64)
		if err == nil {
			return float64(int_value)
		}
	case bool:
		if value {
			return 1.0
		}
	}
	return 0.0
}

// Checks whether the given data is a unique value of removing data.
//
// @param data The data to check.
// @return True if the data is removing data or false if not.
func IsRemoveData(data interface{}) bool {
	switch data := data.(type) {
	case []byte:
		return IsRemoveBytes(data)
	case string:
		return IsRemoveString(data)
	default:
		return false
	}
}

// Checks whether the given bytes are the removing bytes.
//
// @param data The data to check.
// @return True if the data are the removing bytes.
func IsRemoveBytes(data []byte) bool {
	return ((*reflect.SliceHeader)(unsafe.Pointer(&data)).Data ==
		(*reflect.SliceHeader)(unsafe.Pointer(&RemoveBytes)).Data)
}

// Checks whether the given string is the removing string.
//
// @param data The data to check.
// @return True if the data is the removing string, or false if not.
func IsRemoveString(data string) bool {
	return ((*reflect.StringHeader)(unsafe.Pointer(&data)).Data ==
		(*reflect.StringHeader)(unsafe.Pointer(&RemoveString)).Data)
}

// Checks whether the given data is a unique value of any data.
//
// @param data The data to check.
// @return True if the data is any data or false if not.
func IsAnyData(data interface{}) bool {
	switch data := data.(type) {
	case []byte:
		return IsAnyBytes(data)
	case string:
		return IsAnyString(data)
	default:
		return false
	}
}

// Checks whether the given bytes are the any bytes.
//
// @param data The data to check.
// @return True if the data are the any bytes.
func IsAnyBytes(data []byte) bool {
	return ((*reflect.SliceHeader)(unsafe.Pointer(&data)).Data ==
		(*reflect.SliceHeader)(unsafe.Pointer(&AnyBytes)).Data)
}

// Checks whether the given string is the any string.
//
// @param data The data to check.
// @return True if the data is the any string, or false if not.
func IsAnyString(data string) bool {
	return ((*reflect.StringHeader)(unsafe.Pointer(&data)).Data ==
		(*reflect.StringHeader)(unsafe.Pointer(&AnyString)).Data)
}

// Checks whether the given data is a nil-equivalent value.
//
// @param data The data to check.
// @return True if the data is a nil-equivalent value, or false if not.
func IsNilData(data interface{}) bool {
	switch data := data.(type) {
	case []byte:
		return reflect.ValueOf(data).IsNil()
	case string:
		return IsNilString(data)
	}
	return data == nil
}

// Checks whether the given string is the nil string.
//
// @param data The data to check.
// @return True if the data is the nil string, or false if not.
func IsNilString(data string) bool {
	return ((*reflect.StringHeader)(unsafe.Pointer(&data)).Data ==
		(*reflect.StringHeader)(unsafe.Pointer(&NilString)).Data)
}

// Gets the memory capacity of the platform.
//
// @return The memory capacity of the platform in bytes, or -1 on failure.
func GetMemoryCapacity() int64 {
	return get_memory_capacity()
}

// Gets the current memory usage of the process.
//
// @return The current memory usage of the process in bytes, or -1 on failure.
func GetMemoryUsage() int64 {
	return get_memory_usage()
}

// Primary hash function for the hash database.
//
// @param data The data to calculate the hash value for.
// @param num_buckets: The number of buckets of the hash table.
func PrimaryHash(data []byte, num_buckets uint64) uint64 {
	return primary_hash(data, num_buckets)
}

// Secondary hash function for sharding.
//
// @param data The data to calculate the hash value for.
// @aram num_shards The number of shards.
// @return The hash value.
func SecondaryHash(data []byte, num_shards uint64) uint64 {
	return secondary_hash(data, num_shards)
}

// Gets the Levenshtein edit distance of two Unicode strings.
//
// @param a A string.
// @param b The other string.
// @param utf If true, text is treated as UTF-8.  If false, it is treated as raw bytes.
// @return The Levenshtein edit distance of the two strings.
func EditDistanceLev(a string, b string, utf bool) int {
	return edit_distance_lev(a, b, utf)
}

// Parses a parameter string to make a parameter string map.
//
// @param expr A parameter string in "name=value,name=value,..." format.
// @return The string map of the parameters.
func ParseParams(expr string) map[string]string {
	params := make(map[string]string)
	fields := strings.Split(expr, ",")
	for _, field := range fields {
		columns := strings.SplitN(field, "=", 2)
		if len(columns) != 2 {
			continue
		}
		name := strings.TrimSpace(columns[0])
		value := strings.TrimSpace(columns[1])
		if len(name) > 0 {
			params[name] = value
		}
	}
	return params
}

// Storage to make RecordProcessor accessible to the C code.
type RecordProcessorPool struct {
	data  map[unsafe.Pointer]RecordProcessor
	mutex sync.Mutex
}

var recordProcessorPool = RecordProcessorPool{data: make(map[unsafe.Pointer]RecordProcessor)}

// Register a Go function to the storage.
func registerRecordProcessor(proc RecordProcessor) unsafe.Pointer {
	var up unsafe.Pointer = C.malloc(C.size_t(1))
	if up == nil {
		panic("memory allocation failed")
	}
	recordProcessorPool.mutex.Lock()
	recordProcessorPool.data[up] = proc
	recordProcessorPool.mutex.Unlock()
	return up
}

// Deregister a Go function from the storage.
func deregisterRecordProcessor(up unsafe.Pointer) {
	recordProcessorPool.mutex.Lock()
	delete(recordProcessorPool.data, up)
	recordProcessorPool.mutex.Unlock()
	C.free(up)
}

// Call a Go function in the storage.
//export callRecordProcessor
func callRecordProcessor(up unsafe.Pointer, keyPtr unsafe.Pointer, keySize C.int32_t,
	valuePtr unsafe.Pointer, valueSize C.int32_t) (unsafe.Pointer, int32) {
	recordProcessorPool.mutex.Lock()
	proc := recordProcessorPool.data[up]
	recordProcessorPool.mutex.Unlock()
	var key []byte
	if keyPtr == nil {
		key = nil
	} else {
		key = C.GoBytes(keyPtr, keySize)
	}
	var value []byte
	if valuePtr == nil {
		value = nil
	} else {
		value = C.GoBytes(valuePtr, valueSize)
	}
	rv := proc(key, value)
	var retPtr unsafe.Pointer
	var retSize int32
	if IsNilData(rv) {
		retPtr = nil
		retSize = 0
	} else if IsRemoveData(rv) {
		retPtr = unsafe.Pointer(uintptr(1))
		retSize = 0
	} else {
		rv_bytes := ToByteArray(rv)
		retPtr = C.CBytes(rv_bytes)
		retSize = int32(len(rv_bytes))
	}
	return retPtr, retSize
}

// END OF FILE
