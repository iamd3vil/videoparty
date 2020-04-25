package player

import (
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/DexterLB/mpvipc"
)

// mpv will implement Player interface for mpv player
type mpv struct {
	socketName string
	conn       *mpvipc.Connection
	eventChan  chan *mpvipc.Event
	stopChan   chan struct{}
	cb         map[string]EventCallback
	pid        *os.Process
}

// NewMpvPlayer returns a new mpv player instance
func NewMpvPlayer(filePath, socketName string) (Player, error) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		cmd := exec.Command("mpv", "--pause", "--input-ipc-server="+socketName, filePath)
		err := cmd.Start()
		if err != nil {
			return
		}

		wg.Done()

		cmd.Wait()
	}()

	wg.Wait()

	player := &mpv{
		cb:         make(map[string]EventCallback),
		socketName: socketName,
	}

	return player, nil
}

func (m *mpv) Listen() error {
	var (
		conn        *mpvipc.Connection
		eventsChan  = make(chan *mpvipc.Event)
		stopChan    = make(chan struct{}, 1)
		connRetries int
	)
	// Hack to try multiple times so that it can connect when mpv comes up
	for {
		time.Sleep(1 * time.Second)
		connRetries++
		conn = mpvipc.NewConnection(m.socketName)
		err := conn.Open()
		if err != nil {
			if connRetries == 10 {
				return err
			}
			continue
		}
		break
	}

	// Start a go routine to listen on events
	m.eventChan = eventsChan
	m.stopChan = stopChan

	go conn.ListenForEvents(eventsChan, stopChan)

	for {
		select {
		case e := <-eventsChan:
			log.Printf("event: %#v", e)
			m.handleEvent(e)
		}
	}
}

func (m *mpv) PauseCallback(pauseCB EventCallback) error {
	m.cb["pause"] = pauseCB
	return nil
}

func (m *mpv) ExitCallback(exitCB EventCallback) error {
	m.cb["end-file"] = exitCB
	return nil
}

func (m *mpv) Close() {
	m.stopChan <- struct{}{}
	os.Remove(m.socketName)
}

func (m *mpv) handleEvent(e *mpvipc.Event) {
	cb, ok := m.cb[e.Name]
	if !ok {
		return
	}

	cb(nil)
}
