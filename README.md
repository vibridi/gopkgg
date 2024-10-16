# gopkgg

gopkgg (Go package graph) prints the directed graph of the dependencies between packages of a Go module.

This package showcases one practical application of [`nulab/autog`](https://github.com/nulab/autog) (also mainly written by me).

## Install

Install the binary simply with

    $ go install github.com/vibridi/gopkgg@latest

## Usage

Run the command and supply the path of the Go module you wish to analyze. For example:

    gopkgg /Users/me/nulab/autog

The command creates a file named `depgraph.svg` in the same directory as where the executable was run. 
You can open this file with your favorite browser.  

![autog_graph](https://raw.githubusercontent.com/vibridi/gopkgg/refs/heads/main/example_autog_dep_graph.svg)][autog_dep_graph]


## Authors

* **[Gabriele V.](https://github.com/vibridi/)** - *Main contributor*
* Currently, there are no other contributors

## License

This project is licensed under the MIT License. For detailed licensing information, refer to the [LICENSE](LICENSE) file included in the repository.
