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

import "fmt"
import zmq "github.com/pebbe/zmq4"
import "os"
import "os/signal"
import "syscall"
import "time"
import . "github.com/formwork-io/greenline/internal"

const (
	MsgEvent = 1 << iota
	BinEvent = 1 << iota
)

func main() {
	info := "greenline: notoriously unreliable\n" +
		"https://github.com/formwork-io/greenline\n" +
		"This is free software with ABSOLUTELY NO WARRANTY."
	fmt.Printf("%s\n--\n", info)
	var rails []rail
	if len(os.Args) == 2 {
		var err error
		rails, err = ReadConfigFile(os.Args[1])
		if err != nil {
			die(err.Error())
		}
	} else {
		var err error
		rails, err = ReadEnvironment()
		if err != nil {
			die(err.Error())
		}
	}
	pprint("configuring %d rails", len(rails))

	socket_pairs := make(map[*zmq.Socket]*zmq.Socket)
	socket_names := make(map[*zmq.Socket]string)
	poller := zmq.NewPoller()
	for _, rail := range rails {
		pprint("starting rail %s as %s", rail.Name, rail.Pattern)

		var ingress *zmq.Socket
		var egress *zmq.Socket
		switch rail.Pattern {
		case "pubsub":
			ingress, egress = railToPubSub(&rail, poller)
		case "reqrep":
			ingress, egress = railToRouterDealer(&rail, poller)
		default:
			die("The pattern %s is not valid.", rail.Pattern)
		}

		socket_pairs[ingress] = egress
		socket_names[ingress] = fmt.Sprintf("%s (ingress)", rail.Name)

		socket_pairs[egress] = ingress
		socket_names[egress] = fmt.Sprintf("%s (egress)", rail.Name)

		defer ingress.Close()
		defer egress.Close()
	}

	pprint("greenline alive")
	exitchan := make(chan os.Signal, 0)
	signal.Notify(exitchan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-exitchan
		out("received %s signal, exiting.\n", sig.String())
		os.Exit(0)
	}()

	reloadchan := make(chan int)
	go reloader(reloadchan)
	readychan := make(chan bool)
	pollchan := make(chan bool)
	go func() {
		for {
			sockets, err := poller.Poll(-1)
			if err != nil {
				// get EINTR while polling?
				if IsEINTR(err) {
					// Continue polling; EINTR is normal for us...
					continue
					// ... with our use of signals.
				}

				// otherwise shutdown
				readychan <- false
				break
			}
			if len(sockets) != 0 {
				readychan <- true
			}
			// wait to poll once msgs processed
			<-pollchan
		}
	}()
	pprint("greenline ready")
	for {
		select {
		case reloadOp := <-reloadchan:
			if reloadOp&BinReload == BinReload {
				pprint("new binary available, restarting greenline")
				for key, value := range socket_pairs {
					key.Close()
					value.Close()
				}
				zmq.Term()
				// exec or die
				restart()
			} else if reloadOp&ConfigReload == ConfigReload {
				pprint("new configuration available, restarting greenline")
				for key, value := range socket_pairs {
					key.Close()
					value.Close()
				}
				zmq.Term()
				// exec or die
				restart()
			}
		case ready := <-readychan:
			if !ready {
				die("ready set fail")
			}
			// ready set go
			sockets, err := poller.Poll(-1)
			if err != nil {
				die("poll returned err: %s", err.Error())
			}

			for _, polled := range sockets {
				socket := polled.Socket
				paired_socket := socket_pairs[socket]
				name := socket_names[socket]

				pprint("processing message for %s", name)
				for {
					msg, err := socket.Recv(0)
					if err != nil {
						die("failed on receive: %s", err.Error())
					}
					more, err := socket.GetRcvmore()
					if err != nil {
						die("failed on receive more: %s", err.Error())
					}
					if more {
						paired_socket.Send(msg, zmq.SNDMORE)
					} else {
						paired_socket.Send(msg, 0)
						break
					}
				}
			}
			pollchan <- true
		}
	}
}

func railToPubSub(rail *rail, poller *zmq.Poller) (ingress *zmq.Socket, egress *zmq.Socket) {
	// CREATE EACH SOCKET...
	ingress = newSocket(zmq.SUB)
	ingress.SetSubscribe("")

	egress = newSocket(zmq.PUB)
	egress.SetLinger(1)

	// ... AND BIND
	bind(ingress, "tcp", "0.0.0.0", rail.Ingress)
	bind(egress, "tcp", "0.0.0.0", rail.Egress)

	poller.Add(ingress, zmq.POLLIN)
	return
}

func railToRouterDealer(rail *rail, poller *zmq.Poller) (ingress *zmq.Socket, egress *zmq.Socket) {
	// CREATE EACH SOCKET...
	ingress = newSocket(zmq.ROUTER)

	egress = newSocket(zmq.DEALER)
	egress.SetLinger(1)

	// ... AND BIND
	bind(ingress, "tcp", "0.0.0.0", rail.Ingress)
	bind(egress, "tcp", "0.0.0.0", rail.Egress)

	poller.Add(ingress, zmq.POLLIN)
	poller.Add(egress, zmq.POLLIN)
	return
}

func newSocket(ztype zmq.Type) (socket *zmq.Socket) {
	socket, err := zmq.NewSocket(ztype)
	if err != nil {
		die("failed creating socket type %d: %s", ztype, err.Error())
	}
	return
}

func bind(socket *zmq.Socket, transport string, address string, port int) {
	endpoint := fmt.Sprintf("%s://%s:%d", transport, address, port)
	out("Binding socket %d...", port)
	err := socket.Bind(endpoint)
	if err != nil {
		die("failed binding %s: %s", endpoint, err.Error())
	}
	fmt.Println("done.")
}

func makeMsg(msg string, args ...interface{}) string {
	const layout = "%d%02d%02d-%02d-%02d-%02d greenline[%d]: %s"
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	hour := now.Hour()
	minute := now.Minute()
	seconds := now.Second()
	pid := os.Getpid()
	arg := fmt.Sprintf(msg, args...)
	ret := fmt.Sprintf(layout, year, month, day, hour, minute, seconds, pid, arg)
	return ret
}

func pprint(msg string, args ...interface{}) {
	msg = makeMsg(msg, args...)
	fmt.Fprintf(os.Stdout, msg+"\n")
}

func out(msg string, args ...interface{}) {
	msg = makeMsg(msg, args...)
	fmt.Fprintf(os.Stdout, msg)
	os.Stdout.Sync()
}

func die(msg string, args ...interface{}) {
	msg = makeMsg(msg, args...)
	fmt.Fprintf(os.Stderr, msg+"\n")
	os.Exit(1)
}

// vim: ts=4 noexpandtab
