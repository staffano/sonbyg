package builder

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/staffano/sonbyg/pkg/utils"
)

// An Operation is the functionality associated with a task
type Operation func(utils.Variables) *Artifact

// The Builder is the controller responsible for creating a
// deterministic execution order between the tasks registered
// with the Builder. It is also responsible for keeping the Artifacts
// resulted from the execution of the Tasks and deploying them in
// case of their existence.
// The Builder is also responsible for keeping a Task register that
// will be used when creating implicitly declared Tasks.
type Builder struct {
	Verbose        bool
	Variables      utils.Variables
	WorkspaceRoot  string
	tasks          map[string]*Task
	operations     map[string]Operation
	artifacts      map[string]*Artifact
	resolvedTasks  []*Task
	lock           sync.Mutex
	ctxt           context.Context
	workers        []*Worker
	workerDoneChan chan bool
}

// NewBuilder creates a new builder with some workspace level definitions
func NewBuilder(ctxt context.Context,
	workspaceRoot string,
	workerCount int,
	verbose bool,
	build, host, target string,
	downloadDir, toolsDir, cacheDir, stampDir string) *Builder {

	if workspaceRoot == "" {
		workspaceRoot = path.Join(".", ".sonbyg")
	}
	absWspRoot, err := filepath.Abs(workspaceRoot)
	if err != nil {
		log.Panicf("Could create workspace root path %v", err)
	}
	v := utils.NewVariables(os.Environ())
	if runtime.GOOS == "windows" {
		v["PATH"] = v["Path"]
	}
	b := &Builder{
		Verbose:       verbose,
		WorkspaceRoot: absWspRoot,
		ctxt:          ctxt}
	v["WORKSPACE"] = absWspRoot
	if downloadDir == "" {
		downloadDir = path.Join(absWspRoot, "downloads")
	}
	if toolsDir == "" {
		toolsDir = path.Join(absWspRoot, "tools")
	}
	if cacheDir == "" {
		cacheDir = path.Join(absWspRoot, "cache")
	}
	if stampDir == "" {
		stampDir = path.Join(absWspRoot, "stamps")
	}
	v["VERBOSE"] = verbose
	v["BUILD"] = build
	v["HOST"] = host
	v["TARGET"] = target
	v["CACHE_DIR"], _ = filepath.Abs(cacheDir)
	v["TOOLS_DIR"], _ = filepath.Abs(toolsDir)
	v["DOWNLOAD_DIR"], _ = filepath.Abs(downloadDir)
	v["STAMPS_DIR"], _ = filepath.Abs(stampDir)
	b.Variables = v
	b.EstablishPath(755, v["CACHE_DIR"].(string))
	b.EstablishPath(755, v["TOOLS_DIR"].(string))
	b.EstablishPath(755, v["DOWNLOAD_DIR"].(string))
	b.EstablishPath(755, v["STAMPS_DIR"].(string))

	b.workerDoneChan = make(chan bool)
	b.workers = make([]*Worker, workerCount)
	for i := range b.workers {
		b.workers[i] = NewWorker(ctxt, b.workerDoneChan)
	}
	b.tasks = make(map[string]*Task)
	b.operations = make(map[string]Operation)
	b.artifacts = make(map[string]*Artifact)

	return b
}

// EstablishPath will make sure a path exists within the environment.
// If the supplied path is relative, it will be relative to WorkspaceRoot.
// The returned path is the absolute path to the directory.
// It returns the absolute path to the directory
func (b *Builder) EstablishPath(perm os.FileMode, paths ...string) string {
	p := b.WorkspacePath(paths...)
	if !filepath.IsAbs(p) {
		p = path.Join(b.WorkspaceRoot, p)
	}
	if err := os.MkdirAll(p, perm); err != nil {
		log.Fatalf("Could not create directory %s", p)
	}
	return p
}

// WorkspacePath returns a path within the workspace. If paths is a
// absolute path, it will be checked to see if it's within the
// workspace, and panic if not.
func (b *Builder) WorkspacePath(paths ...string) string {
	p := path.Join(paths...)
	if !filepath.IsAbs(p) {
		p = path.Join(b.WorkspaceRoot, p)
	} else {
		if !strings.HasPrefix(p, b.WorkspaceRoot) {
			log.Panicf("The path %s is not within the workspace %s", p, b.WorkspaceRoot)
		}
	}
	return p
}

// Stamp will create a stamp in the stamp directory, marking
// that the task t is done.
func (b *Builder) Stamp(t *Task) {
	stampFile := path.Join(b.Variables["STAMPS_DIR"].(string), t.Signature)
	if err := ioutil.WriteFile(stampFile, []byte{}, 0644); err != nil {
		log.Fatal(err)
	}

}

// IsStamped checks to see if the Task t is already built
func (b *Builder) IsStamped(t *Task) bool {
	stampFile := path.Join(b.Variables["STAMPS_DIR"].(string), t.Signature)
	return utils.FileExists(stampFile)
}

// Add a new task to the builder. If the task already existed
// the existing task will be returned and provided task
// will be nulled, in order to detect erroneous reuse of
// duplicated task.
func (b *Builder) Add(task *Task, op Operation) *Task {
	if t, ok := b.tasks[string(task.Signature)]; ok {

		// Copy over any dependencies t have that task doesn't.
	loop1:
		for _, d := range task.Dependencies {
			for _, d2 := range t.Dependencies {
				if d.Signature == d2.Signature {
					continue loop1
				}
			}
			t.DependsOn(d)
		}
		return t
	}
	b.tasks[string(task.Signature)] = task
	b.operations[task.ID] = op
	return task
}

// Execute a task. If parallell execution is applied
// then the call should be properly embedded in a go routine
// and safe-guarded by mutexes.
func (b *Builder) Execute(t *Task) {
	b.lock.Lock()
	sig := string(t.Signature)
	if b.artifacts[sig] != nil {
		log.Printf("Task %s already done. Applying Artifact instead", t.ID)
	}
	op := b.operations[t.GetID()]
	b.lock.Unlock()
	t.Variables.Set("_TASK_ID", t.ID)
	a := op(t.Variables)
	b.lock.Lock()
	b.artifacts[string(t.Signature)] = a
	t.MarkDone()
	b.lock.Unlock()
}

// An Artifact describes a payload that is built using data
// matching a signature and wich should be deployed at a
// specific location at the workspace.
type Artifact struct {
	Signature string
	Path      string
	Payload   []byte
}

// Apply the artifact to the workdir
func (a Artifact) Apply(workdir string) {
	if len(a.Payload) == 0 {
		return
	}
	dst := path.Join(workdir, a.Path)
	// Payload is stored using libz format
	src, _ := zip.NewReader(bytes.NewReader(a.Payload), int64(len(a.Payload)))
	for _, file := range src.File {
		zippedFile, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer zippedFile.Close()
		if file.FileInfo().IsDir() {
			utils.CreateDir(dst, file.Name, file.Mode())
		} else {
			utils.CreateFile(zippedFile, dst, file.Name, file.Mode())
		}
	}
	log.Printf("Applied artifact [%s] with length %d at %s", a.Signature, len(a.Payload), dst)
}

// Build will use the registered targets and try to build the provided targets,
// which also needs to be part of the environment...
func (b *Builder) Build(tasks ...*Task) {
	parents := []*Task{}
	// Resolve Dependencies
	for _, t := range tasks {
		b.checkDep(parents, t)
	}

	// ResolvedTargets are now listed in dependency order.
	// The first item have no dependency to anyone.

	// We have N build workers. When all dependencies of an
	// unmade target are made, then that target (T) is ready
	// to be sent to a build worker.
	// A search for a free worker is done and selected (N_i).
	// The worker is attached to the target and set to build it.
	// When a worker is done, the target will be set to done
	// and the worker will report on the Worker_Done channel
	// that it's done.

	// for {
	//   for i in ResolvedTargets {
	//     if i.done => continue
	//     targetsLeft = true
	//     if not i.DependenciesDone => continue
	//     if i.AssignedToWorker => continue
	//     while NoFreeWorkers {
	//       wait for Worker_Done
	//     }
	//     for w in workers {
	//       if w.free	{
	//          Assign i to w and build
	//       }
	//     }
	//   }
	//   if ! targetsLeft => DONE!

	if b.Verbose {
		var sb strings.Builder
		sb.WriteString("ResolvedTargets: ")
		for _, rs := range b.resolvedTasks {
			sb.WriteString(rs.ID)
			sb.WriteString(", ")
		}
		log.Println(sb.String())
	}

	// Find tasks that are already done according to stamps directory
	c := 0
	for _, rs := range b.resolvedTasks {
		if b.IsStamped(rs) {
			rs.MarkDone()
			c++
		}
	}
	log.Printf("%d tasks already done", c)

	var targetsLeft bool
	for {
		targetsLeft = false

		// Select eligible target
		var ts *Task
		for _, rt := range b.resolvedTasks {
			if rt.IsDone() {
				continue
			}
			targetsLeft = true
			if !rt.AllDependenciesDone() {
				continue
			}
			if b.IsAssignedToWorker(rt) {
				continue
			}
			ts = rt
			break
		}

		if ts != nil {
			// There is work to be done.
			// Select a free worker, or wait
			// for a worker to be free and then
			// start the task
			w := b.GetFreeWorker()
			w.Work(b, ts)
			continue
		}

		// Check if all targets built and quit.
		if ts == nil && !targetsLeft {
			log.Printf("Done building all targets")
			if b.Verbose {
				b.ReportWorkerUtilization()
			}
			return
		}

		// Tasks left, but no task available for executing
		// This means that there is tasks executing that has
		// dependencies. Wait until a worker reported done
		// before evaluating task execution again.
		<-b.workerDoneChan
	}
}

// ReportWorkerUtilization prints how many tasks each user have executed
func (b *Builder) ReportWorkerUtilization() {
	for i, w := range b.workers {
		log.Printf("Worker %d have execute %d tasks", i, w.utilization)
	}
}

// GetFreeWorker returns the first worker that indicates it has nothing to do.
func (b *Builder) GetFreeWorker() *Worker {
	for {
		for _, w := range b.workers {
			if w.Free() {
				return w
			}
		}
		<-b.workerDoneChan
	}
}

// IsAssignedToWorker checks to see if there is a worker working on the Target
func (b *Builder) IsAssignedToWorker(t *Task) bool {
	for _, w := range b.workers {
		if w.WorkingOn() == t {
			return true
		}
	}
	return false
}

// Task retrievs the MakeTarget based on the ID. It will
// panic in case the target does not exist in the environment
func (b *Builder) Task(id string) *Task {
	for _, t := range b.tasks {
		if t.ID == id {
			return t
		}
	}
	panic(fmt.Sprintf("Target %s not found in environment", id))
}

// checkDep checks the dependecies recursively for a target ID and updares the
// resolved targets
func (b *Builder) checkDep(parents []*Task, t *Task) {
	var sb strings.Builder
	if b.Verbose {
		for _, p := range parents {
			sb.WriteString(fmt.Sprintf("%s-->", p.ID))
		}
		sb.WriteString(t.GetID())
		log.Printf(sb.String())
	}

	// First check if we're already resolved
	for _, rt := range b.resolvedTasks {
		if t == rt {
			if b.Verbose {
				sb.WriteString(": Already Resolved")
				log.Printf(sb.String())
			}
			return
		}
	}

	// Create a new copy of parents and append ourself
	parents = append(parents, t)
	for _, c := range t.GetDependencies() {
		// Check for circular dependencies
		for _, p := range parents {
			if c.Signature == p.Signature {
				var sb strings.Builder
				for _, t := range parents {
					sb.WriteString(t.DumpStr(false))
					sb.WriteString("->")
				}
				panic(fmt.Sprintf("Circular dependency detected. [%s] <-> [%s](%s)", c.IDSig(), p.IDSig(), sb.String()))
			}
		}
		b.checkDep(parents, c)
	}

	// All children collected, we're resolved
	b.resolvedTasks = append(b.resolvedTasks, t)
	if b.Verbose {
		sb.WriteString(": Resolved")
		log.Printf(sb.String())
	}
}

// DumpTasks will show all registered tasks in the builder
func (b *Builder) DumpTasks() {
	for _, v := range b.tasks {
		log.Println(v.DumpStr(true))
	}
}

// Message creates an operation that only prints a message
func Message(v utils.Variables, msg string) Operation {
	return func(vars utils.Variables) *Artifact {
		v.Printf(msg)
		return nil
	}
}
