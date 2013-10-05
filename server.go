package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

const (
	RECV_BUF_LEN = 1024
	PORT         = "8008"
)

func main() {
	fmt.Println("Started the logging server")

	listener, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		println("error listening:", err.Error())
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accept:", err.Error())
			return
		}
		//creates a go routine to execute grep in shell
		go grepMyLog(conn)
	}
}

/*
 * invokes execGrep in the shell and returns the result through a buffer
 * @param conn socket through which the server communicates with the client
 */
func grepMyLog(conn net.Conn) {
	recvBuf := make([]byte, RECV_BUF_LEN)
	_, err := conn.Read(recvBuf)

	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

	//convert byte array to a string
	n := bytes.Index(recvBuf, []byte{0})
	s := string(recvBuf[:n])

	//check to see if this is a request from a unit test
	env := "production"
	if strings.HasPrefix(s, "test") {
		s = s[4:len(s)]
		env = "test"
	}
	//send the results back
	var results string
	if strings.EqualFold(env, "production") {
		results = execGrep(s, os.Args[1], os.Args[2])
	} else {
		//generate random logs with rare, frequent and lines appearing in average frequency
		genLogs()
		results = execGrep(s, "test.log", "test.1")
	}
	sendBuf := make([]byte, len(results))
	copy(sendBuf, string(results))
	conn.Write(sendBuf)
	conn.Close()
}

/*
 * executes grep in unix shell
 * @param s           the query string
 * @param logName     name of the log file
 * @param machineName name of the machine
 */
func execGrep(s string, logName string, machineName string) string {
	cmd := exec.Command("grep", s, logName)
	cmdOut, cmdErr := cmd.Output()

	results := ""
	//check if there is any error in our grep
	if cmdErr != nil {
		fmt.Println("ERROR WHILE READING")
		fmt.Println(cmdErr)
	}

	if len(cmdOut) > 0 {
		results = machineName + "\n" + string(cmdOut)
	} else {
		results = "No matching patterns found in " + machineName
	}
	return results
}

/*
 * Invokes the helper function with the strings we need to append to the test.log
 */
func genLogs() {
	rare := "i:only appear once"
	frequent := "frequent: line they call me"
	sometimes := "they:call me sometimes"

	var lines = []string{}

	for i := 0; i < 20; i++ {
		lines = append(lines, frequent)
		if i%4 == 0 {
			lines = append(lines, sometimes)
		}
	}
	lines = append(lines, rare)
	writeLines(lines, "test.log")
}

/*
 * Creates a test.log file and appends the required strings
 * @param lines a slice of lines that we will append to the given file name
 * @param fileName of the test log file we are going to create
 */
func writeLines(lines []string, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}
