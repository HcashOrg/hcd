// Copyright (c) 2016 The btcsuite developers
// Copyright (c) 2017 The Decred developers
// Copyright (c) 2018-2020 The Hc developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package rpctest

import (
	"crypto/elliptic"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/james-ray/hcd/wire"

	"github.com/james-ray/hcd/hcutil"
	rpc "github.com/james-ray/hcrpcclient"
)

// nodeConfig contains all the args, and data required to launch a hcd process
// and connect the rpc client to it.
type nodeConfig struct {
	rpcUser    string
	rpcPass    string
	listen     string
	rpcListen  string
	rpcConnect string
	dataDir    string
	logDir     string
	profile    string
	debugLevel string
	extra      []string
	prefix     string

	exe          string
	endpoint     string
	certFile     string
	keyFile      string
	certificates []byte
}

// newConfig returns a newConfig with all default values.
func newConfig(prefix, certFile, keyFile string, extra []string) (*nodeConfig, error) {
	a := &nodeConfig{
		listen:    "127.0.0.1:18555",
		rpcListen: "127.0.0.1:18556",
		rpcUser:   "user",
		rpcPass:   "pass",
		extra:     extra,
		prefix:    prefix,

		exe:      "hcd",
		endpoint: "ws",
		certFile: certFile,
		keyFile:  keyFile,
	}
	if err := a.setDefaults(); err != nil {
		return nil, err
	}
	return a, nil
}

// setDefaults sets the default values of the config. It also creates the
// temporary data, and log directories which must be cleaned up with a call to
// cleanup().
func (n *nodeConfig) setDefaults() error {
	n.dataDir = filepath.Join(n.prefix, "data")
	n.logDir = filepath.Join(n.prefix, "logs")
	cert, err := ioutil.ReadFile(n.certFile)
	if err != nil {
		return err
	}
	n.certificates = cert
	return nil
}

// arguments returns an array of arguments that be used to launch the hcd
// process.
func (n *nodeConfig) arguments() []string {
	args := []string{}
	// --simnet
	args = append(args, fmt.Sprintf("--%s", strings.ToLower(wire.SimNet.String())))
	if n.rpcUser != "" {
		// --rpcuser
		args = append(args, fmt.Sprintf("--rpcuser=%s", n.rpcUser))
	}
	if n.rpcPass != "" {
		// --rpcpass
		args = append(args, fmt.Sprintf("--rpcpass=%s", n.rpcPass))
	}
	if n.listen != "" {
		// --listen
		args = append(args, fmt.Sprintf("--listen=%s", n.listen))
	}
	if n.rpcListen != "" {
		// --rpclisten
		args = append(args, fmt.Sprintf("--rpclisten=%s", n.rpcListen))
	}
	if n.rpcConnect != "" {
		// --rpcconnect
		args = append(args, fmt.Sprintf("--rpcconnect=%s", n.rpcConnect))
	}
	// --rpccert
	args = append(args, fmt.Sprintf("--rpccert=%s", n.certFile))
	// --rpckey
	args = append(args, fmt.Sprintf("--rpckey=%s", n.keyFile))
	// --txindex
	args = append(args, "--txindex")
	// --addrindex
	args = append(args, "--addrindex")
	if n.dataDir != "" {
		// --datadir
		args = append(args, fmt.Sprintf("--datadir=%s", n.dataDir))
	}
	if n.logDir != "" {
		// --logdir
		args = append(args, fmt.Sprintf("--logdir=%s", n.logDir))
	}
	if n.profile != "" {
		// --profile
		args = append(args, fmt.Sprintf("--profile=%s", n.profile))
	}
	if n.debugLevel != "" {
		// --debuglevel
		args = append(args, fmt.Sprintf("--debuglevel=%s", n.debugLevel))
	}
	args = append(args, n.extra...)
	return args
}

// command returns the exec.Cmd which will be used to start the hcd process.
func (n *nodeConfig) command() *exec.Cmd {
	return exec.Command(n.exe, n.arguments()...)
}

// rpcConnConfig returns the rpc connection config that can be used to connect
// to the hcd process that is launched via Start().
func (n *nodeConfig) rpcConnConfig() rpc.ConnConfig {
	return rpc.ConnConfig{
		Host:                 n.rpcListen,
		Endpoint:             n.endpoint,
		User:                 n.rpcUser,
		Pass:                 n.rpcPass,
		Certificates:         n.certificates,
		DisableAutoReconnect: true,
	}
}

// String returns the string representation of this nodeConfig.
func (n *nodeConfig) String() string {
	return n.prefix
}

// node houses the necessary state required to configure, launch, and manage a
// hcd process.
type node struct {
	config *nodeConfig

	cmd     *exec.Cmd
	pidFile string

	dataDir string
}

// newNode creates a new node instance according to the passed config. dataDir
// will be used to hold a file recording the pid of the launched process, and
// as the base for the log and data directories for hcd.
func newNode(config *nodeConfig, dataDir string) (*node, error) {
	return &node{
		config:  config,
		dataDir: dataDir,
		cmd:     config.command(),
	}, nil
}

// start creates a new hcd process, and writes its pid in a file reserved for
// recording the pid of the launched process. This file can be used to
// terminate the process in case of a hang, or panic. In the case of a failing
// test case, or panic, it is important that the process be stopped via stop(),
// otherwise, it will persist unless explicitly killed.
func (n *node) start() error {
	if err := n.cmd.Start(); err != nil {
		return err
	}

	pid, err := os.Create(filepath.Join(n.config.String(), "hcd.pid"))
	if err != nil {
		return err
	}

	n.pidFile = pid.Name()
	if _, err = fmt.Fprintf(pid, "%d\n", n.cmd.Process.Pid); err != nil {
		return err
	}

	if err := pid.Close(); err != nil {
		return err
	}

	return nil
}

// stop interrupts the running hcd process process, and waits until it exits
// properly. On windows, interrupt is not supported, so a kill signal is used
// instead
func (n *node) stop() error {
	if n.cmd == nil || n.cmd.Process == nil {
		// return if not properly initialized
		// or error starting the process
		return nil
	}
	defer n.cmd.Wait()
	if runtime.GOOS == "windows" {
		return n.cmd.Process.Signal(os.Kill)
	}
	return n.cmd.Process.Signal(os.Interrupt)
}

// cleanup cleanups process and args files. The file housing the pid of the
// created process will be deleted, as well as any directories created by the
// process.
func (n *node) cleanup() error {
	if n.pidFile != "" {
		if err := os.Remove(n.pidFile); err != nil {
			log.Printf("unable to remove file %s: %v", n.pidFile,
				err)
			return err
		}
	}

	return nil
}

// shutdown terminates the running hcd process, and cleans up all
// file/directories created by node.
func (n *node) shutdown() error {
	if err := n.stop(); err != nil {
		return err
	}
	if err := n.cleanup(); err != nil {
		return err
	}
	return nil
}

// genCertPair generates a key/cert pair to the paths provided.
func genCertPair(certFile, keyFile string) error {
	org := "rpctest autogenerated cert"
	validUntil := time.Now().Add(10 * 365 * 24 * time.Hour)
	cert, key, err := hcutil.NewTLSCertPair(elliptic.P521(), org,
		validUntil, nil)
	if err != nil {
		return err
	}

	// Write cert and key files.
	if err = ioutil.WriteFile(certFile, cert, 0666); err != nil {
		return err
	}
	if err = ioutil.WriteFile(keyFile, key, 0600); err != nil {
		os.Remove(certFile)
		return err
	}

	return nil
}
