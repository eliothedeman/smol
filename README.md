#idea
## Goal

A network of small neural nets and tools working together to solve generic problems.

Smol must fit within a 2GB memory budget and be executable on CPU. It may do computation offloading, but must use generic compute APIs provided by pytorch or similar libraries. No direct dependency on a spesific piece of hardware or vendor may be taken.


Assume the host system has the installed dependencies
1. python version >= 3.12
2. go - version >= 1.24

## Overview
Similar to a CPU, smol is comprised of various logical units that can be combined to execute instructions.

These logic units connected to each other via a message passing interface.

The main process where units are executed is written in go. Units that require access to python, such as to run pytorch models, do so by using the embed package to store their python scripts in the binary, then write the scripts to disk and execute them with subprocesses. Models themselves must be embedded in a similar manor.

The host system must provide all necessary runtime dependencies. If not already in a python virtual environment at execution time, one will be created using the systems installed python.

### Unit Interfaces
```go

type UnitDesc struct {
	Name string
	Proxy Unit
}

type Ctx interface {
	Units() []UnitDesc
	Spawn(name string, f UnitFactory) UnitRef
	Self() UnitRef()
	Subscribe(other Unit)
}

type UnitFactory func() Unit

type UnitRef  interface {
	Name() string
	Send(msg any)
	Stop() 
}

type Unit interface {
	Init(ctx Context)
	Handle(ctx Context, from UnitRef, message any) error
}
```


### RequiredUnits
* Control
	* InstructionExecutor - Decides the next instruction and executes it via tools and neural nets
	* Lifecycle
* Tools
	* Math - Execute simple math expressions. May reference numbers stored in registers via $register_name string substitution.
	* Storage - Persistent storage using basic filesystem CRUD API
	* Registers - ephemeral key value pairs for cooperating between control units and tools
	* Code Execution Sandbox - Execute the given python code in a subprocess and return the stdout and stderr as strings
	* Memory
		* Named memories. Stored as markdown files. Internal links supported using obsidian link formats.
	* OpenAI compatible API server
	* NeuralNet prediction
		* Neural Nets
			* Image classification
			* Object detection

