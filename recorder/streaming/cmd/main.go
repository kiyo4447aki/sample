package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"streaming/config"
	"streaming/internal"
	"sync"
	"syscall"
	"time"
)

func main(){
	ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

	sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
    var wg sync.WaitGroup

	

	probe_timeout := 60 * time.Second
	probe_interval := 5 * time.Second

	log.Printf("設定を読み込んでいます")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("設定の読み込みに失敗しました：%v", err)
	}

	for name := range cfg.Devices {
		wg.Add(1)
		go func (deviceName string) {
			defer wg.Done()
			err := startStreamAndRecord(ctx, probe_timeout, probe_interval, deviceName, cfg)
			if err != nil {
				log.Fatalf("[%s] エラーが発生しました：%v",deviceName, err)
				cancel()
			}
		}(name)
	}
	<-sigCh
    log.Println("プロセスを停止します")
    cancel()
	wg.Wait()
    log.Println("プロセスが正常に停止されました")
}

func startStreamAndRecord(
	ctx context.Context, 
	probeTimeout time.Duration, 
	probeInterval time.Duration, 
	deviceName string, 
	cfg *config.Config,
	) error {
	
	log.Printf("[%s] パイプラインを生成します", deviceName)

	pl, err := internal.NewPipeLine(deviceName, 47000)
	if err != nil {
		return fmt.Errorf("[%s] パイプラインの生成に失敗しました: %v",deviceName , err)
	}

	defer func() {
		log.Printf("[%s] パイプライン停止処理開始", deviceName)
		pl.Stop()
		pl.Cleanup()
		log.Printf("[%s] パイプラインが正常に停止されました", deviceName)

	}()

	
	appsink, err := pl.Pipeline.GetElementByName("appsink")
	if err != nil {
		return fmt.Errorf("[%s] appsinkの取得に失敗しました：%v",deviceName, err)
	}

	log.Printf("[%s] janusへのストリームを開始します", deviceName)
	
	if err := internal.StartStream(
		ctx, 
		appsink, 
		cfg.JanusURL, 
		cfg.RoomID, 
		cfg.RoomPass,
		deviceName,
	); err != nil {
		return fmt.Errorf("[%s] janusへの接続に失敗しました：%v", deviceName, err)
	}

	log.Printf("[%s] 監視プローブを開始します", deviceName)

	internal.StartProbe(ctx, appsink, probeTimeout, probeInterval,func() {
		//TODO: 再接続処理の実装
	})

	log.Printf("[%s] パイプラインをスタートします", deviceName)
	if err := pl.Start(); err != nil {
		return fmt.Errorf("[%s] パイプラインをスタートできませんでした：%v",deviceName , err.Error())
	}

	log.Printf("[%s] 正常にスタートされました", deviceName)

	<-ctx.Done()
    log.Printf("[%s] 停止処理が要求されました", deviceName)
    return nil
}