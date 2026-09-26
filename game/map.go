package game

import (
	"fmt"
	"math/rand/v2"

	"Sapper/ui"
)

type CellState int

const (
	Closed CellState = iota
	Opened
	Flagged
	Questioned
)

type Cell struct {
	Mine   bool
	Around int
	State  CellState
}

type Map struct {
	W, H  int
	Mines int
	Cells [][]Cell // Cells[y][x]

	Lost      bool
	ExplodedX int
	ExplodedY int
	placed    bool
	opened    int
}

// New создаёт закрытое поле w×h. Мины раскидываются при первом Open,
// чтобы первый клик никогда не попадал на мину.
func New(w int, h int, mines int) *Map {
	if mines > w*h-1 {
		mines = w*h - 1
	}
	cells := make([][]Cell, h)
	for y := range cells {
		cells[y] = make([]Cell, w)
	}
	return &Map{W: w, H: h, Mines: mines, Cells: cells}
}

func (m *Map) inBounds(x, y int) bool {
	return x >= 0 && x < m.W && y >= 0 && y < m.H
}

// neighbours вызывает f для каждого соседа клетки (x, y), который есть на поле.
func (m *Map) neighbours(x, y int, f func(nx, ny int)) {
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if (dx != 0 || dy != 0) && m.inBounds(x+dx, y+dy) {
				f(x+dx, y+dy)
			}
		}
	}
}

// placeMines раскидывает мины, не трогая клетку (safeX, safeY) и её соседей,
// так что первый клик всегда открывает пустую область. Если поле для этого
// слишком плотное, свободной остаётся только сама клетка.
func (m *Map) placeMines(safeX, safeY int) {
	near := func(x, y int) bool {
		return abs(x-safeX) <= 1 && abs(y-safeY) <= 1
	}
	if m.W*m.H-9 < m.Mines {
		near = func(x, y int) bool { return x == safeX && y == safeY }
	}

	var free [][2]int
	for y := 0; y < m.H; y++ {
		for x := 0; x < m.W; x++ {
			if !near(x, y) {
				free = append(free, [2]int{x, y})
			}
		}
	}
	rand.Shuffle(len(free), func(i, j int) { free[i], free[j] = free[j], free[i] })
	for _, p := range free[:m.Mines] {
		m.Cells[p[1]][p[0]].Mine = true
	}

	m.countAround()
	m.placed = true
}

func (m *Map) countAround() {
	for y := 0; y < m.H; y++ {
		for x := 0; x < m.W; x++ {
			n := 0
			m.neighbours(x, y, func(nx, ny int) {
				if m.Cells[ny][nx].Mine {
					n++
				}
			})
			m.Cells[y][x].Around = n
		}
	}
}

// Open открывает клетку (левый клик). Пустые клетки открывают соседей по цепочке.
// Возвращает true, если под клеткой была мина.
func (m *Map) Open(x, y int) (boom bool) {
	if m.Over() || !m.inBounds(x, y) {
		return false
	}
	if !m.placed {
		m.placeMines(x, y)
	}
	c := &m.Cells[y][x]
	if c.State == Opened || c.State == Flagged {
		return false
	}
	if c.Mine {
		c.State = Opened
		m.Lost = true
		m.ExplodedX, m.ExplodedY = x, y
		return true
	}

	// Разлив пустых областей. Стек вместо рекурсии, чтобы не упираться в глубину на больших полях.
	stack := [][2]int{{x, y}}
	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		c := &m.Cells[p[1]][p[0]]
		if c.State == Opened || c.State == Flagged || c.Mine {
			continue
		}
		c.State = Opened
		m.opened++
		if c.Around == 0 {
			m.neighbours(p[0], p[1], func(nx, ny int) {
				stack = append(stack, [2]int{nx, ny})
			})
		}
	}
	return false
}

// Chord — клик по открытой цифре: если вокруг уже стоит столько же флагов,
// открывает всех остальных закрытых соседей. Возвращает true, если среди них была мина.
func (m *Map) Chord(x, y int) (boom bool) {
	if m.Over() || !m.inBounds(x, y) {
		return false
	}
	c := m.Cells[y][x]
	if c.State != Opened || c.Around == 0 {
		return false
	}
	flags := 0
	m.neighbours(x, y, func(nx, ny int) {
		if m.Cells[ny][nx].State == Flagged {
			flags++
		}
	})
	if flags != c.Around {
		return false
	}
	m.neighbours(x, y, func(nx, ny int) {
		if m.Open(nx, ny) {
			boom = true
		}
	})
	return boom
}

// Started — мины уже расставлены, то есть был первый клик.
func (m *Map) Started() bool {
	return m.placed
}

// ToggleMark переключает метку на закрытой клетке (правый клик):
// Closed → Flagged → Questioned → Closed.
func (m *Map) ToggleMark(x, y int) {
	if m.Over() || !m.inBounds(x, y) {
		return
	}
	c := &m.Cells[y][x]
	switch c.State {
	case Closed:
		c.State = Flagged
	case Flagged:
		c.State = Questioned
	case Questioned:
		c.State = Closed
	}
}

// Won — все клетки без мин открыты.
func (m *Map) Won() bool {
	return !m.Lost && m.opened == m.W*m.H-m.Mines
}

// Over — игра закончена, победой или поражением.
func (m *Map) Over() bool {
	return m.Lost || m.Won()
}

// MinesLeft — число для счётчика мин: мины минус флаги. Может быть меньше нуля.
func (m *Map) MinesLeft() int {
	n := m.Mines
	for y := range m.Cells {
		for _, c := range m.Cells[y] {
			if c.State == Flagged {
				n--
			}
		}
	}
	return n
}

// Tile возвращает спрайт для клетки (x, y) с учётом конца игры.
func (m *Map) Tile(x, y int) ui.Tile {
	c := m.Cells[y][x]
	if m.Lost {
		switch {
		case c.Mine && x == m.ExplodedX && y == m.ExplodedY:
			return ui.TileMineExploded
		case c.Mine && c.State != Flagged:
			return ui.TileMine
		case !c.Mine && c.State == Flagged:
			return ui.TileMineWrong
		}
	}
	if m.Won() && c.Mine {
		return ui.TileFlag
	}
	switch c.State {
	case Opened:
		if c.Around == 0 {
			return ui.TileOpen
		}
		return ui.TileNumber(c.Around)
	case Flagged:
		return ui.TileFlag
	case Questioned:
		return ui.TileQuestion
	}
	return ui.TileClosed
}

// Print выводит поле в консоль целиком, с минами: * — мина, . — пусто, цифра — мин вокруг.
func (m *Map) Print() {
	for y := m.H - 1; y >= 0; y-- {
		for x := 0; x < m.W; x++ {
			c := m.Cells[y][x]
			switch {
			case c.Mine:
				fmt.Print("* ")
			case c.Around == 0:
				fmt.Print(". ")
			default:
				fmt.Print(c.Around, " ")
			}
		}
		fmt.Println()
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
