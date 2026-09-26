package ui

// Anim — id анимации в карте из helper.Build_sprites.
// Длительности кадров и Loop берутся из atlas.json.
type Anim string

// animations.png
const (
	AnimReveal      Anim = "reveal"       // открытие клетки, один раз
	AnimFlagPlant   Anim = "flag_plant"   // установка флага, один раз
	AnimFlagWave    Anim = "flag_wave"    // флаг развевается, цикл
	AnimQuestionBob Anim = "question_bob" // знак вопроса покачивается, цикл
	AnimMineSpin    Anim = "mine_spin"    // мина крутится, цикл
	AnimExplosion   Anim = "explosion"    // взрыв, один раз; последний кадр как tiles/mine_exploded
	AnimFaceBlink   Anim = "face_blink"   // смайлик моргает, цикл (26×26)
	AnimFacePress   Anim = "face_press"   // смайлик при нажатии, один раз (26×26)
)

// shimmer.png, glint: по клетке проходит полоска света, потом пауза. Все зациклены.
const (
	AnimGlintClosed   Anim = "glint_closed"
	AnimGlintFlag     Anim = "glint_flag"
	AnimGlintQuestion Anim = "glint_question"
	AnimGlintMine     Anim = "glint_mine"
	AnimGlintFaceCool Anim = "glint_face_cool" // 26×26
	AnimGlintDigit8   Anim = "glint_digit_8"   // 13×23
)

// GlintWaveStep — сдвиг блика между соседними клетками, чтобы он шёл волной:
// кадр для клетки (col, row) = t - (col+row)*GlintWaveStep.
const GlintWaveStep = 8

// shimmer.png, wave: непрерывные диагональные полосы с периодом 16 px (= TileSize),
// поэтому все клетки могут показывать один и тот же кадр.
const (
	AnimWaveClosed   Anim = "wave_closed"
	AnimWaveFlag     Anim = "wave_flag"
	AnimWaveQuestion Anim = "wave_question"
)
