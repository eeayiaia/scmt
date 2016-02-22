package main

import "fmt"
import "golang.org/x/crypto/ssh"
import "bytes"

func main() {
	fmt.Println("Stub!")

	sshConfig := &ssh.ClientConfig{
		User: "selund",
		Auth: []ssh.AuthMethod{
			ssh.Password("galenanka1"),
		},
	}

	connection, err := ssh.Dial("tcp", "129.16.22.6:2222", sshConfig)
	if err != nil {
		fmt.Println("Error when connecting: {}", err)
	}

	session, err := connection.NewSession()
	if err != nil {
		fmt.Println("Error when creating a session: {}", err)
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf

	stdin, err := session.StdinPipe()
	if err != nil {
		fmt.Println("Error when creating stdin-pipe: {}", err)
	}

	//session.Run("hostname -f; pwd; ssh odroid@10.46.0.101")
	session.Shell()

	stdin.Write([]byte("hostname -f\n"))
	stdin.Write([]byte("ls -la\n"))
	stdin.Write([]byte("exit\n"))

	session.Wait()

	//	stdin.Write([]byte("odroid\n"))
	//	stdin.Write([]byte("hostname -f; exit\n"))
	//session.Run("odroid")
	//session.Run("hostname -f")

	fmt.Println("Result: " + stdoutBuf.String())
}
