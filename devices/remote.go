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
func (conn *RemoteConnection) RunInShell(query string) string {
	session, err := conn.Connection.NewSession()
	if err != nil {
		fmt.Println("[DeviceRemote] could not open a new session towards ", conn.Device.IpAddress, ": ", err)
		return err.Error()
	}

	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf

	session.Run(query)

	return stdoutBuf.String()
}
