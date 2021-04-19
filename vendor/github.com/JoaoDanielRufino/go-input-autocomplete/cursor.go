package input_autocomplete

import "fmt"

type Cursor struct {
	position int
}

func NewCursor() *Cursor {
	return &Cursor{
		position: 0,
	}
}

func (c *Cursor) canMoveLeft(n int) bool {
	return c.position-n >= 0
}

func (c *Cursor) IncrementPosition() {
	c.position++
}

func (c *Cursor) SetPosition(newPos int) {
	c.position = newPos
}

func (c *Cursor) MoveRight() {
	fmt.Print("\033[1C")
	c.position++
}

func (c *Cursor) MoveLeft() {
	if c.canMoveLeft(1) {
		fmt.Print("\033[1D")
		c.position--
	}
}

func (c *Cursor) MoveLeftNPos(x int) {
	if c.canMoveLeft(x) && x > 0 {
		fmt.Printf("\033[%dD", x)
		c.position -= x
	}
}

func (c *Cursor) GetPosition() int {
	return c.position
}
