package utils

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"runtime/debug"
	"strings"
)

var resolveVarRe *regexp.Regexp

func init() {
	resolveVarRe = regexp.MustCompile(`\$\{([a-zA-Z0-9-_]+)\}`)
}

// Variables represents a map of variables.
type Variables map[string]interface{}

// IsNull checks if the variable list is empty
func (v Variables) IsNull() bool {
	return v == nil
}

// Resolve the string and replace all variable references with
// values from the Variables
func (v Variables) Resolve(str string) string {
	// Get all variable references.
	ss := resolveVarRe.FindAllStringSubmatch(str, -1)
	if len(ss) == 0 {
		return str
	}
	for _, m := range ss {
		// If there are circular dependencies, the program will crash...
		val, ok := v[m[1]]
		if !ok {
			debug.PrintStack()
			log.Fatalf("Variable [%s] not found in Variables", m[1])
		}
		str = strings.ReplaceAll(str, m[0], v.Resolve(fmt.Sprintf("%v", val)))
	}
	return str
}

// Get will get the resolved value of key
func (v Variables) Get(key string) string {
	val, ok := v[key]
	if !ok {
		debug.PrintStack()
		log.Fatalf("Variable [%s] not found in Variables", key)
	}
	return v.Resolve(fmt.Sprintf("%v", val))
}

// Set the value of the variable
func (v *Variables) Set(key string, val interface{}) {
	(*v)[key] = val
}

// Append a set of variables
func (v *Variables) Append(vals Variables) {
	for k2, v2 := range vals {
		(*v)[k2] = v2
	}
}

// ImportEnv  imports environment variables in the form "key=val"
func (v *Variables) ImportEnv(env []string) {
	for _, vi := range env {
		pair := strings.Split(vi, "=")
		(*v)[pair[0]] = pair[1]
	}
}

// ExportEnv exports environment variables in the form "key=val"
// The value will use the default string conversion of the underlying type
func (v Variables) ExportEnv() []string {
	res := []string{}
	for k, x := range v {
		res = append(res, fmt.Sprintf("%s=%v", k, x))
	}
	return res
}

// AppendEnv takes an existing env and merges in (overwrites) the variables
func (v Variables) AppendEnv(env []string) []string {
	t := env
	te := v.ExportEnv()
	for k, v := range te {
		t[k] = v
	}
	return t
}

// PrependEnv prepends the value to the variable with key.
// If the key variable is not an environment list variable, the
// result is unknown
func (v *Variables) PrependEnv(key string, value string) {
	vv, ok := (*v)[key].(string)
	if !ok {
		(*v)[key] = value
		return
	}
	(*v)[key] = fmt.Sprintf("%s%c%s", value, os.PathListSeparator, vv)
	log.Println((*v)[key])
}

// Copy will select the keys from variables and create a new
// variables. if "all" is specified all variables will be copied.
func (v Variables) Copy(keys ...string) Variables {
	if len(keys) == 0 {
		return Variables{}
	}
	if keys[0] == "all" {
		c := DeepCopy(v)
		return c.(Variables)
	}
	r := Variables{}
	for _, k := range keys {
		// First we do a shallow copy and then
		// create the deep one
		// Check if it's options (*VAR)
		key := strings.TrimPrefix(k, "*")
		optional := k[0] == '*'
		val, ok := v[key]
		if !ok {
			if !optional {
				debug.PrintStack()
				log.Fatalf("Variable %s not found", k)
			}
		} else {
			t := DeepCopy(val)
			r[key] = t
		}
	}
	return r
}

// Names returns a list of all variable name
func (v Variables) Names() []string {
	res := make([]string, len(v))
	i := 0
	for k := range v {
		res[i] = k
		i = i + 1
	}
	return res
}

// Dump prints the variables
func (v Variables) Dump() {
	for k, v := range v {
		log.Printf("%s=%v", k, v)
	}
}

// Printf will check the existence of VERBOSE and log accordingly
// If "DEBUG_BUILD" is set, then all variables will be printed
func (v Variables) Printf(format string, args ...interface{}) {
	var verbose, debug, ok bool
	var iface interface{}
	if iface, ok = v["VERBOSE"]; ok {
		verbose = iface.(bool)
		if iface, ok = v["DEBUG_BUILD"]; ok {
			debug = iface.(bool)
		}
		if verbose {
			log.Printf(format, args...)
		}
		if debug {
			log.Println("Variables:")
			for k, v := range v {
				log.Printf("  %s=%v", k, v)
			}
			log.Println("---")
		}
	}
}

// NewVariables creates a new Variables object and import
// the environment variables
func NewVariables(vars []string) Variables {
	v := Variables{}
	v.ImportEnv(vars)
	return v
}

// ResolveAll will resolve all variable references within Variables
func (v *Variables) ResolveAll() {
	// Wrap the original in a reflect.Value
	original := reflect.ValueOf(v)

	copy := reflect.New(original.Type()).Elem()
	v.translateRecursive(copy, original)

	// Remove the reflection wrapper
	vv := copy.Interface().(*Variables)
	for k, vvv := range *vv {
		(*v)[k] = vvv
	}
}

func (v *Variables) translateRecursive(copy, original reflect.Value) {
	if !original.CanInterface() {
		log.Fatalf("Not allowed values as variables")
	}
	switch original.Kind() {
	// The first cases handle nested structures and translate them recursively

	// If it is a pointer we need to unwrap and call once again
	case reflect.Ptr:
		// To get the actual value of the original we have to call Elem()
		// At the same time this unwraps the pointer so we don't end up in
		// an infinite recursion
		originalValue := original.Elem()
		// Check if the pointer is nil
		if !originalValue.IsValid() {
			return
		}
		// Allocate a new object and set the pointer to it
		copy.Set(reflect.New(originalValue.Type()))
		// Unwrap the newly created pointer
		v.translateRecursive(copy.Elem(), originalValue)

	// If it is an interface (which is very similar to a pointer), do basically the
	// same as for the pointer. Though a pointer is not the same as an interface so
	// note that we have to call Elem() after creating a new object because otherwise
	// we would end up with an actual pointer
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()
		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()
		v.translateRecursive(copyValue, originalValue)
		copy.Set(copyValue)

	// If it is a struct we translate each field
	case reflect.Struct:
		for i := 0; i < original.NumField(); i++ {
			v.translateRecursive(copy.Field(i), original.Field(i))
		}

	// If it is a slice we create a new slice and translate each element
	case reflect.Slice:
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i++ {
			v.translateRecursive(copy.Index(i), original.Index(i))
		}

	// If it is a map we create a new map and translate each value
	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			// New gives us a pointer, but again we want the value
			copyValue := reflect.New(originalValue.Type()).Elem()
			v.translateRecursive(copyValue, originalValue)
			copy.SetMapIndex(key, copyValue)
		}

	// Otherwise we cannot traverse anywhere so this finishes the the recursion

	// If it is a string translate it (yay finally we're doing what we came for)
	case reflect.String:
		if original.CanInterface() {
			translatedString := v.Resolve(original.Interface().(string))
			copy.SetString(translatedString)
		}
	// And everything else will simply be taken from the original
	default:
		if original.CanInterface() {
			copy.Set(original)
		}
	}

}
