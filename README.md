## Goal

A network of small neural nets and tools working together to solve generic problems.

Smol must fit within a 2GB memory budget and be executable on CPU. It may do computation offloading, but must use generic compute APIs provided by pytorch or similar libraries. No direct dependency on a spesific piece of hardware or vendor may be taken.

Assume the host system has the installed dependencies
1. python version >= 3.12
2. go - version >= 1.24