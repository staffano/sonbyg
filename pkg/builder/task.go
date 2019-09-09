package builder

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/staffano/sonbyg/pkg/utils"
)

// A Task is a combination of an Operation and a set of Variables
// to call the opration with
type Task struct {
	// The ID must match the Id of the operation in builder
	ID           string
	Dependencies []*Task
	Variables    utils.Variables
	Signature    string
	RunAlways    bool
	done         bool
}

// GetID ...
func (t *Task) GetID() string { return t.ID }

// GetDependencies ...
func (t *Task) GetDependencies() []*Task { return t.Dependencies }

// IsDone ...
func (t *Task) IsDone() bool { return t.done }

// MarkDone marks the Target as done
func (t *Task) MarkDone() { t.done = true }

// AllDependenciesDone checks if all the dependencies
func (t *Task) AllDependenciesDone() bool {
	for _, d := range t.Dependencies {
		if !d.IsDone() {
			return false
		}
	}
	return true
}

// AssignDefaultSignature will generate an md5 signature based on
// the id and a sorted list of the task variables.
func (t *Task) AssignDefaultSignature() {

	var sb strings.Builder
	// Sort variables
	names := t.Variables.Names()
	sort.Strings(names)

	sb.WriteString(t.ID)
	for _, v := range names {
		sb.WriteString(fmt.Sprintf("%v", t.Variables[v]))
	}
	strSig := sb.String()
	h := md5.New()
	io.WriteString(h, strSig)
	t.Signature = hex.EncodeToString(h.Sum(nil))
}

// RunAlwaysSignature generates a random signature that
// will make sure it is always run
func (t *Task) RunAlwaysSignature() {
	b := make([]byte, 8)
	rand.Read(b)
	t.Signature = fmt.Sprintf("%x", b)
}

// DependsOn states a dependency to a another task.
func (t *Task) DependsOn(d ...*Task) {
append_loop:
	for _, a := range d {
		for _, b := range t.Dependencies {
			if a.Signature == b.Signature {
				continue append_loop
			}
		}
		t.Dependencies = append(t.Dependencies, a)
	}
}

// DependsOnSerial makes t depend on d_N, d_N on d_N-1, etc
func (t *Task) DependsOnSerial(d ...*Task) {
	t.DependsOn(d[len(d)-1])
	for i := len(d) - 1; i > 0; i-- {
		d[i].DependsOn(d[i-1])
	}

}

// DumpStr will show interesting stuff in da string...
func (t *Task) DumpStr(deps bool) string {
	var sb strings.Builder
	sb.WriteRune('[')
	sb.WriteString(t.ID)
	sb.WriteRune(':')
	sb.WriteString(t.IDSig())
	sb.WriteRune(']')
	if deps {
		sb.WriteString(" depends on: ")
		for _, d := range t.Dependencies {
			sb.WriteString(d.DumpStr(false))
		}
	}
	return sb.String()
}

// IDSig returns id:Signature
func (t Task) IDSig() string {
	var sb strings.Builder
	sb.WriteString(t.ID)
	sb.WriteRune(':')
	if len(t.Signature) > 6 {
		sb.WriteString(t.Signature[0:6])
	} else {
		sb.WriteString(t.Signature)
	}
	return sb.String()
}

// NewTask creates a new Task of type TaskImpl
func NewTask(b *Builder, id string) *Task {
	t := Task{}
	t.ID = id
	return &t
}
