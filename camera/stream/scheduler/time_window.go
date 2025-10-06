package scheduler

import (
	"fmt"
	"strconv"
	"strings"
)

//時間帯を分形式に変換
type TimeWindow struct {
	StartMin  int
	EndMin    int
}

func ParseTimeWindow(s string) (TimeWindow, error){
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return TimeWindow{}, fmt.Errorf("時間帯の形式が不正です")
	}
	toMin := func(hm string)(int, error){
		p := strings.Split(hm, ":")
		if len(p) != 2 {
			return 0, fmt.Errorf("時刻の形式が不正です")
		}

		h, err := strconv.Atoi(p[0])
		if err != nil {
			return 0, err
		}

		m, err := strconv.Atoi(p[1])
		if err != nil {
			return 0, err
		}

		if h < 0 || h > 23 || m < 0 || m > 59{
			return 0, fmt.Errorf("不正な時刻が入力されました")
		}

		return h * 60 + m, nil
	}

	startMin ,err := toMin(parts[0])
	if err != nil {
		return TimeWindow{}, err
	}

	endMin ,err := toMin(parts[1])
	if err != nil {
		return TimeWindow{}, err
	}

	return TimeWindow{
		StartMin: startMin,
		EndMin: endMin,
	}, nil
}

//現在時刻（分形式）が保持している時間帯と合致するか返す
func (w TimeWindow) Contains(min int) bool{
	if w.StartMin <= w.EndMin {
		return min >= w.StartMin && min <= w.EndMin
	}
	return min >= w.StartMin || min <= w.EndMin
}
