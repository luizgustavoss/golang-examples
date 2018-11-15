# Concurrent File Process

To execute this program you need to create a directory called ***poc*** in the temp directory of your system. In linux it would be:

```
/tmp/poc
```

The program creates a random number of files (max 50) in the directory, writing a random number in each file (max value 3).
After that the program process the files in batches of 5, delegating each file to a goroutine, which will read the file content, and make the goroutine sleep for that amount of time (seconds). At the end of the process of the file the goroutine writes a message to a channel, and at the end of thr program all the messages are shown.