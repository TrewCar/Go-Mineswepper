package ui

import "fmt"

// Системные спрайты из sprites/system. Build_sprites их пока не загружает,
// поэтому здесь пути к файлам относительно sprites/system.

const IconSize = 16

// IconTheme — набор иконок.
type IconTheme string

const (
	IconDark  IconTheme = "dark"  // тёмные иконки для светлого фона
	IconLight IconTheme = "light" // светлые иконки для тёмного фона
)

// Icon — кадр из icons/icons_<тема>.png. Значение совпадает с номером кадра в png.
type Icon int

const (
	IconMine       Icon = iota // мина
	IconFlag                   // флаг
	IconQuestion               // знак вопроса
	IconClock                  // часы (таймер)
	IconGear                   // шестерёнка (настройки)
	IconRestart                // круговая стрелка (заново)
	IconTrophy                 // кубок (рекорды)
	IconStar                   // звезда
	IconHeart                  // сердце
	IconSkull                  // череп
	IconSoundOn                // звук включён
	IconSoundOff               // звук выключен
	IconPlay                   // играть
	IconPause                  // пауза
	IconClose                  // крестик (закрыть)
	IconMenu                   // три полоски (меню)
	IconHome                   // домик (главное меню)
	IconHint                   // лампочка (подсказка)
	IconCheck                  // галочка
	IconDiffEasy               // сложность: лёгкая
	IconDiffMedium             // сложность: средняя
	IconDiffHard               // сложность: сложная
)

var iconNames = [...]string{
	IconMine:       "mine",
	IconFlag:       "flag",
	IconQuestion:   "question",
	IconClock:      "clock",
	IconGear:       "gear",
	IconRestart:    "restart",
	IconTrophy:     "trophy",
	IconStar:       "star",
	IconHeart:      "heart",
	IconSkull:      "skull",
	IconSoundOn:    "sound_on",
	IconSoundOff:   "sound_off",
	IconPlay:       "play",
	IconPause:      "pause",
	IconClose:      "close",
	IconMenu:       "menu",
	IconHome:       "home",
	IconHint:       "hint",
	IconCheck:      "check",
	IconDiffEasy:   "diff_easy",
	IconDiffMedium: "diff_medium",
	IconDiffHard:   "diff_hard",
}

// String возвращает имя кадра, например "sound_on".
func (i Icon) String() string {
	if i < 0 || int(i) >= len(iconNames) {
		return fmt.Sprintf("Icon(%d)", int(i))
	}
	return iconNames[i]
}

// Path возвращает путь к отдельному файлу иконки, например "icons/dark/gear.png".
func (i Icon) Path(theme IconTheme) string {
	return "icons/" + string(theme) + "/" + i.String() + ".png"
}

// Element — элемент интерфейса, значение — путь к файлу.
type Element string

// Панели и кнопки 8×8, растягиваются по 9-slice.
const (
	PanelRaised   Element = "ui/panel_raised.png"   // выпуклая панель
	PanelSunken   Element = "ui/panel_sunken.png"   // вдавленная панель (поле, табло)
	LedPanel      Element = "ui/led_panel.png"      // вдавленная панель с чёрным центром под цифры
	Button        Element = "ui/button.png"         // кнопка
	ButtonHover   Element = "ui/button_hover.png"   // кнопка под курсором
	ButtonPressed Element = "ui/button_pressed.png" // нажатая кнопка
)

// Чекбоксы и радиокнопки 12×12, без 9-slice.
const (
	CheckboxOff Element = "ui/checkbox_off.png"
	CheckboxOn  Element = "ui/checkbox_on.png"
	RadioOff    Element = "ui/radio_off.png"
	RadioOn     Element = "ui/radio_on.png"
)

const (
	ElementSize = 8  // панели и кнопки
	ToggleSize  = 12 // чекбоксы и радиокнопки
)

// NineSlice возвращает ширину рамки для 9-slice или 0, если элемент не растягивается.
func (e Element) NineSlice() int {
	switch e {
	case PanelRaised, PanelSunken, LedPanel:
		return 2
	case Button, ButtonHover, ButtonPressed:
		return 3
	}
	return 0
}

// Курсор.
const (
	CursorArrow    = "cursors/arrow.png"
	CursorSize     = 16
	CursorHotspotX = 0 // точка клика, отсчёт от левого верхнего угла
	CursorHotspotY = 0
)
