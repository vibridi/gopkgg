package main

import (
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/nulab/autog"
	"github.com/nulab/autog/graph"
	"github.com/vibridi/graphify"
	"golang.org/x/mod/modfile"
)

func main() {
	dir := targetDir()
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	gomod, err := os.ReadFile(dir + "go.mod")
	if err != nil {
		fmt.Println("failed to parse go.mod:", err.Error())
		os.Exit(1)
	}

	modpath := modfile.ModulePath(gomod)

	var edges [][]string

	err = filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, fs.SkipDir) {
				return nil
			}
			return err
		}
		if !entry.IsDir() {
			return nil
		}
		if strings.HasPrefix(entry.Name(), ".") || strings.HasPrefix(entry.Name(), "_") || path == "testfiles" {
			return fs.SkipDir
		}

		fset := token.NewFileSet()
		pkgs, err := parser.ParseDir(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			fmt.Println("failed to parse files:", err.Error())
			os.Exit(1)
		}

		pkgname := strings.TrimPrefix(path, dir)
		if pkgname == "" || pkgname == modpath {
			s := strings.Split(modpath, "/")
			pkgname = s[len(s)-1]
		}

		for _, pkg := range pkgs {
			for _, f := range pkg.Files {
				for _, i := range f.Imports {
					path := strings.Trim(i.Path.Value, `"`)
					if !strings.HasPrefix(path, modpath) {
						continue
					}

					dep := strings.TrimPrefix(path, modpath+"/")
					edges = append(edges, []string{dep, pkgname})

				}
			}
		}
		return nil
	})

	m := map[[2]string]struct{}{}
	for _, e := range edges {
		m[[2]string{e[0], e[1]}] = struct{}{}
	}

	edges = nil
	for k := range m {
		edges = append(edges, k[:])
	}

	sizes := map[string]graph.Size{}
	for _, e := range edges {
		for _, f := range e {
			if _, ok := sizes[f]; !ok {
				sizes[f] = graph.Size{W: 20.0 * float64(len(f)), H: 40.0}
			}
		}
	}

	layout := autog.Layout(
		graph.EdgeSlice(edges),
		autog.WithNodeSize(sizes),
		autog.WithLayerSpacing(200),
		autog.WithPositioning(autog.PositioningSinkColoring),
		autog.WithEdgeRouting(autog.EdgeRoutingPolyline),
	)

	f, err := os.Create("depgraph.svg")
	if err != nil {
		fmt.Println("failed to create output file:", err.Error())
		os.Exit(1)
	}

	graphify.DrawSVG(
		layout,
		f,
		graphify.WithCanvasPadding(40),
	)
}

func targetDir() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("failed to get working directory:", err.Error())
		os.Exit(1)
	}
	return dir
}
