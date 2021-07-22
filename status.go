package tkrzw

// Type alias for the enumeration of status codes.
type StatusCode int32

// Enumeration of status codes.
const (
	// Success.
	STATUS_SUCCESS                = StatusCode(0)
	// Generic error whose cause is unknown.
	STATUS_UNKNOWN_ERROR          = StatusCode(1)
	// Generic error from underlying systems.
	STATUS_SYSTEM_ERROR           = StatusCode(2)
	// Error that the feature is not implemented.
	STATUS_NOT_IMPLEMENTED_ERROR  = StatusCode(3)
	// Error that a precondition is not met.
	STATUS_PRECONDITION_ERROR     = StatusCode(4)
	// Error that a given argument is invalid.
	STATUS_INVALID_ARGUMENT_ERROR = StatusCode(5)
	// Error that the operation is canceled.
	STATUS_CANCELED_ERROR         = StatusCode(6)
	// Error that a specific resource is not found.
	STATUS_NOT_FOUND_ERROR        = StatusCode(7)
	// Error that the operation is not permitted.
	STATUS_PERMISSION_ERROR       = StatusCode(8)
	// Error that the operation is infeasible.
	STATUS_INFEASIBLE_ERROR       = StatusCode(9)
	// Error that a specific resource is duplicated.
	STATUS_DUPLICATION_ERROR      = StatusCode(10)
	// Error that internal data are broken.
	STATUS_BROKEN_DATA_ERROR      = StatusCode(11)
	// Generic error caused by the application logic.
	STATUS_APPLICATION_ERROR      = StatusCode(12)
)

// Status of operations
type Status struct {
	// The status code.
	code    StatusCode
	// The status message.
	message string
}

// Gets the name of a status code.
//
// @param code The status code.
// @return The name of the status code.
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
// @param args If the first parameter is given, it is treated as the status code.  If the second parameter is given, it is treated as the status message.
// @return The pointer to the created status object.
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

// Makes a new status, with a code.
//
// @param code The status code.
// @return The pointer to the created status object.
func NewStatus1(code StatusCode) *Status {
	return &Status{code, ""}
}

// Makes a new status, with a code and a message.
//
// @param code The status code.
// @param code The status message.
// @return The pointer to the created status object.
func NewStatus2(code StatusCode, message string) *Status {
	return &Status{code, message}
}

// Makes a string representing the status.
//
// @return The string representing the status.
func (self *Status) String() string {
	expr := StatusCodeName(self.code)
	if len(self.message) > 0 {
		expr += ": " + self.message
	}
	return expr
}

// Gets the status code.
//
// @return The status code.
func (self *Status) GetCode() StatusCode {
	return self.code
}

// Gets the status message.
//
// @return The status message.
func (self *Status) GetMessage() string {
	return self.message
}

// Sets the code and the message.
//
// @param code The status code.
// @param message An arbitrary status message.
func (self *Status) Set(args ...interface{}) {
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

// Checks whether the status equal to another status.
//
// @param rhs a status object or a status code.
// @param true for the both operands are equal, or false if not.
func (self *Status) Equals(rhs interface{}) bool {
	switch rhs := rhs.(type) {
	case Status:
		return self.code == rhs.code
	case *Status:
		return self.code == rhs.code
	case StatusCode:
		return self.code == rhs
	}
	return false
}

// Returns true if the status is success.
//
// @return true if the status is success, or false if not.
func (self *Status) IsOK() bool {
	return self.code == STATUS_SUCCESS
}

// Causes a panic if the status is not success.
func (self *Status) OrDie() {
  if self.code != STATUS_SUCCESS {
		panic(self.String())
	}
}
