# Concurrent File Process

By default to execute this program you need to create a directory called ***poc*** in the temp directory of your system, but this can be configured to another directory. In linux it would be:

```
/tmp/poc
```


The program creates a random number (configurable) of files in the directory, writing a random number (configurable) in each file. Then the program processes the files, opening a goroutine which process each file. It also has an implementation to process files in batches of a configurable size.

The process of a file consists in opening and reading a value from it that represents a number of seconds to make the goroutine sleep, than sleep, and after waking up write a message on a channel.

At the end all the messages are read in the main func.

Once you've cloned the repository, enter golang-examples and execute:

```
go run concurrent-file-process/main.go
```
