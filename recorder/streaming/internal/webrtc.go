package internal

import (
	"context"
	"log"
	"time"

	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
	"github.com/notedit/janus-go"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media"
)

/*
func watchHandle(handle *janus.Handle){
	for {
		msg := <- handle.Events
		switch msg := msg.(type){
		case *janus.SlowLinkMsg:
			log.Println("SlowLinkMsg type ", handle.ID)
		case *janus.MediaMsg:
            log.Println("MediaEvent type", msg.Type, " receiving ", msg.Receiving)
        case *janus.WebRTCUpMsg:
            log.Println("WebRTCUp type ", handle.ID)
        case *janus.HangupMsg:
            log.Println("HangupEvent type ", handle.ID)
        case *janus.EventMsg:
            log.Printf("EventMsg %+v", msg.Plugindata.Data)
		}
	}
}
*/

func StartStream(
	ctx context.Context,
	appsink *gst.Element,
	janusURL string,
	roomID string,
	roomPass string,
	feedID string,
) error {
	// 再接続トリガー用チャネル
    reconnectCh := make(chan struct{}, 1)

	//STUN・TURNサーバーの設定
	webRTCcfg := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{{
			URLs: []string{"stun:stun.l.google.com:19302"},
		}},
		SDPSemantics: webrtc.SDPSemanticsUnifiedPlanWithFallback,
	}

	

	go func(){
		RESTART:
			for {
				//peer connectionの作成
				pc, err := webrtc.NewPeerConnection(webRTCcfg)
				if err != nil {
					log.Printf("PeerConnectionの作成に失敗しました：%v", err.Error())
				}

				pc.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
					log.Printf("Connection State has changed %s \n", state.String())
					switch state {
					case webrtc.ICEConnectionStateDisconnected, webrtc.ICEConnectionStateFailed:
						select{
						case reconnectCh <- struct{}{}:
						default:
						}
					}
				})

				videoTrack, err := webrtc.NewTrackLocalStaticSample(
					webrtc.RTPCodecCapability{
						MimeType: "video/h264",
					}, "video", "main",
				)
				if err != nil {
					log.Printf("VideoTrackの作成に失敗しました：%v", err.Error())
				} else if _, err = pc.AddTrack(videoTrack); err != nil{
					log.Printf("PeerConnectionにTrackを追加できませんでした：%v", err.Error())
				}

				app.SinkFromElement(appsink).SetCallbacks(&app.SinkCallbacks{
					NewSampleFunc: func(sink *app.Sink) gst.FlowReturn{
						sample := sink.PullSample()
						if sample == nil {
							return gst.FlowEOS
						}

						buffer := sample.GetBuffer()
						if buffer == nil {
							return gst.FlowError
						}

						samples := buffer.Map(gst.MapRead).Bytes()
						defer buffer.Unmap()

						if err := videoTrack.WriteSample(media.Sample{
							Data: samples, 
							Duration: *buffer.Duration().AsDuration(),
						}); err != nil {
							return gst.FlowError
						}
						return gst.FlowOK
					},
				})

				offer, err := pc.CreateOffer(nil)
				if err != nil {
					log.Printf("Offerの作成に失敗しました：%v", err.Error())
					time.Sleep(5 * time.Second)
					continue
				}
				gatherComplete := webrtc.GatheringCompletePromise(pc)
				if err = pc.SetLocalDescription(offer); err != nil {
					log.Printf("Offerをセットできませんでした：%v", err.Error())
					time.Sleep(5 * time.Second)
					continue
				}
				<- gatherComplete

				gateway, err := janus.Connect(janusURL)
				if err != nil {
					log.Printf("janusへの接続に失敗しました：%v", err.Error())
					time.Sleep(5 * time.Second)
					continue	
				}

				session, err := gateway.Create()
				if err != nil {
					log.Printf("janusセッションの作成に失敗しました：%v", err.Error())	
					time.Sleep(5 * time.Second)
					continue
				}

				handle, err := session.Attach("janus.plugin.videoroom")
				if err != nil {
					log.Printf("video room pluginへのアタッチに失敗しました：%v", err.Error())	
					time.Sleep(5 * time.Second)
					continue
				}
				
				//KeepAlive停止用チャネル
				stopKA := make(chan struct{})
				//KeepAlive送信用Goroutine
				go func(){
					ticker := time.NewTicker(5 * time.Second)
					defer ticker.Stop()
					for {
						select{
						case <- ticker.C:
							if _, err := session.KeepAlive(); err != nil {
								log.Printf("KeepAliveの送信時にエラーが発生しました: %v", err)
								select {
								case reconnectCh <- struct{}{}:
								default:
								}
								return
							}
						case <- stopKA:
							return
						}
						
						
					}
				}()


				_, err = handle.Message(map[string]interface{}{
					"request": "join",
					"ptype":   "publisher",
					"room":    roomID,
					"id":      "cam-01",
				}, nil)
				if err != nil {
					log.Printf("Roomへの参加に失敗しました：%v", err.Error())
					close(stopKA)
					time.Sleep(5 * time.Second)
					continue
				}
				
				msg, err := handle.Message(map[string]interface{}{
					"request": "publish",
					"audio":   false,
					"video":   true,
					"data":    false,
				}, map[string]interface{}{
					"type":    "offer",
					"sdp":     pc.LocalDescription().SDP,
					"trickle": false,
				})
				if err != nil {
					log.Printf("Publishメッセージの送信に失敗しました")
					close(stopKA)
					time.Sleep(5 * time.Second)
					continue
				}

				if msg.Jsep != nil {
					sdpVal, ok := msg.Jsep["sdp"].(string)
					if !ok {
						log.Printf("remoteDescriptionの取得に失敗しました")
						close(stopKA)
						time.Sleep(5 * time.Second)
						continue
					}
					err = pc.SetRemoteDescription(webrtc.SessionDescription{
						Type: webrtc.SDPTypeAnswer,
						SDP: sdpVal,
						
					})
					if err != nil {
						log.Printf("remoteDescriptionの登録に失敗しました")
						close(stopKA)
						time.Sleep(5 * time.Second)
						continue
					}
				}

				log.Println("publishに成功しました")

				for {
					select {
					case <- ctx.Done():
						close(stopKA)
						pc.Close()
						return
					case <- reconnectCh:
						log.Println("janusサーバへ再接続します")
							close(stopKA)
							pc.Close()
							continue RESTART
						
					case ev, ok :=  <-handle.Events:
						if !ok {
							log.Println("janusサーバへ再接続します")
							close(stopKA)
							pc.Close()
							continue RESTART
						}
						switch ev := ev.(type){
						case *janus.SlowLinkMsg:
							log.Println("SlowLinkMsg type ", handle.ID)
						case *janus.MediaMsg:
							log.Println("MediaEvent type", ev.Type, " receiving ", ev.Receiving)
						case *janus.WebRTCUpMsg:
							log.Println("WebRTCUp type ", handle.ID)
						case *janus.HangupMsg:
							log.Println("HangupEvent type ", handle.ID)
						case *janus.EventMsg:
							log.Printf("EventMsg %+v", msg.Plugindata.Data)
						}
					}
					time.Sleep(5 * time.Second)
				}
					
			}
	}()

	

	
	return nil
}


