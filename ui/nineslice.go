package ui

import "github.com/gopxl/pixel/v2"

// Nine_slice рисует панель или кнопку любого размера из маленькой картинки:
// углы как есть, стороны растягиваются в одну сторону, центр — в обе.
type Nine_slice struct {
	border float64
	parts  [3][3]*pixel.Sprite // [ряд снизу вверх][колонка слева направо]
}

func New_nine_slice(pic pixel.Picture, border int) *Nine_slice {
	b := float64(border)
	r := pic.Bounds()
	xs := [4]float64{r.Min.X, r.Min.X + b, r.Max.X - b, r.Max.X}
	ys := [4]float64{r.Min.Y, r.Min.Y + b, r.Max.Y - b, r.Max.Y}
	n := &Nine_slice{border: b}
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			n.parts[row][col] = pixel.NewSprite(pic, pixel.R(xs[col], ys[row], xs[col+1], ys[row+1]))
		}
	}
	return n
}

// Draw рисует панель так, чтобы она заняла прямоугольник r.
func (n *Nine_slice) Draw(t pixel.Target, r pixel.Rect) {
	b := n.border
	xs := [4]float64{r.Min.X, r.Min.X + b, r.Max.X - b, r.Max.X}
	ys := [4]float64{r.Min.Y, r.Min.Y + b, r.Max.Y - b, r.Max.Y}
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			part := n.parts[row][col]
			dst := pixel.R(xs[col], ys[row], xs[col+1], ys[row+1])
			src := part.Frame()
			scale := pixel.V(dst.W()/src.W(), dst.H()/src.H())
			part.Draw(t, pixel.IM.ScaledXY(pixel.ZV, scale).Moved(dst.Center()))
		}
	}
}
