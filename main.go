package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

// :globals constants
const (
	MAX_ENTITY_COUNT               = 1024
	MAX_HAND_COUNT                 = 40
	entitySelectionRadius  float32 = 10.0
	PLAYER_HEALTH                  = 10
	TROLL_HEALTH                   = 3
	GOBLIN_HEALTH                  = 3
	PLAYER_MOVEMENT_RADIUS float32 = 1
)

// :enum EntityArchType
type EntityArchType int

const (
	ARCH_NIL           EntityArchType = 0
	ARCH_TROLL         EntityArchType = 1
	ARCH_GOBLIN        EntityArchType = 2
	ARCH_PLAYER        EntityArchType = 3
	ARCH_CARD_FIREBALL EntityArchType = 4
)

type Sprite struct {
	Image rl.Texture2D
}

// :enum SpriteId
type SpriteId int

const (
	SPRITE_NIL SpriteId = iota
	SPRITE_PLAYER
	SPRITE_GOBLIN
	SPRITE_TROLL
	SPRITE_CARD_FIREBALL
	SPRITE_MAX
)

type Entity struct {
	Position rl.Vector2
	IsValid  bool
	Type     EntityArchType
	SpriteId SpriteId
	Health   int
}

type Hand struct {
	Cards [MAX_HAND_COUNT]Entity
}

type World struct {
	Entities [MAX_ENTITY_COUNT]Entity
}

type WorldFrame struct {
	SelectedEntity *Entity
}

// :globals structs
var worldFrame WorldFrame
var world *World = nil
var hand *Hand = nil

// :helpers engine functions
func boolToInt(x bool) int32 {
	if x {
		return 1
	} else {
		return 0
	}
}

func assert(condition bool, error string) {
	if condition == false {
		fmt.Print(error)
		panic("")
	}
}

func almostEquals(a, b, epsilon float32) bool {
	return float32(math.Abs(float64(a-b))) <= epsilon
}

func animateF32ToTarget(value *float32, target, delta_t, rate float32) bool {
	*value += (target - *value) * (1.0 - float32(math.Pow(2.0, float64(-rate*delta_t))))
	if almostEquals(*value, target, 0.001) {
		*value = target
		return true
	}

	return false
}

func animateV2ToTarget(value *rl.Vector2, target rl.Vector2, delta_t, rate float32) {
	animateF32ToTarget(&value.X, target.X, delta_t, rate)
	animateF32ToTarget(&value.Y, target.Y, delta_t, rate)
}

const tileWidth int32 = 8

func worldPositionToTilePosition(worldPosition float32) float32 {
	return float32(math.Round(float64(worldPosition) / float64(tileWidth)))
}

func tilePositionToWorldPosition(tilePosition float32) float32 {
	return float32(tilePosition) * float32(tileWidth)
}

func roundV2ToTile(worldPosition rl.Vector2) rl.Vector2 {
	worldPosition.X = tilePositionToWorldPosition(worldPositionToTilePosition(worldPosition.X))
	worldPosition.Y = tilePositionToWorldPosition(worldPositionToTilePosition(worldPosition.Y))
	return worldPosition
}

// :helper game functions
func createCardInHand() *Entity {
	var entityFound *Entity = nil
	for i := 0; i < MAX_HAND_COUNT; i++ {
		var existingEntity *Entity = &(hand.Cards[i])
		if !existingEntity.IsValid {
			entityFound = existingEntity
			break

		}
	}
	// :TODO assert here
	assert(entityFound != nil, "max # of cards in hand reached")
	entityFound.IsValid = true
	return entityFound
}

func createEntity() *Entity {
	var entityFound *Entity = nil
	for i := 0; i < MAX_ENTITY_COUNT; i++ {

		var existingEntity *Entity = &world.Entities[i]
		if !existingEntity.IsValid {

			entityFound = existingEntity
			break
		}
	}

	assert(entityFound != nil, "max # of entities reached")
	entityFound.IsValid = true
	return entityFound
}

func destroyEntity(en *Entity) {
	en = nil
}

func setupTroll(en *Entity) {
	en.Type = ARCH_TROLL
	en.SpriteId = SPRITE_TROLL
	en.Health = TROLL_HEALTH
}

func setupPlayer(en *Entity) {
	en.Type = ARCH_PLAYER
	en.SpriteId = SPRITE_PLAYER
	en.Health = TROLL_HEALTH
}

func setupGoblin(en *Entity) {
	en.Type = ARCH_GOBLIN
	en.SpriteId = SPRITE_GOBLIN
	en.Health = GOBLIN_HEALTH
}

var sprites [SPRITE_MAX]Sprite

func getSprite(id SpriteId) *Sprite {
	if id >= 0 && id < SPRITE_MAX {
		return &sprites[id]
	}
	return &sprites[0]
}

func setupCardFireball(en *Entity) {
	en.Type = ARCH_CARD_FIREBALL
	en.SpriteId = SPRITE_CARD_FIREBALL
}

/**/
func main() {
	rl.SetConfigFlags(rl.FlagVsyncHint | rl.FlagWindowHighdpi)

	const screenWidth int32 = 800
	const screenHeight int32 = 450

	rl.InitWindow(screenWidth, screenHeight, "Dueling Monsters")

	//setup globals

	world = &World{}
	assert(world != nil, "world not correctly initialized")
	hand = &Hand{}
	assert(hand != nil, "hand not correctly initialized")

	// initalze t
	var cardFireballTest *Entity = createEntity()
	setupCardFireball(cardFireballTest)

	sprites[SPRITE_PLAYER] = Sprite{Image: rl.LoadTexture("./resources/player.png")}
	sprites[SPRITE_GOBLIN] = Sprite{Image: rl.LoadTexture("./resources/goblin.png")}
	sprites[SPRITE_TROLL] = Sprite{Image: rl.LoadTexture("./resources/troll.png")}
	sprites[SPRITE_CARD_FIREBALL] =
		Sprite{Image: rl.LoadTexture("./resources/troll.png")}

	for i := 0; i < 2; i++ {
		var en *Entity = createEntity()
		setupTroll(en)
		var en2 *Entity = createEntity()
		setupGoblin(en2)
		en.Position = rl.Vector2{X: float32(rl.GetRandomValue(0, 200)), Y: float32(rl.GetRandomValue(0, 200))}
		en.Position = roundV2ToTile(en.Position)

		en2.Position = rl.Vector2{X: float32(rl.GetRandomValue(0, 200)), Y: float32(rl.GetRandomValue(0, 200))}
		en2.Position = roundV2ToTile(en2.Position)
	}

	var playerEntity *Entity = createEntity()
	setupPlayer(playerEntity)

	playerEntity.Position = rl.Vector2{X: 0, Y: 0}

	// :camera initialze

	var camera rl.Camera2D = rl.Camera2D{}

	camera.Zoom = 3.0
	camera.Offset = rl.Vector2{X: float32(float32(screenWidth) / 2.0), Y: float32(float32(screenHeight) / 2.0)}
	camera.Rotation = 0
	camera.Target = rl.Vector2{
		X: playerEntity.Position.X + (float32(sprites[SPRITE_PLAYER].Image.Width) / 2.0),
		Y: playerEntity.Position.Y + (float32(sprites[SPRITE_PLAYER].Image.Height) / 2.0),
	}

	// :input movement variables
	var inputAxis rl.Vector2 = rl.Vector2{X: 0, Y: 0}
	var terminalPoint rl.Vector2

	/// remove bellow
	rl.SetTargetFPS(60)

	defer rl.CloseWindow()

	for !rl.WindowShouldClose() {
		fmt.Println(camera.Target)
		worldFrame = WorldFrame{}

		// :input
		{

			if rl.Vector2Distance(playerEntity.Position, terminalPoint) <
				PLAYER_MOVEMENT_RADIUS {
				inputAxis = rl.Vector2{X: 0, Y: 0}
				terminalPoint = rl.Vector2{X: 0, Y: 0}
			}
		}

		// :camera
		{
			var target rl.Vector2 = playerEntity.Position
			target.X = target.X + (float32(sprites[SPRITE_PLAYER].Image.Width))/2.0
			target.Y = target.Y + (float32(sprites[SPRITE_PLAYER].Image.Height))/2.0
			animateV2ToTarget(&camera.Target, target, rl.GetFrameTime(), 30.0)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.LightGray)

		// :world entities positions
		{
			rl.BeginMode2D(camera)
			var mousePositionWorld = rl.GetScreenToWorld2D(rl.GetMousePosition(), camera)

			// var mouseTileX int = int(worldPositionToTilePosition(mousePositionWorld.X))
			// var mouseTileY int = int(worldPositionToTilePosition(mousePositionWorld.Y))

			// :tile rendering
			{

				var playerTileX = int32(worldPositionToTilePosition(playerEntity.Position.X))
				var playerTileY = int32(worldPositionToTilePosition(playerEntity.Position.Y))

				const tileRadiusX int32 = 40
				const tileRadiusY int32 = 30

				for x := playerTileX - tileRadiusX; x < playerTileX+tileRadiusX; x++ {
					for y := playerTileY - tileRadiusY; y < playerTileY+tileRadiusY; y++ {
						if (x+boolToInt(y%2 == 0))%2 == 0 {

							var xPosition float32 = float32(x * tileWidth)
							var yPosition float32 = float32(y * tileWidth)
							var tileColor rl.Color = rl.White

							rl.DrawRectangle(int32(xPosition)+int32(float32(tileWidth)*-0.5), int32(yPosition+(float32(tileWidth)*-0.5)), tileWidth, tileWidth, tileColor)
						}
					}
				}

			}

			// :mouse :selector
			{
				var smallestDistance float32 = math.MaxFloat32

				for i := 0; i < MAX_ENTITY_COUNT; i++ {
					var en *Entity = &world.Entities[i]
					if en.IsValid {
						// var sprite *Sprite = getSprite(en.SpriteId)
						var distance float32 = float32(math.Abs(float64(rl.Vector2Distance(en.Position, mousePositionWorld))))
						if distance < entitySelectionRadius {
							if worldFrame.SelectedEntity == nil || (distance < smallestDistance) {
								worldFrame.SelectedEntity = en
								smallestDistance = distance
							}
						}
					}
				}

			}

			// :click handler
			{

				var isMouseLeftPressed bool = rl.IsMouseButtonPressed(rl.MouseButtonLeft)
				var selectedEntity *Entity = worldFrame.SelectedEntity

				if isMouseLeftPressed {
					isMouseLeftPressed = false
					if selectedEntity != nil {
						selectedEntity.Health -= 1
						if selectedEntity.Health <= 0 {
							destroyEntity(selectedEntity)
						}
					}

				} else if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
					terminalPoint = mousePositionWorld
					inputAxis = rl.Vector2Subtract(terminalPoint, playerEntity.Position)
				}

			}

			inputAxis = rl.Vector2Normalize(inputAxis)

			fmt.Println(inputAxis)

			playerEntity.Position = rl.Vector2Add(playerEntity.Position, rl.Vector2Scale(inputAxis, 100*rl.GetFrameTime()))

			// :render
			{

				for i := 0; i < MAX_ENTITY_COUNT; i++ {
					var entity *Entity = &world.Entities[i]
					if entity.IsValid {
						switch entity.Type {

						default:

							var sprite *Sprite = getSprite(entity.SpriteId)
							var entityColor rl.Color = rl.White
							if worldFrame.SelectedEntity == entity {
								entityColor = rl.Red
							}
							rl.DrawTexture(sprite.Image, int32(entity.Position.X-float32(sprite.Image.Width/2)), int32(entity.Position.Y-float32(sprite.Image.Height/2)), entityColor)

						}
					}

				}

			}

			rl.EndMode2D()
		}

		rl.EndDrawing()
	}

	rl.UnloadTexture(sprites[SPRITE_TROLL].Image)
	rl.UnloadTexture(sprites[SPRITE_PLAYER].Image)
	rl.UnloadTexture(sprites[SPRITE_GOBLIN].Image)

}
