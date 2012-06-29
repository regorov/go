package minecraft

import (
	"fmt"
	"os"
)

// This function runs in the background and listens for incoming packets. It then dispatches the
// handling of them to a handle*Packet method (see packet-handlers.go)
func (client *Client) Receiver() {
	defer func() {
		if client.conn != nil {
			client.Leave()
		}
	}()

	for {
		id, err := client.RecvAnyPacket()
		if err != nil {
			client.ErrChan <- err
			continue
		}

		err = nil
		switch id {
		case 0x00:
			err = client.handleKeepAlivePacket()
		case 0x03:
			err = client.handleChatMessagePacket()
		case 0x04:
			err = client.handleTimeUpdatePacket()
		case 0x05:
			err = client.handleEntityEquipmentPacket()
		case 0x06:
			err = client.handleSpawnPositionPacket()
		case 0x08:
			err = client.handleUpdateHealthPacket()
		case 0x09:
			err = client.handleRespawnPacket()
		case 0x0D:
			err = client.handlePlayerPositionLookPacket()
		case 0x11:
			err = client.handleUseBedPacket()
		case 0x12:
			err = client.handleAnimationPacket()
		case 0x14:
			err = client.handleSpawnNamedEntityPacket()
		case 0x15:
			err = client.handleSpawnDroppedItemPacket()
		case 0x16:
			err = client.handleCollectItemPacket()
		case 0x17:
			err = client.handleSpawnObjectPacket()
		case 0x18:
			err = client.handleSpawnMobPacket()
		case 0x19:
			err = client.handleSpawnPaintingPacket()
		case 0x1A:
			err = client.handleSpawnExperienceOrbPacket()
		case 0x1C:
			err = client.handleEntityVelocityPacket()
		case 0x1D:
			err = client.handleDestroyEntityPacket()
		case 0x1E:
			err = client.handleEntityPacket()
		case 0x1F:
			err = client.handleEntityRelativeMovePacket()
		case 0x20:
			err = client.handleEntityLookPacket()
		case 0x21:
			err = client.handleEntityLookRelativeMovePacket()
		case 0x22:
			err = client.handleEntityTeleportPacket()
		case 0x23:
			err = client.handleEntityHeadLookPacket()
		case 0x26:
			err = client.handleEntityStatusPacket()
		case 0x27:
			err = client.handleAttachEntityPacket()
		case 0x28:
			err = client.handleEntityMetadataPacket()
		case 0x29:
			err = client.handleEntityEffectPacket()
		case 0x2A:
			err = client.handleRemoveEntityEffectPacket()
		case 0x2B:
			err = client.handleSetExperiencePacket()
		case 0x32:
			err = client.handleMapColumnAllocationPacket()
		case 0x33:
			err = client.handleMapChunksPacket()
		case 0x34:
			err = client.handleMultiBlockChangePacket()
		case 0x35:
			err = client.handleBlockChangePacket()
		case 0x36:
			err = client.handleBlockActionPacket()
		case 0x3C:
			err = client.handleExplosionPacket()
		case 0x3D:
			err = client.handleSoundParticleEffectPacket()
		case 0x46:
			err = client.handleChangeGameStatePacket()
		case 0x47:
			err = client.handleThunderboltPacket()
		case 0x64:
			err = client.handleOpenWindowPacket()
		case 0x65:
			err = client.handleCloseWindowPacket()
		case 0x67:
			err = client.handleSetSlotPacket()
		case 0x68:
			err = client.handleSetWindowItemsPacket()
		case 0x69:
			err = client.handleUpdateWindowPropertyPacket()
		case 0x6A:
			err = client.handleConfirmTransactionPacket()
		case 0x6B:
			err = client.handleCreativeInventoryActionPacket()
		case 0x82:
			err = client.handleUpdateSignPacket()
		case 0x83:
			err = client.handleItemDataPacket()
		case 0x84:
			err = client.handleUpdateTileEntityPacket()
		case 0xC8:
			err = client.handleIncrementStatisticPacket()
		case 0xC9:
			err = client.handlePlayerListItemPacket()
		case 0xCA:
			err = client.handlePlayerAbilitiesPacket()
		case 0xFA:
			err = client.handlePluginMessagePacket()
		case 0xFF:
			err = client.handleKickPacket()
		default:
			fmt.Fprintf(os.Stderr, "Ignoring unhandled packet with id 0x%02X", id)
		}

		if err == Stop {
			return
		}

		if err != nil {
			client.ErrChan <- err
		}
	}
}
