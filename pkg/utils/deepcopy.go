package utils

import (
	"encoding/gob"
	"log"
	"os"
	"reflect"
)

// DeepCopy copies all public properties of src into the return value,
// that is a value of the same type as src.
func DeepCopy(src interface{}) interface{} {
	var ro reflect.Value
	ptr := false
	if reflect.TypeOf(src).Kind() == reflect.Ptr {
		ro = reflect.New(reflect.ValueOf(src).Elem().Type())
		ptr = true
	} else {
		ro = reflect.New(reflect.TypeOf(src))
	}
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(w)
	err = enc.Encode(src)
	if err != nil {
		return err
	}
	if err := gob.NewDecoder(r).Decode(ro.Interface()); err != nil {
		log.Fatal(err)
	}
	if ptr {
		return ro.Interface()
	}
	// Since we allways get a pointer to the type by reflect.New, we need
	// to unref that pointer before going to the value
	res := ro.Elem().Interface()

	return res
}
