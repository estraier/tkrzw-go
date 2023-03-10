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
	"math"
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
	if IsNilData(want) {
		if !IsNilData(got) {
			println(got, IsNilData(got), got == nil)
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
				t.Errorf("line=%d: not equal: want=%s, got=%s", line,
					StatusCodeName(want), StatusCodeName(got))
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
	if IsNilData(want) {
		if IsNilData(got) {
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
				t.Errorf("line=%d: equal: want=%s, got=%s",
					line, StatusCodeName(want), StatusCodeName(got))
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
	CheckEq(t, interface{}(nil), interface{}(nil))
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

func TestSpecialData(t *testing.T) {
	myAnyBytes := make([]byte, len(AnyBytes))
	myAnyString := string([]byte(AnyString))
	CheckFalse(t, IsAnyData(0))
	CheckFalse(t, IsAnyData(nil))
	CheckFalse(t, IsAnyData(""))
	CheckTrue(t, IsAnyData(AnyBytes))
	CheckTrue(t, IsAnyData(AnyString))
	copy(myAnyBytes, AnyBytes)
	CheckFalse(t, IsAnyData(myAnyBytes))
	CheckFalse(t, IsAnyData(myAnyString))
	CheckTrue(t, IsAnyBytes(AnyBytes))
	CheckFalse(t, IsAnyBytes(myAnyBytes))
	CheckTrue(t, IsAnyString(AnyString))
	CheckFalse(t, IsAnyString(myAnyString))
	myNilString := string([]byte(NilString))
	CheckFalse(t, IsNilData(0))
	CheckTrue(t, IsNilData(nil))
	CheckFalse(t, IsNilData(""))
	CheckTrue(t, IsNilData(NilString))
	CheckFalse(t, IsNilData(myNilString))
	CheckFalse(t, IsNilString(""))
	CheckTrue(t, IsNilString(NilString))
	CheckFalse(t, IsNilString(myNilString))
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
	params := ParseParams("a=A, bb = BBB, ccc=")
	CheckEq(t, 3, len(params))
	CheckEq(t, "A", params["a"])
	CheckEq(t, "BBB", params["bb"])
	CheckEq(t, "", params["ccc"])
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
	s2 := NewStatus2(StatusNotImplementedError, "void")
	s.Join(s2)
	CheckEq(t, "NOT_FOUND_ERROR: bazquux", s.String())
	s.Set(StatusSuccess, "OK")
	s.Join(s2)
	CheckEq(t, "NOT_IMPLEMENTED_ERROR: void", s.String())
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
	status := dbm.Open(filePath, true, ParseParams("truncate=true,num_buckets=5"))
	CheckEq(t, StatusSuccess, status)
	CheckTrue(t, len(dbm.String()) > len(filePath))
	CheckFalse(t, dbm.Check("one"))
	CheckTrue(t, dbm.Set("one", "first", false).IsOK())
	CheckEq(t, StatusDuplicationError, dbm.Set("one", "uno", false))
	CheckTrue(t, dbm.Check("one"))
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
	strValue, status := dbm.GetStr([]byte("three"))
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "third:3", strValue)
	strValue, status = dbm.GetStr([]byte("fourth"))
	CheckEq(t, StatusNotFoundError, status)
	CheckEq(t, "", strValue)
	CheckEq(t, "first", dbm.GetSimple("one", "*"))
	CheckEq(t, "second", dbm.GetStrSimple("two", "*"))
	CheckEq(t, "third:3", dbm.GetStrSimple([]byte("three"), "*"))
	CheckEq(t, "*", dbm.GetStrSimple([]byte("four"), "*"))
	CheckEq(t, StatusSuccess, dbm.Remove("one"))
	CheckEq(t, StatusSuccess, dbm.Remove("two"))
	CheckEq(t, StatusSuccess, dbm.Remove([]byte("three")))
	CheckEq(t, StatusNotFoundError, dbm.Remove([]byte("fourth")))
	CheckEq(t, StatusSuccess, dbm.Set("日本", "東京", true))
	CheckEq(t, "東京", dbm.GetSimple("日本", "*"))
	CheckEq(t, StatusSuccess, dbm.Remove("日本"))
	CheckEq(t, 0, dbm.CountSimple())
	CheckEq(t, StatusSuccess, dbm.CompareExchange("num", nil, "first"))
	CheckEq(t, "first", dbm.GetSimple("num", "*"))
	CheckEq(t, StatusInfeasibleError, dbm.CompareExchange("num", nil, "first"))
	CheckEq(t, StatusSuccess, dbm.CompareExchange("num", "first", "second"))
	CheckEq(t, "second", dbm.GetSimple("num", "*"))
	CheckEq(t, StatusSuccess, dbm.CompareExchange("num", "second", nil))
	CheckEq(t, "*", dbm.GetSimple("num", "*"))
	CheckEq(t, StatusInfeasibleError, dbm.CompareExchange("xyz", AnyString, AnyString))
	CheckEq(t, StatusSuccess, dbm.CompareExchange("xyz", nil, "abc"))
	CheckEq(t, StatusSuccess, dbm.CompareExchange("xyz", AnyBytes, AnyBytes))
	CheckEq(t, "abc", dbm.GetSimple("xyz", "*"))
	CheckEq(t, StatusSuccess, dbm.CompareExchange("xyz", AnyBytes, "def"))
	CheckEq(t, "def", dbm.GetSimple("xyz", "*"))
	CheckEq(t, StatusSuccess, dbm.CompareExchange("xyz", AnyString, nil))
	CheckEq(t, "*", dbm.GetSimple("xyz", "*"))
	actual, status := dbm.CompareExchangeAndGetStr("xyz", nil, "123")
	CheckEq(t, StatusSuccess, status)
	CheckTrue(t, IsNilString(actual))
	actual, status = dbm.CompareExchangeAndGetStr("xyz", AnyString, AnyString)
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "123", actual)
	actual, status = dbm.CompareExchangeAndGetStr("xyz", AnyString, NilString)
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "123", actual)
	incValue, status := dbm.Increment("num", 2, 100)
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, 102, incValue)
	incValue, status = dbm.Increment("num", 3, 100)
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, 105, incValue)
	CheckEq(t, StatusSuccess, dbm.Remove("num"))
	oldValue, status := dbm.SetAndGet("zero", "nil", false)
	CheckEq(t, nil, oldValue)
	CheckEq(t, StatusSuccess, status)
	oldValue, status = dbm.SetAndGet("zero", "nothing", false)
	CheckEq(t, "nil", oldValue)
	CheckEq(t, StatusDuplicationError, status)
	oldStrValue, status := dbm.SetAndGetStr("zero", "void", false)
	CheckEq(t, "nil", *oldStrValue)
	CheckEq(t, StatusDuplicationError, status)
	oldValue, status = dbm.RemoveAndGet("zero")
	CheckEq(t, "nil", oldValue)
	CheckEq(t, StatusSuccess, status)
	oldValue, status = dbm.RemoveAndGet("zero")
	CheckEq(t, nil, oldValue)
	CheckEq(t, StatusNotFoundError, status)
	oldStrValue, status = dbm.SetAndGetStr("zero", "void", false)
	CheckEq(t, nil, oldValue)
	CheckEq(t, StatusSuccess, status)
	oldStrValue, status = dbm.RemoveAndGetStr("zero")
	CheckEq(t, "void", *oldStrValue)
	CheckEq(t, StatusSuccess, status)
	oldStrValue, status = dbm.RemoveAndGetStr("zero")
	CheckEq(t, nil, oldValue)
	CheckEq(t, StatusNotFoundError, status)
	records := map[string]string{"one": "first", "two": "second"}
	CheckEq(t, StatusSuccess, dbm.SetMultiStr(records, false))
	CheckEq(t, StatusSuccess, dbm.AppendMultiStr(records, ":"))
	keys := []string{"one", "two", "three"}
	records = dbm.GetMultiStr(keys)
	CheckEq(t, 2, len(records))
	CheckEq(t, "first:first", records["one"])
	CheckEq(t, "second:second", records["two"])
	rawRecords := dbm.GetMulti(keys)
	CheckEq(t, 2, len(rawRecords))
	CheckEq(t, "first:first", rawRecords["one"])
	CheckEq(t, "second:second", rawRecords["two"])
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
	CheckEq(t, StatusInfeasibleError, dbm.CompareExchangeMulti(
		[]KeyValuePair{{[]byte("xyz"), AnyBytes}},
		[]KeyValuePair{{[]byte("xyz"), []byte("abc")}}))
	CheckEq(t, StatusSuccess, dbm.CompareExchangeMulti(
		[]KeyValuePair{{[]byte("xyz"), nil}},
		[]KeyValuePair{{[]byte("xyz"), []byte("abc")}}))
	CheckEq(t, StatusSuccess, dbm.CompareExchangeMulti(
		[]KeyValuePair{{[]byte("xyz"), []byte("abc")}},
		[]KeyValuePair{{[]byte("xyz"), []byte("def")}}))
	CheckEq(t, "def", dbm.GetSimple("xyz", "*"))
	CheckEq(t, StatusSuccess, dbm.CompareExchangeMulti(
		[]KeyValuePair{{[]byte("xyz"), []byte("def")}},
		[]KeyValuePair{{[]byte("xyz"), nil}}))
	CheckEq(t, "*", dbm.GetSimple("xyz", "*"))
	CheckEq(t, 0, dbm.CountSimple())
	set4 := []KeyValueStrPair{KeyValueStrPair{"one", NilString}, KeyValueStrPair{"two", NilString}}
	set5 := []KeyValueStrPair{KeyValueStrPair{"one", "apple"}, KeyValueStrPair{"two", "orange"}}
	CheckEq(t, StatusSuccess, dbm.CompareExchangeMultiStr(set4, set5))
	CheckEq(t, StatusInfeasibleError, dbm.CompareExchangeMultiStr(set4, set5))
	CheckEq(t, "apple", dbm.GetSimple("one", "*"))
	CheckEq(t, "orange", dbm.GetSimple("two", "*"))
	CheckEq(t, StatusSuccess, dbm.CompareExchangeMultiStr(set5, set4))
	CheckEq(t, StatusInfeasibleError, dbm.CompareExchangeMultiStr(set5, set4))
	CheckEq(t, "*", dbm.GetSimple("one", "*"))
	CheckEq(t, "*", dbm.GetSimple("two", "*"))
	CheckEq(t, StatusInfeasibleError, dbm.CompareExchangeMultiStr(
		[]KeyValueStrPair{{"xyz", AnyString}},
		[]KeyValueStrPair{{"xyz", "abc"}}))
	CheckEq(t, StatusSuccess, dbm.CompareExchangeMultiStr(
		[]KeyValueStrPair{{"xyz", NilString}},
		[]KeyValueStrPair{{"xyz", "abc"}}))
	CheckEq(t, StatusSuccess, dbm.CompareExchangeMultiStr(
		[]KeyValueStrPair{{"xyz", "abc"}},
		[]KeyValueStrPair{{"xyz", "def"}}))
	CheckEq(t, "def", dbm.GetStrSimple("xyz", "*"))
	CheckEq(t, StatusSuccess, dbm.CompareExchangeMultiStr(
		[]KeyValueStrPair{{"xyz", "def"}},
		[]KeyValueStrPair{{"xyz", NilString}}))
	CheckEq(t, "*", dbm.GetStrSimple("xyz", "*"))
	CheckEq(t, 0, dbm.CountSimple())
	fileSize, status := dbm.GetFileSize()
	CheckEq(t, StatusSuccess, status)
	CheckTrue(t, fileSize > 0)
	CheckEq(t, fileSize, dbm.GetFileSizeSimple())
	gotFilePath, status := dbm.GetFilePath()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, filePath, gotFilePath)
	CheckEq(t, filePath, dbm.GetFilePathSimple())
	timestamp, status := dbm.GetTimestamp()
	CheckEq(t, StatusSuccess, status)
	CheckTrue(t, timestamp > 0)
	CheckEq(t, timestamp, dbm.GetTimestampSimple())
	for i := 1; i <= 10; i++ {
		CheckTrue(t, dbm.Set(i, i*i, true).IsOK())
	}
	CheckEq(t, 10, dbm.CountSimple())
	tobe, status := dbm.ShouldBeRebuilt()
	CheckEq(t, StatusSuccess, status)
	CheckTrue(t, tobe)
	CheckTrue(t, dbm.ShouldBeRebuiltSimple())
	CheckEq(t, StatusSuccess, dbm.Rebuild(nil))
	CheckEq(t, StatusSuccess, dbm.Synchronize(true, nil))
	CheckEq(t, StatusSuccess, dbm.CopyFileData(copyPath, false))
	CheckEq(t, StatusSuccess, dbm.Clear())
	CheckEq(t, 0, dbm.CountSimple())
	CheckEq(t, StatusSuccess, dbm.Close())
	CheckEq(t, StatusSuccess, dbm.Open(filePath, true, nil))
	copyDBM := NewDBM()
	CheckEq(t, StatusSuccess, copyDBM.Open(copyPath, false, nil))
	CheckEq(t, 10, copyDBM.CountSimple())
	CheckEq(t, StatusSuccess, copyDBM.Export(dbm))
	CheckEq(t, 10, dbm.CountSimple())
	CheckEq(t, StatusSuccess, copyDBM.Close())
	inspRecords := dbm.Inspect()
	CheckEq(t, "10", inspRecords["num_records"])
	CheckEq(t, "HashDBM", inspRecords["class"])
	iter := dbm.MakeIterator()
	CheckTrue(t, len(iter.String()) > 0)
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
		strKey, strValue, status := iter.GetStr()
		CheckEq(t, StatusSuccess, status)
		CheckEq(t, strKey, string(key))
		CheckEq(t, strValue, string(value))
		records[strKey] = strValue
		oneKey, status := iter.GetKey()
		CheckEq(t, StatusSuccess, status)
		CheckEq(t, strKey, string(oneKey))
		oneStrKey, status := iter.GetKeyStr()
		CheckEq(t, StatusSuccess, status)
		CheckEq(t, strKey, oneStrKey)
		oneValue, status := iter.GetValue()
		CheckEq(t, StatusSuccess, status)
		CheckEq(t, strValue, string(oneValue))
		oneStrValue, status := iter.GetValueStr()
		CheckEq(t, StatusSuccess, status)
		CheckEq(t, strValue, oneStrValue)
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
	strValue, status = iter.GetValueStr()
	CheckEq(t, StatusSuccess, iter.Remove())
	CheckEq(t, 9, dbm.CountSimple())
	CheckEq(t, StatusSuccess, dbm.Set("key1", "value1", true))
	CheckEq(t, StatusSuccess, dbm.Rekey("key1", "key2", false, false))
	CheckEq(t, "*", dbm.GetSimple("key1", "*"))
	CheckEq(t, "value1", dbm.GetSimple("key2", "*"))
	CheckEq(t, StatusSuccess, dbm.Rekey("key2", "key1", false, true))
	CheckEq(t, "value1", dbm.GetSimple("key1", "*"))
	CheckEq(t, "value1", dbm.GetSimple("key2", "*"))
	CheckEq(t, StatusDuplicationError, dbm.Rekey("key1", "key2", false, false))
	CheckEq(t, StatusNotFoundError, dbm.Rekey("key0", "key2", false, false))
	iter.Destruct()
	CheckEq(t, StatusSuccess, dbm.Close())
	os.Remove(copyPath)
	CheckEq(t, StatusSuccess, RestoreDatabase(filePath, copyPath, "", -1, ""))
	CheckEq(t, StatusSuccess, copyDBM.Open(copyPath, false, nil))
	CheckEq(t, 11, copyDBM.CountSimple())
	CheckEq(t, StatusSuccess, copyDBM.Close())
	CheckEq(t, StatusSuccess, dbm.Open(filePath, false, nil))
	CheckEq(t, StatusSuccess, dbm.Close())
}

func TestDBMProcess(t *testing.T) {
	tmpDir := MakeTempDir()
	defer os.RemoveAll(tmpDir)
	filePath := path.Join(tmpDir, "casket.tkh")
	dbm := NewDBM()
	status := dbm.Open(filePath, true, ParseParams("truncate=true,num_buckets=1000"))
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, StatusSuccess,
		dbm.Process("abc", func(k []byte, v []byte) interface{} { return nil }, true))
	CheckEq(t, "*", dbm.GetStrSimple("abc", "*"))
	CheckEq(t, StatusSuccess,
		dbm.Process("abc", func(k []byte, v []byte) interface{} { return RemoveBytes }, true))
	CheckEq(t, "*", dbm.GetStrSimple("abc", "*"))
	CheckEq(t, StatusSuccess,
		dbm.Process("abc", func(k []byte, v []byte) interface{} { return []byte("ABCDE") }, true))
	CheckEq(t, "ABCDE", dbm.GetStrSimple("abc", "*"))
	proc1 := func(k []byte, v []byte) interface{} {
		CheckEq(t, "abc", string(k))
		CheckEq(t, "ABCDE", string(v))
		return nil
	}
	CheckEq(t, StatusSuccess, dbm.Process("abc", proc1, true))
	CheckEq(t, StatusSuccess,
		dbm.Process("abc", func(k []byte, v []byte) interface{} { return RemoveBytes }, true))
	CheckEq(t, "*", dbm.GetStrSimple("abc", "*"))
	proc2 := func(k []byte, v []byte) interface{} {
		CheckEq(t, "abc", string(k))
		CheckEq(t, nil, v)
		return nil
	}
	CheckEq(t, StatusSuccess, dbm.Process("abc", proc2, false))
	for i := 0; i < 10; i++ {
		CheckEq(t, StatusSuccess,
			dbm.Process(ToString(i), func(k []byte, v []byte) interface{} {
				return []byte(ToString(i * i))
			}, true))
	}
	CheckEq(t, 10, dbm.CountSimple())
	count_full := 0
	count_empty := 0
	proc3 := func(k []byte, v []byte) interface{} {
		if k == nil {
			count_empty += 1
		} else {
			count_full += 1
			num := ToInt(k)
			CheckEq(t, num*num, ToInt(v))
		}
		return nil
	}
	CheckEq(t, StatusSuccess, dbm.ProcessEach(proc3, false))
	CheckEq(t, 2, count_empty)
	CheckEq(t, 10, count_full)
	proc4 := func(k []byte, v []byte) interface{} {
		if k == nil {
			return nil
		}
		num := ToInt(v)
		return int(math.Sqrt(float64(num)))
	}
	CheckEq(t, StatusSuccess, dbm.ProcessEach(proc4, true))
	proc5 := func(k []byte, v []byte) interface{} {
		if k == nil {
			return nil
		}
		CheckEq(t, ToInt(k), ToInt(v))
		return RemoveBytes
	}
	CheckEq(t, StatusSuccess, dbm.ProcessEach(proc5, true))
	CheckEq(t, 0, dbm.CountSimple())
	ops := []KeyProcPair{
		{"one", func(k []byte, v []byte) interface{} { return "hop" }},
		{"two", func(k []byte, v []byte) interface{} { return "step" }},
		{"three", func(k []byte, v []byte) interface{} { return "jump" }},
	}
	CheckEq(t, StatusSuccess, dbm.ProcessMulti(ops, true))
	proc6 := func(k []byte, v []byte) interface{} {
		if v == nil {
			return "x"
		}
		return ToString(v) + ToString(v)
	}
	ops = []KeyProcPair{
		{"one", func(k []byte, v []byte) interface{} { return RemoveBytes }},
		{"two", func(k []byte, v []byte) interface{} { return RemoveBytes }},
		{"three", proc6},
		{"four", proc6},
		{"three", proc6},
		{"four", proc6},
	}
	CheckEq(t, StatusSuccess, dbm.ProcessMulti(ops, true))
	CheckEq(t, 2, dbm.CountSimple())
	CheckEq(t, "*", dbm.GetStrSimple("one", "*"))
	CheckEq(t, "*", dbm.GetStrSimple("two", "*"))
	CheckEq(t, "jumpjumpjumpjump", dbm.GetStrSimple("three", "*"))
	CheckEq(t, "xx", dbm.GetStrSimple("four", "*"))
	CheckEq(t, StatusSuccess, dbm.Close())
}

func TestDBMIterator(t *testing.T) {
	tmpDir := MakeTempDir()
	defer os.RemoveAll(tmpDir)
	filePath := path.Join(tmpDir, "casket.tkt")
	dbm := NewDBM()
	status := dbm.Open(filePath, true, ParseParams("truncate=true"))
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
	step_count := 0
	CheckEq(t, StatusSuccess, iter.First())
	for {
		if step_count%2 == 0 {
			key, value, status := iter.Step()
			if !status.IsOK() {
				CheckEq(t, StatusNotFoundError, status)
				break
			}
			CheckEq(t, ToInt(key)*ToInt(key), ToInt(value))
		} else {
			strKey, strValue, status := iter.StepStr()
			if !status.IsOK() {
				CheckEq(t, StatusNotFoundError, status)
				break
			}
			CheckEq(t, ToInt(strKey)*ToInt(strKey), ToInt(strValue))
		}
		CheckEq(t, StatusSuccess, status)
		step_count++
	}
	CheckEq(t, count, step_count)
	pop_count := 0
	for {
		if pop_count%2 == 0 {
			key, value, status := dbm.PopFirst()
			if !status.IsOK() {
				CheckEq(t, StatusNotFoundError, status)
				break
			}
			CheckEq(t, ToInt(key)*ToInt(key), ToInt(value))
		} else {
			strKey, strValue, status := dbm.PopFirstStr()
			if !status.IsOK() {
				CheckEq(t, StatusNotFoundError, status)
				break
			}
			CheckEq(t, ToInt(strKey)*ToInt(strKey), ToInt(strValue))
		}
		CheckEq(t, StatusSuccess, status)
		pop_count++
	}
	CheckEq(t, step_count, pop_count)
	CheckEq(t, 0, dbm.CountSimple())
	CheckEq(t, StatusSuccess, dbm.PushLast("foo", 0))
	strKey, strValue, status := dbm.PopFirstStr()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "\x00\x00\x00\x00\x00\x00\x00\x00", strKey)
	CheckEq(t, "foo", strValue)
	iter.Destruct()
	CheckEq(t, StatusSuccess, dbm.Close())
}

func TestDBMThread(t *testing.T) {
	tmpDir := MakeTempDir()
	defer os.RemoveAll(tmpDir)
	filePath := path.Join(tmpDir, "casket.tkh")
	dbm := NewDBM()
	status := dbm.Open(filePath, true, ParseParams("truncate=true"))
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
				iter.Destruct()
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
	CheckEq(t, StatusSuccess, dbm.Open(filePath, true, ParseParams("truncate=true")))
	CheckEq(t, StatusSuccess, dbm.Set("one", "first", true))
	CheckEq(t, StatusSuccess, dbm.Set("two", "second", true))
	CheckEq(t, 2, dbm.CountSimple())
	copyFile := NewFile()
	CheckEq(t, StatusSuccess, copyFile.Open(copyPath, true, ParseParams("truncate=true")))
	CheckEq(t, StatusSuccess, dbm.ExportToFlatRecords(copyFile))
	CheckEq(t, StatusSuccess, dbm.Clear())
	CheckEq(t, 0, dbm.CountSimple())
	CheckEq(t, StatusSuccess, dbm.ImportFromFlatRecords(copyFile))
	CheckEq(t, 2, dbm.CountSimple())
	CheckEq(t, "first", dbm.GetSimple("one", "*"))
	CheckEq(t, "second", dbm.GetSimple("two", "*"))
	CheckEq(t, StatusSuccess, copyFile.Close())
	CheckEq(t, StatusSuccess, copyFile.Open(copyPath, true, ParseParams("truncate=true")))
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
	status := dbm.Open(filePath, true, ParseParams("truncate=true"))
	CheckEq(t, StatusSuccess, status)
	for i := 1; i <= 100; i++ {
		key := fmt.Sprintf("%08d", i)
		value := fmt.Sprintf("%d", i*i)
		CheckEq(t, StatusSuccess, dbm.Set(key, value, false))
	}
	CheckEq(t, StatusSuccess, dbm.Synchronize(false, ParseParams("reducer=ReduceToFirst")))
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

func TestAsyncDBM(t *testing.T) {
	tmpDir := MakeTempDir()
	defer os.RemoveAll(tmpDir)
	filePath := path.Join(tmpDir, "casket.tkh")
	copyPath := path.Join(tmpDir, "casket-copy.tkh")
	dbm := NewDBM()
	status := dbm.Open(filePath, true, ParseParams("truncate=true,num_buckets=5"))
	CheckEq(t, StatusSuccess, status)
	async := NewAsyncDBM(dbm, 4)
	CheckTrue(t, len(async.String()) > 0)
	future := async.Set("one", "hop", false)
	CheckTrue(t, len(future.String()) > 0)
	future.Wait(0)
	CheckTrue(t, future.Wait(-1))
	CheckEq(t, StatusSuccess, future.Get())
	CheckEq(t, StatusSuccess, async.Set("two", "step", false).Get())
	CheckEq(t, StatusSuccess, async.Set("three", "jump", false).Get())
	CheckEq(t, StatusDuplicationError, async.Set("three", "jump", false).Get())
	async.Set("three", "jump", false).Destruct()
	CheckEq(t, 3, dbm.CountSimple())
	CheckEq(t, StatusSuccess, async.Append("one", "1", ":").Get())
	CheckEq(t, StatusSuccess, async.Append("two", "2", ":").Get())
	value_bytes, status := async.Get("one").GetBytes()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "hop:1", value_bytes)
	strValue, status := async.Get("two").GetStr()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "step:2", strValue)
	CheckEq(t, StatusSuccess, async.Remove("one").Get())
	CheckEq(t, StatusSuccess, async.Remove("two").Get())
	CheckEq(t, StatusSuccess, async.Remove("three").Get())
	CheckEq(t, 0, dbm.CountSimple())
	records := map[string]string{"one": "first", "two": "second"}
	CheckEq(t, StatusSuccess, async.SetMultiStr(records, false).Get())
	CheckEq(t, StatusSuccess, async.AppendMultiStr(records, ":").Get())
	keys := []string{"one", "two", "three"}
	records, status = async.GetMulti(keys).GetMapStr()
	CheckEq(t, StatusNotFoundError, status)
	CheckEq(t, 2, len(records))
	CheckEq(t, "first:first", records["one"])
	CheckEq(t, "second:second", records["two"])
	rawRecords, status := async.GetMulti(keys).GetMap()
	CheckEq(t, StatusNotFoundError, status)
	CheckEq(t, 2, len(rawRecords))
	CheckEq(t, "first:first", rawRecords["one"])
	CheckEq(t, "second:second", rawRecords["two"])
	CheckEq(t, StatusNotFoundError, async.RemoveMulti(keys).Get())
	keys = []string{"three"}
	CheckEq(t, 0, dbm.CountSimple())
	CheckEq(t, StatusSuccess, async.CompareExchange("num", nil, "first").Get())
	CheckEq(t, "first", dbm.GetSimple("num", "*"))
	CheckEq(t, StatusInfeasibleError, async.CompareExchange("num", nil, "first").Get())
	CheckEq(t, StatusSuccess, async.CompareExchange("num", "first", "second").Get())
	CheckEq(t, "second", dbm.GetSimple("num", "*"))
	CheckEq(t, StatusSuccess, async.CompareExchange("num", "second", nil).Get())
	CheckEq(t, "*", dbm.GetSimple("num", "*"))
	CheckEq(t, StatusInfeasibleError, async.CompareExchange("xyz", AnyBytes, AnyBytes).Get())
	CheckEq(t, StatusSuccess, async.CompareExchange("xyz", NilString, "abc").Get())
	CheckEq(t, StatusSuccess, async.CompareExchange("xyz", AnyString, AnyString).Get())
	CheckEq(t, "abc", dbm.GetSimple("xyz", "*"))
	CheckEq(t, StatusSuccess, async.CompareExchange("xyz", AnyString, "def").Get())
	CheckEq(t, "def", dbm.GetSimple("xyz", "*"))
	CheckEq(t, StatusSuccess, async.CompareExchange("xyz", AnyString, NilString).Get())
	CheckEq(t, "*", dbm.GetSimple("xyz", "*"))
	incValue, status := async.Increment("num", 2, 100).GetInt()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, 102, incValue)
	incValue, status = async.Increment("num", 3, 100).GetInt()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, 105, incValue)
	CheckEq(t, StatusSuccess, async.Remove("num").Get())
	set1 := []KeyValuePair{KeyValuePair{[]byte("one"), []byte(nil)},
		KeyValuePair{[]byte("two"), []byte(nil)}}
	set2 := []KeyValuePair{KeyValuePair{[]byte("one"), []byte("ichi")},
		KeyValuePair{[]byte("two"), []byte("ni")}}
	set3 := []KeyValuePair{KeyValuePair{[]byte("one"), []byte("uno")},
		KeyValuePair{[]byte("two"), []byte("dos")}}
	CheckEq(t, StatusSuccess, async.CompareExchangeMulti(set1, set2).Get())
	CheckEq(t, "ichi", dbm.GetSimple("one", "*"))
	CheckEq(t, "ni", dbm.GetSimple("two", "*"))
	CheckEq(t, StatusInfeasibleError, async.CompareExchangeMulti(set1, set2).Get())
	CheckEq(t, StatusSuccess, async.CompareExchangeMulti(set2, set3).Get())
	CheckEq(t, "uno", dbm.GetSimple("one", "*"))
	CheckEq(t, "dos", dbm.GetSimple("two", "*"))
	CheckEq(t, StatusSuccess, async.CompareExchangeMulti(set3, set1).Get())
	CheckEq(t, "*", dbm.GetSimple("one", "*"))
	CheckEq(t, "*", dbm.GetSimple("two", "*"))
	CheckEq(t, 0, dbm.CountSimple())
	set4 := []KeyValueStrPair{KeyValueStrPair{"one", NilString}, KeyValueStrPair{"two", NilString}}
	set5 := []KeyValueStrPair{KeyValueStrPair{"one", "apple"}, KeyValueStrPair{"two", "orange"}}
	CheckEq(t, StatusSuccess, async.CompareExchangeMultiStr(set4, set5).Get())
	CheckEq(t, StatusInfeasibleError, async.CompareExchangeMultiStr(set4, set5).Get())
	CheckEq(t, "apple", dbm.GetSimple("one", "*"))
	CheckEq(t, "orange", dbm.GetSimple("two", "*"))
	CheckEq(t, StatusSuccess, async.CompareExchangeMultiStr(set5, set4).Get())
	CheckEq(t, StatusInfeasibleError, async.CompareExchangeMultiStr(set5, set4).Get())
	CheckEq(t, "*", dbm.GetSimple("one", "*"))
	CheckEq(t, "*", dbm.GetSimple("two", "*"))
	CheckEq(t, 0, dbm.CountSimple())
	CheckEq(t, StatusSuccess, async.Set("hello", "world", false).Get())
	CheckEq(t, StatusSuccess, async.Rebuild(nil).Get())
	CheckEq(t, StatusSuccess, async.Synchronize(false, nil).Get())
	CheckEq(t, StatusSuccess, async.CopyFileData(copyPath, false).Get())
	CheckEq(t, StatusSuccess, async.Clear().Get())
	CheckEq(t, StatusInfeasibleError, async.CompareExchangeMultiStr(
		[]KeyValueStrPair{KeyValueStrPair{"xyz", AnyString}},
		[]KeyValueStrPair{KeyValueStrPair{"xyz", AnyString}}).Get())
	CheckEq(t, StatusSuccess, async.CompareExchangeMultiStr(
		[]KeyValueStrPair{KeyValueStrPair{"xyz", NilString}},
		[]KeyValueStrPair{KeyValueStrPair{"xyz", "abc"}}).Get())
	CheckEq(t, "abc", dbm.GetSimple("xyz", "*"))
	CheckEq(t, StatusSuccess, async.CompareExchangeMultiStr(
		[]KeyValueStrPair{KeyValueStrPair{"xyz", AnyString}},
		[]KeyValueStrPair{KeyValueStrPair{"xyz", "def"}}).Get())
	CheckEq(t, "def", dbm.GetSimple("xyz", "*"))
	CheckEq(t, StatusSuccess, async.CompareExchangeMultiStr(
		[]KeyValueStrPair{KeyValueStrPair{"xyz", AnyString}},
		[]KeyValueStrPair{KeyValueStrPair{"xyz", NilString}}).Get())
	CheckEq(t, "*", dbm.GetSimple("xyz", "*"))
	CheckEq(t, 0, dbm.CountSimple())
	copyDBM := NewDBM()
	CheckEq(t, StatusSuccess, copyDBM.Open(copyPath, true, nil))
	CheckEq(t, 1, copyDBM.CountSimple())
	CheckEq(t, "world", copyDBM.GetSimple("hello", "*"))
	CheckEq(t, StatusSuccess, async.Set("japan", "tokyo", false).Get())
	CheckEq(t, StatusSuccess, async.Export(copyDBM).Get())
	CheckEq(t, 2, copyDBM.CountSimple())
	CheckEq(t, "tokyo", copyDBM.GetSimple("japan", "*"))
	CheckEq(t, StatusSuccess, copyDBM.Close())
	CheckEq(t, StatusSuccess, async.Set("hello", "good-bye", true).Get())
	CheckEq(t, StatusSuccess, async.Set("hi", "bye", true).Get())
	CheckEq(t, StatusSuccess, async.Set("chao", "adios", true).Get())
	values, status := async.Search("begin", "h", 0).GetArrayStr()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, 2, len(values))
	rawValues, status := async.Search("begin", "h", 0).GetArray()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, 2, len(rawValues))
	CheckEq(t, StatusSuccess, async.Clear().Get())
	CheckEq(t, StatusSuccess, async.Set("aa", "AAA", false).Get())
	CheckEq(t, StatusSuccess, async.Rekey("aa", "bb", true, false).Get())
	CheckEq(t, "*", dbm.GetSimple("aa", "*"))
	CheckEq(t, "AAA", dbm.GetSimple("bb", "*"))
	key, value, status := async.PopFirst().GetPair()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "bb", key)
	CheckEq(t, "AAA", value)
	CheckEq(t, StatusSuccess, async.Set("cc", "CCC", false).Get())
	strKey, strValue, status := async.PopFirst().GetPairStr()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "cc", strKey)
	CheckEq(t, "CCC", strValue)
	CheckEq(t, StatusSuccess, async.PushLast("foo", 0).Get())
	strKey, strValue, status = async.PopFirst().GetPairStr()
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "\x00\x00\x00\x00\x00\x00\x00\x00", strKey)
	CheckEq(t, "foo", strValue)
	async.Destruct()
	CheckEq(t, StatusSuccess, dbm.Close())
}

func TestFile(t *testing.T) {
	tmpDir := MakeTempDir()
	defer os.RemoveAll(tmpDir)
	filePath := path.Join(tmpDir, "casket.txt")
	file := NewFile()
	status := file.Open(filePath, true, ParseParams("truncate=true"))
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
	strData, status := file.ReadStr(3, 5)
	CheckEq(t, StatusSuccess, status)
	CheckEq(t, "defgh", strData)
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
