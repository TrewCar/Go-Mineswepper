package helper

import (
	"time"

	"github.com/gopxl/pixel/v2"
)

type Sprite struct {
	*pixel.Sprite
}

type Animate struct {
	Frames    []*pixel.Sprite
	Durations []time.Duration
	Loop      bool
}

type palette struct {
	K  string `json:"K"`
	Dk string `json:"d"`
	D  string `json:"D"`
	G  string `json:"G"`
	S  string `json:"S"`
	W  string `json:"W"`
}

type animation struct {
	X            int   `json:"x"`
	Y            int   `json:"y"`
	Frame_w      int   `json:"frame_w"`
	Frame_h      int   `json:"frame_h"`
	Frames       int   `json:"frames"`
	Durations_ms []int `json:"durations_ms"`
	Loop         bool  `json:"loop"`
}

type sheet struct {
	Frame_w    int                  `json:"frame_w,omitempty"`
	Frame_h    int                  `json:"frame_h,omitempty"`
	Frames     []string             `json:"frames,omitempty"`
	Nine_slice int                  `json:"nine_slice,omitempty"`
	Hotspot    *[2]int              `json:"hotspot,omitempty"`
	Animations map[string]animation `json:"animations,omitempty"`
}

type notes struct {
	Icons   string `json:"icons"`
	Ui      string `json:"ui"`
	Cursors string `json:"cursors"`
}

type atlas struct {
	Palette palette          `json:"palette"`
	Sheets  map[string]sheet `json:"sheets"`
	Notes   notes            `json:"notes"`
}

type Sprite_style struct {
	Style     string           `json:"style"`
	Tile_size int              `json:"tile_size"`
	Palette   palette          `json:"palette"`
	Sheets    map[string]sheet `json:"sheets"`
}
