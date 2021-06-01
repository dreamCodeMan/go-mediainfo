package mediainfo

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os/exec"
	"strings"
)

var mediainfoBinary = flag.String("mediainfo-bin", "mediainfo", "the path to the mediainfo binary if it is not in the system $PATH")

type mediainfo struct {
	XMLName xml.Name `xml:"MediaInfo"`
	Media    media     `xml:"media"`
	FilePath string  `xml:"ref,attr"`
}

type track struct {
	XMLName                   xml.Name `xml:"track"`
	Type                      string   `xml:"type,attr"`
	FormatProfile            string   `xml:"Format_Profile"`
	FileExtension            string   `xml:"FileExtension"`
	ColorSpace               string   `xml:"ColorSpace"`
	ChromaSubsampling        string   `xml:"ChromaSubsampling"`
	EncodedApplication       string   `xml:"Encoded_Application"`
	StreamSizeProportion string   `xml:"StreamSize_Proportion"`
	Width                     []string `xml:"Width"`
	Height                    []string `xml:"Height"`
	Format                    []string `xml:"Format"`
	Duration                  []string `xml:"Duration"`
	BitRate                  []string `xml:"BitRate"`
	BitDepth                 []string `xml:"BitDepth"`
	ScanType                 []string `xml:"ScanType"`
	FileSize                 []string `xml:"FileSize"`
	FrameRate                []string `xml:"FrameRate"`
	Channels                []string `xml:"Channels"`
	StreamSize               []string `xml:"StreamSize"`
	BitRateMode             []string `xml:"BitRate_Mode"`
	SamplingRate             []string `xml:"SamplingRate"`
	FrameRateMode           []string `xml:"FrameRate_Mode"`
	OverallBitRate          []string `xml:"OverallBitRate"`
	DisplayAspectRatio      []string `xml:"Display_aspect_ratio"`
	OverallBitRateMode     []string `xml:"OverallBitRate_Mode"`
	FormatSettingsCABAC    []string `xml:"Format_Settings_CABAC"`
	FormatSettingsRefFrames []string `xml:"Format_Settings_RefFrames"`
}

type media struct {
	XMLName xml.Name `xml:"media"`
	Tracks  []track  `xml:"track"`
}

type MediaInfo struct {
	General general `json:"general,omitempty"`
	Video   video   `json:"video,omitempty"`
	Audio   audio   `json:"audio,omitempty"`
	Menu    menu    `json:"menu,omitempty"`
}

type general struct {
	Format                string `json:"format"`
	Duration              string `json:"duration"`
	File_size             string `json:"file_size"`
	Overall_bit_rate_mode string `json:"overall_bit_rate_mode"`
	Overall_bit_rate      string `json:"overall_bit_rate"`
	File_extension        string `json:"file_extension"`
	Frame_rate            string `json:"frame_rate"`
	Stream_size           string `json:"stream_size"`
}

type video struct {
	Width                     string `json:"width"`
	Height                    string `json:"height"`
	Format                    string `json:"format"`
	Bit_rate                  string `json:"bitrate"`
	Duration                  string `json:"duration"`
	Format_profile            string `json:"format_profile"`
	Format_settings__CABAC    string `json:"format_settings_cabac"`
	Format_settings__ReFrames string `json:"format_settings_reframes"`
	Frame_rate                string `json:"frame_rate"`
	Bit_depth                 string `json:"bit_depth"`
	Scan_type                 string `json:"scan_type"`
}

type audio struct {
	Format         string `json:"format"`
	Duration       string `json:"duration"`
	Bit_rate       string `json:"bitrate"`
	Channel_s_     string `json:"channels"`
	Frame_rate     string `json:"frame_rate"`
	Sampling_rate  string `json:"sampling_rate"`
	Format_profile string `json:"format_profile"`
}

type menu struct {
	Format   string `json:"format"`
	Duration string `json:"duration"`
}

func IsInstalled() bool {
	cmd := exec.Command(*mediainfoBinary)
	err := cmd.Run()
	if err != nil {
		if strings.HasSuffix(err.Error(), "no such file or directory") ||
			strings.HasSuffix(err.Error(), "executable file not found in %PATH%") ||
			strings.HasSuffix(err.Error(), "executable file not found in $PATH") {
			return false
		} else if strings.HasPrefix(err.Error(), "exit status 255") {
			return true
		}
	}
	return true
}

func (info MediaInfo) IsMedia() bool {
	return info.Video.Duration != "" && info.Audio.Duration != ""
}

func GetMediaInfo(fname string) (MediaInfo, error) {
	info := MediaInfo{}
	mInfo := mediainfo{}
	general := general{}
	video := video{}
	audio := audio{}
	menu := menu{}

	if !IsInstalled() {
		return info, fmt.Errorf("Must install mediainfo")
	}
	out, err := exec.Command(*mediainfoBinary, "--Output=XML", "-f", fname).Output()

	if err != nil {
		return info, err
	}
	
	if err := xml.Unmarshal(out, &mInfo); err != nil {
		return info, err
	}

	for _, v := range mInfo.Media.Tracks {
		if v.Type == "General" {
			general.Duration = v.Duration[0]
			general.Format = v.Format[0]
			general.File_size = v.FileSize[0]
			if len(v.OverallBitRateMode) > 0 {
				general.Overall_bit_rate_mode = v.OverallBitRateMode[0]
			}
			general.Overall_bit_rate = v.OverallBitRate[0]
			general.File_extension = v.FileExtension
			general.Frame_rate = v.FrameRate[0]
			general.Stream_size = v.StreamSize[0]
		} else if v.Type == "Video" {
			video.Width = v.Width[0]
			video.Height = v.Height[0]
			video.Format = v.Format[0]
			video.Bit_rate = v.BitRate[0]
			video.Duration = v.Duration[0]
			video.Bit_depth = v.BitDepth[0]
			video.Scan_type = v.ScanType[0]
			video.Frame_rate = v.FrameRate[0]
			video.Format_profile = v.FormatProfile
			video.Format_settings__CABAC = v.FormatSettingsCABAC[0]
			video.Format_settings__ReFrames = v.FormatSettingsRefFrames[0]
		} else if v.Type == "Audio" {
			audio.Format = v.Format[0]
			audio.Channel_s_ = v.Channels[0]
			audio.Duration = v.Duration[0]
			audio.Bit_rate = v.BitRate[0]
			audio.Frame_rate = v.FrameRate[0]
			audio.Sampling_rate = v.SamplingRate[0]
			audio.Format_profile = v.FormatProfile
		} else if v.Type == "Menu" {
			menu.Duration = v.Duration[0]
			menu.Format = v.Format[0]
		}
	}
	info = MediaInfo{General: general, Video: video, Audio: audio, Menu: menu}

	return info, nil
}
