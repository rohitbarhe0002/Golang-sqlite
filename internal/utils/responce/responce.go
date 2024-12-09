package responce

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Responce struct{
	Status string `json:"status"`
	Error string  `json:"error"`
}

const (
	StatusOk = "OK"
	StatusError = "Error"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Responce{
return Responce{
	Status:StatusError ,
	Error: err.Error(),
}
}

func ValidationError(errs validator.ValidationErrors)Responce{
	var errMsgs []string
	for _, err := range errs{
		switch err.ActualTag(){
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("filed %s is required field",err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("filed %s is invalid ",err.Field()))
		}
	}
	return Responce{
		Status: StatusError,
		Error: strings.Join(errMsgs,","),
	}
}