/*************************************************************************************************
 * Test cases
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
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"
)

func CheckEq(t *testing.T, want interface{}, got interface{}) {
	_, _, line, _ := runtime.Caller(1)
	if want == nil {
		if got != nil {
			t.Errorf("line=%d: not equal: want=%q, got=%q", line, want, got)
		}
		return
	}
	switch want := want.(type) {
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		if ToInt(want) != ToInt(got) {
			t.Errorf("line=%d: not equal: want=%d, got=%d", line, want, got)
		}
	case float32, float64, complex64, complex128:
		if ToFloat(want) != ToFloat(got) {
			t.Errorf("line=%d: not equal: want=%d, got=%d", line, want, got)
		}
	case string:
		if want != ToString(got) {
			t.Errorf("line=%d: not equal: want=%s, got=%s", line, want, got)
		}
	case []byte:
		if !reflect.DeepEqual(want, ToByteArray(got)) {
			t.Errorf("line=%d: not equal: want=%q, got=%q", line, want, got)
		}
	case Status:
		if !want.Equals(got) {
			t.Errorf("line=%d: not equal: want=%s, got=%s", line, want.String(), got)
		}
	case *Status:
		if !want.Equals(got) {
			t.Errorf("line=%d: not equal: want=%s, got=%s", line, want.String(), got)
		}
	case StatusCode:
		switch got := got.(type) {
		case Status:
			if !got.Equals(want) {
				t.Errorf("line=%d: not equal: want=%s, got=%s", line, StatusCodeName(want), got.String())
			}
		case *Status:
			if !got.Equals(want) {
				t.Errorf("line=%d: not equal: want=%s, got=%s", line, StatusCodeName(want), got.String())
			}
		case StatusCode:
			if want != got {
				t.Errorf("line=%d: not equal: want=%s, got=%s", line, StatusCodeName(want), StatusCodeName(got))
			}
		default:
			t.Errorf("line=%d: not comparable: want=%s, got=%q", line, StatusCodeName(want), got)
		}
	default:
		if want != got {
			t.Errorf("line=%d: not equal: want=%q, got=%q", line, want, got)
		}
	}
}

func CheckNe(t *testing.T, want interface{}, got interface{}) {
	_, _, line, _ := runtime.Caller(1)
	if want == nil {
		if got == nil {
			t.Errorf("line=%d: equal: want=%q, got=%q", line, want, got)
		}
		return
	}
	switch want := want.(type) {
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		if ToInt(want) == ToInt(got) {
			t.Errorf("line=%d: equal: want=%d, got=%d", line, want, got)
		}
	case float32, float64, complex64, complex128:
		if ToFloat(want) == ToFloat(got) {
			t.Errorf("line=%d: equal: want=%d, got=%d", line, want, got)
		}
	case string:
		if want == ToString(got) {
			t.Errorf("line=%d: equal: want=%s, got=%s", line, want, got)
		}
	case []byte:
		if reflect.DeepEqual(want, ToByteArray(got)) {
			t.Errorf("line=%d: equal: want=%q, got=%q", line, want, got)
		}
	case Status:
		if want.Equals(got) {
			t.Errorf("line=%d: equal: want=%s, got=%s", line, want.String(), got)
		}
	case *Status:
		if want.Equals(got) {
			t.Errorf("line=%d: equal: want=%s, got=%s", line, want.String(), got)
		}
	case StatusCode:
		switch got := got.(type) {
		case Status:
			if got.Equals(want) {
				t.Errorf("line=%d: equal: want=%s, got=%s", line, StatusCodeName(want), got.String())
			}
		case *Status:
			if got.Equals(want) {
				t.Errorf("line=%d: equal: want=%s, got=%s", line, StatusCodeName(want), got.String())
			}
		case StatusCode:
			if want == got {
				t.Errorf("line=%d: equal: want=%s, got=%s", line, StatusCodeName(want), StatusCodeName(got))
			}
		default:
			t.Errorf("line=%d: not comparable: want=%s, got=%q", line, StatusCodeName(want), got)
		}
	default:
		if want == got {
			t.Errorf("line=%d: equal: want=%q, got=%q", line, want, got)
		}
	}
}

func CheckTrue(t *testing.T, got bool) {
	_, _, line, _ := runtime.Caller(1)
	if !got {
		t.Errorf("line=%d: not true", line)
	}
}

func CheckFalse(t *testing.T, got bool) {
	_, _, line, _ := runtime.Caller(1)
	if got {
		t.Errorf("line=%d: true", line)
	}
}

func MakeTempDir() string {
	tmpPath := path.Join(os.TempDir(), fmt.Sprintf(
		"tkrzw-test-%04x%08x", os.Getpid()%(1<<16), time.Now().Unix()%(1<<32)))
	error := os.MkdirAll(tmpPath, 0755)
	if error != nil {
		panic(fmt.Sprintf("cannot create directory: %s", error))
	}
	return tmpPath
}

func TestAssertion(t *testing.T) {
	CheckEq(t, nil, nil)
	CheckNe(t, nil, 0)
	CheckEq(t, 2, 2)
	CheckEq(t, 2.0, 2.0)
	CheckEq(t, "two", "two")
	CheckEq(t, []byte("two"), []byte("two"))
	CheckEq(t, nil, nil)
	CheckEq(t, 2, 2.0)
	CheckEq(t, 2, "2")
	CheckEq(t, 2.0, 2)
	CheckEq(t, 2.0, "2")
	CheckEq(t, "2", 2)
	CheckEq(t, []byte("2"), 2)
	CheckNe(t, 2, 3)
	CheckNe(t, 2.0, 3.0)
	CheckNe(t, "two", "three")
	CheckNe(t, []byte("two"), []byte("three"))
	CheckNe(t, nil, 0)
	CheckTrue(t, true)
	CheckTrue(t, 1 > 0)
	CheckFalse(t, false)
	CheckFalse(t, 1 < 0)
}

type Person struct {
	Name string
}

func (self Person) String() string {
	return fmt.Sprintf("I'm %s.", self.Name)
}

func TestToString(t *testing.T) {
	CheckEq(t, "123", ToString("123"))
	CheckEq(t, "123", ToString([]byte("123")))
	CheckEq(t, "123", ToString(123))
	CheckEq(t, "123.000000", ToString(123.0))
	CheckEq(t, "true", ToString(true))
	CheckEq(t, "false", ToString(false))
	CheckEq(t, "Boom", ToString(errors.New("Boom")))
	CheckEq(t, "I'm Alice.", ToString(Person{"Alice"}))
	CheckEq(t, "I'm Bob.", ToString(&Person{"Bob"}))
}

func TestToByteArray(t *testing.T) {
	CheckEq(t, []byte("123"), ToByteArray("123"))
	CheckEq(t, []byte("123"), ToByteArray([]byte("123")))
	CheckEq(t, []byte("123"), ToByteArray(123))
	CheckEq(t, []byte("123.000000"), ToByteArray(123.0))
	CheckEq(t, []byte("true"), ToByteArray(true))
	CheckEq(t, []byte("false"), ToByteArray(false))
	CheckEq(t, []byte("Boom"), ToByteArray(errors.New("Boom")))
	CheckEq(t, []byte("I'm Alice."), ToByteArray(Person{"Alice"}))
	CheckEq(t, []byte("I'm Bob."), ToByteArray(&Person{"Bob"}))
}

func TestToInt(t *testing.T) {
	CheckEq(t, -123, ToInt("-123"))
	CheckEq(t, -123, ToInt("-123.0"))
	CheckEq(t, -123, ToInt(int8(-123)))
	CheckEq(t, -123, ToInt(int16(-123)))
	CheckEq(t, -123, ToInt(int32(-123)))
	CheckEq(t, -123, ToInt(int64(-123)))
	CheckEq(t, 255, ToInt(uint8(255)))
	CheckEq(t, 255, ToInt(uint16(255)))
	CheckEq(t, 255, ToInt(uint32(255)))
	CheckEq(t, 255, ToInt(uint64(255)))
	CheckEq(t, -255, ToInt(float32(-255)))
	CheckEq(t, -255, ToInt(float64(-255)))
}

func TestToFloat(t *testing.T) {
	CheckEq(t, -123.0, ToFloat("-123"))
	CheckEq(t, -123.5, ToFloat("-123.5"))
	CheckEq(t, -123.0, ToFloat(int8(-123)))
	CheckEq(t, -123.0, ToFloat(int16(-123)))
	CheckEq(t, -123.0, ToFloat(int32(-123)))
	CheckEq(t, -123.0, ToFloat(int64(-123)))
	CheckEq(t, 255.0, ToFloat(uint8(255)))
	CheckEq(t, 255.0, ToFloat(uint16(255)))
	CheckEq(t, 255.0, ToFloat(uint32(255)))
	CheckEq(t, 255.0, ToFloat(uint64(255)))
	CheckEq(t, -255.5, ToFloat(float32(-255.5)))
	CheckEq(t, -255.5, ToFloat(float64(-255.5)))
}

func TestConstants(t *testing.T) {
	CheckTrue(t, len(Version) > 3)
	CheckTrue(t, len(OSName) > 0)
	CheckTrue(t, PageSize > 0)
	CheckEq(t, int64(^uint64(0)>>1)*-1-1, Int64Min)
	CheckEq(t, ^uint64(0)>>1, Int64Max)
}

func TestMiscUtils(t *testing.T) {
	if OSName == "Linux" {
		CheckTrue(t, GetMemoryCapacity() > 0)
		CheckTrue(t, GetMemoryUsage() > 0)
	}
	CheckEq(t, 3042090208, PrimaryHash([]byte("abc"), (1<<32)-1))
	CheckEq(t, uint64(16973900370012003622), PrimaryHash([]byte("abc"), ^uint64(0)))
	CheckEq(t, 702176507, SecondaryHash([]byte("abc"), (1<<32)-1))
	CheckEq(t, uint64(1765794342254572867), SecondaryHash([]byte("abc"), ^uint64(0)))
	CheckEq(t, 0, EditDistanceLev("", "", true))
	CheckEq(t, 1, EditDistanceLev("ac", "abc", true))
	CheckEq(t, 1, EditDistanceLev("あいう", "あう", true))
	CheckEq(t, 3, EditDistanceLev("あいう", "あう", false))
}

func TestStatus(t *testing.T) {
	s := NewStatus()
	CheckEq(t, StatusSuccess, s.GetCode())
	CheckEq(t, "", s.GetMessage())
	CheckTrue(t, s.Equals(s))
	CheckTrue(t, s.Equals(*s))
	CheckTrue(t, s.Equals(StatusSuccess))
	CheckFalse(t, s.Equals(StatusNotFoundError))
	CheckFalse(t, s.Equals(100))
	CheckEq(t, "SUCCESS", s)
	CheckTrue(t, s.IsOK())
	s.OrDie()
	s = NewStatus(StatusNotFoundError, "foobar")
	CheckEq(t, StatusNotFoundError, s.GetCode())
	CheckEq(t, "foobar", s.GetMessage())
	CheckEq(t, "NOT_FOUND_ERROR: foobar", s.String())
	CheckEq(t, "NOT_FOUND_ERROR", s)
	CheckEq(t, "NOT_FOUND_ERROR", s.Error())
	CheckTrue(t, s.Equals(s))
	CheckTrue(t, s.Equals(*s))
	CheckTrue(t, s.Equals(StatusNotFoundError))
	CheckFalse(t, s.Equals(StatusSuccess))
	CheckFalse(t, s.IsOK())
	CheckFalse(t, s.Equals(100))
	s = NewStatus1(StatusSuccess)
	CheckEq(t, StatusSuccess, s.GetCode())
	CheckEq(t, "", s.GetMessage())
	s = NewStatus2(StatusNotFoundError, "bazquux")
	CheckEq(t, StatusNotFoundError, s.GetCode())
	CheckEq(t, "bazquux", s.GetMessage())
	CheckEq(t, StatusSuccess, StatusSuccess)
	CheckEq(t, StatusSuccess, NewStatus1(StatusSuccess))
	CheckEq(t, NewStatus1(StatusSuccess), StatusSuccess)
	CheckEq(t, NewStatus1(StatusSuccess), NewStatus1(StatusSuccess))
	CheckEq(t, StatusNotFoundError, NewStatus1(StatusNotFoundError))
	CheckNe(t, StatusNotFoundError, StatusSuccess)
	CheckNe(t, StatusNotFoundError, NewStatus1(StatusSuccess))
	CheckNe(t, NewStatus1(StatusNotFoundError), StatusSuccess)
	CheckNe(t, NewStatus1(StatusNotFoundError), NewStatus1(StatusSuccess))
	CheckNe(t, StatusNotFoundError, NewStatus1(StatusUnknownError))
}

func TestDBMBasic(t *testing.T) {
	tmpDir := MakeTempDir()
	defer os.RemoveAll(tmpDir)
	filePath := path.Join(tmpDir, "casket.tkh")
	copyPath := path.Join(tmpDir, "casket-copy.tkh")
	dbm := NewDBM()
	status := dbm.Open(filePath, true, "truncate=true,num_buckets=5")
	CheckEq(t, StatusSuccess, status)
	CheckTrue(t, len(dbm.String()) > len(filePath))
	CheckTrue(t, dbm.Set("one", "first", false).IsOK())
	CheckEq(t, StatusDuplicationError, dbm.Set("one", "uno", false))
	CheckTrue(t, dbm.Set("two", "second", false).IsOK())
	CheckTrue(t, dbm.Set("three", "third", false).IsOK())
	CheckTrue(t, dbm.Append("three", "3", ":").IsOK())
	count, status := dbm.Count()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, 3, count)
	CheckEq(t, 3, dbm.CountSimple())
	value, status := dbm.Get("one")
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "first", value)
	value, status = dbm.Get([]byte("two"))
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "second", value)
	value_str, status := dbm.GetStr([]byte("three"))
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "third:3", value_str)
	value_str, status = dbm.GetStr([]byte("fourth"))
	CheckEq(t, StatusNotFoundError, status)
	CheckEq(t, "", value_str)
	CheckEq(t, "first", dbm.GetSimple("one", "*"))
	CheckEq(t, "second", dbm.GetStrSimple("two", "*"))
	CheckEq(t, "third:3", dbm.GetStrSimple([]byte("three"), "*"))
	CheckEq(t, "*", dbm.GetStrSimple([]byte("four"), "*"))
	CheckEq(t, StatusSuccess, dbm.Remove("one"))
	CheckEq(t, StatusSuccess, dbm.Remove("two"))
	CheckEq(t, StatusSuccess, dbm.Remove([]byte("three")))
	CheckEq(t, StatusNotFoundError, dbm.Remove([]byte("fourth")))
	CheckEq(t, 0, dbm.CountSimple())
	CheckEq(t, StatusSuccess, dbm.CompareExchange("num", nil, "first"))
	CheckEq(t, "first", dbm.GetSimple("num", "*"))
	CheckEq(t, StatusInfeasibleError, dbm.CompareExchange("num", nil, "first"))
	CheckEq(t, StatusSuccess, dbm.CompareExchange("num", "first", "second"))
	CheckEq(t, "second", dbm.GetSimple("num", "*"))
	CheckEq(t, StatusSuccess, dbm.CompareExchange("num", "second", nil))
	CheckEq(t, "*", dbm.GetSimple("num", "*"))
	inc_value, status := dbm.Increment("num", 2, 100)
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, 102, inc_value)
	inc_value, status = dbm.Increment("num", 3, 100)
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, 105, inc_value)
	CheckEq(t, StatusSuccess, dbm.Remove("num"))
	old_value, status := dbm.SetAndGet("zero", "nil", false)
	CheckTrue(t, old_value == nil)
	CheckEq(t, StatusSuccess, status)
	old_value, status = dbm.SetAndGet("zero", "nothing", false)
	CheckEq(t, "nil", old_value)
	CheckEq(t, StatusDuplicationError, status)
	old_value_str, status := dbm.SetAndGetStr("zero", "void", false)
	CheckEq(t, "nil", *old_value_str)
	CheckEq(t, StatusDuplicationError, status)
	old_value, status = dbm.RemoveAndGet("zero")
	CheckEq(t, "nil", old_value)
	CheckEq(t, StatusSuccess, status)
	old_value, status = dbm.RemoveAndGet("zero")
	CheckTrue(t, old_value == nil)
	CheckEq(t, StatusNotFoundError, status)
	old_value_str, status = dbm.SetAndGetStr("zero", "void", false)
	CheckTrue(t, old_value == nil)
	CheckEq(t, StatusSuccess, status)
	old_value_str, status = dbm.RemoveAndGetStr("zero")
	CheckEq(t, "void", *old_value_str)
	CheckEq(t, StatusSuccess, status)
	old_value_str, status = dbm.RemoveAndGetStr("zero")
	CheckTrue(t, old_value == nil)
	CheckEq(t, StatusNotFoundError, status)
	records := map[string]string{"one": "first", "two": "second"}
	CheckEq(t, StatusSuccess, dbm.SetMultiStr(records, false))
	keys := []string{"one", "two", "three"}
	records = dbm.GetMultiStr(keys)
	CheckEq(t, 2, len(records))
	CheckEq(t, "first", records["one"])
	CheckEq(t, "second", records["two"])
	CheckEq(t, StatusNotFoundError, dbm.RemoveMulti(keys))
	set1 := []KeyValuePair{KeyValuePair{[]byte("one"), []byte(nil)},
		KeyValuePair{[]byte("two"), []byte(nil)}}
	set2 := []KeyValuePair{KeyValuePair{[]byte("one"), []byte("ichi")},
		KeyValuePair{[]byte("two"), []byte("ni")}}
	set3 := []KeyValuePair{KeyValuePair{[]byte("one"), []byte("uno")},
		KeyValuePair{[]byte("two"), []byte("dos")}}
	CheckEq(t, StatusSuccess, dbm.CompareExchangeMulti(set1, set2))
	CheckEq(t, "ichi", dbm.GetSimple("one", "*"))
	CheckEq(t, "ni", dbm.GetSimple("two", "*"))
	CheckEq(t, StatusInfeasibleError, dbm.CompareExchangeMulti(set1, set2))
	CheckEq(t, StatusSuccess, dbm.CompareExchangeMulti(set2, set3))
	CheckEq(t, "uno", dbm.GetSimple("one", "*"))
	CheckEq(t, "dos", dbm.GetSimple("two", "*"))
	CheckEq(t, StatusSuccess, dbm.CompareExchangeMulti(set3, set1))
	CheckEq(t, "*", dbm.GetSimple("one", "*"))
	CheckEq(t, "*", dbm.GetSimple("two", "*"))
	CheckEq(t, 0, dbm.CountSimple())
	set4 := []KeyValueStrPair{KeyValueStrPair{"one", ""}, KeyValueStrPair{"two", ""}}
	set5 := []KeyValueStrPair{KeyValueStrPair{"one", "apple"}, KeyValueStrPair{"two", "orange"}}
	CheckEq(t, StatusSuccess, dbm.CompareExchangeMultiStr(set4, set5))
	CheckEq(t, StatusInfeasibleError, dbm.CompareExchangeMultiStr(set4, set5))
	CheckEq(t, "apple", dbm.GetSimple("one", "*"))
	CheckEq(t, "orange", dbm.GetSimple("two", "*"))
	CheckEq(t, StatusSuccess, dbm.CompareExchangeMultiStr(set5, set4))
	CheckEq(t, StatusInfeasibleError, dbm.CompareExchangeMultiStr(set5, set4))
	CheckEq(t, "*", dbm.GetSimple("one", "*"))
	CheckEq(t, "*", dbm.GetSimple("two", "*"))
	CheckEq(t, 0, dbm.CountSimple())
	fileSize, status := dbm.GetFileSize()
	CheckEq(t, StatusSuccess, status)
	CheckTrue(t, fileSize > 0)
	CheckEq(t, fileSize, dbm.GetFileSizeSimple())
	gotFilePath, status := dbm.GetFilePath()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, filePath, gotFilePath)
	CheckEq(t, filePath, dbm.GetFilePathSimple())
	for i := 1; i <= 10; i++ {
		CheckTrue(t, dbm.Set(i, i*i, true).IsOK())
	}
	CheckEq(t, 10, dbm.CountSimple())
	tobe, status := dbm.ShouldBeRebuilt()
	CheckEq(t, StatusSuccess, status)
	CheckTrue(t, tobe)
	CheckTrue(t, dbm.ShouldBeRebuiltSimple())
	CheckEq(t, StatusSuccess, dbm.Rebuild(""))
	CheckEq(t, StatusSuccess, dbm.Synchronize(true, ""))
	CheckEq(t, StatusSuccess, dbm.CopyFileData(copyPath))
	CheckEq(t, StatusSuccess, dbm.Clear())
	CheckEq(t, 0, dbm.CountSimple())
	CheckEq(t, StatusSuccess, dbm.Close())
	CheckEq(t, StatusSuccess, dbm.Open(filePath, true, ""))
	copyDBM := NewDBM()
	CheckEq(t, StatusSuccess, copyDBM.Open(copyPath, false, ""))
	CheckEq(t, 10, copyDBM.CountSimple())
	CheckEq(t, StatusSuccess, copyDBM.Export(dbm))
	CheckEq(t, 10, dbm.CountSimple())
	CheckEq(t, StatusSuccess, copyDBM.Close())
	inspRecords := dbm.Inspect()
	CheckEq(t, "10", inspRecords["num_records"])
	CheckEq(t, "HashDBM", inspRecords["class"])
	iter := dbm.MakeIterator()
	CheckEq(t, StatusSuccess, iter.First())
	CheckTrue(t, len(iter.String()) > 1)
	count = 0
	records = make(map[string]string)
	for {
		key, value, status := iter.Get()
		if !status.IsOK() {
			CheckEq(t, StatusNotFoundError, status)
			break
		}
		key_str, value_str, status := iter.GetStr()
		CheckEq(t, StatusSuccess, status)
		CheckEq(t, key_str, string(key))
		CheckEq(t, value_str, string(value))
		records[key_str] = value_str
		one_key, status := iter.GetKey()
		CheckEq(t, StatusSuccess, status)
		CheckEq(t, key_str, string(one_key))
		one_key_str, status := iter.GetKeyStr()
		CheckEq(t, StatusSuccess, status)
		CheckEq(t, key_str, one_key_str)
		one_value, status := iter.GetValue()
		CheckEq(t, StatusSuccess, status)
		CheckEq(t, value_str, string(one_value))
		one_value_str, status := iter.GetValueStr()
		CheckEq(t, StatusSuccess, status)
		CheckEq(t, value_str, one_value_str)
		CheckEq(t, StatusSuccess, iter.Next())
		count++
	}
	CheckEq(t, 10, count)
	for i := 1; i <= 10; i++ {
		CheckEq(t, ToString(i*i), records[ToString(i)])
	}
	CheckEq(t, StatusSuccess, iter.Jump("5"))
	key, value, status := iter.Get()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "5", key)
	CheckEq(t, "25", value)
	CheckEq(t, StatusSuccess, iter.Set("foobar"))
	value_str, status = iter.GetValueStr()
	CheckEq(t, StatusSuccess, iter.Remove())
	CheckEq(t, 9, dbm.CountSimple())
	iter.Destruct()
	CheckEq(t, StatusSuccess, dbm.Close())
	os.Remove(copyPath)
	CheckEq(t, StatusSuccess, RestoreDatabase(filePath, copyPath, "", -1))
	CheckEq(t, StatusSuccess, copyDBM.Open(copyPath, false, ""))
	CheckEq(t, 9, copyDBM.CountSimple())
	CheckEq(t, StatusSuccess, copyDBM.Close())
}

func TestDBMIterator(t *testing.T) {
	tmpDir := MakeTempDir()
	defer os.RemoveAll(tmpDir)
	filePath := path.Join(tmpDir, "casket.tkt")
	dbm := NewDBM()
	status := dbm.Open(filePath, true, "truncate=true")
	CheckEq(t, StatusSuccess, status)
	for i := 1; i <= 100; i++ {
		key := fmt.Sprintf("%08d", i)
		value := fmt.Sprintf("%d", i*i)
		CheckEq(t, StatusSuccess, dbm.Set(key, value, false))
	}
	CheckEq(t, 100, dbm.CountSimple())
	iter := dbm.MakeIterator()
	CheckEq(t, StatusSuccess, iter.Jump("00000050"))
	CheckEq(t, StatusSuccess, iter.Remove())
	CheckEq(t, StatusSuccess, iter.Jump("00000050"))
	key, status := iter.GetKeyStr()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "00000051", key)
	CheckEq(t, StatusSuccess, iter.JumpLower("00000051", true))
	key, status = iter.GetKeyStr()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "00000051", key)
	CheckEq(t, StatusSuccess, iter.JumpLower("00000051", false))
	key, status = iter.GetKeyStr()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "00000049", key)
	CheckEq(t, StatusSuccess, iter.Next())
	key, status = iter.GetKeyStr()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "00000051", key)
	CheckEq(t, StatusSuccess, iter.JumpUpper("00000049", true))
	key, status = iter.GetKeyStr()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "00000049", key)
	CheckEq(t, StatusSuccess, iter.JumpUpper("00000049", false))
	key, status = iter.GetKeyStr()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "00000051", key)
	CheckEq(t, StatusSuccess, iter.Previous())
	key, status = iter.GetKeyStr()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "00000049", key)
	iter.Destruct()
	count := 0
	for record := range dbm.Each() {
		CheckEq(t, dbm.GetSimple(record.Key, ""), record.Value)
		count++
	}
	CheckEq(t, dbm.CountSimple(), count)
	count = 0
	for record := range dbm.EachStr() {
		CheckEq(t, dbm.GetSimple(record.Key, ""), record.Value)
		count++
	}
	CheckEq(t, dbm.CountSimple(), count)
	CheckEq(t, StatusSuccess, dbm.Close())
}

func TestDBMThread(t *testing.T) {
	tmpDir := MakeTempDir()
	defer os.RemoveAll(tmpDir)
	filePath := path.Join(tmpDir, "casket.tkh")
	dbm := NewDBM()
	status := dbm.Open(filePath, true, "truncate=true")
	CheckEq(t, StatusSuccess, status)
	numIterations := 5000
	numThreads := 5
	recordMaps := make([]map[string]string, 0, numThreads)
	mutexes := make([]sync.Mutex, 0, numThreads)
	for i := 0; i < numThreads; i++ {
		recordMaps = append(recordMaps, make(map[string]string))
		mutexes = append(mutexes, sync.Mutex{})
	}
	task := func(thid int, done chan<- bool) {
		random := rand.New(rand.NewSource(int64(thid)))
		for i := 0; i < numIterations; i++ {
			keyNum := random.Intn(numIterations * numThreads)
			valueNum := random.Intn(numIterations * numThreads)
			key := fmt.Sprintf("%d", keyNum)
			value := fmt.Sprintf("%d", valueNum*valueNum)
			groupIndex := keyNum % numThreads
			recordMap := &recordMaps[groupIndex]
			mutex := &mutexes[groupIndex]
			mutex.Lock()
			if random.Intn(5) == 0 {
				gotValue, status := dbm.Get(key)
				if status.IsOK() {
					CheckEq(t, (*recordMap)[key], gotValue)
				} else {
					CheckEq(t, StatusNotFoundError, status)
				}
			} else if random.Intn(5) == 0 {
				status := dbm.Remove(key)
				CheckTrue(t, status.Equals(StatusSuccess) || status.Equals(StatusNotFoundError))
				delete(*recordMap, key)
			} else {
				CheckEq(t, StatusSuccess, dbm.Set(key, value, true))
				(*recordMap)[key] = value
			}
			mutex.Unlock()
			if random.Intn(10) == 0 {
				iter := dbm.MakeIterator()
				iter.Jump(key)
				_, _, status := iter.Get()
				CheckTrue(t, status.Equals(StatusSuccess) || status.Equals(StatusNotFoundError))
			}
		}
		done <- true
	}
	dones := make([]chan bool, 0)
	for i := 0; i < numThreads; i++ {
		done := make(chan bool)
		go task(i, done)
		dones = append(dones, done)
	}
	for _, done := range dones {
		<-done
	}
	numRecords := 0
	for _, recordMap := range recordMaps {
		numRecords += len(recordMap)
		for key, value := range recordMap {
			gotValue, status := dbm.Get(key)
			CheckEq(t, StatusSuccess, status)
			CheckEq(t, value, gotValue)
		}
	}
	CheckEq(t, numRecords, dbm.CountSimple())
	CheckEq(t, StatusSuccess, dbm.Close())
}

func TestDBMExport(t *testing.T) {
	tmpDir := MakeTempDir()
	defer os.RemoveAll(tmpDir)
	filePath := path.Join(tmpDir, "casket.tkh")
	copyPath := path.Join(tmpDir, "casket-copy.dat")
	dbm := NewDBM()
	CheckEq(t, StatusSuccess, dbm.Open(filePath, true, "truncate=true"))
	CheckEq(t, StatusSuccess, dbm.Set("one", "first", true))
	CheckEq(t, StatusSuccess, dbm.Set("two", "second", true))
	CheckEq(t, 2, dbm.CountSimple())
	copyFile := NewFile()
	CheckEq(t, StatusSuccess, copyFile.Open(copyPath, true, "truncate=true"))
	CheckEq(t, StatusSuccess, dbm.ExportToFlatRecords(copyFile))
	CheckEq(t, StatusSuccess, dbm.Clear())
	CheckEq(t, 0, dbm.CountSimple())
	CheckEq(t, StatusSuccess, dbm.ImportFromFlatRecords(copyFile))
	CheckEq(t, 2, dbm.CountSimple())
	CheckEq(t, "first", dbm.GetSimple("one", "*"))
	CheckEq(t, "second", dbm.GetSimple("two", "*"))
	CheckEq(t, StatusSuccess, copyFile.Close())
	CheckEq(t, StatusSuccess, copyFile.Open(copyPath, true, "truncate=true"))
	CheckEq(t, StatusSuccess, dbm.ExportKeysAsLines(copyFile))
	lines := copyFile.Search("contain", "o", 0)
	CheckEq(t, 2, len(lines))
	CheckEq(t, StatusSuccess, copyFile.Close())
	CheckEq(t, StatusSuccess, dbm.Close())
}

func TestDBMSearch(t *testing.T) {
	tmpDir := MakeTempDir()
	defer os.RemoveAll(tmpDir)
	filePath := path.Join(tmpDir, "casket.tks")
	dbm := NewDBM()
	status := dbm.Open(filePath, true, "truncate=true")
	CheckEq(t, StatusSuccess, status)
	for i := 1; i <= 100; i++ {
		key := fmt.Sprintf("%08d", i)
		value := fmt.Sprintf("%d", i*i)
		CheckEq(t, StatusSuccess, dbm.Set(key, value, false))
	}
	CheckEq(t, StatusSuccess, dbm.Synchronize(false, "reducer=ReduceToFirst"))
	CheckEq(t, 100, dbm.CountSimple())
	keys := dbm.Search("contain", "99", 0)
	CheckEq(t, 1, len(keys))
	CheckEq(t, "00000099", keys[0])
	keys = dbm.Search("edit", "00000100", 2)
	CheckEq(t, 2, len(keys))
	CheckEq(t, "00000100", keys[0])
	CheckEq(t, "00000001", keys[1])
	keys = dbm.Search("begin", "0000005", 0)
	CheckEq(t, 10, len(keys))
	CheckEq(t, "00000050", keys[0])
	keys = dbm.Search("end", "0", 0)
	CheckEq(t, 10, len(keys))
	CheckEq(t, "00000010", keys[0])
	CheckEq(t, StatusSuccess, dbm.Close())
}

func TestFile(t *testing.T) {
	tmpDir := MakeTempDir()
	defer os.RemoveAll(tmpDir)
	filePath := path.Join(tmpDir, "casket.txt")
	file := NewFile()
	status := file.Open(filePath, true, "truncate=true")
	CheckEq(t, StatusSuccess, status)
	CheckTrue(t, len(file.String()) > len(filePath))
	CheckEq(t, StatusSuccess, file.Write(3, "defg"))
	CheckEq(t, StatusSuccess, file.Write(0, "abc"))
	off, status := file.Append("hij")
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, 7, off)
	size, status := file.GetSize()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, 10, size)
	data, status := file.Read(0, 10)
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "abcdefghij", data)
	data_str, status := file.ReadStr(3, 5)
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "defgh", data_str)
	data, status = file.Read(8, 4)
	CheckEq(t, StatusInfeasibleError, status)
	CheckEq(t, StatusSuccess, file.Truncate(5))
	size, status = file.GetSize()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, 5, size)
	CheckEq(t, StatusSuccess, file.Synchronize(false, 0, 5))
	CheckEq(t, StatusSuccess, file.Truncate(0))
	for i := 1; i <= 100; i++ {
		_, status = file.Append(fmt.Sprintf("%08d\n", i))
		CheckEq(t, StatusSuccess, status)
	}
	CheckEq(t, 19, len(file.Search("contain", "9", 0)))
	CheckEq(t, 1, len(file.Search("contain", "100", 0)))
	CheckEq(t, 3, len(file.Search("end", "0", 3)))
	CheckEq(t, StatusSuccess, file.Close())
}

// END OF FILE
