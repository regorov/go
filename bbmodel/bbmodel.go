package bbmodel

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"time"
)

var StringType = reflect.TypeOf("")
var BoolType = reflect.TypeOf(false)
var IntType = reflect.TypeOf(int(0))
var UintType = reflect.TypeOf(uint(0))
var Float64Type = reflect.TypeOf(float64(0.0))
var TimeType = reflect.TypeOf(time.Now())

const IDLength = 32
const IDChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func CreateID() (id string) {
	for i := 0; i < IDLength; i++ {
		id += IDChars[rand.Intn(len(IDChars))]
	}

	return id
}

func handleResponse(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		st, ok := err.(HTTPStatus)
		if !ok {
			log.Printf("Error: %s", err.Error())
			st = HTTPStatus(http.StatusInternalServerError)
		}

		msg := http.StatusText(int(st))
		log.Printf("%s %s -> %d %s", r.Method, r.RequestURI, st, msg)
		w.WriteHeader(int(st))
		fmt.Fprintln(w, msg)

	} else {
		log.Printf("%s %s -> 200 OK", r.Method, r.RequestURI)
	}
}

func HandleModel(path string, model Model) {
	http.Handle(path, NewModelHandler(model))
}

func HandleCollection(path string, collection *Collection) {
	http.Handle(path, NewCollectionHandler(collection))
	http.Handle(path+"/", NewCollectionModelHandler(path+"/", collection))
}

type HTTPStatus int

func (e HTTPStatus) Error() (s string) {
	return fmt.Sprintf("%d", int(e))
}

type ModelSpec map[string]reflect.Type

type Model map[string]interface{}

func NewModel(spec ModelSpec) (model Model) {
	model = make(Model)
	model["id"] = CreateID()

	for name, t := range spec {
		model[name] = reflect.Zero(t)
	}

	return model
}

type Collection struct {
	Models map[uint64]Model
	Spec   ModelSpec
}

func NewCollection(spec ModelSpec) (c *Collection) {
	return &Collection{
		Models: make(map[uint64]Model),
		Spec:   spec,
	}
}

func (c *Collection) New() (model Model) {
	model = NewModel(c.Spec)
	c.Add(model)
	return model
}

func (c *Collection) Add(model Model) {
	c.Models[model["id"]] = model
}

func (c *Collection) Remove(model Model) {
	delete(c.Models, model["id"])
}

func (c *Collection) Get(id string) (model Model, ok bool) {
	model, ok = c.Models[id]
	return model, ok
}

type ModelHandler struct {
	Model Model
}

func NewModelHandler(model Model) (h ModelHandler) {
	return &ModelHandler{
		Model: *model,
	}
}

func (h ModelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleResponse(w, r, h.serve(w, r))
}

func (h ModelHandler) serve(w http.ResponseWriter, r *http.Request) (err error) {
	switch r.Method {
	case "GET":
		err = json.NewEncoder(w).Encode(h.Model)
		if err != nil {
			return err
		}

	case "POST", "PUT":
		err = json.NewDecoder(r.Body).Decode(&h.Model)
		if err != nil {
			return HTTPStatus(http.StatusBadRequest)
		}

		err = json.NewEncoder(w).Encode(h.Model)
		if err != nil {
			return err
		}

	default:
		return HTTPStatus(http.StatusMethodNotAllowed)
	}

	return nil
}

type CollectionHandler struct {
	Collection *Collection
}

func NewCollectionHandler(c *Collection) (h *CollectionHandler) {
	return &CollectionHandler{
		Collection: c,
	}
}

func (h *CollectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleResponse(w, r, h.serve(w, r))
}

func (h *CollectionHandler) serve(w http.RespnoseWriter, r *http.Request) (err error) {
	switch r.Method {
	case "GET":
		var models []Model

		for _, model := range h.Collection.Models {
			models = append(models, model)
		}

		err = json.NewEncoder(w).Encode(models)
		if err != nil {
			return err
		}

	case "POST", "PUT":
		model := h.Collection.New()

		err = json.NewDecoder(r.Body).Decode(&model)
		if err != nil {
			return HTTPStatus(http.StatusBadRequest)
		}

		err = json.NewEncoder(w).Encode(model)
		if err != nil {
			return err
		}

	default:
		return HTTPStatus(http.StatusMethodNotAllowed)
	}

	return nil
}

type CollectionModelHandler struct {
	BaseURL    string
	Collection *Collection
}

func NewCollectionModelHandler(baseURL string, c *Collection) (h *CollectionModelHandler) {
	return &CollectionModelHandler{
		BaseURL:    baseURL,
		Collection: c,
	}
}

func (h *CollectionModelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleResponse(w, r, h.serve(w, r))
}

func (h *CollectionModelHandler) serve(w http.RespnoseWriter, r *http.Request) (err error) {
	id := r.URL.Path[len(h.BaseURL):]

	switch r.Method {
	case "GET":
		model, ok := h.Collection.Get(id)
		if !ok {
			return HTTPStatus(http.StatusNotFound)
		}

		err = json.NewEncoder(w).Encode(model)
		if err != nil {
			return err
		}

	case "POST", "PUT":
		model, ok := h.Collection.Get(id)
		if !ok {
			model = h.Collection.New()
		}

		err = json.NewDecoder(r.Body).Decode(&model)
		if err != nil {
			return HTTPStatus(http.StatusBadRequest)
		}

		err = json.NewEncoder(w).Encode(model)
		if err != nil {
			return err
		}

	case "DELETE":
		model, ok := h.Collection.Get(id)
		if !ok {
			return HTTPStatus(http.StatusNotFound)
		}

		h.Collection.Delete(model)

	default:
		return HTTPStatus(http.StatusMethodNotAllowed)
	}

	return nil
}
