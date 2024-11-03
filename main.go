package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 640
	screenHeight = 480
	ballSpeed    = 3
	paddleSpeed  = 6
)

type Object struct {
	X, Y, W, H int
}

type Paddle struct {
	Object
}

type Ball struct {
	Object
	dxdt, dydt int
}

type Game struct {
	paddle           Paddle
	ball             Ball
	score, highScore int
}

func main() {
	ebiten.SetWindowTitle("Pong in golang")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	paddle := Paddle{
		Object: Object{
			X: 600,
			Y: 200,
			W: 15,
			H: 100,
		},
	}
	ball := Ball{
		Object: Object{
			X: 0,
			Y: 0,
			W: 15,
			H: 15,
		},
		dxdt: ballSpeed,
		dydt: ballSpeed,
	}

	g := &Game{
		paddle: paddle,
		ball:   ball,
	}

	err := ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	faceSource, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		log.Println("couldnt load font face")
	}

	vector.DrawFilledRect(
		screen,
		float32(g.paddle.X),
		float32(g.paddle.Y),
		float32(g.paddle.W),
		float32(g.paddle.H),
		color.White,
		false,
	)
	vector.DrawFilledRect(
		screen,
		float32(g.ball.X),
		float32(g.ball.Y),
		float32(g.ball.W),
		float32(g.ball.H),
		color.White,
		false,
	)

	scoreStr := "Score: " + fmt.Sprint(g.score)
	drawOptions := &text.DrawOptions{}
	drawOptions.ColorScale.ScaleWithColor(color.White)
	drawOptions.GeoM.Translate(10, 10)
	text.Draw(
		screen,
		scoreStr,
		&text.GoTextFace{
			Source: faceSource,
			Size:   16,
		},
		drawOptions,
	)

	highScoreStr := "High Score: " + fmt.Sprint(g.score)
	drawOptions = &text.DrawOptions{}
	drawOptions.ColorScale.ScaleWithColor(color.White)
	drawOptions.GeoM.Translate(10, 40)
	text.Draw(screen, highScoreStr, &text.GoTextFace{
		Source: faceSource,
		Size:   16,
	}, drawOptions)
}

func (g *Game) Update() error {
	g.paddle.MoveOnKeyPress()
	g.ball.Move()
	g.CollideWithWall()
	g.CollideWithPaddle()
	return nil
}

func (p *Paddle) MoveOnKeyPress() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		if p.Y >= screenHeight {
			p.Y = screenHeight
		}
		p.Y += paddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		if p.Y <= 0 {
			p.Y = paddleSpeed
		}
		p.Y -= paddleSpeed
	}
}

func (b *Ball) Move() {
	b.X += b.dxdt
	b.Y += b.dydt
}

func (g *Game) Reset() {
	g.ball.X = 0
	g.ball.Y = 0
	g.score = 0
}

func (g *Game) CollideWithWall() {
	if g.ball.X >= screenWidth {
		g.Reset()
	} else if g.ball.X <= 0 {
		g.ball.dxdt = ballSpeed
	} else if g.ball.Y <= 0 {
		g.ball.dydt = ballSpeed
	} else if g.ball.Y >= screenHeight {
		g.ball.dydt = -ballSpeed
	}
}

func (g *Game) CollideWithPaddle() {
	if g.ball.X >= g.paddle.X && g.ball.Y >= g.paddle.Y && g.ball.Y <= g.paddle.Y+g.paddle.H {
		g.ball.dxdt = -g.ball.dxdt
		g.score++
		if g.score > g.highScore {
			g.highScore = g.score
		}
	}
}
