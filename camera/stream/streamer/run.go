package streamer

import (
	"camera/config"
	"camera/scheduler"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
)


func (s *SwitchStreamer) Run(cfg config.Config)error{
	ticker := time.NewTicker(cfg.SwitchEvery)
	stop := make(chan struct{})

	//終了処理の重複防止
	var shutdownOnce sync.Once
	shutdown := func(){
		shutdownOnce.Do(
			func() {
				ticker.Stop()
				close(stop)
			})
	}

	loop := glib.NewMainLoop(glib.MainContextDefault(), false)

	bus := s.Pipeline.GetBus()
	bus.AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEOS:
			log.Println("EOS detected")
			shutdown()
			loop.Quit()
		case gst.MessageError:
			gErr := msg.ParseError()
			log.Printf("%sでエラーが発生しました：%v", msg.Source(), gErr)
			shutdown()
			loop.Quit()
		case gst.MessageStateChanged:
			srcName := msg.Source()
			if srcName == s.Pipeline.GetName(){
				oldSt, newSt := msg.ParseStateChanged()
				log.Printf("state: %s -> %s", oldSt, newSt)
			}
		}
		return true
	})

	if err := s.Pipeline.SetState(gst.StateReady); err != nil {
		return fmt.Errorf("stateの変更に失敗しました：%w", err)
	}
	

	loc, err := time.LoadLocation(cfg.Tz)
	if err != nil {
		return fmt.Errorf("タイムゾーンの読み込みに失敗しました：%w", err)
	}

	now := time.Now().In(loc)

	tw, err := scheduler.ParseTimeWindow(cfg.NightVisionTime)
	if err != nil {
		return err
	}

	min := now.Hour() * 60 + now.Minute()
	if tw.Contains(min){
		err := s.setActivePadSync("night")
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		err := s.setActivePadSync("normal")
		if err != nil {
			log.Println(err.Error())
		}
	}

	if err := s.Pipeline.SetState(gst.StatePlaying); err != nil {
		return fmt.Errorf("stateの変更に失敗しました：%w", err)
	}

	go func() {
		for {
			select {
			case <- ticker.C:
				n := time.Now().In(loc)
				m := n.Hour() * 60 + n.Minute()
				if tw.Contains(m){
					err := s.setActivePad("night")
					if err != nil {
						log.Println(err.Error())
					}
				} else {
					err := s.setActivePad("normal")
					if err != nil {
						log.Println(err.Error())
					}
				}
			case <- stop:
				return 
			}
		}
	}()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	go func(){
		<- sigc
		log.Println("signal detected：stopping...")
		shutdown()
		s.Pipeline.SendEvent(gst.NewEOSEvent())
	}()

	loop.Run()
	_ = s.Pipeline.SetState(gst.StateNull)
	return nil
}
