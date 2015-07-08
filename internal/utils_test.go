package internal

import "syscall"
import "testing"
import . "github.com/smartystreets/goconvey/convey"
import zmq "github.com/pebbe/zmq4"

func TestIsEINTR(t *testing.T) {

    Convey("Returns true when given", t, func() {

        Convey("syscall.EINTR", func() {
            So(IsEINTR(syscall.EINTR), ShouldEqual, true)
        })

        Convey("zmq.Errno", func() {
            So(IsEINTR(zmq.Errno(syscall.EINTR)), ShouldEqual, true)
        })

    })
}

