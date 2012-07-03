package minecraft

import (
	"fmt"
)

// This function listens for incoming packets and dispatches the handling of them to a
// handle*Packet method (see packet-handlers.go)
func (client *Client) Run() (kickMessage string) {
	defer func() {
		if client.conn != nil {
			client.Leave()
		}
	}()

	for client.conn != nil {
		id, err := client.RecvAnyPacket()
		if err != nil {
			client.ErrChan <- err
			continue
		}

		err = client.dispatchPacket(id)
		if err != nil {
			if kick, ok := err.(Kick); ok {
				return string(kick)
			}

			client.ErrChan <- err
		}
	}

	return ""
}

func (client *Client) dispatchPacket(id byte) (err error) {
	/*
		defer func() {
			if err == nil {
				e := recover()

				if e != nil {
					var ok bool
					err, ok = e.(error)

					if !ok {
						err = fmt.Errorf("%v", e)
					}
				}
			}
		}()
	*/

	println(id)

	switch id {
	case 0x00:
		return client.handleKeepAlivePacket()
	case 0x03:
		return client.handleChatMessagePacket()
	case 0x04:
		return client.handleTimeUpdatePacket()
	case 0x05:
		return client.handleEntityEquipmentPacket()
	case 0x06:
		return client.handleSpawnPositionPacket()
	case 0x08:
		return client.handleUpdateHealthPacket()
	case 0x09:
		return client.handleRespawnPacket()
	case 0x0D:
		return client.handlePlayerPositionLookPacket()
	case 0x11:
		return client.handleUseBedPacket()
	case 0x12:
		return client.handleAnimationPacket()
	case 0x14:
		return client.handleSpawnNamedEntityPacket()
	case 0x15:
		return client.handleSpawnDroppedItemPacket()
	case 0x16:
		return client.handleCollectItemPacket()
	case 0x17:
		return client.handleSpawnObjectPacket()
	case 0x18:
		return client.handleSpawnMobPacket()
	case 0x19:
		return client.handleSpawnPaintingPacket()
	case 0x1A:
		return client.handleSpawnExperienceOrbPacket()
	case 0x1C:
		return client.handleEntityVelocityPacket()
	case 0x1D:
		return client.handleDestroyEntityPacket()
	case 0x1E:
		return client.handleEntityPacket()
	case 0x1F:
		return client.handleEntityRelativeMovePacket()
	case 0x20:
		return client.handleEntityLookPacket()
	case 0x21:
		return client.handleEntityLookRelativeMovePacket()
	case 0x22:
		return client.handleEntityTeleportPacket()
	case 0x23:
		return client.handleEntityHeadLookPacket()
	case 0x26:
		return client.handleEntityStatusPacket()
	case 0x27:
		return client.handleAttachEntityPacket()
	case 0x28:
		return client.handleEntityMetadataPacket()
	case 0x29:
		return client.handleEntityEffectPacket()
	case 0x2A:
		return client.handleRemoveEntityEffectPacket()
	case 0x2B:
		return client.handleSetExperiencePacket()
	case 0x32:
		return client.handleMapColumnAllocationPacket()
	case 0x33:
		return client.handleMapChunksPacket()
	case 0x34:
		return client.handleMultiBlockChangePacket()
	case 0x35:
		return client.handleBlockChangePacket()
	case 0x36:
		return client.handleBlockActionPacket()
	case 0x3C:
		return client.handleExplosionPacket()
	case 0x3D:
		return client.handleSoundParticleEffectPacket()
	case 0x46:
		return client.handleChangeGameStatePacket()
	case 0x47:
		return client.handleThunderboltPacket()
	case 0x64:
		return client.handleOpenWindowPacket()
	case 0x65:
		return client.handleCloseWindowPacket()
	case 0x67:
		return client.handleSetSlotPacket()
	case 0x68:
		return client.handleSetWindowItemsPacket()
	case 0x69:
		return client.handleUpdateWindowPropertyPacket()
	case 0x6A:
		return client.handleConfirmTransactionPacket()
	case 0x6B:
		return client.handleCreativeInventoryActionPacket()
	case 0x82:
		return client.handleUpdateSignPacket()
	case 0x83:
		return client.handleItemDataPacket()
	case 0x84:
		return client.handleUpdateTileEntityPacket()
	case 0xC8:
		return client.handleIncrementStatisticPacket()
	case 0xC9:
		return client.handlePlayerListItemPacket()
	case 0xCA:
		return client.handlePlayerAbilitiesPacket()
	case 0xFA:
		return client.handlePluginMessagePacket()
	case 0xFF:
		return client.handleKickPacket()

	default:
		if client.DebugWriter != nil {
			fmt.Fprintf(client.DebugWriter, "Ignoring unhandled packet with id 0x%02X\n", id)
		}
	}

	return nil
}
