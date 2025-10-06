package streamer

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
)

type SwitchStreamer struct {
	Pipeline            *gst.Pipeline //パイプライン本体

	Sel                 *gst.Element  //input-selector
	SelPadNormal        *gst.Pad      //通常カメラの入力
	SelPadNight   *gst.Pad      //暗視カメラの入力
	ActiveName          string        //"normal", "night"
	mu sync.Mutex
}

func (s *SwitchStreamer) setActivePad(name string)error {
	s.mu.Lock()
	if name == s.ActiveName { 
		s.mu.Unlock()
		return nil
	} 
	s.mu.Unlock()
	var pad *gst.Pad
	var padVal *glib.Value
	var err error
	switch name {
	case "normal":
		pad = s.SelPadNormal
	case "night":
		pad = s.SelPadNight
	default:
		return fmt.Errorf("setActivePadの引数が不正です：%v", name)
	}

	if pad == nil {
        return fmt.Errorf("対象の Pad が nil です: %s", name)
    }

	padVal, err = pad.ToGValue()
	if err != nil {
		return fmt.Errorf("sinkPadをGValueに変換できませんでした：%w", err)
	}

	done := make(chan error, 1)
	glib.IdleAdd(func()bool{
		done <- s.Sel.SetProperty("active-pad", padVal)
		return false
	})

	select {
	case err := <-done:
    	if err != nil {
			return fmt.Errorf("active-padプロパティをセットできませんでした：%w", err)
		} 
	case <- time.After(2 * time.Second):
		return fmt.Errorf("active-pad設定がタイムアウトしました")
	}

	s.mu.Lock()
	s.ActiveName = name
	s.mu.Unlock()
	log.Printf("[switch] active = %s", name)
	return nil
}

func (s *SwitchStreamer) setActivePadSync(name string) error {
    var pad *gst.Pad
    switch name {
    case "normal":
        pad = s.SelPadNormal
    case "night":
        pad = s.SelPadNight
    default:
        return fmt.Errorf("setActivePadの引数が不正です：%v", name)
    }
    if pad == nil {
        return fmt.Errorf("対象の Pad が nil です: %s", name)
    }
    gv, err := pad.ToGValue()
    if err != nil {
        return fmt.Errorf("sinkPadをGValueに変換できませんでした：%w", err)
    }
    if err := s.Sel.SetProperty("active-pad", gv); err != nil {
        return fmt.Errorf("active-padプロパティをセットできませんでした：%w", err)
    }
	s.mu.Lock()
    s.ActiveName = name
	s.mu.Unlock()
    log.Printf("[switch] active = %s", name)
    return nil
}
