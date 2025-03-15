## Lecture 8: IO and interfaces

How to read and write files in Go? One easier way is the "bufio" package. It provides buffered I/O. It reads and writes data in chunks, which is more efficient than reading and writing one byte at a time.

An instance of `bufio.Scanner` have a `Split` method, which takes a function that returns a boolean. This function is used to split the input into tokens. The default split function is `bufio.ScanLines`, which splits the input into lines.

### Interface

### Multithreading and mutual exclusion

In Go, threads are called goroutines. A goroutine is a lightweight thread managed by the Go runtime. Go runtime schedules goroutines on multiple OS threads. The Go runtime manages the goroutines, and the programmer does not need to worry about the threads.

"Go doc" is a tool that shows the documentation of a package.

To see the documentation of the `sync` package, run `go doc sync`. The `sync` package provides synchronization primitives.
