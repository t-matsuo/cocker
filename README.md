# Cocker

Cocker is pre processer for Dockerfile.

[![Go Report Card](https://goreportcard.com/badge/github.com/t-matsuo/cocker)](https://goreportcard.com/report/github.com/t-matsuo/cocker)

It provides these features.

* merge RUN command 
* split RUN command
* include another Dockerfile

## Usage

Merge (use -m option)  
Keyword `#break` breaks merging RUN

```
$ cat Dockerfile
FROM centos:7
RUN echo 1
RUN echo 2
RUN echo 3
#break
RUN echo 4
```
```
$ cocker -m Dockerfile
FROM centos:7
RUN echo 1 && \
    echo 2 && \
    echo 3
RUN echo 4
```

Split (use -s option)

```
$ cat Dockerfile
FROM centos:7
RUN echo 1 && \
    echo 2 && \
    echo 3
```
```
$ cocker -s Dockerfile
FROM centos:7
RUN echo 1
RUN echo 2
RUN echo 3
```

Include (use -i option) 

Note : It includes files recursively, but it cannot detect loop. You can use "ifdef" or "ifndef" option to switch condition.

```
$ cat Dockerfile
FROM centos:7
RUN echo 1
RUN echo 2
#include Dockerfile.inc1
#include Dockerfile.inc2 ifdef MY_ENV1
#include Dockerfile.inc3 ifdef MY_ENV2
#include Dockerfile.inc4 ifndef MY_ENV1
```
```
$ cat Dockerfile.inc 
RUN echo a
RUN echo b
$ cat Dockerfile.inc2
RUN echo c
RUN echo d
$ cat Dockerfile.inc3
RUN echo e
RUN echo f
$ cat Dockerfile.inc4
RUN echo g
RUN echo h
```
```
$ export MY_ENV2=""
$ cocker -i Dockerfile
FROM centos:7
RUN echo 1
RUN echo 2
RUN echo a
RUN echo b
RUN echo e
RUN echo f
RUN echo g
RUN echo h
```

Condition (use -c option)

```
$ export MY_ENV=""
$ cat Dockerfile
FROM centos:7
RUN echo 1
RUN echo 2
#ifndef MY_ENV
RUN echo 3
RUN echo 4
#endif
#ifdef MY_ENV
RUN echo 5
RUN echo 6
#endif
```
```
$ cocker -c Dockerfile
FROM centos:7
RUN echo 1
RUN echo 2
RUN echo 5
RUN echo 6
```

## Help

```
$ cocker -h

Cocker is pre processor for Dockerfile.

Usage:
  $ cocker [options...] filename
  $ cat Dockerfile | cocker [options...]

Options:
   -m --merge   : Merge RUN mode (cannot use -s option)
   -s --split   : Split RUN mode (cannot use -m option)
   -i --include : Include another Dockerfile
   -d --debug   : Print debug messages
   -c --condition : Handle ifdef/ifndef condition only
   --version    : Show version number
   -h --help    : Show help

Environment Variables:
   CC_DEBUG=true : Print debug messages (=--debug option)

```
