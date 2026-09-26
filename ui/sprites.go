package ui

import "fmt"

// Размеры кадров в пикселях.
const (
	TileSize = 16 // tiles.png и большинство анимаций
	FaceSize = 26 // faces.png, face_blink, face_press
	DigitW   = 13 // digits.png
	DigitH   = 23
)

// Tile — кадр из tiles.png. Значение совпадает с номером кадра в png.
type Tile int

const (
	TileClosed Tile = iota // закрытая клетка (выпуклая)
	TileOpen               // открытая пустая клетка
	Tile1                  // Tile1…Tile8 — число мин вокруг
	Tile2
	Tile3
	Tile4
	Tile5
	Tile6
	Tile7
	Tile8
	TileMine         // мина на открытой клетке (конец игры)
	TileMineExploded // мина, на которой подорвался игрок
	TileMineWrong    // зачёркнутая мина: флаг стоял не на мине
	TileFlag         // флаг на закрытой клетке
	TileQuestion     // знак вопроса на закрытой клетке
	TileQuestionOpen // бледный знак вопроса на открытой клетке
	TilePressed      // клетка, пока на неё нажимают
)

var tileNames = [...]string{
	TileClosed:       "closed",
	TileOpen:         "open",
	Tile1:            "1",
	Tile2:            "2",
	Tile3:            "3",
	Tile4:            "4",
	Tile5:            "5",
	Tile6:            "6",
	Tile7:            "7",
	Tile8:            "8",
	TileMine:         "mine",
	TileMineExploded: "mine_exploded",
	TileMineWrong:    "mine_wrong",
	TileFlag:         "flag",
	TileQuestion:     "question",
	TileQuestionOpen: "question_open",
	TilePressed:      "pressed",
}

// String возвращает id спрайта для карты из helper.Build_sprites, например "tiles/mine".
func (t Tile) String() string {
	if t < 0 || int(t) >= len(tileNames) {
		return fmt.Sprintf("Tile(%d)", int(t))
	}
	return "tiles/" + tileNames[t]
}

// TileNumber возвращает тайл с числом n (1…8).
func TileNumber(n int) Tile {
	return Tile1 + Tile(n-1)
}

// Face — кадр из faces.png.
type Face int

const (
	FaceSmile        Face = iota // обычное состояние
	FaceSmilePressed             // на смайлик нажали (перезапуск)
	FaceWow                      // игрок зажал кнопку мыши на клетке
	FaceDead                     // проигрыш
	FaceCool                     // победа
)

var faceNames = [...]string{
	FaceSmile:        "smile",
	FaceSmilePressed: "smile_pressed",
	FaceWow:          "wow",
	FaceDead:         "dead",
	FaceCool:         "cool",
}

func (f Face) String() string {
	if f < 0 || int(f) >= len(faceNames) {
		return fmt.Sprintf("Face(%d)", int(f))
	}
	return "faces/" + faceNames[f]
}

// Digit — кадр из digits.png. Digit0…Digit9 равны самим цифрам, так что Digit(n) тоже работает.
type Digit int

const (
	Digit0 Digit = iota
	Digit1
	Digit2
	Digit3
	Digit4
	Digit5
	Digit6
	Digit7
	Digit8
	Digit9
	DigitMinus // мин осталось меньше нуля
	DigitBlank // все сегменты погашены
)

func (d Digit) String() string {
	switch {
	case d >= Digit0 && d <= Digit9:
		return fmt.Sprintf("digits/%d", int(d))
	case d == DigitMinus:
		return "digits/minus"
	case d == DigitBlank:
		return "digits/blank"
	}
	return fmt.Sprintf("Digit(%d)", int(d))
}
