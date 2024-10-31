package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// :globals constants
const (
	MAX_ENTITY_COUNT               = 1024
	MAX_HAND_COUNT                 = 40
	entitySelectionRadius  float32 = 10.0
	PLAYER_HEALTH                  = 100
	TROLL_HEALTH                   = 10
	GOBLIN_HEALTH                  = 10
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
	ARCH_CARD          EntityArchType = 5
	ARCH_ATTACK        EntityArchType = 6
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
	SPRITE_ATTACK_FIREBALL
	SPRITE_ATTACK_BASIC
	SPRITE_ATTACK_SWORD
	SPRITE_MAX
)

type Entity struct {
	Position           rl.Vector2
	IsValid            bool
	Type               EntityArchType
	SpriteId           SpriteId
	Health             int32
	inputAxis          rl.Vector2
	CollisionRectangle rl.Rectangle

	// for cards
	Range  int32
	Width  int32
	Damage int32
	Speed  int32

	// for attacks
	MaxPosition  rl.Vector2
	Angle        float32
	isMelee      bool
	isProjectile bool
}

type Card struct {
	base Entity
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
var sprites [SPRITE_MAX]Sprite

// :helpers :engine functions
func IsInRange(entity1 *Entity, entity2 *Entity) bool {
	distance := rl.Vector2Distance(entity1.Position, entity2.Position)
	return distance <= float32(entity1.Range)
}

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
func getSprite(id SpriteId) *Sprite {
	if id >= 0 && id < SPRITE_MAX {
		return &sprites[id]
	}
	return &sprites[0]
}

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

	entityFound.inputAxis = rl.Vector2{X: 0, Y: 0}
	return entityFound
}

func destroyEntity(en *Entity) {
	*en = Entity{}
	en.IsValid = false
}

func setupTroll(en *Entity, position *rl.Vector2) {
	en.Type = ARCH_TROLL
	en.SpriteId = SPRITE_TROLL
	en.Health = TROLL_HEALTH
	en.Damage = 30
	en.Speed = 50
	en.Range = 100

	if position != nil {
		en.Position = *position
	} else {
		en.Position = rl.Vector2{X: 0, Y: 0}

	}

	var sprite *Sprite = getSprite(en.SpriteId)
	en.CollisionRectangle.X = (en.Position.X - float32(sprite.Image.Width)/2)
	en.CollisionRectangle.Y = en.Position.Y - float32(sprite.Image.Height/2)
	en.CollisionRectangle.Width = float32(sprite.Image.Width)
	en.CollisionRectangle.Height = float32(sprite.Image.Height)
}

func setupPlayer(en *Entity, position *rl.Vector2) {
	en.Type = ARCH_PLAYER
	en.SpriteId = SPRITE_PLAYER
	en.Health = PLAYER_HEALTH
	en.Speed = 100

	if position != nil {
		en.Position = *position
	} else {
		en.Position = rl.Vector2{X: 0, Y: 0}

	}

	var sprite *Sprite = getSprite(en.SpriteId)
	en.CollisionRectangle.X = (en.Position.X - float32(sprite.Image.Width)/2)
	en.CollisionRectangle.Y = en.Position.Y - float32(sprite.Image.Height/2)
	en.CollisionRectangle.Width = float32(sprite.Image.Width)
	en.CollisionRectangle.Height = float32(sprite.Image.Height)
}

func setupGoblin(en *Entity, position *rl.Vector2) {
	en.Type = ARCH_GOBLIN
	en.SpriteId = SPRITE_GOBLIN
	en.Health = GOBLIN_HEALTH
	en.Damage = 10
	en.Speed = 50
	en.Range = 100

	if position != nil {
		en.Position = *position
	} else {
		en.Position = rl.Vector2{X: 0, Y: 0}

	}

	var sprite *Sprite = getSprite(en.SpriteId)
	en.CollisionRectangle.X = (en.Position.X - float32(sprite.Image.Width)/2)
	en.CollisionRectangle.Y = en.Position.Y - float32(sprite.Image.Height/2)
	en.CollisionRectangle.Width = float32(sprite.Image.Width)
	en.CollisionRectangle.Height = float32(sprite.Image.Height)
}

func setupCardFireball(en *Entity) {
	en.Type = ARCH_CARD
	en.SpriteId = SPRITE_CARD_FIREBALL
	en.Range = 100
	en.Width = 5
	en.Damage = 2
	en.Health = 1
}

func setupAttackFireball(en *Entity) {
	en.Type = ARCH_ATTACK
	en.SpriteId = SPRITE_ATTACK_FIREBALL
	en.Damage = 4
	en.Speed = 200
	en.Health = 1
	en.Range = 100
	en.isProjectile = true
}

func setupAttackBasic(en *Entity) {
	en.Type = ARCH_ATTACK
	en.SpriteId = SPRITE_ATTACK_BASIC
	en.Damage = 1
	en.Speed = 150
	en.Health = 1
	en.Range = 85
	en.isProjectile = true
}

func setupAttackSword(en *Entity, position, maxPosition rl.Vector2) {
	en.Type = ARCH_ATTACK
	en.SpriteId = SPRITE_ATTACK_SWORD
	en.Damage = 3
	en.Speed = 40
	en.Health = 1
	en.Range = 85
	en.isProjectile = true
	en.MaxPosition = maxPosition

}

func main() {

	rl.SetConfigFlags(rl.FlagVsyncHint | rl.FlagWindowHighdpi)

	const screenWidth int32 = 800
	const screenHeight int32 = 450

	var spawnTrollRate float32 = 4
	var spawnGlobinRate float32 = 2
	var spawnTrollPosition rl.Vector2 = rl.Vector2{X: 30, Y: 40}
	var spawnGoblinPosition rl.Vector2 = rl.Vector2{X: 20, Y: 20}
	var spawnTrollAmount int32 = 1
	var spawnGoblinAmount int32 = 2
	//var maxEnemyToSpawn int32 = 12
	var elapsedTimeGoblin float32 = 0
	var elapsedTimeTroll float32 = 0

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
		Sprite{Image: rl.LoadTexture("./resources/card_fireball.png")}
	sprites[SPRITE_ATTACK_FIREBALL] =
		Sprite{Image: rl.LoadTexture("./resources/attack_fireball.png")}
	sprites[SPRITE_ATTACK_BASIC] =
		Sprite{Image: rl.LoadTexture("./resources/basic_attack.png")}
	sprites[SPRITE_ATTACK_SWORD] =
		Sprite{Image: rl.LoadTexture("./resources/sword.png")}

	/* for i := 0; i < 2; i++ { */
	/* 	var en *Entity = createEntity() */
	/* 	setupTroll(en) */
	/* 	var en2 *Entity = createEntity() */
	/* 	setupGoblin(en2) */
	/* 	en.Position = rl.Vector2{X: float32(rl.GetRandomValue(0, 200)), Y: float32(rl.GetRandomValue(0, 200))} */
	/* 	en.Position = roundV2ToTile(en.Position) */
	/**/
	/* 	en2.Position = rl.Vector2{X: float32(rl.GetRandomValue(0, 200)), Y: float32(rl.GetRandomValue(0, 200))} */
	/* 	en2.Position = roundV2ToTile(en2.Position) */
	/* } */

	var playerEntity *Entity = createEntity()

	setupPlayer(playerEntity, &rl.Vector2{X: 0, Y: 0})

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

	/// remove bellow
	rl.SetTargetFPS(60)

	defer rl.CloseWindow()

	var runningMultiplier float32 = 1

	// grabbed entity
	var grabbedEntity *Entity = nil

	for !rl.WindowShouldClose() {
		// :clean :reset

		worldFrame = WorldFrame{}
		var delta_t float32 = rl.GetFrameTime()
		elapsedTimeGoblin += delta_t
		elapsedTimeTroll += delta_t

		// :input
		{
			playerEntity.inputAxis = rl.Vector2{X: 0, Y: 0}
			if rl.IsKeyDown(rl.KeyLeftShift) {
				runningMultiplier = 1.5
			} else {
				runningMultiplier = 1
			}

			if rl.IsKeyDown(rl.KeyS) {
				playerEntity.inputAxis.Y += 1
			}
			if rl.IsKeyDown(rl.KeyW) {
				playerEntity.inputAxis.Y -= 1
			}
			if rl.IsKeyDown(rl.KeyD) {
				playerEntity.inputAxis.X += 1
			}
			if rl.IsKeyDown(rl.KeyA) {
				playerEntity.inputAxis.X -= 1
			}

		}
		// :spawn Enemies
		{
			if elapsedTimeGoblin >= spawnGlobinRate {
				for i := int32(0); i < spawnGoblinAmount; i++ {
					var goblinEntity *Entity = createEntity()
					var goblinPosition = (rl.Vector2AddValue(spawnGoblinPosition, float32(i*20)))
					setupGoblin(goblinEntity, &goblinPosition)
				}
				elapsedTimeGoblin = 0
			}

			if elapsedTimeTroll >= spawnTrollRate {
				for i := int32(0); i < spawnTrollAmount; i++ {
					var trollEntity *Entity = createEntity()
					var trollPosition = (rl.Vector2AddValue(spawnTrollPosition, float32(i*20)))
					setupTroll(trollEntity, &trollPosition)
				}
				elapsedTimeTroll = 0
			}
		}

		// :camera
		{
			var target rl.Vector2 = playerEntity.Position
			var sprite = getSprite(playerEntity.SpriteId)
			target.X = target.X + (float32(sprite.Image.Width))/2.0
			target.Y = target.Y + (float32(sprite.Image.Height))/2.0
			animateV2ToTarget(&camera.Target, target, delta_t, 30.0)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.LightGray)

		// :world entities positions
		{
			rl.BeginMode2D(camera)
			var mousePositionScreen rl.Vector2 = rl.GetMousePosition()
			var mousePositionWorld rl.Vector2 = rl.GetScreenToWorld2D(mousePositionScreen, camera)

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

			// :mouse :click handler
			{

				var top20Percent float32 = float32(screenHeight) - (float32(screenHeight) * 0.20)

				var IsMouseButtonLeftDown bool = rl.IsMouseButtonDown(rl.MouseButtonLeft)
				var IsMouseButtonLeftRelased bool = rl.IsMouseButtonReleased(rl.MouseButtonLeft)
				var IsMouseButtonLeftPressed bool = rl.IsMouseButtonPressed(rl.MouseButtonLeft)
				var selectedEntity *Entity = worldFrame.SelectedEntity

				if IsMouseButtonLeftDown {
					IsMouseButtonLeftDown = false
					if selectedEntity != nil && selectedEntity.Type == ARCH_CARD {
						grabbedEntity = selectedEntity
					}

					if selectedEntity != nil {
						if mousePositionScreen.Y < top20Percent {
							// :TODO show attack direction league like
						}
					}

				} else if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
					/* inputAxis = rl.Vector2Subtract(terminalPoint, playerEntity.Position) */
				}
				if IsMouseButtonLeftPressed {
					IsMouseButtonLeftPressed = false
					/* var basicAttack *Entity = createEntity() */
					/* setupAttackBasic(basicAttack) */
					/* basicAttack.Position = playerEntity.Position */
					/* basicAttack.inputAxis = rl.Vector2Normalize((rl.Vector2Subtract(mousePositionWorld, playerEntity.Position))) */
					/* basicAttack.MaxPosition = rl.Vector2AddValue(playerEntity.Position, float32(basicAttack.Range)) */
					var swordAttack *Entity = createEntity()

					// Calculate direction vector
					deltaX := float64(mousePositionWorld.X - playerEntity.Position.X)
					deltaY := float64(mousePositionWorld.Y - playerEntity.Position.Y)

					// Calculate angle using arctangent
					angle := float32(math.Atan2(deltaY, deltaX))

					swordAttack.MaxPosition = rl.Vector2AddValue(playerEntity.Position, float32(swordAttack.Range))
					swordAttack.Angle = angle
					setupAttackSword(swordAttack, playerEntity.Position, swordAttack.MaxPosition)
					var playerSprite *Sprite = getSprite(playerEntity.SpriteId)
					var swingDirectionStart rl.Vector2 = rl.Vector2Normalize(rl.NewVector2(float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle)))))
					swordAttack.inputAxis = rl.Vector2{X: swingDirectionStart.Y, Y: -swingDirectionStart.X}

					swordAttack.Position = rl.Vector2{
						X: (playerEntity.Position.X + swingDirectionStart.X*float32(playerSprite.Image.Width)) + (swordAttack.inputAxis.X * 10),
						Y: playerEntity.Position.Y + swingDirectionStart.Y*float32(playerSprite.Image.Height) + (swordAttack.inputAxis.Y * 10)}

					swordAttack.MaxPosition = rl.Vector2AddValue(playerEntity.Position, float32(swordAttack.Range))
				}

				if IsMouseButtonLeftRelased {
					IsMouseButtonLeftRelased = false
					if grabbedEntity != nil {
						if mousePositionScreen.Y < top20Percent {
							// :TODO do the attack on the direction mouse poing to
							var fireballAttack *Entity = createEntity()
							setupAttackFireball(fireballAttack)
							fireballAttack.Position = playerEntity.Position
							fireballAttack.MaxPosition = rl.Vector2AddValue(playerEntity.Position, float32(fireballAttack.Range))
							fireballAttack.inputAxis = rl.Vector2Normalize((rl.Vector2Subtract(mousePositionWorld, playerEntity.Position)))
							destroyEntity(grabbedEntity)
						}

						grabbedEntity = nil
					}

				}

			}

			// :update
			{
				// :update grabbedEntity position

				for i := 0; i < MAX_ENTITY_COUNT; i++ {
					var entity *Entity = &world.Entities[i]
					if entity.Type == ARCH_ATTACK {
						if entity.Position.X > entity.MaxPosition.X || entity.Position.Y > entity.MaxPosition.Y {
							destroyEntity(entity)
							continue
						}
					}
					// :update :positions

					if entity.Type == ARCH_PLAYER {
						entity.Position = rl.Vector2Add(entity.Position, rl.Vector2Scale(entity.inputAxis, (float32(entity.Speed)*delta_t)*runningMultiplier))
					} else if entity.Type == ARCH_TROLL || entity.Type == ARCH_GOBLIN {
						/* entity.inputAxis = rl.Vector2Normalize((rl.Vector2Subtract(playerEntity.Position, entity.Position))) */
						/* entity.Position = rl.Vector2Add(entity.Position, rl.Vector2Scale(entity.inputAxis, (float32(entity.Speed)*delta_t))) */
					} else {
						entity.Position = rl.Vector2Add(entity.Position, rl.Vector2Scale(entity.inputAxis, (float32(entity.Speed)*delta_t)))
					}

					var sprite *Sprite = getSprite(entity.SpriteId)
					entity.CollisionRectangle.X = (entity.Position.X - float32(sprite.Image.Width)/2)
					entity.CollisionRectangle.Y = entity.Position.Y - float32(sprite.Image.Height/2)

					// :update :existance
					if entity.Health <= 0 {
						destroyEntity(entity)
					}
				}

				if grabbedEntity != nil {
					grabbedEntity.Position.X = mousePositionWorld.X
					grabbedEntity.Position.Y = mousePositionWorld.Y
				}
			}

			// :enemy :attack
			{
				for i := 0; i < MAX_ENTITY_COUNT; i++ {
					var entity *Entity = &world.Entities[i]
					if entity.Type == ARCH_TROLL || entity.Type == ARCH_GOBLIN {
						if IsInRange(entity, playerEntity) {
							// TODO enemy attacks
							/* var sword *Entity = createEntity() */
							/* var enemySprite = getSprite(entity.SpriteId) */
							/* var initPosition = rl.Vector2{X: entity.Position.X + float32(enemySprite.Image.Width), Y: entity.Position.Y + float32(enemySprite.Image.Height)} */
							/**/
						}

					}

				}
			}

			// :collision
			{

				for i := 0; i < MAX_ENTITY_COUNT; i++ {

					var firstEntity *Entity = &world.Entities[i]
					if firstEntity.Type == ARCH_ATTACK {

						var didGlobalCollisionHappen bool = false
						for j := 0; j < MAX_ENTITY_COUNT; j++ {

							var didLocalCollisionHappend bool = false
							var secondEntity *Entity = &world.Entities[j]
							if firstEntity == secondEntity {
								continue
							}

							if secondEntity.Type == ARCH_ATTACK || secondEntity.Type == ARCH_CARD || secondEntity.Type == ARCH_PLAYER {
								continue
							}

							didLocalCollisionHappend = rl.CheckCollisionRecs(firstEntity.CollisionRectangle, secondEntity.CollisionRectangle)
							if didLocalCollisionHappend {
								secondEntity.Health -= firstEntity.Damage
								didGlobalCollisionHappen = didLocalCollisionHappend
							}

						}

						if didGlobalCollisionHappen {
							destroyEntity(firstEntity)
						}
					}

				}
			}

			// :render
			{

				var numberOfCards = 0
				for i := 0; i < MAX_ENTITY_COUNT; i++ {
					var entity *Entity = &world.Entities[i]
					if entity.IsValid {

						var entityColor rl.Color = rl.White
						if worldFrame.SelectedEntity == entity {
							entityColor = rl.Red
						}
						switch entity.Type {

						case ARCH_CARD:
							var sprite *Sprite = getSprite(entity.SpriteId)
							var entityColor rl.Color = rl.White
							if worldFrame.SelectedEntity == entity {
								entityColor = rl.Red
							}
							xPosition := int32(camera.Target.X)
							// move to bottom
							yPosition := int32(camera.Target.Y) - (sprite.Image.Height / 2)

							yPosition = yPosition + ((screenHeight / 2) / 3)

							if entity == grabbedEntity {

								xPosition = int32(entity.Position.X)
								// move to bottom
								yPosition = int32(entity.Position.Y)

							} else {
								entity.Position.X = float32(xPosition)
								entity.Position.Y = float32(yPosition)
							}

							rl.DrawTexture(sprite.Image, xPosition-(sprite.Image.Width/2), yPosition-(sprite.Image.Height/2), entityColor)

							numberOfCards += 1
							// do nothig for now

						default:
							var sprite *Sprite = getSprite(entity.SpriteId)
							// show collisions
							//rl.DrawRectangle(int32(entity.CollisionRectangle.X), int32(entity.CollisionRectangle.Y), int32(entity.CollisionRectangle.Width), int32(entity.CollisionRectangle.Height), rl.Blue)
							rl.DrawTexture(sprite.Image, int32(entity.Position.X-float32(sprite.Image.Width/2)), int32(entity.Position.Y-float32(sprite.Image.Height/2)), entityColor)

						}
					}

				}

			}

			// :render ui
			{
				for i := 0; i < MAX_HAND_COUNT; i++ {
					var entity *Entity = &hand.Cards[i]
					if entity.IsValid {
						switch entity.Type {
						default:

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
