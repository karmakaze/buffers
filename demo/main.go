package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/karmakaze/buffers"
)

func main() {
	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	r, w := buffers.NewPipe(1300)

	go write(w)
	read(r)
}

func read(r buffers.PipeReader) {
	alphabet := "abcdefghijklmnopqrstuvwxyz"

	buf := make([]byte, 260)
	received := buffers.NewCircular(520)
	for {
		rnd := rand.Intn(len(buf) + 1)
		n, err := r.Read(buf[0:rnd])
		errorExit(err, "Error on read: %v", err)

		c, err := received.Write(buf[0:n])
		errorExit(err, "Error on write to 'received': %v", err)
		if c != n {
			err = fmt.Errorf("Only wrote %d of %d bytes.", c, n)
			errorExit(err, "Error on write to 'received': %v", err)
		}

		for received.Buffered() >= 26 {
			n, err = received.Read(buf[:26])
			errorExit(err, "Error on read from 'received': %v", err)
			if n != 26 {
				err = fmt.Errorf("Only read %d of %d bytes.", n, 26)
				errorExit(err, "Error on read of 'received': %v", err)
			}

			fmt.Printf(" read %v\n", string(buf[0:n]))

			if string(buf[0:n]) != alphabet {
				err = fmt.Errorf("expecting '%v', got '%v'", alphabet, string(buf[0:n]))
				errorExit(err, "Error: %v", err)
			}
		}
	}
}

func write(w buffers.PipeWriter) {
	alphabet3 := "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"

	var i int
	for {
		c := rand.Intn(26 + 1)
		n, err := w.Write([]byte(alphabet3[i : i+c]))
		errorExit(err, "Error on write: %v", err)
		fmt.Printf("wrote %s\n", alphabet3[i:i+n])
		i += n
		if i >= 26 {
			i -= 26
		}
		time.Sleep(50)
	}
}

func errorExit(err error, message string, args ...interface{}) {
	if err != nil {
		fmt.Printf(message, args)
		os.Exit(1)
	}
}
