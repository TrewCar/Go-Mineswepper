package ui

import "image/color"

// Палитра из atlas.json (ключ palette). Структуры в Go не бывают const,
// поэтому цвета объявлены как var.
var (
	ColorBlack     = color.RGBA{0x11, 0x11, 0x11, 0xff} // K, почти чёрный
	ColorDarkGray  = color.RGBA{0x33, 0x33, 0x33, 0xff} // d, тёмно-серый
	ColorGray      = color.RGBA{0x80, 0x80, 0x80, 0xff} // D, серый, тени
	ColorLightGray = color.RGBA{0xc0, 0xc0, 0xc0, 0xff} // G, светло-серый, фон клеток
	ColorLight     = color.RGBA{0xe0, 0xe0, 0xe0, 0xff} // S, светлый
	ColorWhite     = color.RGBA{0xff, 0xff, 0xff, 0xff} // W, белый, блики
)
