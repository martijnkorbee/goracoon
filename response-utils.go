package goracoon

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"
)

func (gr *Goracoon) ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only have only a single json value")
	}

	return nil
}

func (gr *Goracoon) WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (gr *Goracoon) WriteXML(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := xml.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (gr *Goracoon) ShowFile(w http.ResponseWriter, r *http.Request, pathToFile, fileName string) error {
	fp := path.Join(pathToFile, fileName)
	fileToServe := filepath.Clean(fp)
	w.Header().Set("Content-Type", fmt.Sprintf("attachement; file=\"%s\"", fileName))
	http.ServeFile(w, r, fileToServe)

	return nil
}

func (gr *Goracoon) DownloadFile(w http.ResponseWriter, r *http.Request, pathToFile, fileName string) error {
	fp := path.Join(pathToFile, fileName)
	fileToServe := filepath.Clean(fp)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachement; file=\"%s\"", fileName))
	http.ServeFile(w, r, fileToServe)

	return nil
}

// response errors

func (gr *Goracoon) Error404(w http.ResponseWriter, r *http.Request) {
	gr.ErrorStatus(w, http.StatusNotFound)
}

func (gr *Goracoon) Error500(w http.ResponseWriter, r *http.Request) {
	gr.ErrorStatus(w, http.StatusInternalServerError)
}

func (gr *Goracoon) ErrorUnAuthorized(w http.ResponseWriter, r *http.Request) {
	gr.ErrorStatus(w, http.StatusUnauthorized)
}

func (gr *Goracoon) ErrorForbidden(w http.ResponseWriter, r *http.Request) {
	gr.ErrorStatus(w, http.StatusForbidden)
}

func (gr *Goracoon) ErrorStatus(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
