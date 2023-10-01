package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	_ "image/png"
)

type scrollDemo struct {
	player          *ebiten.Image
	xloc            int
	yloc            int
	background      *ebiten.Image
	backgroundXView int
	eggPict         *ebiten.Image
	eggs            []Shot
}

type Shot struct {
	pict   *ebiten.Image
	xShot  int
	yShot  int
	deltaX int
}

func PlayerInput(demo *scrollDemo) {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && demo.yloc > 0 {
		demo.yloc -= 4
	}
	//window height is 1000 pixels and dragon is 100 pixels
	if ebiten.IsKeyPressed(ebiten.KeyDown) && demo.yloc < 900 {
		demo.yloc += 4
	}

	//projectile
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		newEgg := NewProjectile(demo.eggPict, demo)
		demo.eggs = append(demo.eggs, newEgg)
		updateShots(demo)
	}
}

func NewProjectile(picture *ebiten.Image, demo *scrollDemo) Shot {
	return Shot{
		pict:   picture,
		xShot:  int(demo.xloc + 100),
		yShot:  int(demo.yloc),
		deltaX: 8,
	}
}

func updateShots(demo *scrollDemo) {
	for i := 0; i < len(demo.eggs); i++ {
		demo.eggs[i].xShot += demo.eggs[i].deltaX
		//remove shots off-screen here
	}
}

func (demo *scrollDemo) Update() error {
	//background scroll
	backgroundWidth := demo.background.Bounds().Dx()
	maxX := backgroundWidth * 2
	demo.backgroundXView -= 4
	demo.backgroundXView %= maxX

	//player input
	PlayerInput(demo)

	//update projectiles
	updateShots(demo)
	return nil
}

func (demo *scrollDemo) Draw(screen *ebiten.Image) {
	drawOps := ebiten.DrawImageOptions{}
	//draw background
	const repeat = 5
	backgroundWidth := demo.background.Bounds().Dx()
	for count := 0; count < repeat; count += 1 {
		drawOps.GeoM.Reset()
		drawOps.GeoM.Translate(float64(backgroundWidth*count),
			float64(-1000))
		drawOps.GeoM.Translate(float64(demo.backgroundXView), 0)
		screen.DrawImage(demo.background, &drawOps)
	}

	//draw player
	drawOps.GeoM.Reset()
	drawOps.GeoM.Translate(float64(demo.xloc), float64(demo.yloc))
	screen.DrawImage(demo.player, &drawOps)

	//draw shots
	for _, shot := range demo.eggs {
		drawOps.GeoM.Reset()
		drawOps.GeoM.Translate(float64(shot.xShot), float64(shot.yShot))
		screen.DrawImage(shot.pict, &drawOps)
	}
}

func (s scrollDemo) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("Scroller Example")
	//New image from file returns image as image.Image (_) and ebiten.Image
	//background image
	backgroundPict, _, err := ebitenutil.NewImageFromFile("background.png")
	if err != nil {
		fmt.Println("Unable to load background image:", err)
	}
	//player image
	playerPict, _, err := ebitenutil.NewImageFromFile("dragon.png")
	if err != nil {
		fmt.Println("Unable to load player image:", err)
	}
	//egg image
	eggPict, _, err := ebitenutil.NewImageFromFile("EggBlue.png")
	if err != nil {
		fmt.Println("Unable to load egg projectile image:", err)
	}
	demo := scrollDemo{
		player:     playerPict,
		background: backgroundPict,
		eggPict:    eggPict,
		xloc:       0,
	}
	err = ebiten.RunGame(&demo)
	if err != nil {
		fmt.Println("Failed to run game", err)
	}
}
