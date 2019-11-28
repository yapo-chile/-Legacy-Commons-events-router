package infrastructure

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/Yapo/goutils"
	"gopkg.in/gorilla/mux.v1"

	"github.mpi-internal.com/Yapo/events-router/pkg/interfaces/handlers"
)

// InputSource defines the type for an input source
type InputSource string

const (
	// BODY defines the constant for the body params
	BODY InputSource = "body"
	// RAWBODY defines the constant for process the raw body in bytes
	RAWBODY InputSource = "raw"
	// PATH defines the constant for the path params
	PATH InputSource = "path"
	// QUERY defines the constant for the query params
	QUERY InputSource = "query"
	// HEADERS defines the constant for the headers params
	HEADERS InputSource = "headers"
	// COOKIES defines the constant for the cookies params
	COOKIES InputSource = "cookies"
	// FORM defines the constant for the FORM params
	FORM InputSource = "form"

	// NotSeteable defines the error string of this error
	NotSeteable string = "PROVIDED_INPUT_IS_NOT_SETEABLE"
	// NotPointer defines the error string of this error
	NotPointer string = "PROVIDED_INPUT_IS_NOT_A_POINTER"
	// NotStruct defines the error string of this error
	NotStruct string = "PROVIDED_INPUT_IS_NOT_A_STRUCT"
)

// ErrNotSeteable describes an error that occurs when the var is not seteable
var ErrNotSeteable = errors.New(NotSeteable)

// ErrNotPointer describes an error that occurs when the var is not a pointer
var ErrNotPointer = errors.New(NotPointer)

// ErrNotStruct describes an error that occurs when the var is not a struct
var ErrNotStruct = errors.New(NotStruct)

type inputRequest struct {
	output      interface{}
	httpRequest *http.Request
	sources     []InputSource
}

func (ri *inputRequest) Set(output interface{}) handlers.InputRequest {
	ri.output = output
	return ri
}

func (ri *inputRequest) FromJSONBody() handlers.InputRequest {
	ri.sources = append(ri.sources, BODY)
	return ri
}

func (ri *inputRequest) FromRawBody() handlers.InputRequest {
	ri.sources = append(ri.sources, RAWBODY)
	return ri
}

func (ri *inputRequest) FromPath() handlers.InputRequest {
	ri.sources = append(ri.sources, PATH)
	return ri
}

func (ri *inputRequest) FromQuery() handlers.InputRequest {
	ri.sources = append(ri.sources, QUERY)
	return ri
}

func (ri *inputRequest) FromHeaders() handlers.InputRequest {
	ri.sources = append(ri.sources, HEADERS)
	return ri
}

// FromCookies sets request cookies as handler input
func (ri *inputRequest) FromCookies() handlers.InputRequest {
	ri.sources = append(ri.sources, COOKIES)
	return ri
}

// FromForm sets request form as handler input
func (ri *inputRequest) FromForm() handlers.InputRequest {
	ri.sources = append(ri.sources, FORM)
	return ri
}

type inputHandler struct {
	inputRequest *inputRequest
	output       handlers.HandlerInput
}

// NewInputHandler returns a new InputHandler
func NewInputHandler() handlers.InputHandler {
	return &inputHandler{}
}

// NewInputRequest returns a new InputRequest based on a http.request
func (ih *inputHandler) NewInputRequest(r *http.Request) handlers.InputRequest {
	return &inputRequest{httpRequest: r}
}

// SetInputRequest sets the input request and the request input
func (ih *inputHandler) SetInputRequest(ri handlers.InputRequest, hi handlers.HandlerInput) {
	ih.inputRequest = ri.(*inputRequest)
	ih.output = hi
}

// Input does the actual process of getting the input
func (ih *inputHandler) Input() (handlers.HandlerInput, *goutils.Response) {
	if ih.inputRequest.output == nil {
		return ih.output, &goutils.Response{
			Code: http.StatusInternalServerError,
			Body: goutils.GenericError{
				ErrorMessage: "Output was not set correctly",
			},
		}
	}

	hasError := false
	reflectedOutput := reflect.ValueOf(ih.inputRequest.output)
	for _, source := range ih.inputRequest.sources {
		switch source {
		case BODY:
			hasError = hasError ||
				goutils.ParseJSONBody(
					ih.inputRequest.httpRequest,
					ih.inputRequest.output,
				) != nil
		case RAWBODY:
			rawBody, err := ioutil.ReadAll(ih.inputRequest.httpRequest.Body)
			ih.inputRequest.httpRequest.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
			hasError = hasError || err != nil ||
				ih.parseInput(
					map[string]string{"body": string(rawBody)},
					source,
					reflectedOutput,
				) != nil
		case QUERY:
			hasError = hasError ||
				ih.parseInput(
					ih.httpValuesToMap(ih.inputRequest.httpRequest.URL.Query()),
					source,
					reflectedOutput,
				) != nil
		case PATH:
			hasError = hasError ||
				ih.parseInput(
					mux.Vars(ih.inputRequest.httpRequest),
					source,
					reflectedOutput,
				) != nil
		case HEADERS:
			hasError = hasError ||
				ih.parseInput(
					ih.httpValuesToMap(ih.inputRequest.httpRequest.Header),
					source,
					reflectedOutput,
				) != nil
		case COOKIES:
			hasError = hasError ||
				ih.parseInput(
					ih.httpCookiesToMap(ih.inputRequest.httpRequest.Cookies()),
					source,
					reflectedOutput,
				) != nil
		case FORM:
			hasError = hasError ||
				ih.parseInput(
					ih.formToMap(ih.inputRequest.httpRequest.Body),
					source,
					reflectedOutput,
				) != nil
		}
	}

	if hasError {
		return ih.output, &goutils.Response{
			Code: http.StatusInternalServerError,
			Body: goutils.GenericError{
				ErrorMessage: "Output was not set correctly",
			},
		}
	}
	return ih.output, nil
}

func (ih *inputHandler) httpValuesToMap(values map[string][]string) map[string]string {
	outValues := make(map[string]string)
	for k, v := range values {
		outValues[k] = strings.Join(v, ",")
	}
	return outValues
}

func (ih *inputHandler) httpCookiesToMap(cookies []*http.Cookie) map[string]string {
	outValues := make(map[string]string)
	for _, cookie := range cookies {
		outValues[cookie.Name] = cookie.Value
	}
	return outValues
}

func (ih *inputHandler) formToMap(values io.Reader) map[string]string {
	rawBody, _ := ioutil.ReadAll(values)
	params, _ := url.ParseQuery(string(rawBody))
	mapBody := make(map[string]string)
	for key, value := range params {
		mapBody[key] = value[0]
	}
	return mapBody
}

func (ih *inputHandler) parseInput(vars map[string]string, inputTag InputSource, input reflect.Value) error {
	if input.Kind() != reflect.Ptr {
		return ErrNotPointer
	}
	reflectedInput := reflect.Indirect(input)
	// We should only keep going if we can set values
	if !reflectedInput.IsValid() || !reflectedInput.CanSet() {
		return ErrNotSeteable
	}
	// this part of the function is made just for structs
	if reflectedInput.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	// Recursively load inner struct fields
	for i := 0; i < reflectedInput.NumField(); i++ {
		if tag, ok := reflectedInput.Type().Field(i).Tag.Lookup(string(inputTag)); ok {
			switch reflectedInput.Field(i).Kind() {
			case reflect.Struct:
				if ih.parseInput(vars, inputTag, reflectedInput.Field(i).Addr()) != nil {
					continue
				}
			case reflect.String:
				reflectedInput.Field(i).SetString(vars[tag])
			case reflect.Int:
				if value, err := strconv.Atoi(vars[tag]); err == nil {
					reflectedInput.Field(i).Set(reflect.ValueOf(value))
				}
			case reflect.Slice:
				values := strings.Split(vars[tag], ",")
				switch reflectedInput.Field(i).Interface().(type) {
				case []string:
					trimValues := []string{}
					for _, value := range values {
						trimValues = append(trimValues, strings.TrimSpace(value))
					}
					reflectedInput.Field(i).Set(reflect.ValueOf(trimValues))
				case []int:
					intValues := []int{}
					for _, value := range values {
						if val, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
							intValues = append(intValues, val)
						}
					}
					reflectedInput.Field(i).Set(reflect.ValueOf(intValues))
				case []byte:
					reflectedInput.Field(i).SetBytes([]byte(vars[tag]))
				}
			}
		}
	}
	return nil
}
