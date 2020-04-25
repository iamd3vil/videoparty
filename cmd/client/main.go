package main

import (
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/iamd3vil/videoparty/pkg/player"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Printf("got termination signal")
		done <- true
	}()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("couldn't get current working directory")
	}
	newPlayer, err := player.NewMpvPlayer(os.Args[1], path.Join(wd, "mpv_socket"))
	if err != nil {
		log.Fatalf("error starting the player: %v", err)
	}

	defer newPlayer.Close()

	newPlayer.PauseCallback(func(data interface{}) error {
		log.Println("paused....")
		return err
	})
	newPlayer.ExitCallback(func(data interface{}) error {
		log.Println("player exited...")
		newPlayer.Close()
		os.Exit(0)
		return err
	})

	go newPlayer.Listen()
	<-done
}
