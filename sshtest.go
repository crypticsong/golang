// https://blogs.oracle.com/janp/entry/how_the_scp_protocol_works
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"
)

func main() {

	username := ""
	password := ""
	hostname := ""
	port := "22"

	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", hostname+":"+port, clientConfig)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()
	file, err := os.Open("users.xml")
	if err != nil {
		fmt.Println("ERROR in OPENing FILE", err)
	}
	fmt.Println("file is : ", file)
	contentfile, _ := ioutil.ReadAll(file)
	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()

		//fmt.Println("CONTENT OF USERS.XML \n", string(contentfile))
		content := "1234567890000\n"
		fmt.Fprintln(w, "D0755", 0, "testdir") // mkdir

		fmt.Fprintln(w, "C0644", len(content), "testfile12")
		fmt.Fprint(w, string(content))
		fmt.Fprint(w, "\x00") // transfer end with \x00
		//fmt.Fprintln(w, "D0755", 0, "1testdir")
		fmt.Fprintln(w, "C0644", len(contentfile), "testfile22")
		fmt.Fprint(w, string(contentfile))
		fmt.Fprint(w, "\x00")
	}()
	if err := session.Run("/usr/bin/scp -tr ./"); err != nil {
		panic("Failed to run: " + string(err.Error()))
	}
}
