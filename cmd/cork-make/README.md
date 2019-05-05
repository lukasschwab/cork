# cork-make

`cork-make` is a file event handler CLI akin to [`watchman-make`](https://facebook.github.io/watchman/docs/watchman-make.html). Some key differences:

+ `cork-make` has no special treatment for `make` targets.
+ `cork-make` allows the specification of multiple script targets (resolves [watchman:#688](https://github.com/facebook/watchman/issues/688)).

## Getting started

In this subdirectory, run:

```
$ go install
```

## Usage

`cork-make` maps at least one pattern to each command.

```bash
$ cork-make -p [patterns...] -r [command]
$ # Basic example: one pattern group to one command.
$ cork-make -p README.md -r "cat README.md"
$ # Multiple pattern groups to different commands.
$ cork-make -p *.go -r "go build" \
            -p README.md -r "cat README.md"
$ # Multiple pattern groups to one command.
$ cork-make -p *.go *.md -r date
```
