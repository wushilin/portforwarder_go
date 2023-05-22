package worker

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/wushilin/portforwarder_go/logging"
)

type WorkerConfig struct {
	BindAddress  string
	BindPort     int
	TargetHost   string
	TargetPort   int
	Uploaded     uint64
	Downloaded   uint64
	TotalHandled int64
	Active       int64
}

var ID_GEN uint64 = 0

func (v *WorkerConfig) Start(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	listener := fmt.Sprintf("%s:%d", v.BindAddress, v.BindPort)
	server, err := net.Listen("tcp4", listener)
	if err != nil {
		ERROR("Failed to LISTEN to %s:%s", listener, err)
		os.Exit(1)
	}
	defer server.Close()
	INFO("Listening %s:%d -> %s:%d", v.BindAddress, v.BindPort, v.TargetHost, v.TargetPort)
	for {
		connection, err := server.Accept()
		if err != nil {
			ERROR("Error accepting connection: %s ", err.Error())
			os.Exit(1)
		}
		var connection_id = atomic.AddUint64(&ID_GEN, 1)
		INFO("%d Accepted from %v to %v", connection_id, connection.RemoteAddr(), connection.LocalAddr())
		go v.processClient(connection, connection_id)
	}
}

func (v *WorkerConfig) processClient(connection net.Conn, conn_id uint64) {
	atomic.AddInt64(&v.Active, 1)
	start := time.Now()
	var uploaded uint64 = 0
	var downloaded uint64 = 0
	defer func() {
		connection.Close()
		INFO("%d Done. Uptime: %v Uploaded: %d bytes Downloaded: %d bytes", conn_id, time.Since(start), uploaded, downloaded)
		atomic.AddUint64(&v.Uploaded, uploaded)
		atomic.AddUint64(&v.Downloaded, downloaded)
		atomic.AddInt64(&v.Active, -1)
		atomic.AddInt64(&v.TotalHandled, 1)
	}()

	// client hello must be read in 30 seconds
	dest, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", v.TargetHost, v.TargetPort), 10*time.Second)
	if err != nil {
		WARN("%d Client %v can't connect to host %s: %s", conn_id, connection.RemoteAddr(), v.TargetHost, err)
		return
	}
	INFO("%d Client %v connected to host %s via %s", conn_id, connection.RemoteAddr(), v.TargetHost, dest.LocalAddr())
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go pipe(conn_id, &uploaded, &downloaded, connection, dest, wg)
	wg.Wait()
}

func pipe(conn_id uint64, uploaded *uint64, downloaded *uint64, src net.Conn, dest net.Conn, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		written, _ := io.Copy(src, dest)
		atomic.AddUint64(downloaded, uint64(written))
		src.Close()
	}()
	go func() {
		defer wg.Done()
		written, _ := io.Copy(dest, src)
		atomic.AddUint64(uploaded, uint64(written))
		dest.Close()
	}()
}
