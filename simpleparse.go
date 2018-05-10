// This program needs to perform the following tasks:
// 1. Get `info commandstats` output from command line
// 2. Parse the numbers from the output (file?)
// 3. Compare current values to previous (delta)

package main 

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	// "reflect"
	"strings"
	_"log"
	"errors"
)

func main() {

	// TODO: Have the IP address managed by the program so it can used on different machines

	cmd := exec.Command("redis-cli", "-h", "172.17.0.2", "info", "commandstats")

	// TODO: Output to temp file https://www.devdungeon.com/content/working-files-go#temp_files

	outfile, err := os.Create("./out.txt")
	if err != nil {
		panic(err)
	}
	defer outfile.Close()
	cmd.Stdout = outfile

    err = cmd.Start(); if err != nil {
        panic(err)
    }
    cmd.Wait()

	// fmt.Printf("%s\n", outfile)
	// fmt.Println(reflect.TypeOf(outfile))
	// fmt.Println(reflect.TypeOf(cmd))

	removeLines("./out.txt", 1, 1)
	readLines("./out.txt")

}

func readLines(fn string) (err error) {

    file, err := os.Open(fn)
    defer file.Close()

    if err != nil {
        return err
    }

    // Start reading from the file with a reader.
    reader := bufio.NewReader(file)

    for {
        var buffer bytes.Buffer

        var l []byte
        var isPrefix bool
        for {
            l, isPrefix, err = reader.ReadLine()
            buffer.Write(l)

            // If we've reached the end of the line, stop reading.
            if !isPrefix {
                break
            }

            // If we're just at the EOF, break
            if err != nil {
                break
            }
        }

        if err == io.EOF {
            break
        }

        line := buffer.String()

	    line1 := strings.Replace(line, ",", "=", -1)
		new_line := strings.Split(line1, "=")

		calls, usec, usec_per_call := new_line[1], new_line[3], new_line[5]

		fmt.Println(calls, usec, usec_per_call)
		// fmt.Println("Original line")
		// fmt.Println(line)
		// fmt.Println("Replaced comma w equals")
		// fmt.Println(line1)
		// fmt.Println("Split line")
		// fmt.Println(new_line)

        // Process the line here.
    }

    if err != io.EOF {
        fmt.Printf(" > Failed!: %v\n", err)
    }

    return
}

func removeLines(fn string, start, n int) (err error) {
    if start < 1 {
        return errors.New("invalid request.  line numbers start at 1.")
    }
    if n < 0 {
        return errors.New("invalid request.  negative number to remove.")
    }
    var f *os.File
    if f, err = os.OpenFile(fn, os.O_RDWR, 0); err != nil {
        return
    }
    defer func() {
        if cErr := f.Close(); err == nil {
            err = cErr
        }
    }()
    var b []byte
    if b, err = ioutil.ReadAll(f); err != nil {
        return
    }
    cut, ok := skip(b, start-1)
    if !ok {
        return fmt.Errorf("less than %d lines", start)
    }
    if n == 0 {
        return nil
    }
    tail, ok := skip(cut, n)
    if !ok {
        return fmt.Errorf("less than %d lines after line %d", n, start)
    }
    t := int64(len(b) - len(cut))
    if err = f.Truncate(t); err != nil {
        return
    }
    if len(tail) > 0 {
        _, err = f.WriteAt(tail, t)
    }
    return
}

func skip(b []byte, n int) ([]byte, bool) {
    for ; n > 0; n-- {
        if len(b) == 0 {
            return nil, false
        }
        x := bytes.IndexByte(b, '\n')
        if x < 0 {
            x = len(b)
        } else {
            x++
        }
        b = b[x:]
    }
    return b, true
}