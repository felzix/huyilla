package engine

import (
    "path"
    "runtime"
    "testing"
)


func getContent(t *testing.T) *Content {
    directory := getDirectory(t)

    content, err := LoadContent(directory)
    if err != nil {
        t.Fatal(err)
    }

    return content
}


func getDirectory (t *testing.T) string {
    _, filename, _, ok := runtime.Caller(0)
    if !ok {
        t.Fatal("Failed to discover current directory")
    }

    return path.Dir(filename) + "/../content"
}