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
	File    file     `xml:"media"`
}

type track struct {
	XMLName                   xml.Name `xml:"track"`
	Type                      string   `xml:"type,attr"`
	File_name                 string   `xml:"File_name"`
	Format_Info               string   `xml:"Format_Info"`
	Color_space               string   `xml:"Color_space"`
	Complete_name             string   `xml:"Complete_name"`
	Format_profile            string   `xml:"Format_profile"`
	File_extension            string   `xml:"File_extension"`
	Chroma_subsampling        string   `xml:"Chroma_subsampling"`
	Writing_application       string   `xml:"Writing_application"`
	Proportion_of_this_stream string   `xml:"Proportion_of_this_stream"`
	Width                     []string `xml:"Width"`
	Height                    []string `xml:"Height"`
	Format                    []string `xml:"Format"`
	Duration                  []string `xml:"Duration"`
	Bit_rate                  []string `xml:"Bit_rate"`
	Bit_depth                 []string `xml:"Bit_depth"`
	Scan_type                 []string `xml:"Scan_type"`
	File_size                 []string `xml:"File_size"`
	Frame_rate                []string `xml:"Frame_rate"`
	Channel_s_                []string `xml:"Channel_s_"`
	Stream_size               []string `xml:"Stream_size"`
	Interlacement             []string `xml:"Interlacement"`
	Bit_rate_mode             []string `xml:"Bit_rate_mode"`
	Sampling_rate             []string `xml:"Sampling_rate"`
	Writing_library           []string `xml:"Writing_library"`
	Frame_rate_mode           []string `xml:"Frame_rate_mode"`
	Overall_bit_rate          []string `xml:"Overall_bit_rate"`
	Display_aspect_ratio      []string `xml:"Display_aspect_ratio"`
	Overall_bit_rate_mode     []string `xml:"Overall_bit_rate_mode"`
	Format_settings__CABAC    []string `xml:"Format_settings__CABAC"`
	Format_settings__ReFrames []string `xml:"Format_settings__ReFrames"`
}

type file struct {
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
	Complete_name         string `json:"complete_name"`
	File_name             string `json:"file_name"`
	File_extension        string `json:"file_extension"`
	Frame_rate            string `json:"frame_rate"`
	Stream_size           string `json:"stream_size"`
	Writing_application   string `json:"writing_application"`
}

type video struct {
	Width                     string `json:"width"`
	Height                    string `json:"height"`
	Format                    string `json:"format"`
	Bit_rate                  string `json:"bitrate"`
	Duration                  string `json:"duration"`
	Format_Info               string `json:"format_info"`
	Format_profile            string `json:"format_profile"`
	Format_settings__CABAC    string `json:"format_settings_cabac"`
	Format_settings__ReFrames string `json:"format_settings__reframes"`
	Frame_rate                string `json:"frame_rate"`
	Bit_depth                 string `json:"bit_depth"`
	Scan_type                 string `json:"scan_type"`
	Interlacement             string `json:"interlacement"`
	Writing_library           string `json:"writing_library"`
}

type audio struct {
	Format         string `json:"format"`
	Duration       string `json:"duration"`
	Bit_rate       string `json:"bitrate"`
	Channel_s_     string `json:"channels"`
	Frame_rate     string `json:"frame_rate"`
	Format_Info    string `json:"format_Info"`
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
	minfo := mediainfo{}
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

	if err := xml.Unmarshal(out, &minfo); err != nil {
		return info, err
	}

	for _, v := range minfo.File.Tracks {
		if v.Type == "General" {
			general.Duration = v.Duration[0]
			general.Format = v.Format[0]
			general.File_size = v.File_size[0]
			if len(v.Overall_bit_rate_mode) > 0 {
				general.Overall_bit_rate_mode = v.Overall_bit_rate_mode[0]
			}
			general.Overall_bit_rate = v.Overall_bit_rate[0]
			general.Complete_name = v.Complete_name
			general.File_name = v.File_name
			general.File_extension = v.File_extension
			general.Frame_rate = v.Frame_rate[0]
			general.Stream_size = v.Stream_size[0]
			general.Writing_application = v.Writing_application
		} else if v.Type == "Video" {
			video.Width = v.Width[0]
			video.Height = v.Height[0]
			video.Format = v.Format[0]
			video.Bit_rate = v.Bit_rate[0]
			video.Duration = v.Duration[0]
			video.Bit_depth = v.Bit_depth[0]
			video.Scan_type = v.Scan_type[0]
			video.Format_Info = v.Format_Info
			video.Frame_rate = v.Frame_rate[0]
			video.Format_profile = v.Format_profile
			video.Interlacement = v.Interlacement[1]
			video.Writing_library = v.Writing_library[0]
			video.Format_settings__CABAC = v.Format_settings__CABAC[0]
			video.Format_settings__ReFrames = v.Format_settings__ReFrames[0]
		} else if v.Type == "Audio" {
			audio.Format = v.Format[0]
			audio.Channel_s_ = v.Channel_s_[0]
			audio.Duration = v.Duration[0]
			audio.Bit_rate = v.Bit_rate[0]
			audio.Format_Info = v.Format_Info
			audio.Frame_rate = v.Frame_rate[0]
			audio.Sampling_rate = v.Sampling_rate[0]
			audio.Format_profile = v.Format_profile
		} else if v.Type == "Menu" {
			menu.Duration = v.Duration[0]
			menu.Format = v.Format[0]
		}
	}
	info = MediaInfo{General: general, Video: video, Audio: audio, Menu: menu}

	return info, nil
}
