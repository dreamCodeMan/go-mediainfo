# mediainfo
Golang binding for [mediainfo](https://mediaarea.net/en/MediaInfo)

Duration, Bitrate, Codec, Streams and a lot of other meta-information about media files can be extracted through it.

Your advices and suggestions are welcome!

## Example
```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/dreamCodeMan/go_mediainfo"
)

func main() {
	mediainfo, _ := mediainfo.GetMediaInfo("/Users/Fang/Movies/11.flv")
	info, _ := json.Marshal(mediainfo)
	fmt.Println(string(info), err)
}

```