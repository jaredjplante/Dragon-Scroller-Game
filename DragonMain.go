package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	_ "github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	_ "image/png"
	"os"
)

type scrollDemo struct {
	player          *ebiten.Image
	xloc            int
	yloc            int
	background      *ebiten.Image
	backgroundXView int
	eggPict         *ebiten.Image
	eggs            []Shot
	popSound        sound
	shatterSound    sound
}

type Shot struct {
	pict   *ebiten.Image
	xShot  int
	yShot  int
	deltaX int
}

type sound struct {
	audioContext *audio.Context
	soundPlayer  *audio.Player
}

const (
	WINDOW_WIDTH      = 1000
	WINDOW_HEIGHT     = 1000
	DRAGON_WIDTH      = 100
	SOUND_SAMPLE_RATE = 48000
)

func PlayerInput(demo *scrollDemo) {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && demo.yloc > 0 {
		demo.yloc -= 4
	}
	//window height is 1000 pixels and dragon is 100 pixels
	if ebiten.IsKeyPressed(ebiten.KeyDown) && demo.yloc < WINDOW_WIDTH-DRAGON_WIDTH {
		demo.yloc += 4
	}

	//projectile
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		newEgg := NewProjectile(demo.eggPict, demo)
		demo.eggs = append(demo.eggs, newEgg)
		updateShots(demo)
		//play egg sound
		demo.popSound.soundPlayer.Rewind()
		demo.popSound.soundPlayer.Play()
	}
}

func NewProjectile(picture *ebiten.Image, demo *scrollDemo) Shot {
	return Shot{
		pict:   picture,
		xShot:  int(demo.xloc + DRAGON_WIDTH),
		yShot:  int(demo.yloc),
		deltaX: 8,
	}
}

func updateShots(demo *scrollDemo) {
	for i := 0; i < len(demo.eggs); i++ {
		demo.eggs[i].xShot += demo.eggs[i].deltaX
		//shift elements to remove projectile off-screen
		if demo.eggs[i].xShot > WINDOW_WIDTH {
			demo.eggs = append(demo.eggs[:i], demo.eggs[i+1:]...)
			i--
		}
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
	ebiten.SetWindowSize(WINDOW_WIDTH, WINDOW_HEIGHT)
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

	//handle sound
	soundContext := audio.NewContext(SOUND_SAMPLE_RATE)
	popSound := sound{
		audioContext: soundContext,
		soundPlayer:  LoadWav("pop.wav", soundContext),
	}
	shatterSound := sound{
		audioContext: soundContext,
		soundPlayer:  LoadWav("shatter.wav", soundContext),
	}

	//setup game and run
	demo := scrollDemo{
		player:       playerPict,
		background:   backgroundPict,
		eggPict:      eggPict,
		popSound:     popSound,
		shatterSound: shatterSound,
		xloc:         0,
	}
	err = ebiten.RunGame(&demo)
	if err != nil {
		fmt.Println("Failed to run game", err)
	}
}

func LoadWav(name string, context *audio.Context) *audio.Player {
	File, err := os.Open(name)
	if err != nil {
		fmt.Println("Error Loading sound: ", err)
	}
	Sound, err := wav.DecodeWithoutResampling(File)
	if err != nil {
		fmt.Println("Error interpreting sound file: ", err)
	}
	Player, err := context.NewPlayer(Sound)
	if err != nil {
		fmt.Println("Couldn't create sound player: ", err)
	}
	return Player
}
