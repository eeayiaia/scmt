package devices

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
)

type RemoteConnection struct {
	Device     *Slave
	Connection *ssh.Client
}

/*
	Remote execution "service"
*/

func NewRemoteConnection(device *Slave) (*RemoteConnection, error) {
	sshConfig := &ssh.ClientConfig{
		User: device.UserName,
		Auth: []ssh.AuthMethod{
			ssh.Password(device.Password),
		},
	}

	connection, err := ssh.Dial("tcp", device.IpAddress, sshConfig)
	rc := &RemoteConnection{
		Device:     device,
		Connection: connection,
	}

	return rc, err
}

func (conn *RemoteConnection) RunInShell(query string, sudo bool) string {
	session, err := conn.Connection.NewSession()
	if err != nil {
		Log.Error("could not open a new session towards ", conn.Device.IpAddress, ": ", err)
		return err.Error()
	}

	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf

	var q string = query
	if sudo {
		q = fmt.Sprintf("echo %s | sudo -S %s", conn.Device.Password, query)
	}

	e := session.Run(q)
	if e != nil {
		Log.Error(" could not run command: ", e)
		return e.Error()
	}

	return stdoutBuf.String()
}

/* NOTE: entirely experimental at this moment .. */
func (conn *RemoteConnection) RunInShellAsync(query string, sudo bool) (chan string, error) {
	ch := make(chan string)

	session, err := conn.Connection.NewSession()
	if err != nil {
		Log.Error("could not open a new session towards ", conn.Device.IpAddress, ": ", err)
		return nil, err
	}

	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		Log.Error("could not open remote stdout: ", err)
		return nil, err
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		Log.Error("[DeviceRemote] could not open remote stdin: ", err)
		return nil, err
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()

			trimmedLine := strings.Trim(line, "\n ")

			ch <- trimmedLine
		}
	}()

	var q string
	if sudo {
		q = "echo " + conn.Device.Password + " | sudo -S " + query
	} else {
		q = query
	}

	e := session.Shell()
	if e != nil {
		Log.Error(e.Error())
		return nil, e
	}

	stdin.Write([]byte(q + "\n"))
	stdin.Write([]byte("exit\n"))

	defer session.Wait()

	return ch, nil
}

func (conn *RemoteConnection) RunScript(scriptpath string) (chan string, error) {
	ch := make(chan string)

	go func() {
		lines, err := readScript(scriptpath)
		if err != nil {
			Log.Error("could not open ", scriptpath, ": ", err)
			return
		}

		session, err := conn.Connection.NewSession()
		if err != nil {
			Log.Error("could not open a new session towards ", conn.Device.IpAddress, ": ", err)
			return
		}

		stdout, err := session.StdoutPipe()
		if err != nil {
			Log.Error("when creating stdout-pipe: ", err)
			return
		}

		go func() {
			var read bool = false

			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				line := scanner.Text()
				trimmedLine := strings.Trim(line, "\n ")

				// Since some dists have MOTDs on shells ...
				if strings.Compare(trimmedLine, "BEGIN") == 0 {
					read = true
					continue
				}

				if read {
					ch <- trimmedLine
				}
			}
		}()

		stdin, err := session.StdinPipe()
		if err != nil {
			Log.Error("when creating stdin-pipe: ", err)
			return
		}

		session.Shell()

		// Set start-caret
		stdin.Write([]byte("echo BEGIN\n"))

		// Enable sudo elevation for later entire session (this is a uglyhack in case
		// elevation is needed)
		stdin.Write([]byte(fmt.Sprintf("echo %s | sudo -S echo boo >/dev/null\n", conn.Device.Password)))
		stdin.Write([]byte("while true; do sudo echo boo >/dev/null && sleep 10; done &")) // The '&' at the end creates a job

		for _, line := range lines {
			stdin.Write([]byte(line + "\n"))
		}

		// Kill all jobs (if any) and exit
		stdin.Write([]byte("kill $(jobs -p) && exit\n"))

		session.Wait()
		session.Close()

		close(ch)
	}()

	return ch, nil
}

func readScript(scriptpath string) ([]string, error) {
	file, err := os.Open(scriptpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}
