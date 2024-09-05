package ffmpeg

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type ProbeData struct {
	Streams []Data `json:"streams"`
}

type Data struct {
	CodecType string `json:"codec_type"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Duration  string `json:"duration"`
}

func ProbeWithContextExec(fileName string, ctx context.Context) (*Data, error) {
	cmd := exec.CommandContext(ctx, "ffprobe", "-show_streams", "-of", "json", fileName)
	buf := bytes.NewBuffer(nil)
	cmd.Stdout = buf

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var resp ProbeData
	err = json.Unmarshal(buf.Bytes(), &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Streams) > 0 {
		return &resp.Streams[0], nil
	}

	return nil, errors.New("streams not found")
}

func ProbeFromBytes(ctx context.Context, file *[]byte) (*int, *string, error) {
	tmp, err := os.CreateTemp("", "*")
	if err != nil {
		return nil, nil, err
	}
	defer os.Remove(tmp.Name())

	_, err = tmp.Write(*file)
	if err != nil {
		return nil, nil, err
	}

	meta, err := ProbeWithContextExec(tmp.Name(), ctx)
	if err != nil {
		return nil, nil, err
	}

	if meta != nil {
		durationString, _ := strconv.ParseFloat(meta.Duration, 32)
		duration := int(durationString)
		resolution := fmt.Sprintf("%dx%d", meta.Width, meta.Height)

		if duration > 0 {
			return &duration, &resolution, nil
		}

		return nil, &resolution, nil
	}

	return nil, nil, nil
}
