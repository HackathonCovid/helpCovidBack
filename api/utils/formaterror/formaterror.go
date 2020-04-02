package formaterror

import (
	//"errors"
	"strings"
)

var errorMessages = make(map[string]string)
var err error

// FormatError : handle errors when they occur
func FormatError(errString string) map[string]string {
	if strings.Contains(errString, "email") {
		errorMessages["Taken_email"] = "Email Already Taken"
	}

	if strings.Contains(errString, "title") {
		errorMessages["Taken_title"] = "Title Already Taken"
	}

	if strings.Contains(errString, "hashedPassword") {
		errorMessages["Incorrect_password"] = "Incorrect Password"
	}

	if strings.Contains(errString, "record not found") {
		errorMessages["No_record"] = "No Record Found"
	}

	if strings.Contains(errString, "You already applied to this mission") {
		errorMessages["Double_apply"] = "You cannot apply to this mission twice"
	}

	if len(errorMessages) > 0 {
		return errorMessages
	}

	if len(errorMessages) == 0 {
		errorMessages["Incorrect_details"] = "Incorrect Details"
		return errorMessages
	}

	return nil
}
