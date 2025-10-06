package streamer

import (
	"camera/config"
	"fmt"

	"github.com/go-gst/go-gst/gst"
)

func buildPipelineString(cfg config.Config)string{
	return fmt.Sprintf(`
input-selector name=sel sync-streams=true sync-mode=clock cache-buffers=true
    sel. ! queue name=q_out leaky=2 max-size-buffers=1 !
        mpph264enc name=enc bitrate=%d !
        h264parse name=h264parse config-interval=1 !
        rtph264pay name=pay pt=96 config-interval=1 mtu=1400 !
        udpsink name=sink host=%s port=%d sync=false

v4l2src device=%s !
    image/jpeg, width=%d, height=%d, framerate=%d/1 !
    jpegparse ! mppjpegdec !
    video/x-raw, format=NV12, width=%d, height=%d, framerate=%d/1 !
    queue name=q_normal leaky=2 max-size-buffers=2 !
    sel.

v4l2src device=%s io-mode=4 !
    video/x-raw, format=NV12, width=%d, height=%d, framerate=%d/1 !
    queue name=q_night leaky=2 max-size-buffers=2 !
    sel.
	`,
	cfg.Bitrate, cfg.Host, cfg.Port,
	cfg.NormalDev, cfg.Width, cfg.Height, cfg.Fps, cfg.Width, cfg.Height, cfg.Fps,
	cfg.NightDev, cfg.Width, cfg.Height, cfg.Fps,)
}

func BuildPipeline(cfg config.Config) (*SwitchStreamer, error){
	gst.Init(nil)

	ps := buildPipelineString(cfg)
	pipeline, err := gst.NewPipelineFromString(ps)
	if err != nil {
		return nil, fmt.Errorf("パイプラインの作成に失敗しました：%w", err)
	}

	sel, err := pipeline.GetElementByName("sel")
	if err != nil {
		return nil, fmt.Errorf("セレクターのエレメントを取得できませんでした：%w", err)
	}

	sinkPads, err := sel.GetSinkPads()
	if err != nil {
		return nil, fmt.Errorf("sinkpadsを取得できませんでした：%w", err)
	}
	var padNormal, padNight *gst.Pad
	for _, sp := range sinkPads{
		if peer := sp.GetPeer(); peer != nil{
			if parent := peer.GetParentElement(); parent != nil {
				switch parent.GetName(){
				case "q_normal":
					padNormal = sp
				case "q_night":
					padNight = sp
				}
			}
		}
	}

	if padNormal == nil || padNight == nil {
		return nil, fmt.Errorf("セレクターのsinkPadを取得できませんでした（normal：%v, night：%v）", padNormal, padNight)
	}

	return &SwitchStreamer{
		Pipeline: pipeline,
		Sel: sel,
		SelPadNormal: padNormal,
		SelPadNight: padNight,
	}, nil
}

