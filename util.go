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

import (
	"fmt"
	"strconv"
)

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

func GetMemoryCapacity() int64 {
	return get_memory_capacity()
}

func GetMemoryUsage() int64 {
	return get_memory_usage()
}

// END OF FILE
