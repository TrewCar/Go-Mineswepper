package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"Sapper/game"
	"Sapper/helper"
	"Sapper/ui"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
)

type player struct {
	helper.Animate
	delay float64
	t     float64
	i     int
	done  bool
	hold  bool
}

func (p *player) update(dt float64) {
	if p.done {
		return
	}
	if p.delay > 0 {
		p.delay -= dt
		return
	}
	p.t += dt
	for p.t >= p.Durations[p.i].Seconds() {
		p.t -= p.Durations[p.i].Seconds()
		if p.i == len(p.Frames)-1 {
			if !p.Loop {
				p.done = true
				return
			}
			p.i = 0
		} else {
			p.i++
		}
	}
}

func (p *player) sprite() *pixel.Sprite {
	return p.Frames[p.i]
}

type App struct {
	win *opengl.Window
	imd *imdraw.IMDraw

	sprites map[string]helper.Sprite  // "tiles/mine", "faces/cool", "digits/7" ...
	anims   map[string]helper.Animate // "reveal", "explosion" ...

	players map[[2]int]*player

	timer     float64
	startGame bool

	lPressed bool
	lUp      bool

	rPressed bool
	rUp      bool

	m    game.Map
	won  bool
	over bool

	padding     float64
	padding_top float64
	scale       float64
}

func NewApp(win *opengl.Window) (*App, error) {
	styles, err := helper.Get_styles()
	if err != nil {
		return nil, err
	}
	sprites, anims, err := helper.Build_sprites(styles[1])
	if err != nil {
		return nil, err
	}

	a := &App{
		win:         win,
		imd:         imdraw.New(nil),
		sprites:     sprites,
		anims:       anims,
		players:     make(map[[2]int]*player),
		padding:     6,
		padding_top: 100,
		scale:       2,
	}
	a.NewGame()
	return a, nil
}

func (a *App) NewGame() {
	a.m = *game.New(20, 20, 10)
	a.players = make(map[[2]int]*player)
	a.over = false
	a.won = false
	a.startGame = false
	a.timer = 0

	a.win.SetBounds(pixel.R(0, 0,
		a.padding*2+float64(a.m.W)*16*a.scale, a.padding*2+a.padding_top+float64(a.m.H)*16*a.scale))
}

func (a *App) sprite(id fmt.Stringer) *pixel.Sprite {
	s, ok := a.sprites[id.String()]
	if !ok {
		panic("нет спрайта " + id.String())
	}
	return s.Sprite
}

func (a *App) play(name ui.Anim, delay float64) *player {
	an, ok := a.anims[string(name)]
	if !ok {
		panic("нет анимации " + string(name))
	}
	return &player{Animate: an, delay: delay}
}

func (a *App) drawAt(s *pixel.Sprite, pos pixel.Vec) {
	half := s.Frame().Size().Scaled(0.5 * a.scale)
	s.Draw(a.win, pixel.IM.Scaled(pixel.ZV, a.scale).Moved(pos.Add(half)))
}

// ---------- ввод и время ----------

func (a *App) Update(dt float64) {

	a.onClickSetup()
	a.updatePlayers(dt)
	a.click()
}

func (a *App) onClickSetup() {
	w := a.win

	a.lUp = false
	a.rUp = false
	if w.JustPressed(pixel.MouseButtonLeft) {
		a.lPressed = true
	}
	if w.JustReleased(pixel.MouseButtonLeft) {
		a.lPressed = false
		a.lUp = true
	}
	if w.JustPressed(pixel.MouseButtonRight) {
		a.rPressed = true
	}
	if w.JustReleased(pixel.MouseButtonRight) {
		a.rPressed = false
		a.rUp = true
	}
}

func (a *App) click() {

	X, Y, in := a.mouseAtGrid()
	fX, fY := a.faceBound()

	temp_map := make([][]game.CellState, len(a.m.Cells))
	for i := range temp_map {
		temp_map[i] = make([]game.CellState, len(a.m.Cells[i]))
		for j := range temp_map[i] {
			temp_map[i][j] = a.m.Cells[i][j].State
		}
	}

	isBoom := false
	if a.lUp {
		if in {
			a.startGame = true
			if a.m.Cells[Y][X].State == game.Closed {
				isBoom = a.m.Open(X, Y)
			} else if a.m.Cells[Y][X].State == game.Opened {
				isBoom = a.m.Chord(X, Y)
			}
		} else if a.mouseInBound(fX, fY) {
			a.NewGame()
			return
		}
	} else if a.rUp {
		if in {
			if a.m.Cells[Y][X].State == game.Closed {
				a.m.ToggleMark(X, Y)
			} else if a.m.Cells[Y][X].State == game.Flagged {
				a.m.ToggleMark(X, Y)
			} else if a.m.Cells[Y][X].State == game.Questioned {
				a.m.ToggleMark(X, Y)
			}
		}
	}

	wasOver := a.over
	if isBoom {
		a.over = true
	}
	if wasOver {
		return
	}

	bx, by := float64(a.m.ExplodedX), float64(a.m.ExplodedY)

	for i := range temp_map {
		for j := range temp_map[i] {
			key := [2]int{j, i}
			c := a.m.Cells[i][j]

			switch {
			case isBoom && c.Mine && c.State != game.Flagged:
				p := a.play(ui.AnimExplosion, a.distance(float64(j), float64(i), bx, by)/10)
				p.hold = true
				a.players[key] = p
			case temp_map[i][j] == c.State:
			case c.State == game.Opened:
				a.players[key] = a.play(ui.AnimReveal, a.distance(float64(j), float64(i), float64(X), float64(Y))/10)
			case c.State == game.Flagged:
				a.players[key] = a.play(ui.AnimFlagPlant, 0)
			default:
				delete(a.players, key)
			}
		}
	}
}

func (a *App) updatePlayers(dt float64) {
	for i, anm := range a.players {
		if anm.done && !anm.hold {
			delete(a.players, i)
			continue
		}
		anm.update(dt)
	}
	if a.startGame && !a.over && !a.won {
		a.timer += dt
	}
}

// ---------- позиционирование ----------

func (a *App) distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}

func (a *App) inBound(pos_t pixel.Vec, pos_min pixel.Vec, pos_max pixel.Vec) bool {
	return pos_t.X >= pos_min.X && pos_t.Y >= pos_min.Y &&
		pos_t.X < pos_max.X && pos_t.Y < pos_max.Y
}

func (a *App) mouseInBound(pos_min pixel.Vec, pos_max pixel.Vec) bool {
	return a.inBound(a.win.MousePosition(), pos_min, pos_max)
}

func (a *App) mouseInGrid() bool {
	min, max := a.gridBound()
	return a.mouseInBound(min, max)
}

func (a *App) mouseAtGrid() (int, int, bool) {
	if !a.mouseInGrid() {
		return -1, -1, false
	}

	pad := pixel.V(-a.padding, -a.padding)

	mouse_pos := a.win.MousePosition().Add(pad)
	sb := a.sprite(ui.Tile1).Frame()

	X := int(mouse_pos.X) / int(sb.W()*a.scale)
	Y := int(mouse_pos.Y) / int(sb.W()*a.scale)

	return X, Y, true
}

func (a *App) gridBound() (pixel.Vec, pixel.Vec) {
	sb := a.sprite(ui.Tile1).Frame()
	pos := pixel.V(a.padding, a.padding)
	pos_max := pos.Add(pixel.V(sb.W()*a.scale*float64(a.m.W), sb.H()*a.scale*float64(a.m.H)))
	return pos, pos_max
}

func (a *App) faceBound() (pixel.Vec, pixel.Vec) {
	wb := a.win.Bounds()
	fb := a.sprite(ui.FaceSmile).Frame()
	pos := pixel.V(wb.W()/2-fb.W()*a.scale/2, wb.H()-a.padding-fb.H()*a.scale)
	pos_max := pos.Add(pixel.V(fb.W()*a.scale, fb.H()*a.scale))
	return pos, pos_max
}

// ---------- отрисовка ----------

func (a *App) Draw() {
	a.drawBackground()
	a.drawFace()
	a.drawDigest()
	a.drawBoard()
}

func (a *App) drawBackground() {
	open, closed := a.sprite(ui.TilePressed).Frame(), a.sprite(ui.TileClosed).Frame()
	face := spriteColor(a.sprite(ui.TilePressed), open.Center())
	light := spriteColor(a.sprite(ui.TileClosed), pixel.V(closed.Min.X+0.5, closed.Max.Y-0.5))
	shadow := spriteColor(a.sprite(ui.TileClosed), pixel.V(closed.Max.X-0.5, closed.Min.Y+0.5))

	a.win.Clear(face)

	wb := a.win.Bounds()
	gmin, gmax := a.gridBound()
	fmin, _ := a.faceBound()

	a.imd.Clear()
	a.sunkenFrame(pixel.R(gmin.X, gmin.Y, gmax.X, gmax.Y), a.padding, shadow, light)
	a.sunkenFrame(pixel.R(a.padding, fmin.Y, wb.W()-a.padding, wb.H()-a.padding), a.padding, shadow, light)
	a.imd.Draw(a.win)
}

func (a *App) sunkenFrame(r pixel.Rect, t float64, shadow, light pixel.RGBA) {
	o := pixel.R(r.Min.X-t, r.Min.Y-t, r.Max.X+t, r.Max.Y+t)
	side := func(c pixel.RGBA, pts ...pixel.Vec) {
		a.imd.Color = c
		a.imd.Push(pts...)
		a.imd.Polygon(0)
	}
	side(shadow, o.Min, r.Min, pixel.V(r.Min.X, r.Max.Y), pixel.V(o.Min.X, o.Max.Y)) // слева
	side(shadow, pixel.V(o.Min.X, o.Max.Y), pixel.V(r.Min.X, r.Max.Y), r.Max, o.Max) // сверху
	side(light, o.Max, r.Max, pixel.V(r.Max.X, r.Min.Y), pixel.V(o.Max.X, o.Min.Y))  // справа
	side(light, pixel.V(o.Max.X, o.Min.Y), pixel.V(r.Max.X, r.Min.Y), r.Min, o.Min)  // снизу
}

func spriteColor(s *pixel.Sprite, at pixel.Vec) pixel.RGBA {
	pd, ok := s.Picture().(*pixel.PictureData)
	if !ok {
		return pixel.RGB(0, 0, 0)
	}
	return pd.Color(at)
}

func (a *App) drawFace() {
	pos, pos_max := a.faceBound()

	face_sprite := a.sprite(ui.FaceSmile)
	if a.lPressed && a.mouseInBound(pos, pos_max) {
		face_sprite = a.sprite(ui.FaceSmilePressed)
	} else if a.over {
		face_sprite = a.sprite(ui.FaceDead)
	} else if a.won {
		face_sprite = a.sprite(ui.FaceCool)
	} else if a.lPressed && a.mouseInGrid() {
		face_sprite = a.sprite(ui.FaceWow)
	}

	a.drawAt(face_sprite, pos)
}

func (a *App) drawDigest() {
	w := a.win
	sprite_digest := a.sprite(ui.Digit0)

	pos_left_digest := pixel.V(a.padding, w.Bounds().H()-a.padding-sprite_digest.Frame().H()*a.scale)

	pos_right_digest := pixel.V(w.Bounds().W()-a.padding-sprite_digest.Frame().W()*a.scale, w.Bounds().H()-a.padding-sprite_digest.Frame().H()*a.scale)

	mines := strings.Split(strconv.Itoa(a.m.MinesLeft()), "")
	l := len(mines)
	for i := 0; i < 3-l; i++ {
		mines = append([]string{""}, mines...)
	}
	timer := strings.Split(strconv.Itoa(int(a.timer)), "")
	l = len(timer)
	for i := 0; i < 3-l; i++ {
		timer = append([]string{""}, timer...)
	}

	for i := 0; i < 3; i++ {
		pos_left := pos_left_digest.Add(pixel.V(sprite_digest.Frame().W()*a.scale*float64(i), 0))
		pos_right := pos_right_digest.Add(pixel.V(-sprite_digest.Frame().W()*a.scale*float64(i), 0))

		lSprite := a.getDigest(mines[i])
		rSprite := a.getDigest(timer[2-i])

		a.drawAt(lSprite, pos_left)
		a.drawAt(rSprite, pos_right)
	}
}

func (a *App) getDigest(char string) *pixel.Sprite {
	switch char {
	case "-":
		return a.sprite(ui.DigitMinus)
	case "":
		return a.sprite(ui.DigitBlank)
	default:
		n, _ := strconv.Atoi(char)
		return a.sprite(ui.Digit(n))
	}
}

func (a *App) drawBoard() {
	mx, my, in := a.mouseAtGrid()
	isOpen := in && a.m.Cells[my][mx].State == game.Opened
	for y := 0; y < a.m.H; y++ {
		for x := 0; x < a.m.W; x++ {
			sprite := a.sprite(a.m.Tile(x, y))
			if p, ok := a.players[[2]int{x, y}]; ok {
				sprite = p.sprite()
			}
			if a.lPressed && in && x == mx && y == my && a.m.Cells[y][x].State == game.Closed {
				sprite = a.sprite(ui.TilePressed)
			}

			if a.lPressed && a.m.Cells[y][x].State == game.Closed && isOpen {
				if !(x == mx && y == my) && (math.Abs(float64(y-my)) <= 1 && math.Abs(float64(x-mx)) <= 1) {
					sprite = a.sprite(ui.TilePressed)
				}
			}
			a.drawAt(sprite, pixel.V(float64(x)*sprite.Frame().H()*a.scale+a.padding, float64(y)*sprite.Frame().W()*a.scale+a.padding))
		}
	}
}
