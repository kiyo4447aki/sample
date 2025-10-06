package internal

import (
	"fmt"
	"time"

	"github.com/go-gst/go-gst/gst"
)

type Pipeline struct {
	Pipeline *gst.Pipeline
}


//udpsrcからappsinkとsplitmuxsinkするパイプラインを作成
func NewPipeLine(id string, port int)(*Pipeline, error){
	recordDir := "/tmp/id"
	//TODO カメラより先に起動しておくとnoPTSとなって落ちる問題の対応
	pipelineStr := fmt.Sprintf(`udpsrc port=%d caps="application/x-rtp,media=video,clock-rate=90000,encoding-name=H264,payload=96" ! rtpjitterbuffer latency=10 ! rtph264depay ! tee name=t ! queue ! h264parse config-interval=-1 ! video/x-h264,stream-format=byte-stream,alignment=au ! appsink name=appsink sync=false t. ! queue ! h264parse config-interval=1 ! video/x-h264,stream-format=avc,alignment=au ! splitmuxsink name="splitmuxsink" muxer-factory=mp4mux muxer-properties="properties,streamable=true,faststart=true"`, port)

	gst.Init(nil)

	pipeline, err := gst.NewPipelineFromString(pipelineStr)
	if err != nil {
		return nil, fmt.Errorf("パイプラインの作成に失敗しました：%w", err)
	}


	
	//splitmuxsinkの設定
	splitmux, err := pipeline.GetElementByName("splitmuxsink")
	if err != nil {
		return nil, fmt.Errorf("splitmuxsinkを取得できませんでした：%w", err)
	}
	splitmux.SetProperty("max-size-time", uint64(time.Minute.Nanoseconds()))
	// format-location シグナルでファイル名を動的生成
	splitmux.Connect("format-location",
		func(self *gst.Element, fragmentID uint) string {
			// YYYYmmdd-HHMMSS 形式のタイムスタンプ
			t := time.Now().Format("20060102-150405")
			// 例: /record/20250419-123045.mp4
			return fmt.Sprintf("%s/%s-%s.mp4",recordDir , t , id)
		},
	)
	


	return &Pipeline{
		Pipeline: pipeline,
	}, nil
}	

func(p *Pipeline) Start()error{
	if err := p.Pipeline.SetState(gst.StatePlaying); err != nil {
		return fmt.Errorf("パイプライン PLAYING への遷移に失敗しました：%w", err)
	}
	return nil
}

func (p *Pipeline) Stop() {
	//bus := p.Pipeline.GetBus()
	p.Pipeline.SendEvent(gst.NewEOSEvent())
	time.Sleep(time.Second * 5)
	p.Pipeline.SetState(gst.StateNull)
}

func (p *Pipeline) Cleanup() {
	p.Stop()
	p.Pipeline = nil
}





