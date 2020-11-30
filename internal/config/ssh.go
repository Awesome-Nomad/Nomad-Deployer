package config

import (
	"fmt"
	"github.com/phayes/freeport"
	"github.com/rgzr/sshtun"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SSHConnection is holding a ssh tunnel to remote address
// Dont use this in critical section
type SSHConnection struct {
	Config    *SSHConnectionConfig
	tunnel    *sshtun.SSHTun
	localPort int
	remoteUrl *url.URL
}

func (conn *SSHConnection) Init(address string) error {
	httpAddress := HTTPURLEnhancer(address)
	u, err := url.Parse(httpAddress)
	if err != nil {
		return err
	}
	localPort, err := freeport.GetFreePort()
	if err != nil {
		return err
	}
	remotePort, err := strconv.Atoi(u.Port())
	if err != nil {
		return err
	}

	// Get ssh address
	localEndpoint := &sshtun.Endpoint{
		Host: "127.0.0.1",
		Port: localPort,
	}
	sshEndpoint, err := conn.getSSHAddress(address)
	if err != nil {
		return err
	}
	targetEndpoint, err := createEndpoint(address, remotePort)
	if err != nil {
		return err
	}
	log.Printf("sshEndpoint: %+v", sshEndpoint)
	tunnel := newSSHTunnel(localEndpoint, sshEndpoint, targetEndpoint)
	// Setting up instance
	conn.remoteUrl = u
	conn.localPort = localPort
	conn.tunnel = tunnel
	// Start and wait until tunnel available or failed
	startTunnel(tunnel)
	return nil
}

func newSSHTunnel(localEndpoint, sshEndpoint, targetEndpoint *sshtun.Endpoint) *sshtun.SSHTun {
	tunnel := sshtun.New(localEndpoint.Port, sshEndpoint.Host, targetEndpoint.Port)
	tunnel.SetRemoteHost(targetEndpoint.Host)
	tunnel.SetPort(sshEndpoint.Port)
	tunnel.SetLocalHost(localEndpoint.Host)
	tunnel.SetDebug(viper.GetBool("verbose"))
	return tunnel
}

func (conn *SSHConnection) Destroy() error {
	conn.tunnel.Stop()
	return nil
}

func (conn *SSHConnection) GetAddress() (string, error) {
	localAddress := fmt.Sprintf("%s://%s:%d", conn.remoteUrl.Scheme, "127.0.0.1", conn.localPort)
	return localAddress, nil
}

func (conn *SSHConnection) getSSHAddress(targetAddress string) (endpoint *sshtun.Endpoint, err error) {
	var addr string
	if conn.Config.Address != "" {
		addr = conn.Config.Address
	} else {
		addr = strings.Split(targetAddress, ":")[0]
	}
	return createEndpoint(addr, SSHPort)
}

func createEndpoint(addr string, defaultPort int) (endpoint *sshtun.Endpoint, err error) {
	endpoint = &sshtun.Endpoint{}
	if len(addr) == 0 {
		endpoint.Host = "localhost"
		endpoint.Port = defaultPort
		return
	}
	pieces := strings.Split(addr, ":")
	endpoint.Host = pieces[0]
	if len(pieces) > 1 {
		endpoint.Port, err = strconv.Atoi(pieces[1])
	} else {
		endpoint.Port = defaultPort
	}
	return
}

func startTunnel(tunnel *sshtun.SSHTun) {
	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(1)
	// We set a callback to know when the tunnel is ready
	tunnel.SetConnState(func(_ *sshtun.SSHTun, state sshtun.ConnState) {
		switch state {
		case sshtun.StateStarting:
			log.Printf("STATE is Starting")
		case sshtun.StateStarted:
			wg.Done()
			log.Printf("STATE is Started")
		case sshtun.StateStopped:
			log.Printf("STATE is Stopped")
		}
	})
	// We start the tunnel (and restart it every time it is stopped)
	go func() {
		for {
			if err := tunnel.Start(); err != nil {
				log.Printf("SSH tunnel stopped: %s", err.Error())
				time.Sleep(time.Second) // don't flood if there's a start error :)
			}
		}
	}()
}
