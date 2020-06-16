package api

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
)

type Remote struct {
	ssh *ssh.Client
}

func (r *Remote) Close() {
	r.ssh.Close()
}

func (r *Remote) Run(cmd string) error {
	session, err := r.ssh.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	var outBuf, errBuf bytes.Buffer
	session.Stdout = &outBuf
	session.Stderr = &errBuf
	exe := "source /etc/profile;" + cmd // non-login形式默认不读/etc/profile
	err = session.Run(exe)
	if err != nil {
		log.Printf("%s error :\n%s", cmd, errBuf.String())
		return err
	}
	log.Printf("%s output :\n%s", cmd, outBuf.String())
	return nil
}

func (r *Remote) Start(cmd string) error {
	session, err := r.ssh.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin
	exe := "source /etc/profile;" + cmd // non-login形式默认不读/etc/profile
	err = session.Start(exe)
	if err != nil {
		return err
	}
	err = session.Wait()
	if err != nil {
		return err
	}
	return nil
}

// NewRemote 通过用户名密码创建实例
func NewRemote(host, usr, pwd string) (*Remote, error) {
	auth := ssh.Password(pwd)
	client, err := newSSHClient(host, usr, auth)
	if err != nil {
		return nil, err
	}
	return &Remote{ssh: client}, nil
}

// NewRemoteByDefaultKey 通过默认免密key创建实例
func NewRemoteByDefaultKey(host, usr string) (*Remote, error) {
	var keyPath string
	if usr == "root" {
		keyPath = "/root/.ssh/id_rsa"
	} else {
		keyPath = "/home/" + usr + "/.ssh/id_rsa"
	}
	return NewRemoteByKey(host, usr, keyPath)
}

// NewRemoteByKey 通过指定免密key创建实例
func NewRemoteByKey(host, usr, keyPath string) (*Remote, error) {
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	auth := ssh.PublicKeys(signer)
	client, err := newSSHClient(host, usr, auth)
	if err != nil {
		return nil, err
	}
	return &Remote{ssh: client}, nil
}

// NewRemoteByKeyWithPwd 通过带密码的key创建实例
func NewRemoteByKeyWithPwd(host, usr, keyPath, pwd string) (*Remote, error) {
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKeyWithPassphrase(key, []byte(pwd))
	if err != nil {
		return nil, err
	}
	auth := ssh.PublicKeys(signer)
	client, err := newSSHClient(host, usr, auth)
	if err != nil {
		return nil, err
	}
	return &Remote{ssh: client}, nil
}

func newSSHClient(host, usr string, auth ssh.AuthMethod) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: usr,
		Auth: []ssh.AuthMethod{
			auth,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 接受所有hostkey
	}
	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil, err
	}
	return client, nil
}
