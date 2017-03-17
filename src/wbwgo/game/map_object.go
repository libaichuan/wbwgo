package game

import (
	"log"
	"wbwgo/misc"
)

type MapObjectBase struct {
	id int64

	x, y float32

	speedx, speedy float32
}

func (s *MapObjectBase) ComputeSpeed(targe_x float32, targe_y float32, speed float32) {
	d_x := targe_x - s.x
	d_y := targe_y - s.y

	ticks := misc.InvSqrt(d_x*d_x+d_y*d_y) / speed
	ticks *= 10

	s.speedx = d_x / ticks
	s.speedy = d_y / ticks
}

func (s *MapObjectBase) Move() {
	s.x += s.speedx
	s.y += s.speedy
}

func (s *MapObjectBase) Print() {
	log.Printf("x:%f,y:%f", s.x, s.y)
}
