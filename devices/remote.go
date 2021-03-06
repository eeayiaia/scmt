package devices

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"

	log "github.com/Sirupsen/logrus"
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

	connection, err := ssh.Dial("tcp", (device.IPAddress + ":" + device.Port), sshConfig)
	rc := &RemoteConnection{
		Device:     device,
		Connection: connection,
	}

	return rc, err
}

func (conn *RemoteConnection) RunInShell(query string, sudo bool) string {
	session, err := conn.Connection.NewSession()
	if err != nil {
		Log.WithFields(log.Fields{
			"IP":    conn.Device.IPAddress,
			"MAC":   conn.Device.HardwareAddress,
			"error": err,
		}).Error("could not open session")

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
		Log.WithFields(log.Fields{
			"IP":    conn.Device.IPAddress,
			"MAC":   conn.Device.HardwareAddress,
			"error": e,
		}).Error("could not run command")

		return e.Error()
	}

	return stdoutBuf.String()
}

/*
	Copies a folder to a remote NewRemoteConnection
    Example: CopyFolder("/home/xxxx/SuperK/", "/tmp/") will copy SuperK to /tmp/SuperK
*/
func (conn *RemoteConnection) CopyFolder(folderpath string, destination string) error {
	session, err := conn.Connection.NewSession()
	if err != nil {
		Log.Error("could not open a new session", err)
		return err
	}
	defer session.Close()

	f, err := os.Open(folderpath)
	if err != nil {
		return err
	}

	defer f.Close()

	if err != nil {
		return err
	}

	if folderpath[len(folderpath)-1] == '/' {
		folderpath = folderpath[:len(folderpath)-1]
	}

	path, foldername := path.Split(folderpath)
	fmt.Println(path)

	tmpPath := fmt.Sprintf("/tmp/%s.tar.gz", foldername)
	cmd := exec.Command("tar", "-C", path, "-zcf", tmpPath, foldername)

	Log.WithFields(log.Fields{
		"target":  conn.Device.IPAddress,
		"command": cmd,
	}).Info("Running command")

	if err := cmd.Run(); err != nil {
		Log.WithFields(log.Fields{
			"target":  conn.Device.IPAddress,
			"command": cmd,
			"error":   err,
		}).Fatal("Failed to run command")
	}

	err = conn.CopyFile(tmpPath, tmpPath)

	if err != nil {
		Log.WithFields(log.Fields{
			"target":  conn.Device.IPAddress,
			"command": cmd,
			"error":   err,
		}).Warn("Failed to copy file")
		return err
	}

	Log.WithFields(log.Fields{
		"target":      conn.Device.IPAddress,
		"folderpath":  folderpath,
		"destination": destination,
	}).Info("copied file")

	shellCMD := fmt.Sprintf("/bin/tar -xf %s -C %s", tmpPath, destination)

	/*If possible run without sudo*/
	writeCheck := "if [ -w \"" + destination + "\" ]; then echo \"WRITABLE\"; fi"
	outp := conn.RunInShell(writeCheck, false)

	if strings.Contains(outp, "WRITABLE") {
		conn.RunInShell(shellCMD, false)
	} else {
		conn.RunInShell(shellCMD, true)
	}

	shellCMD = "rm " + tmpPath
	conn.RunInShell(shellCMD, true)

	return nil
}

func (conn *RemoteConnection) CopyFile(filepath string, destination string) error {
	session, err := conn.Connection.NewSession()
	if err != nil {
		Log.Error("could not open a new session", err)
		return err
	}
	defer session.Close()

	f, err := os.Open(filepath)
	if err != nil {
		Log.Error("could not open file", err)
		return err
	}

	defer f.Close()

	s, err := f.Stat()
	if err != nil {
		Log.Error("could not stat file", err)
		return err
	}

	Log.WithFields(log.Fields{
		"filepath":    filepath,
		"destination": destination,
	}).Info("Copying file to device")

	go func() {
		var stdoutBuf bytes.Buffer
		session.Stdout = &stdoutBuf

		w, e := session.StdinPipe()
		if e != nil {
			Log.WithFields(log.Fields{
				"filepath":    filepath,
				"destination": destination,
			}).Error("could not get stdin pipe!")
		}

		fileName := path.Base(filepath)
		mode := s.Mode().Perm()
		size := s.Size()

		fmt.Fprintf(w, "C%#o %d %s\n", mode, size, fileName)
		io.Copy(w, f)
		fmt.Fprint(w, "\x00 \r\n")

		Log.WithFields(log.Fields{
			"target":      conn.Device.IPAddress,
			"filepath":    filepath,
			"destination": destination,
			"size":        size,
		}).Info("copied file")
	}()

	cmd := fmt.Sprintf("/usr/bin/scp -qtr %s", destination)
	session.Run(cmd)

	return nil
}

/* NOTE: entirely experimental at this moment .. */
func (conn *RemoteConnection) RunInShellAsync(query string, sudo bool) (chan string, error) {
	ch := make(chan string)

	session, err := conn.Connection.NewSession()
	if err != nil {
		Log.Error("could not open a new session towards ", conn.Device.IPAddress, ": ", err)
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

/*
   Note: This function no longer transfer given file, it only runs given script.
*/
func (conn *RemoteConnection) RunScript(scriptpath string, env map[string]string) (chan string, error) {
	ch := make(chan string)
	exit := false
	var wg sync.WaitGroup

	go func() {
		session, err := conn.Connection.NewSession()
		if err != nil {
			Log.WithFields(log.Fields{
				"IP":  conn.Device.IPAddress,
				"MAC": conn.Device.HardwareAddress,
				"ERR": err,
			}).Error("could not open session")
			return
		}

		stdout, err := session.StdoutPipe()
		if err != nil {
			Log.Error("when creating stdout-pipe: ", err)
			return
		}

		stderr, err := session.StderrPipe()

		// Read stdout to channel
		go func() {
			wg.Add(1)
			defer wg.Done()

			var read bool = false
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() && !exit {
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

		// Read stderr to channel
		go func() {
			wg.Add(1)
			defer wg.Done()

			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() && !exit {
				line := scanner.Text()
				trimmedLine := strings.Trim(line, "\n ")

				ch <- trimmedLine
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

		if env != nil {
			for k, v := range env {
				Log.Debug("export " + k + "=" + v)
				stdin.Write([]byte("export " + k + "=" + v + "\n"))
			}
		}

		sudo := fmt.Sprintf("echo %s | sudo -S ", conn.Device.Password)

		Log.WithFields(log.Fields{
			"scriptpath": scriptpath,
		}).Debug("chmod")
		stdin.Write([]byte(sudo + "chmod +x " + scriptpath + "\n"))

		Log.WithFields(log.Fields{
			"scriptpath": scriptpath,
		}).Debug("sudo -E bash -C ..")
		stdin.Write([]byte(sudo + "-E bash -C '" + scriptpath + "'\n"))

		stdin.Write([]byte("exit\n"))

		session.Wait()
		session.Close()

		exit = true
		// Do not close until readers are done
		wg.Wait()
		close(ch)
	}()

	return ch, nil
}
