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
	default:
		if want != got {
			t.Errorf("line=%d: not equal: want=%q, got=%q", line, want, got)
		}
	}
}

func CheckNe(t *testing.T, want interface{}, got interface{}) {
	_, _, line, _ := runtime.Caller(1)
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
	CheckEq(t, "NOT_FOUND_ERROR: foobar", s)
	CheckEq(t, "NOT_FOUND_ERROR: foobar", s.Error())
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
}

func TestVersion(t *testing.T) {
	CheckTrue(t, len(Version) > 3)
	CheckTrue(t, len(OSName) > 0)
	CheckTrue(t, PageSize > 0)
}

func TestDBMBasic(t *testing.T) {
	tmpDir := MakeTempDir()
	defer os.RemoveAll(tmpDir)
	filePath := path.Join(tmpDir, "casket.tkh")
	copyPath := path.Join(tmpDir, "casket-copy.tkh")
	dbm := NewDBM()
	status := dbm.Open(filePath, true, "truncate=true,num_buckets=5")
	CheckTrue(t, status.Equals(StatusSuccess))
	CheckTrue(t, len(dbm.String()) > len(filePath))
	CheckTrue(t, dbm.Set("one", "first", false).IsOK())
	CheckTrue(t, dbm.Set("one", "uno", false).Equals(StatusDuplicationError))
	CheckTrue(t, dbm.Set("two", "second", false).IsOK())
	CheckTrue(t, dbm.Set("three", "third", false).IsOK())
	CheckTrue(t, dbm.Append("three", "3", ":").IsOK())
	count, status := dbm.Count()
	CheckTrue(t, status.Equals(StatusSuccess))
	CheckEq(t, 3, count)
	CheckEq(t, 3, dbm.CountSimple())
	value, status := dbm.Get("one")
	CheckTrue(t, status.Equals(StatusSuccess))
	CheckEq(t, "first", value)
	value, status = dbm.Get([]byte("two"))
	CheckTrue(t, status.Equals(StatusSuccess))
	CheckEq(t, "second", value)
	value_str, status := dbm.GetStr([]byte("three"))
	CheckTrue(t, status.Equals(StatusSuccess))
	CheckEq(t, "third:3", value_str)
	value_str, status = dbm.GetStr([]byte("fourth"))
	CheckTrue(t, status.Equals(StatusNotFoundError))
	CheckEq(t, "", value_str)
	CheckEq(t, "first", dbm.GetSimple("one", "*"))
	CheckEq(t, "second", dbm.GetStrSimple("two", "*"))
	CheckEq(t, "third:3", dbm.GetStrSimple([]byte("three"), "*"))
	CheckEq(t, "*", dbm.GetStrSimple([]byte("four"), "*"))
	CheckTrue(t, dbm.Remove("one").Equals(StatusSuccess))
	CheckTrue(t, dbm.Remove("two").Equals(StatusSuccess))
	CheckTrue(t, dbm.Remove([]byte("three")).Equals(StatusSuccess))
	CheckTrue(t, dbm.Remove([]byte("fourth")).Equals(StatusNotFoundError))
	CheckEq(t, 0, dbm.CountSimple())
	CheckTrue(t, dbm.CompareExchange("num", nil, "first").Equals(StatusSuccess))
	CheckEq(t, "first", dbm.GetSimple("num", "*"))
	CheckTrue(t, dbm.CompareExchange("num", nil, "first").Equals(StatusInfeasibleError))
	CheckTrue(t, dbm.CompareExchange("num", "first", "second").Equals(StatusSuccess))
	CheckEq(t, "second", dbm.GetSimple("num", "*"))
	CheckTrue(t, dbm.CompareExchange("num", "second", nil).Equals(StatusSuccess))
	CheckEq(t, "*", dbm.GetSimple("num", "*"))
	inc_value, status := dbm.Increment("num", 2, 100)
	CheckTrue(t, status.Equals(StatusSuccess))
	CheckEq(t, 102, inc_value)
	inc_value, status = dbm.Increment("num", 3, 100)
	CheckTrue(t, status.Equals(StatusSuccess))
	CheckEq(t, 105, inc_value)
	CheckTrue(t, dbm.Remove("num").Equals(StatusSuccess))
	old_value, status := dbm.SetAndGet("zero", "nil", false)
	CheckTrue(t, old_value == nil)
	CheckTrue(t, status.Equals(StatusSuccess))
	old_value, status = dbm.SetAndGet("zero", "nothing", false)
	CheckEq(t, "nil", old_value)
	CheckTrue(t, status.Equals(StatusDuplicationError))
	old_value_str, status := dbm.SetAndGetStr("zero", "void", false)
	CheckEq(t, "nil", *old_value_str)
	CheckTrue(t, status.Equals(StatusDuplicationError))
	old_value, status = dbm.RemoveAndGet("zero")
	CheckEq(t, "nil", old_value)
	CheckTrue(t, status.Equals(StatusSuccess))
	old_value, status = dbm.RemoveAndGet("zero")
	CheckTrue(t, old_value == nil)
	CheckTrue(t, status.Equals(StatusNotFoundError))
	old_value_str, status = dbm.SetAndGetStr("zero", "void", false)
	CheckTrue(t, old_value == nil)
	CheckTrue(t, status.Equals(StatusSuccess))
	old_value_str, status = dbm.RemoveAndGetStr("zero")
	CheckEq(t, "void", *old_value_str)
	CheckTrue(t, status.Equals(StatusSuccess))
	old_value_str, status = dbm.RemoveAndGetStr("zero")
	CheckTrue(t, old_value == nil)
	CheckTrue(t, status.Equals(StatusNotFoundError))
	records := map[string]string{"one": "first", "two": "second"}
	CheckTrue(t, dbm.SetMultiStr(records, false).Equals(StatusSuccess))
	keys := []string{"one", "two", "three"}
	records = dbm.GetMultiStr(keys)
	CheckEq(t, 2, len(records))
	CheckEq(t, "first", records["one"])
	CheckEq(t, "second", records["two"])
	CheckTrue(t, dbm.RemoveMulti(keys).Equals(StatusNotFoundError))
	set1 := []KeyValuePair{KeyValuePair{[]byte("one"), []byte(nil)},
		KeyValuePair{[]byte("two"), []byte(nil)}}
	set2 := []KeyValuePair{KeyValuePair{[]byte("one"), []byte("ichi")},
		KeyValuePair{[]byte("two"), []byte("ni")}}
	set3 := []KeyValuePair{KeyValuePair{[]byte("one"), []byte("uno")},
		KeyValuePair{[]byte("two"), []byte("dos")}}
	CheckTrue(t, dbm.CompareExchangeMulti(set1, set2).Equals(StatusSuccess))
	CheckEq(t, "ichi", dbm.GetSimple("one", "*"))
	CheckEq(t, "ni", dbm.GetSimple("two", "*"))
	CheckTrue(t, dbm.CompareExchangeMulti(set1, set2).Equals(StatusInfeasibleError))
	CheckTrue(t, dbm.CompareExchangeMulti(set2, set3).Equals(StatusSuccess))
	CheckEq(t, "uno", dbm.GetSimple("one", "*"))
	CheckEq(t, "dos", dbm.GetSimple("two", "*"))
	CheckTrue(t, dbm.CompareExchangeMulti(set3, set1).Equals(StatusSuccess))
	CheckEq(t, "*", dbm.GetSimple("one", "*"))
	CheckEq(t, "*", dbm.GetSimple("two", "*"))
	CheckEq(t, 0, dbm.CountSimple())
	fileSize, status := dbm.GetFileSize()
	CheckTrue(t, status.Equals(StatusSuccess))
	CheckTrue(t, fileSize > 0)
	CheckEq(t, fileSize, dbm.GetFileSizeSimple())
	gotFilePath, status := dbm.GetFilePath()
	CheckTrue(t, status.Equals(StatusSuccess))
	CheckEq(t, filePath, gotFilePath)
	CheckEq(t, filePath, dbm.GetFilePathSimple())
	for i := 1; i <= 10; i++ {
		CheckTrue(t, dbm.Set(i, i*i, true).IsOK())
	}
	CheckEq(t, 10, dbm.CountSimple())
	tobe, status := dbm.ShouldBeRebuilt()
	CheckTrue(t, status.Equals(StatusSuccess))
	CheckTrue(t, tobe)
	CheckTrue(t, dbm.ShouldBeRebuiltSimple())
	CheckTrue(t, dbm.Rebuild("").Equals(StatusSuccess))
	CheckTrue(t, dbm.Synchronize(true, "").Equals(StatusSuccess))
	CheckTrue(t, dbm.CopyFileData(copyPath).Equals(StatusSuccess))
	CheckTrue(t, dbm.Clear().Equals(StatusSuccess))
	CheckEq(t, 0, dbm.CountSimple())
	CheckTrue(t, dbm.Close().Equals(StatusSuccess))
	CheckTrue(t, dbm.Open(copyPath, true, "").Equals(StatusSuccess))
	copyDBM := NewDBM()
	CheckTrue(t, copyDBM.Open(copyPath, false, "").Equals(StatusSuccess))
	CheckEq(t, 10, copyDBM.CountSimple())
	CheckTrue(t, copyDBM.Export(dbm).Equals(StatusSuccess))
	CheckEq(t, 10, dbm.CountSimple())
	CheckTrue(t, copyDBM.Close().Equals(StatusSuccess))
	inspRecords := dbm.Inspect()
	CheckEq(t, "10", inspRecords["num_records"])
	CheckEq(t, "HashDBM", inspRecords["class"])
	iter := dbm.MakeIterator()
	CheckTrue(t, iter.First().Equals(StatusSuccess))
	CheckTrue(t, len(iter.String()) > 1)
	count = 0
	records = make(map[string]string)
	for {
		key, value, status := iter.Get()
		if !status.IsOK() {
			CheckTrue(t, status.Equals(StatusNotFoundError))
			break
		}
		key_str, value_str, status := iter.GetStr()
		CheckTrue(t, status.Equals(StatusSuccess))
		CheckEq(t, key_str, string(key))
		CheckEq(t, value_str, string(value))
		records[key_str] = value_str
		one_key, status := iter.GetKey()
		CheckTrue(t, status.Equals(StatusSuccess))
		CheckEq(t, key_str, string(one_key))
		one_key_str, status := iter.GetKeyStr()
		CheckTrue(t, status.Equals(StatusSuccess))
		CheckEq(t, key_str, one_key_str)
		one_value, status := iter.GetValue()
		CheckTrue(t, status.Equals(StatusSuccess))
		CheckEq(t, value_str, string(one_value))
		one_value_str, status := iter.GetValueStr()
		CheckTrue(t, status.Equals(StatusSuccess))
		CheckEq(t, value_str, one_value_str)
		CheckTrue(t, iter.Next().Equals(StatusSuccess))
		count++
	}
	CheckEq(t, 10, count)
	for i := 1; i <= 10; i++ {
		CheckEq(t, ToString(i*i), records[ToString(i)])
	}
	CheckTrue(t, iter.Jump("5").Equals(StatusSuccess))
	key, value, status := iter.Get()
	CheckTrue(t, status.Equals(StatusSuccess))
	CheckEq(t, "5", key)
	CheckEq(t, "25", value)
	CheckTrue(t, iter.Set("foobar").Equals(StatusSuccess))
	value_str, status = iter.GetValueStr()
	CheckTrue(t, iter.Remove().Equals(StatusSuccess))
	CheckEq(t, 9, dbm.CountSimple())
	iter.Destruct()
	CheckTrue(t, dbm.Close().Equals(StatusSuccess))
}

func TestDBMIterator(t *testing.T) {
	tmpDir := MakeTempDir()
	defer os.RemoveAll(tmpDir)
	filePath := path.Join(tmpDir, "casket.tkt")
	dbm := NewDBM()
	status := dbm.Open(filePath, true, "truncate=true")
	CheckTrue(t, status.Equals(StatusSuccess))
	for i := 1; i <= 100; i++ {
		key := fmt.Sprintf("%08d", i)
		value := fmt.Sprintf("%d", i*i)
		CheckTrue(t, dbm.Set(key, value, false).Equals(StatusSuccess))
	}
	CheckEq(t, 100, dbm.CountSimple())
	iter := dbm.MakeIterator()
	CheckTrue(t, iter.Jump("00000050").Equals(StatusSuccess))
	CheckTrue(t, iter.Remove().Equals(StatusSuccess))
	CheckTrue(t, iter.Jump("00000050").Equals(StatusSuccess))
	key, status := iter.GetKeyStr()
	CheckTrue(t, status.Equals(StatusSuccess))
	CheckEq(t, "00000051", key)
	CheckTrue(t, iter.JumpLower("00000051", true).Equals(StatusSuccess))
	key, status = iter.GetKeyStr()
	CheckTrue(t, status.Equals(StatusSuccess))
	CheckEq(t, "00000051", key)
	CheckTrue(t, iter.JumpLower("00000051", false).Equals(StatusSuccess))
	key, status = iter.GetKeyStr()
	CheckTrue(t, status.Equals(StatusSuccess))
	CheckEq(t, "00000049", key)
	CheckTrue(t, iter.Next().Equals(StatusSuccess))
	key, status = iter.GetKeyStr()
	CheckTrue(t, status.Equals(StatusSuccess))
	CheckEq(t, "00000051", key)
	CheckTrue(t, iter.JumpUpper("00000049", true).Equals(StatusSuccess))
	key, status = iter.GetKeyStr()
	CheckTrue(t, status.Equals(StatusSuccess))
	CheckEq(t, "00000049", key)
	CheckTrue(t, iter.JumpUpper("00000049", false).Equals(StatusSuccess))
	key, status = iter.GetKeyStr()
	CheckTrue(t, status.Equals(StatusSuccess))
	CheckEq(t, "00000051", key)
	CheckTrue(t, iter.Previous().Equals(StatusSuccess))
	key, status = iter.GetKeyStr()
	CheckTrue(t, status.Equals(StatusSuccess))
	CheckEq(t, "00000049", key)
	iter.Destruct()
	CheckTrue(t, dbm.Close().Equals(StatusSuccess))
}

func TestThread(t *testing.T) {
	tmpDir := MakeTempDir()
	defer os.RemoveAll(tmpDir)
	filePath := path.Join(tmpDir, "casket.tkh")
	dbm := NewDBM()
	status := dbm.Open(filePath, true, "truncate=true")
	CheckTrue(t, status.Equals(StatusSuccess))
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
					CheckTrue(t, status.Equals(StatusNotFoundError))
				}
			} else if random.Intn(5) == 0 {
				status := dbm.Remove(key)
				CheckTrue(t, status.Equals(StatusSuccess) || status.Equals(StatusNotFoundError))
				delete(*recordMap, key)
			} else {
				CheckTrue(t, dbm.Set(key, value, true).Equals(StatusSuccess))
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
			CheckTrue(t, status.Equals(StatusSuccess))
			CheckEq(t, value, gotValue)
		}
	}
	CheckEq(t, numRecords, dbm.CountSimple())
	CheckTrue(t, dbm.Close().Equals(StatusSuccess))
}

func TestSearch(t *testing.T) {
	tmpDir := MakeTempDir()
	defer os.RemoveAll(tmpDir)
	filePath := path.Join(tmpDir, "casket.tks")
	dbm := NewDBM()
	status := dbm.Open(filePath, true, "truncate=true")
	CheckTrue(t, status.Equals(StatusSuccess))
	for i := 1; i <= 100; i++ {
		key := fmt.Sprintf("%08d", i)
		value := fmt.Sprintf("%d", i*i)
		CheckTrue(t, dbm.Set(key, value, false).Equals(StatusSuccess))
	}
	CheckTrue(t, dbm.Synchronize(false, "reducer=ReduceToFirst").Equals(StatusSuccess))
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
	CheckTrue(t, dbm.Close().Equals(StatusSuccess))
}

// END OF FILE
