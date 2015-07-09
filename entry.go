// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
//
// See http://formwork-io.github.io/ for more.

package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

// Return true if we were executed using whatever binary was found in the path,
// false otherwise.
func StartedUsingEnvironment() bool {
	arg0 := os.Args[0]
	dir, _ := filepath.Split(arg0)
	if dir == "" {
		return true
	}
	return false
}

// Return true if
func StartedUsingPath() bool {
	arg0 := os.Args[0]
	dir, _ := filepath.Split(arg0)
	if dir != "" {
		return true
	}
	return false
}

func Arg0Dir() (path string, err error) {
	arg0 := os.Args[0]
	if StartedUsingPath() {
		path, _ = filepath.Split(arg0)
		path, err = filepath.Abs(path)
		if err != nil {
			return
		}
	} else {
		path, err = exec.LookPath(arg0)
		if err != nil {
			return
		}
		path, _ = filepath.Split(path)
	}
	return
}

func Arg0Base() (command string) {
	arg0 := os.Args[0]
	_, command = filepath.Split(arg0)
	return
}

func CfgPath() (path string, err error) {
	if len(os.Args) == 1 {
		return "", nil
	}
	cfgfile := os.Args[1]
	path, err = filepath.Abs(cfgfile)
	if err != nil {
		return "", err
	}
	return path, nil
}
