/*************************************************************************************************
 * Status interface
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

// Type alias for the enumeration of status codes.
type StatusCode int32

// Enumeration of status codes.
const (
	// Success.
	StatusSuccess = StatusCode(0)
	// Generic error whose cause is unknown.
	StatusUnknownError = StatusCode(1)
	// Generic error from underlying systems.
	StatusSystemError = StatusCode(2)
	// Error that the feature is not implemented.
	StatusNotImplementedError = StatusCode(3)
	// Error that a precondition is not met.
	StatusPreconditionError = StatusCode(4)
	// Error that a given argument is invalid.
	StatusInvalidArgumentError = StatusCode(5)
	// Error that the operation is canceled.
	StatusCanceledError = StatusCode(6)
	// Error that a specific resource is not found.
	StatusNotFoundError = StatusCode(7)
	// Error that the operation is not permitted.
	StatusPermissionError = StatusCode(8)
	// Error that the operation is infeasible.
	StatusInfeasibleError = StatusCode(9)
	// Error that a specific resource is duplicated.
	StatusDuplicationError = StatusCode(10)
	// Error that internal data are broken.
	StatusBrokenDataError = StatusCode(11)
	// Generic error caused by the application logic.
	StatusApplicationError = StatusCode(12)
)

// Status of operations
type Status struct {
	// The status code.
	code StatusCode
	// The status message.
	message string
}

// Gets the name of a status code.
//
// @param code The status code.
// @return The name of the status code.
func StatusCodeName(code StatusCode) string {
	switch code {
	case StatusSuccess:
		return "SUCCESS"
	case StatusUnknownError:
		return "UNKNOWN_ERROR"
	case StatusSystemError:
		return "SYSTEM_ERROR"
	case StatusNotImplementedError:
		return "NOT_IMPLEMENTED_ERROR"
	case StatusPreconditionError:
		return "PRECONDITION_ERROR"
	case StatusInvalidArgumentError:
		return "INVALID_ARGUMENT_ERROR"
	case StatusCanceledError:
		return "CANCELED_ERROR"
	case StatusNotFoundError:
		return "NOT_FOUND_ERROR"
	case StatusPermissionError:
		return "PERMISSION_ERROR"
	case StatusInfeasibleError:
		return "INFEASIBLE_ERROR"
	case StatusDuplicationError:
		return "DUPLICATION_ERROR"
	case StatusBrokenDataError:
		return "BROKEN_DATA_ERROR"
	case StatusApplicationError:
		return "APPLICATION_ERROR"
	}
	return "invalid code"
}

// Makes a new status, with variable length arguments.
//
// @param args If the first parameter is given, it is treated as the status code.  If the second parameter is given, it is treated as the status message.
// @return The pointer to the created status object.
func NewStatus(args ...interface{}) *Status {
	code := StatusSuccess
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

// Makes a string representing the status, to be an error.
//
// @return The string representing the code.  As the string contains only the code, comparison
// of the result strings is not affected by difference of the status messages.
// the additiona 
func (self *Status) Error() string {
	return StatusCodeName(self.code)
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
	code := StatusSuccess
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
	return self.code == StatusSuccess
}

// Causes a panic if the status is not success.
func (self *Status) OrDie() {
	if self.code != StatusSuccess {
		panic(self.String())
	}
}

// END OF FILE
