package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type envelope map[string]any

// writeJson encoding data to json format and send it to client
func (app *application) writeJson(w http.ResponseWriter, r *http.Request, status int, data envelope, header http.Header) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// add any potential header
	for key, value := range header {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json)

	return nil
}

// errorResponse sends an error in json-formatted form
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJson(w, r, status, env, nil)
	if err != nil {
		app.logger.Print(err)
		// empty 500 response
		w.WriteHeader(500)
	}
}

// serverErrorResponse sends unexpected error that server encounters at runtime
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Print(err)
	app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
}

// notFoundResponse sends 404 NotFound in form json to client
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusNotFound, "The requested resource could not be found")
}

// methodNotAllowedResponse sends 405 Method Not Allowed in form json to client
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusMethodNotAllowed, fmt.Sprintf("%s method is not support for this resource", r.Method))
}

// readBodyToJSON decode request body into dst
func (app *application) readBodyToJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576 // 2^20
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // return error if there is unknown field

	err := decoder.Decode(dst)

	return err
}

func (app *application) failedValidatorResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
