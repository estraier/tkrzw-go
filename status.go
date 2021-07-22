package tkrzw

type StatusCode int32

const (
	STATUS_SUCCESS                = StatusCode(0)
	STATUS_UNKNOWN_ERROR          = StatusCode(1)
	STATUS_SYSTEM_ERROR           = StatusCode(2)
	STATUS_NOT_IMPLEMENTED_ERROR  = StatusCode(3)
	STATUS_PRECONDITION_ERROR     = StatusCode(4)
	STATUS_INVALID_ARGUMENT_ERROR = StatusCode(5)
	STATUS_CANCELED_ERROR         = StatusCode(6)
	STATUS_NOT_FOUND_ERROR        = StatusCode(7)
	STATUS_PERMISSION_ERROR       = StatusCode(8)
	STATUS_INFEASIBLE_ERROR       = StatusCode(9)
	STATUS_DUPLICATION_ERROR      = StatusCode(10)
	STATUS_BROKEN_DATA_ERROR      = StatusCode(11)
	STATUS_APPLICATION_ERROR      = StatusCode(12)
)

type Status struct {
	code    StatusCode
	message string
}

func StatusCodeName(code StatusCode) string {
	switch code {
	case STATUS_SUCCESS:
		return "SUCCESS"
	case STATUS_UNKNOWN_ERROR:
		return "UNKNOWN_ERROR"
	case STATUS_SYSTEM_ERROR:
		return "SYSTEM_ERROR"
	case STATUS_NOT_IMPLEMENTED_ERROR:
		return "NOT_IMPLEMENTED_ERROR"
	case STATUS_PRECONDITION_ERROR:
		return "PRECONDITION_ERROR"
	case STATUS_INVALID_ARGUMENT_ERROR:
		return "INVALID_ARGUMENT_ERROR"
	case STATUS_CANCELED_ERROR:
		return "CANCELED_ERROR"
	case STATUS_NOT_FOUND_ERROR:
		return "NOT_FOUND_ERROR"
	case STATUS_PERMISSION_ERROR:
		return "PERMISSION_ERROR"
	case STATUS_INFEASIBLE_ERROR:
		return "INFEASIBLE_ERROR"
	case STATUS_DUPLICATION_ERROR:
		return "DUPLICATION_ERROR"
	case STATUS_BROKEN_DATA_ERROR:
		return "BROKEN_DATA_ERROR"
	case STATUS_APPLICATION_ERROR:
		return "APPLICATION_ERROR"
	}
	return "invalid code"
}

// Makes a new status, with variable length arguments.
//
//   Love:
//    Tus
//
// If the first parameter is given, it is treated as the status code.
// If the second parameter is given, it is treated as the status message.
// It returns the pointer to the created object.
func NewStatus(args ...interface{}) *Status {
	code := STATUS_SUCCESS
	if len(args) > 0 {
		code = args[0].(StatusCode)
	}
	message := ""
	if len(args) > 1 {
		message = args[1].(string)
	}
	return &Status{code, message}
}

func NewStatus1(code StatusCode) *Status {
	return &Status{code, ""}
}

func NewStatus2(code StatusCode, message string) *Status {
	return &Status{code, message}
}

func (self Status) String() string {
	expr := StatusCodeName(self.code)
	if len(self.message) > 0 {
		expr += ": " + self.message
	}
	return expr
}

func (self Status) GetCode() StatusCode {
	return self.code
}

func (self Status) GetMessage() string {
	return self.message
}

func (self Status) Set(args ...interface{}) {
	code := STATUS_SUCCESS
	if len(args) > 0 {
		code = args[0].(StatusCode)
	}
	message := ""
	if len(args) > 1 {
		message = args[1].(string)
	}
	self.code = code
	self.message = message
}

func (self Status) Is(x interface{}) bool {
	switch x := x.(type) {
	case Status:
		return self.code == x.code
	case *Status:
		return self.code == x.code
	case StatusCode:
		return self.code == x
	}
	return false
}

func (self Status) IsOK() bool {
	return self.code == STATUS_SUCCESS
}
