// Command mm3ddump is used to demonstrate the functionality of the mm3dmodel package.
package main

import (
	"flag"
	"fmt"
	"github.com/kierdavis/go/amberfell/mm3dmodel"
	"os"
)

func die(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func printMetadata(model *mm3dmodel.Model) {
	fmt.Printf("Metadata:\n")
	for key, value := range model.Metadata() {
		fmt.Printf("  %s: %s\n", key, value)
	}

	if len(model.Metadata()) == 0 {
		fmt.Printf("  None\n\n")
	} else {
		fmt.Printf("\n")
	}
}

func printVertices(model *mm3dmodel.Model) {
	fmt.Printf("Vertices (%d):\n", model.NVertices())
	for i := 0; i < model.NVertices(); i++ {
		vertex := model.Vertex(i)
		fmt.Printf("  %4d: Flags: 0x%04X  (%.3f, %.3f, %.3f)\n", i, vertex.Flags(), vertex.X(),
            vertex.Y(), vertex.Z())
	}

	if model.NVertices() == 0 {
		fmt.Printf("  None\n\n")
	} else {
		fmt.Printf("\n")
	}
}

func printTriangles(model *mm3dmodel.Model) {
	fmt.Printf("Triangles (%d):\n", model.NTriangles())
	for i := 0; i < model.NTriangles(); i++ {
		triangle := model.Triangle(i)
		fmt.Printf("  %4d: Flags: 0x%04X  (%d, %d, %d)\n", i, triangle.Flags(),
            triangle.VertexIndex1(), triangle.VertexIndex2(), triangle.VertexIndex3())
	}

	if model.NTriangles() == 0 {
		fmt.Printf("  None\n\n")
	} else {
		fmt.Printf("\n")
	}
}

func printTriangleNormals(model *mm3dmodel.Model) {
	fmt.Printf("Triangle normals (%d):\n", model.NTriangleNormals())
	for i := 0; i < model.NTriangleNormals(); i++ {
		triangleNormals := model.TriangleNormals(i)
		v1x, v1y, v1z := triangleNormals.Vertex1Normal()
		v2x, v2y, v2z := triangleNormals.Vertex2Normal()
		v3x, v3y, v3z := triangleNormals.Vertex3Normal()
		fmt.Printf("  %4d: Flags: 0x%04X [of triangle %d] (%.3f, %.3f, %.3f), (%.3f, %.3f, %.3f), (%.3f, %.3f, %.3f)\n",
            i, triangleNormals.Flags(), triangleNormals.TriangleIndex(), v1x, v1y, v1z, v2x, v2y,
            v2z, v3x, v3y, v3z)
	}

	if model.NTriangleNormals() == 0 {
		fmt.Printf("  None\n\n")
	} else {
		fmt.Printf("\n")
	}
}

func printTextureCoordinates(model *mm3dmodel.Model) {
	fmt.Printf("Texture Coordinates (%d):\n", model.NTextureCoordinates())
	for i := 0; i < model.NTextureCoordinates(); i++ {
		textureCoordinates := model.TextureCoordinates(i)
		s1, t1 := textureCoordinates.Vertex1Coord()
		s2, t2 := textureCoordinates.Vertex2Coord()
		s3, t3 := textureCoordinates.Vertex3Coord()
		fmt.Printf("  %4d: Flags: 0x%04X [of triangle %d] (%.3f, %.3f), (%.3f, %.3f), (%.3f, %.3f)\n",
            i, textureCoordinates.Flags(), textureCoordinates.TriangleIndex(), s1, t1, s2, t2,
            s3, t3)
	}

	if model.NTextureCoordinates() == 0 {
		fmt.Printf("  None\n\n")
	} else {
		fmt.Printf("\n")
	}
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Not enough arguments.\n\nUsage: %s <file.mm3d>\n", os.Args[0])
		os.Exit(2)
	}

	fname := flag.Arg(0)
	f, err := os.Open(fname)
	die(err)
	defer f.Close()

	model, err := mm3dmodel.Read(f)
	die(err)

	fmt.Printf("Version: %d.%d\n", model.MajorVersion(), model.MinorVersion())
	fmt.Printf("Model flags: 0x%02x\n", model.ModelFlags())
	fmt.Printf("Num dirty segments: %d\n\n", model.NDirtySegments())

	printMetadata(model)
	printVertices(model)
	printTriangles(model)
	printTriangleNormals(model)
	printTextureCoordinates(model)
}
