## Download Manager in go, using concurrency and parallel processing

### How it works
The whole file is divided into 10 chunks for now, which are going to be downloaded concurrently and parallel on multiple processors. If your system has only one core, then it will run on that single core, otherwise it runs on n-1 cores to download the complete file.

### How to compile binary
Binary can be built by 
```
go build main.go
```
The typical usage of the binary is
```
vdlm/vdll <url> <filename>
```
The binary also supports bug reporting and reviewing, for which the typical use case is:-
```
vdll/vdlm <report/review> review
```
for example:-
```
vdlm report absolute file path is an annoying bug, fix it soon, you fuck!
```
