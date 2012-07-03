package minecraft

import (
	"io"
)

type ColumnCoord struct {
	X int
	Z int
}

type Column struct {
	Chunks map[int]*Chunk
	Biomes []byte
}

type Chunk struct {
	BlockTypes    []byte
	BlockMetadata []byte
	BlockLight    []byte
	SkyLight      []byte
	AddTypes      []byte
}

func (client *Client) GetChunk(cx int, cy int, cz int) (chunk *Chunk, ok bool) {
	column, ok := client.Columns[ColumnCoord{cx, cz}]
	if !ok {
		return nil, false
	}

	chunk, ok = column.Chunks[cy]
	return chunk, ok
}

func (client *Client) GetBlock(x int, y int, z int) (blockType byte, blockData byte, blockLight byte, skyLight byte, addType byte, ok bool) {
	chunk, ok := client.GetChunk(x>>4, y>>4, z>>4)
	if !ok {
		return 0, 0, 0, 0, 0, false
	}

	xOffset := x & 15
	yOffset := y & 15
	zOffset := z & 15
	byteOffset := xOffset | (zOffset << 4) | (yOffset << 8)

	if chunk.BlockTypes != nil {
		blockType = chunk.BlockTypes[byteOffset]
	}

	if chunk.BlockMetadata != nil {
		blockData = chunk.BlockMetadata[byteOffset/2]
		if byteOffset%2 == 0 {
			blockData = blockData & 0x0F
		} else {
			blockData = blockData >> 4
		}
	}

	if chunk.BlockLight != nil {
		blockLight = chunk.BlockLight[byteOffset/2]
		if byteOffset%2 == 0 {
			blockLight = blockLight & 0x0F
		} else {
			blockLight = blockLight >> 4
		}
	}

	if chunk.SkyLight != nil {
		skyLight = chunk.SkyLight[byteOffset/2]
		if byteOffset%2 == 0 {
			skyLight = skyLight & 0x0F
		} else {
			skyLight = skyLight >> 4
		}
	}

	if chunk.AddTypes != nil {
		addType = chunk.AddTypes[byteOffset/2]
		if byteOffset%2 == 0 {
			addType = addType & 0x0F
		} else {
			addType = addType >> 4
		}
	}

	return blockType, blockData, blockLight, skyLight, addType, true
}

func (client *Client) readBlockTypes(r io.ReadCloser, column *Column, primaryBitMap uint16) (err error) {
	for cy := 0; cy < 16; cy++ {
		if primaryBitMap&(1<<uint(cy)) != 0 {
			data := make([]byte, 4096)
			_, err = r.Read(data)
			if err != nil {
				return err
			}

			chunk, ok := column.Chunks[cy]
			if !ok {
				chunk = new(Chunk)
				column.Chunks[cy] = chunk
			}

			chunk.BlockTypes = data
		}
	}

	return nil
}

func (client *Client) readBlockMetadata(r io.ReadCloser, column *Column, primaryBitMap uint16) (err error) {
	for cy := 0; cy < 16; cy++ {
		if primaryBitMap&(1<<uint(cy)) != 0 {
			data := make([]byte, 2048)
			_, err = r.Read(data)
			if err != nil {
				return err
			}

			chunk, ok := column.Chunks[cy]
			if !ok {
				chunk = new(Chunk)
				column.Chunks[cy] = chunk
			}

			chunk.BlockMetadata = data
		}
	}

	return nil
}

func (client *Client) readBlockLight(r io.ReadCloser, column *Column, primaryBitMap uint16) (err error) {
	for cy := 0; cy < 16; cy++ {
		if primaryBitMap&(1<<uint(cy)) != 0 {
			data := make([]byte, 2048)
			_, err = r.Read(data)
			if err != nil {
				return err
			}

			chunk, ok := column.Chunks[cy]
			if !ok {
				chunk = new(Chunk)
				column.Chunks[cy] = chunk
			}

			chunk.BlockLight = data
		}
	}

	return nil
}

func (client *Client) readSkyLight(r io.ReadCloser, column *Column, primaryBitMap uint16) (err error) {
	for cy := 0; cy < 16; cy++ {
		if primaryBitMap&(1<<uint(cy)) != 0 {
			data := make([]byte, 2048)
			_, err = r.Read(data)
			if err != nil {
				return err
			}

			chunk, ok := column.Chunks[cy]
			if !ok {
				chunk = new(Chunk)
				column.Chunks[cy] = chunk
			}

			chunk.SkyLight = data
		}
	}

	return nil
}

func (client *Client) readAddTypes(r io.ReadCloser, column *Column, addBitMap uint16) (err error) {
	for cy := 0; cy < 16; cy++ {
		if addBitMap&(1<<uint(cy)) != 0 {
			data := make([]byte, 2048)
			_, err = r.Read(data)
			if err != nil {
				return err
			}

			chunk, ok := column.Chunks[cy]
			if !ok {
				chunk = new(Chunk)
				column.Chunks[cy] = chunk
			}

			chunk.AddTypes = data
		}
	}

	return nil
}

func (client *Client) readBiomeData(r io.ReadCloser, column *Column) (err error) {
	data := make([]byte, 256)
	_, err = r.Read(data)
	if err != nil {
		return err
	}

	column.Biomes = data
	return nil
}
