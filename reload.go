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
	"gopkg.in/fsnotify.v1"
	"os"
	"path/filepath"
	"syscall"
)

const (
	ConfigReload = 1 << iota
	BinReload    = 1 << iota
)

func isReloadEvent(event fsnotify.Event) bool {
	if event.Op&fsnotify.Create == fsnotify.Create {
		return true
	}
	if event.Op&fsnotify.Write == fsnotify.Write {
		return true
	}
	return false
}

func reloader(reload chan int) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		die("failed creating watcher (%s)", err.Error())
	}
	defer watcher.Close()

	bindir, err := Arg0Dir()
	if err != nil {
		die("failed getting path to binary (%s)", err.Error())
	}
	pprint("monitoring %s", bindir)
	me := Arg0Base()
	err = watcher.Add(bindir)
	if err != nil {
		die("failed watching binary (%s)", err.Error())
	}
	watchcfg := false
	cfgfile, cfgdir := "", ""
	if len(os.Args) == 2 {
		cfgpath, err := CfgPath()
		if err != nil {
			die("error getting configuration path (%s)", err.Error())
		}
		cfgdir, cfgfile = filepath.Split(cfgpath)
		err = watcher.Add(cfgdir)
		if err != nil {
			die("failed watching configuration (%s)", err.Error())
		}
		cfgdir = cfgdir[:len(cfgdir)-1]
		pprint("monitoring %s", cfgdir)
		watchcfg = true
	}

	for {
		select {
		case event := <-watcher.Events:
			if isReloadEvent(event) {
				dir, base := filepath.Split(event.Name)
				dir = dir[:len(dir)-1]
				if dir == bindir && base == me {
					reload <- BinReload
				} else if watchcfg && dir == cfgdir && base == cfgfile {
					reload <- ConfigReload
				}
			}
		case err := <-watcher.Errors:
			die("failed gettiing events (%s)", err.Error())
		}
	}
}

func restart() (err error) {
	err = syscall.Exec(os.Args[0], os.Args, os.Environ())
	if err != nil {
		die("failed restarting greenline")
	}
	return
}

// vim: ts=4 noexpandtab
