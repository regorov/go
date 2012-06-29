package minecraft

import (
	"fmt"
)

func (client *Client) handleKeepAlivePacket() (err error) {
	var keepAliveID int32

	err = client.RecvPacketData(&keepAliveID)
	if err != nil {
		return err
	}

	err = client.SendPacket(0x00, keepAliveID)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleChatMessagePacket() (err error) {
	var msg string

	err = client.RecvPacketData(&msg)
	if err != nil {
		return err
	}

	client.HandleMessage(msg)

	return nil
}

func (client *Client) handleTimeUpdatePacket() (err error) {
	var t int64

	err = client.RecvPacketData(&t)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleEntityEquipmentPacket() (err error) {
	var entityId int32
	var slot, itemId, damage int16

	err = client.RecvPacketData(&entityId, &slot, &itemId, &damage)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleSpawnPositionPacket() (err error) {
	var x, y, z int32

	err = client.RecvPacketData(&x, &y, &z)
	if err != nil {
		return err
	}

	client.playerX = float64(x)
	client.playerY = float64(y)
	client.playerZ = float64(z)

	return nil
}

func (client *Client) handleUpdateHealthPacket() (err error) {
	var health, food int16
	var foodSat float32

	err = client.RecvPacketData(&health, &food, &foodSat)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleRespawnPacket() (err error) {
	var dimension int32
	var difficulty, creativeMode int8
	var worldHeight int16
	var levelType string

	err = client.RecvPacketData(&dimension, &difficulty, &creativeMode, &worldHeight, &levelType)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handlePlayerPositionLookPacket() (err error) {
	err = client.RecvPacketData(&client.playerX, &client.playerStance, &client.playerY, &client.playerZ, &client.playerYaw, &client.playerPitch, &client.playerOnGround)
	if err != nil {
		return err
	}

	if client.debug {
		fmt.Printf("Received position update: (%.1f, %.1f, %.1f) (%.1f, %.1f) on ground: %t  stance: %.1f\n", client.playerX, client.playerY, client.playerZ, client.playerYaw, client.playerPitch, client.playerOnGround, client.playerStance)
	}

	err = client.SendPacket(0x0D, client.playerX, client.playerY, client.playerStance, client.playerZ, client.playerYaw, client.playerPitch, client.playerOnGround)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleUseBedPacket() (err error) {
	var entityId, x, z int32
	var unknown, y int8

	err = client.RecvPacketData(&entityId, &unknown, &x, &y, &z)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleAnimationPacket() (err error) {
	var entityId int32
	var animation int8

	err = client.RecvPacketData(&entityId, &animation)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleSpawnNamedEntityPacket() (err error) {
	var entityId, x, y, z int32
	var playerName string
	var yaw, pitch int8
	var currentItem int16

	err = client.RecvPacketData(&entityId, &playerName, &x, &y, &z, &yaw, &pitch, &currentItem)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleSpawnDroppedItemPacket() (err error) {
	var entityId, x, y, z int32
	var item, damage int16
	var count, rotation, pitch, roll int8

	err = client.RecvPacketData(&entityId, &item, &count, &damage, &x, &y, &z, &rotation, &pitch, &roll)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleCollectItemPacket() (err error) {
	var collectedId, collectorId int32

	err = client.RecvPacketData(&collectedId, &collectorId)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleSpawnObjectPacket() (err error) {
	var entityId, x, y, z, throwerId int32
	var entityType int8
	var speedX, speedY, speedZ int16

	err = client.RecvPacketData(&entityId, &entityType, &x, &y, &z, &throwerId)
	if err != nil {
		return err
	}

	if throwerId != 0 {
		err = client.RecvPacketData(&speedX, &speedY, &speedZ)
		if err != nil {
			return err
		}
	}

	return nil
}

func (client *Client) handleSpawnMobPacket() (err error) {
	var entityId, x, y, z int32
	var entityType, yaw, pitch, headYaw int8

	err = client.RecvPacketData(&entityId, &entityType, &x, &y, &z, &yaw, &pitch, &headYaw)
	if err != nil {
		return err
	}

	_, err = client.RecvEntityMetadata()
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleSpawnPaintingPacket() (err error) {
	var entityId, x, y, z, direction int32
	var title string

	err = client.RecvPacketData(&entityId, &title, &x, &y, &z, &direction)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleSpawnExperienceOrbPacket() (err error) {
	var entityId, x, y, z int32
	var count int16

	err = client.RecvPacketData(&entityId, &x, &y, &z, &count)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleEntityVelocityPacket() (err error) {
	var entityId int32
	var vx, vy, vz int16

	err = client.RecvPacketData(&entityId, &vx, &vy, &vz)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleDestroyEntityPacket() (err error) {
	var entityId int32

	err = client.RecvPacketData(&entityId)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleEntityPacket() (err error) {
	var entityId int32

	err = client.RecvPacketData(&entityId)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleEntityRelativeMovePacket() (err error) {
	var entityId int32
	var dx, dy, dz int8

	err = client.RecvPacketData(&entityId, &dx, &dy, &dz)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleEntityLookPacket() (err error) {
	var entityId int32
	var yaw, pitch int8

	err = client.RecvPacketData(&entityId, &yaw, &pitch)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleEntityLookRelativeMovePacket() (err error) {
	var entityId int32
	var dx, dy, dz, yaw, pitch int8

	err = client.RecvPacketData(&entityId, &dx, &dy, &dz, &yaw, &pitch)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleEntityTeleportPacket() (err error) {
	var entityId, x, y, z int32
	var yaw, pitch int8

	err = client.RecvPacketData(&entityId, &x, &y, &z, &yaw, &pitch)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleEntityHeadLookPacket() (err error) {
	var entityId int32
	var headYaw int8

	err = client.RecvPacketData(&entityId, &headYaw)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleEntityStatusPacket() (err error) {
	var entityId int32
	var status int8

	err = client.RecvPacketData(&entityId, &status)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleAttachEntityPacket() (err error) {
	var entityId, vehicleId int32

	err = client.RecvPacketData(&entityId, &vehicleId)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleEntityMetadataPacket() (err error) {
	var entityId int32

	err = client.RecvPacketData(&entityId)
	if err != nil {
		return err
	}

	_, err = client.RecvEntityMetadata()
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleEntityEffectPacket() (err error) {
	var entityId int32
	var effectId, amplifier int8
	var duration int16

	err = client.RecvPacketData(&entityId, &effectId, &amplifier, &duration)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleRemoveEntityEffectPacket() (err error) {
	var entityId int32
	var effectId int8

	err = client.RecvPacketData(&entityId, &effectId)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleSetExperiencePacket() (err error) {
	var xp float32
	var level, totalXp int16

	err = client.RecvPacketData(&xp, &level, &totalXp)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleMapColumnAllocationPacket() (err error) {
	var cx, cz int32
	var initialize bool

	err = client.RecvPacketData(&cx, &cz, &initialize)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleMapChunksPacket() (err error) {
	var cx, cz, compressedSize, unused int32
	var groundUpContinuous bool
	var primaryBitMap, addBitMap uint16

	err = client.RecvPacketData(&cx, &cz, &groundUpContinuous, &primaryBitMap, &addBitMap, &compressedSize, &unused)
	if err != nil {
		return err
	}

	compressedData := make([]byte, compressedSize)
	_, err = client.conn.Read(compressedData)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleMultiBlockChangePacket() (err error) {
	var cx, cz, dataSize int32
	var recordCount int16

	err = client.RecvPacketData(&cx, &cz, &recordCount, &dataSize)
	if err != nil {
		return err
	}

	data := make([]byte, dataSize)
	_, err = client.conn.Read(data)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleBlockChangePacket() (err error) {
	var x, z int32
	var y, blockType, blockMetadata int8

	err = client.RecvPacketData(&x, &y, &z, &blockType, &blockMetadata)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleBlockActionPacket() (err error) {
	var x, z int32
	var y int16
	var b1, b2 int8

	err = client.RecvPacketData(&x, &y, &z, &b1, &b2)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleExplosionPacket() (err error) {
	var x, y, z float64
	var radius float32
	var recordCount int32

	err = client.RecvPacketData(&x, &y, &z, &radius, &recordCount)
	if err != nil {
		return err
	}

	records := make([]ExplosionRecord, recordCount)

	for i := int32(0); i < recordCount; i++ {
		var record ExplosionRecord
		records[i] = record
		err = client.RecvPacketData(&record.X, &record.Y, &record.Z)
		if err != nil {
			return err
		}
	}

	return nil
}

func (client *Client) handleSoundParticleEffectPacket() (err error) {
	var effectId, x, z, data int32
	var y int8

	err = client.RecvPacketData(&effectId, &x, &y, &z, &data)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleChangeGameStatePacket() (err error) {
	var reason, gameMode int8

	err = client.RecvPacketData(&reason, &gameMode)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleThunderboltPacket() (err error) {
	var entityId, x, y, z int32
	var unknown bool

	err = client.RecvPacketData(&entityId, &unknown, &x, &y, &z)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleOpenWindowPacket() (err error) {
	var windowId, inventoryType, numSlots int8
	var windowTitle string

	err = client.RecvPacketData(&windowId, &inventoryType, &windowTitle, &numSlots)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleCloseWindowPacket() (err error) {
	var windowId int8

	err = client.RecvPacketData(&windowId)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleSetSlotPacket() (err error) {
	var windowId int8
	var slot int16
	var data Slot

	err = client.RecvPacketData(&windowId, &slot, &data)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleSetWindowItemsPacket() (err error) {
	var windowId int8
	var count int16

	err = client.RecvPacketData(&windowId, &count)
	if err != nil {
		return err
	}

	slots := make([]Slot, count)

	for i := int16(0); i < count; i++ {
		err = client.RecvPacketData(&slots[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (client *Client) handleUpdateWindowPropertyPacket() (err error) {
	var windowId int8
	var property, value int16

	err = client.RecvPacketData(&windowId, &property, &value)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleConfirmTransactionPacket() (err error) {
	var windowId int8
	var actionNumber int16
	var accepted bool

	err = client.RecvPacketData(&windowId, &actionNumber, &accepted)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleCreativeInventoryActionPacket() (err error) {
	var slot int16
	var data Slot

	err = client.RecvPacketData(&slot, &data)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleUpdateSignPacket() (err error) {
	var x, z int32
	var y int16
	var text1, text2, text3, text4 string

	err = client.RecvPacketData(&x, &y, &z, &text1, &text2, &text3, &text4)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleItemDataPacket() (err error) {
	var itemType, itemId int16
	var length uint8

	err = client.RecvPacketData(&itemType, &itemId, &length)
	if err != nil {
		return err
	}

	data := make([]byte, length)
	_, err = client.conn.Read(data)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleUpdateTileEntityPacket() (err error) {
	var x, z, custom1, custom2, custom3 int32
	var y int16
	var action int8

	err = client.RecvPacketData(&x, &y, &z, &action, &custom1, &custom2, &custom3)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handleIncrementStatisticPacket() (err error) {
	var statisticId int32
	var amount int8

	err = client.RecvPacketData(&statisticId, &amount)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handlePlayerListItemPacket() (err error) {
	var playerName string
	var online bool
	var ping int16

	err = client.RecvPacketData(&playerName, &online, &ping)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handlePlayerAbilitiesPacket() (err error) {
	var invulnerability, isFlying, canFly, instantDestroy bool

	err = client.RecvPacketData(&invulnerability, &isFlying, &canFly, &instantDestroy)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) handlePluginMessagePacket() (err error) {
	var channel string
	var length int16

	err = client.RecvPacketData(&channel, &length)
	if err != nil {
		return err
	}

	data := make([]byte, length)
	_, err = client.conn.Read(data)
	if err != nil {
		return err
	}

	client.HandleMessage(string(data))

	return nil
}

func (client *Client) handleKickPacket() (err error) {
	var message string

	err = client.RecvPacketData(&message)
	if err != nil {
		return err
	}

	fmt.Println("Disconnected: %s\n", message)
	client.LeaveNoKick()

	return Stop
}
