package helper

import (
	"encoding/json"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gopxl/pixel/v2"
)

const PATH_TO_STYLES = "sprites/game"
const PATH_TO_SYSTEM = "sprites/system"
const NAME_ATLAS = "atlas.json"

func Load_system(file string) (*pixel.PictureData, error) {
	return load_picture(filepath.Join(PATH_TO_SYSTEM, file))
}

func Get_styles() ([]string, error) {
	entries, err := os.ReadDir(PATH_TO_STYLES)
	if err != nil {
		return nil, err
	}
	var styles []string
	for _, e := range entries {
		if e.IsDir() {
			styles = append(styles, e.Name())
		}
	}
	return styles, nil
}

func load_sprites(style string) (Sprite_style, error) {
	var sprites Sprite_style
	path := filepath.Join(PATH_TO_STYLES, style, NAME_ATLAS)
	file, err := os.Open(path)
	if err != nil {
		return sprites, err
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(&sprites); err != nil {
		return sprites, fmt.Errorf("decode %s: %w", path, err)
	}
	return sprites, nil
}

func load_picture(path string) (*pixel.PictureData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}
	return pixel.PictureDataFromImage(img), nil
}

func frame_rect(pic *pixel.PictureData, x, y, w, h int) pixel.Rect {
	top := pic.Bounds().H() - float64(y)
	return pixel.R(float64(x), top-float64(h), float64(x+w), top)
}

func Build_sprites(style string) (map[string]Sprite, map[string]Animate, error) {
	atlas, err := load_sprites(style)
	if err != nil {
		return nil, nil, err
	}
	sprites := make(map[string]Sprite)
	animates := make(map[string]Animate)
	for file, sh := range atlas.Sheets {
		pic, err := load_picture(filepath.Join(PATH_TO_STYLES, style, file))
		if err != nil {
			return nil, nil, err
		}
		name := strings.TrimSuffix(file, filepath.Ext(file))
		for i, frame := range sh.Frames {
			rect := frame_rect(pic, i*sh.Frame_w, 0, sh.Frame_w, sh.Frame_h)
			sprites[name+"/"+frame] = Sprite{pixel.NewSprite(pic, rect)}
		}
		for id, an := range sh.Animations {
			if len(an.Durations_ms) != an.Frames {
				return nil, nil, fmt.Errorf("%s: %s: frames=%d, durations_ms=%d", file, id, an.Frames, len(an.Durations_ms))
			}
			a := Animate{Loop: an.Loop}
			for i := 0; i < an.Frames; i++ {
				rect := frame_rect(pic, an.X+i*an.Frame_w, an.Y, an.Frame_w, an.Frame_h)
				a.Frames = append(a.Frames, pixel.NewSprite(pic, rect))
				a.Durations = append(a.Durations, time.Duration(an.Durations_ms[i])*time.Millisecond)
			}
			animates[id] = a
		}
	}
	return sprites, animates, nil
}
