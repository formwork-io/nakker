/*
Permission is hereby granted, free of charge, to any person
obtaining a copy of this software and associated documentation
files (the "Software"), to deal in the Software without
restriction, including without limitation the rights to use,
copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following
conditions:
The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

See http://formwork-io.github.io/ for more.
*/

package main

import "errors"
import "fmt"
import toml "github.com/BurntSushi/toml"
import "io/ioutil"
import "os"
import "strconv"
import "strings"

// Rail ...
type Rail struct {
	Name     string
	Protocol string
	Ingress  int
	Egress   int
}

// Rails ...
type Rails struct {
	rail []Rail
}

// ReadConfigFile ...
func ReadConfigFile(path string) ([]Rail, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("configuration file error: %s", err.Error())
	}

	var rails []Rail
	_, err = toml.Decode(string(data), &rails)
	if err != nil {
		return nil, fmt.Errorf("configuration file error: %s", err.Error())
	}

	_, err = validateRails(rails)
	if err != nil {
		return nil, fmt.Errorf("configuration file error: " + err.Error())
	}

	return rails, nil
}

// ReadEnvironment ...
func ReadEnvironment() ([]Rail, error) {
	nameTemplate := "GL_RAIL_%d_NAME"
	protocolTemplate := "GL_RAIL_%d_PROTOCOL"
	ingressTemplate := "GL_RAIL_%d_INGRESS"
	egressTemplate := "GL_RAIL_%d_EGRESS"

	var rails []Rail
	index := 0
	for {
		name, err := getenv(fmt.Sprintf(nameTemplate, index))
		if err != nil {
			if index != 0 {
				break
			}
			return nil, err
		}
		if name == "" {
			break
		}

		protocol, err := getenv(fmt.Sprintf(protocolTemplate, index))
		if err != nil {
			return nil, err
		}
		ingressStr, err := getenv(fmt.Sprintf(ingressTemplate, index))
		if err != nil {
			return nil, err
		}
		egressStr, err := getenv(fmt.Sprintf(egressTemplate, index))
		if err != nil {
			return nil, err
		}
		ingress, err := asPort(ingressStr)
		if err != nil {
			return nil, err
		}
		egress, err := asPort(egressStr)
		if err != nil {
			return nil, err
		}

		rails = append(rails, Rail{
			Name:     name,
			Protocol: protocol,
			Ingress:  ingress,
			Egress:   egress,
		})

		index++
	}

	_, err := validateRails(rails)
	if err != nil {
		return nil, fmt.Errorf("configuration file error: " + err.Error())
	}

	return rails, nil
}

func validProtocol(protocol string) bool {
	switch protocol {
	case "broadcast":
		return true
	case "request":
		return true
	}
	return false
}

func validateRails(rails []Rail) (*Rail, error) {
	for i, rail := range rails {
		protocol := strings.ToLower(rail.Protocol)
		if !validProtocol(protocol) {
			msg := fmt.Sprintf("%s: unsupported rail type", rail.Protocol)
			return &rail, errors.New(msg)
		}
		rails[i].Protocol = protocol
	}

	return nil, nil
}

func getenv(env string) (string, error) {
	_env := os.Getenv(env)
	if len(_env) == 0 {
		return "", fmt.Errorf("no %s is set", env)
	}
	return _env, nil
}

func asPort(env string) (int, error) {
	port, err := strconv.Atoi(env)
	if err != nil {
		die("invalid port: %s", env)
		return -1, fmt.Errorf("invalid port: %v - %s", env, err.Error())
	} else if port < 1 || port > 65535 {
		die("invalid port: %s", env)
		return -1, fmt.Errorf("invalid port: %v - %s", env, err.Error())
	}
	return port, nil
}
