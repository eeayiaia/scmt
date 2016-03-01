package devices

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
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

/* NOTE: this is a blocking function .. */
func (conn *RemoteConnection) RunInShell(query string, sudo bool) string {
	session, err := conn.Connection.NewSession()
	if err != nil {
		fmt.Println("[DeviceRemote] could not open a new session towards ", conn.Device.IpAddress, ": ", err)
		return err.Error()
	}

	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	//	session.Stdin = strings.NewReader(conn.Device.Password)

	var q string
	if sudo {
		q = "echo " + conn.Device.Password + " | sudo -S " + query
	} else {
		q = query
	}

	e := session.Run(q)
	if e != nil {
		fmt.Println(e.Error())
	}

	return stdoutBuf.String()
}
