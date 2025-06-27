package ssh

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"github.com/povsister/scp"
	"golang.org/x/crypto/ssh"
	"os"
	"path/filepath"
)

func newSSHClient(host string, port int32, username string, password string, key string, keyPassphrase string) (*ssh.Client, error) {
	if host == "" {
		host = "localhost"
	}
	if port == 0 {
		port = 22
	}
	var authMethod ssh.AuthMethod
	if key != "" {
		var signer ssh.Signer
		var err error

		if keyPassphrase != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(key), []byte(keyPassphrase))
		} else {
			signer, err = ssh.ParsePrivateKey([]byte(key))
		}

		if err != nil {
			return nil, err
		}
		authMethod = ssh.PublicKeys(signer)
	} else {
		authMethod = ssh.Password(password)
	}

	return ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
}

func execSSHCommand(sshCli *ssh.Client, command string) (string, string, error) {
	session, err := sshCli.NewSession()
	if err != nil {
		return "", "", err
	}
	defer session.Close()
	stdoutBuf := bytes.NewBuffer(nil)
	session.Stdout = stdoutBuf
	stderrBuf := bytes.NewBuffer(nil)
	session.Stderr = stderrBuf
	err = session.Run(command)
	if err != nil {
		return stdoutBuf.String(), stderrBuf.String(), errors.Wrap(err, "failed to execute ssh command")
	}
	return stdoutBuf.String(), stderrBuf.String(), nil
}

func writeFile(sshCli *ssh.Client, useSCP bool, path string, data []byte) error {
	if useSCP {
		return writeFileWithSCP(sshCli, path, data)
	}
	return writeFileWithSFTP(sshCli, path, data)
}

func writeFileWithSCP(sshCli *ssh.Client, path string, data []byte) error {
	scpClient, err := scp.NewClientFromExistingSSH(sshCli, &scp.ClientOption{})
	if err != nil {
		return errors.Wrap(err, "failed to create scp client")
	}
	defer scpClient.Close()

	reader := bytes.NewReader(data)
	err = scpClient.CopyToRemote(reader, path, &scp.FileTransferOption{})
	if err != nil {
		return errors.Wrap(err, "failed to write to remote file")
	}
	return nil
}

func writeFileWithSFTP(sshCli *ssh.Client, path string, data []byte) error {
	sftpCli, err := sftp.NewClient(sshCli)
	if err != nil {
		return errors.Wrap(err, "failed to create sftp client")
	}
	defer sftpCli.Close()

	if err := sftpCli.MkdirAll(filepath.Dir(path)); err != nil {
		return errors.Wrap(err, "failed to create remote directory")
	}

	file, err := sftpCli.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return errors.Wrap(err, "failed to open remote file")
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return errors.Wrap(err, "failed to write to remote file")
	}
	return nil
}
