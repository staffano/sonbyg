# The build system.
The role of the build system is to apply a series of Tasks in order to shape the
contents of a Workspace on the behalf of a client. The build system controller
is called Builder and is responsible for the task registry and for arranging the
execution of tasks in a way to guarantee a consistent and fast update of the
workspace.

## The Workspace.
The client wants to shape a filesystem according to a root task. In order to
that a number of tools, support directories and state files might be needed.
These are of little interest to the client. For this reason we partition the
Workspace into the InstallationDir and the WorkDir, where InstallationDir is
what the client wants and WorkDir is the workshop where the shaping tools do
their jobs.

### Externa Files
Patch files, configurations and other files required by the tasks are located
in the FileDir (FILES) within the workspace. They should be put under
`${FILES}/${NAMESPACE}/$ID` or similar

### Variables
The build variables for a task define the behaviour of the task. They are defined
when the task is created and resolved when the task is executed.
A task is potentially created provided a great number of variables, defining
it's context. It might however be so that only a few of those variables are
actually needed. Therefore the task creation should start by selecting 
the variables needed for build an only record those in the task context.

A task that depends on another task, might also need access to variables in that
task, for instance if a task is dependent on that the bash script executor is
installed, the task might need the path of the bash command. This could 
be _exported_ to the depening task by letting the bash-installer task update
the depending tasks variable list.


## Phases.
The build system undergoes a number of phases.

### Phase 1 - Task declaration
 During the task declaration, the client initiate the creation of new tasks by
 creating the root task that describes the overall intent of the shaping of the
 workspace. This task in turn will declare dependencies it needs to have
 fulfilled in order to consider itself done. This is a recursive operation and
 once the creation of the root task returns to the client, all explicitly
 declared tasks neccessary to shape the workspace shall been recorded in the
 Builder.

#### Phase 1.1 Task Call Signature
A _Task Call Signature_ is the combination of a task id and a set of variables
to call this task with. The map between a task signature and a behaviour must be
completely deterministic and represents a unique operation.

### Phase 2. Deciding task execution order.
Once all tasks are recorded in the Builder, the builder will decide on a scheme
for executing the tasks.

### Phase 3. Task execution
The client decides how many jobs in parallell that should be used by the
builder and then the builder calls each and every task in the order of
dependencies.

#### Phase 3.1 - Cache
Before a task is selected to be executed, the Builder will check if an artifact
for the task is already present. If so, the artifact will be applied to the
workspace instead of executing the task. All tasks return upon their execution,
either an artifact or a Null_Artifact, indicating that no artifact was
generated. The task execution plan is not altered due to the existence of
artifacts, only the execution of the task.

### Phase 4. Handover.
When the task execution is done the workspace is shaped according to the
requirements of the root task. A database is created with the tasks executed and
their artifacts.

## Task correctness.
The task is completely defined by its build variable and the code it is
executing. The code is defined by a version in a code repository and the build
variables is defined by the root task and the operating system environment when
executing the task. These requirements are recorded in the build variables that
are associated with the execution of a task. The role of the code is to define a
one-to-one mapping of the current state of the workspace and the build variables
to the result of the operation. There are also _sync tasks_ that coordinate work, they
don't provide new information to the workspace, but are used as dependency
holders and indicators within the build system. All communication between the tasks
during execution is done via the workspace and via the variables. This provides
us with means to cache result. By associating the signature with an artefact,
when considering a task, we can check if there is an artefact associated with
the signature and then just apply the artefact, without knowing much more about
the task. This requires the task execution order to be determant and that only
tasks update the workspace. An even more robust version would include the
current state of the workspace in the signature of the task.

## Workspace monitoring
When building the workspace, the builder could keep track of all the updates
made to the workspace during the task execution. This could be used to recreate
a single file, or a section of the workspace, but reapplying the artefacts in
dependency order for each of the tasks whose execution modified the file. It
could also be used to enable packaging of in-workspace builds. The recommended
way of building a set of coherent tasks is to create a staging directory where
the package is installed and when the package is built, an artefact (archive) is
created and recorded. This artefact is then applied to the InstallationDir and
the Task is done. By monitoring the workspace, installation into the
InstallationDir directly can be monitored and an artefact from the created files
could after the task has executed be created from the list of files the task
created. If the task installation fails and a roll-back is required, the
previous artefacts will need to be reapplied for the files changed by the
failing task.
