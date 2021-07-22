package tkrzw

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func check_eq(t *testing.T, want interface{}, got interface{}) {
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

func check_ne(t *testing.T, want interface{}, got interface{}) {
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

func check_true(t *testing.T, got bool) {
	_, _, line, _ := runtime.Caller(1)
	if !got {
		t.Errorf("line=%d: not true", line)
	}
}

func check_false(t *testing.T, got bool) {
	_, _, line, _ := runtime.Caller(1)
	if got {
		t.Errorf("line=%d: true", line)
	}
}

func TestAssertion(t *testing.T) {
	check_eq(t, 2, 2)
	check_eq(t, 2.0, 2.0)
	check_eq(t, "two", "two")
	check_eq(t, []byte("two"), []byte("two"))
	check_eq(t, nil, nil)
	check_eq(t, 2, 2.0)
	check_eq(t, 2, "2")
	check_eq(t, 2.0, 2)
	check_eq(t, 2.0, "2")
	check_eq(t, "2", 2)
	check_eq(t, []byte("2"), 2)
	check_ne(t, 2, 3)
	check_ne(t, 2.0, 3.0)
	check_ne(t, "two", "three")
	check_ne(t, []byte("two"), []byte("three"))
	check_ne(t, nil, 0)
	check_true(t, true)
	check_true(t, 1 > 0)
	check_false(t, false)
	check_false(t, 1 < 0)
}

type Person struct {
	Name string
}

func (self Person) String() string {
	return fmt.Sprintf("I'm %s.", self.Name)
}

func TestToString(t *testing.T) {
	check_eq(t, "123", ToString("123"))
	check_eq(t, "123", ToString([]byte("123")))
	check_eq(t, "123", ToString(123))
	check_eq(t, "123.000000", ToString(123.0))
	check_eq(t, "true", ToString(true))
	check_eq(t, "false", ToString(false))
	check_eq(t, "Boom", ToString(errors.New("Boom")))
	check_eq(t, "I'm Alice.", ToString(Person{"Alice"}))
	check_eq(t, "I'm Bob.", ToString(&Person{"Bob"}))
}

func TestToByteArray(t *testing.T) {
	check_eq(t, []byte("123"), ToByteArray("123"))
	check_eq(t, []byte("123"), ToByteArray([]byte("123")))
	check_eq(t, []byte("123"), ToByteArray(123))
	check_eq(t, []byte("123.000000"), ToByteArray(123.0))
	check_eq(t, []byte("true"), ToByteArray(true))
	check_eq(t, []byte("false"), ToByteArray(false))
	check_eq(t, []byte("Boom"), ToByteArray(errors.New("Boom")))
	check_eq(t, []byte("I'm Alice."), ToByteArray(Person{"Alice"}))
	check_eq(t, []byte("I'm Bob."), ToByteArray(&Person{"Bob"}))
}

func TestToInt(t *testing.T) {
	check_eq(t, -123, ToInt("-123"))
	check_eq(t, -123, ToInt("-123.0"))
	check_eq(t, -123, ToInt(int8(-123)))
	check_eq(t, -123, ToInt(int16(-123)))
	check_eq(t, -123, ToInt(int32(-123)))
	check_eq(t, -123, ToInt(int64(-123)))
	check_eq(t, 255, ToInt(uint8(255)))
	check_eq(t, 255, ToInt(uint16(255)))
	check_eq(t, 255, ToInt(uint32(255)))
	check_eq(t, 255, ToInt(uint64(255)))
	check_eq(t, -255, ToInt(float32(-255)))
	check_eq(t, -255, ToInt(float64(-255)))
}

func TestToFloat(t *testing.T) {
	check_eq(t, -123.0, ToFloat("-123"))
	check_eq(t, -123.5, ToFloat("-123.5"))
	check_eq(t, -123.0, ToFloat(int8(-123)))
	check_eq(t, -123.0, ToFloat(int16(-123)))
	check_eq(t, -123.0, ToFloat(int32(-123)))
	check_eq(t, -123.0, ToFloat(int64(-123)))
	check_eq(t, 255.0, ToFloat(uint8(255)))
	check_eq(t, 255.0, ToFloat(uint16(255)))
	check_eq(t, 255.0, ToFloat(uint32(255)))
	check_eq(t, 255.0, ToFloat(uint64(255)))
	check_eq(t, -255.5, ToFloat(float32(-255.5)))
	check_eq(t, -255.5, ToFloat(float64(-255.5)))
}

func TestStatus(t *testing.T) {
	s := NewStatus()
	check_eq(t, STATUS_SUCCESS, s.GetCode())
	check_eq(t, "", s.GetMessage())
	check_true(t, s.Is(s))
	check_true(t, s.Is(*s))
	check_true(t, s.Is(STATUS_SUCCESS))
	check_false(t, s.Is(STATUS_NOT_FOUND_ERROR))
	check_false(t, s.Is(100))
	check_eq(t, "SUCCESS", s)
	check_true(t, s.IsOK())
	s = NewStatus(STATUS_NOT_FOUND_ERROR, "foobar")
	check_eq(t, STATUS_NOT_FOUND_ERROR, s.GetCode())
	check_eq(t, "foobar", s.GetMessage())
	check_eq(t, "NOT_FOUND_ERROR: foobar", s)
	check_true(t, s.Is(s))
	check_true(t, s.Is(*s))
	check_true(t, s.Is(STATUS_NOT_FOUND_ERROR))
	check_false(t, s.Is(STATUS_SUCCESS))
	check_false(t, s.IsOK())
	check_false(t, s.Is(100))
	s = NewStatus1(STATUS_SUCCESS)
	check_eq(t, STATUS_SUCCESS, s.GetCode())
	check_eq(t, "", s.GetMessage())
	s = NewStatus2(STATUS_NOT_FOUND_ERROR, "bazquux")
	check_eq(t, STATUS_NOT_FOUND_ERROR, s.GetCode())
	check_eq(t, "bazquux", s.GetMessage())
}
