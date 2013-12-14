package main 

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/codegangsta/martini"
)

// GetFoods returns the list of foods (possibly filtered).
func GetFoods(r *http.Request, enc Encoder, db DB) string {
	// Get the query string arguments, if any
	qs := r.URL.Query()
	name, description, venue := qs.Get("name"), qs.Get("description"), qs.Get("venue")
	if name != "" || description != "" || venue != "" {
		// At least one filter, use Find()
		return Must(enc.Encode(toIface(db.Find(name, description, venue))...))
	}
	// Otherwise, return all foods
	return Must(enc.Encode(toIface(db.GetAll())...))
}

// GetFood returns the requested food.
func GetFood(enc Encoder, db DB, params martini.Params) (int, string) {
	id, err := strconv.Atoi(params["id"])
	foo := db.Get(id)
	if err != nil || foo == nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, Must(enc.Encode(
				NewError(ErrCodeNotExist, fmt.Sprintf("the food with id %s does not exist", params["id"]))))
	}
	return http.StatusOK, Must(enc.Encode(foo))
}

// PostFood creates the posted food.
func PostFood(w http.ResponseWriter, r *http.Request, enc Encoder, db DB) (int, string) {
	foo := getPostFood(r)
	id, err := db.Create(foo)
	switch err {
	case ErrAlreadyExists:
		// Duplicate
		return http.StatusConflict, Must(enc.Encode(
			NewError(ErrCodeAlreadyExists, fmt.Sprintf("the food '%s' from '%s' already exists", foo.Name, foo.Venue))))
	case nil:
		// TODO : Location is expected to be an absolute URI, as per the RFC2616
		w.Header().Set("Location", fmt.Sprintf("/foods/%d", id))
		return http.StatusCreated, Must(enc.Encode(foo))
	default:
		panic(err)
	}
}

// PutFood changes the specified food.
func PutFood(r *http.Request, enc Encoder, db DB, params martini.Params) (int, string) {
	foo, err := getPutFood(r, params)
	if err != nil {
		// Invalid id, 404
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the food with id %s does not exist", params["id"]))))
	}
	err = db.Update(foo)
	switch err {
	case ErrAlreadyExists:
		return http.StatusConflict, Must(enc.Encode(
			NewError(ErrCodeAlreadyExists, fmt.Sprintf("the food '%s' from '%s' already exists", foo.Name, foo.Venue))))
	case nil:
		return http.StatusOK, Must(enc.Encode(foo))
	default:
		panic(err)
	}
}

// Parse the request body, load into a Food structure.
func getPostFood(r *http.Request) *Food {
	name, description, venue := r.FormValue("name"), r.FormValue("description"), r.FormValue("venue")
	return &Food{
		Name: name,
		Description: description,
		Venue: venue,
	}
}

// Like getPostFood, but additionally, parse and store the `id` query string.
func getPutFood(r *http.Request, params martini.Params) (*Food, error) {
	foo := getPostFood(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return nil, err
	}
	foo.Id = id 
	return foo, nil
}

// Martini requires that 2 parameters are returned to treat the first one as the
// status code. Delete is an idempotent action, but this does not mean it should
// always return 204 - No content, idempotence relates to the state of the server
// after the request, not the returned status code. So I return a 404 - Not found 
// if the id does not exist.
func DeleteFood(enc Encoder, db DB, params martini.Params) (int, string) {
	id, err := strconv.Atoi(params["id"])
	foo := db.Get(id)
	ir err != nil || foo == nil {
		return http.StatusNotFound, Must(enc.Encode(
			NewError(ErrCodeNotExist, fmt.Sprintf("the food with the id %s does not exist", params["id"]))))
	}
	db.Delete(id)
	return http.StatusNoContent, ""
}

func toIface(v []*Food) []interface{} {
	if len(v) == 0 {
		return nil
	}
	ifs := make([]interface{}, len(v))
	for i, v := range v {
		ifs[i] = v
	}
	return ifs
}
