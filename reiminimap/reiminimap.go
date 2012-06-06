package reiminimap

import (
	"fmt"
	"github.com/kierdavis/goutil"
	"image/color"
	"io"
	"strconv"
	"strings"
)

type Waypoint struct {
	Name	string
	X	int
	Y	int
	Z	int
	Visible	bool
	Color	color.Color
}

func ReadWaypoints(r io.Reader) (waypoints []*Waypoint, err error) {
	lineChan, errChan := util.IterLines(r)

	for line := range lineChan {
		parts := strings.Split(line, ":")

		x, err := strconv.ParseInt(parts[1], 10, 0)
		if err != nil {
			return waypoints, err
		}

		y, err := strconv.ParseInt(parts[2], 10, 0)
		if err != nil {
			return waypoints, err
		}

		z, err := strconv.ParseInt(parts[3], 10, 0)
		if err != nil {
			return waypoints, err
		}

		colorNum, err := strconv.ParseUint(parts[5], 16, 32)
		if err != nil {
			return waypoints, err
		}

		colorRed := uint8(colorNum >> 16)
		colorGreen := uint8(colorNum >> 8)
		colorBlue := uint8(colorNum)

		waypoint := &Waypoint{
			Name:		parts[0],
			X:		int(x),
			Y:		int(y),
			Z:		int(z),
			Visible:	parts[4] == "true",
			Color:		color.RGBA{colorRed, colorGreen, colorBlue, 0},
		}

		waypoints = append(waypoints, waypoint)
	}

	return waypoints, <-errChan
}

func WriteWaypoints(w io.Writer, waypoints []*Waypoint) (err error) {
	for _, waypoint := range waypoints {
		r, g, b, _ := waypoint.Color.RGBA()

		_, err = fmt.Fprintf(w, "%s:%d:%d:%d:%t:%02X%02X%02X\n", waypoint.Name, waypoint.X, waypoint.Y, waypoint.Z, waypoint.Visible, r, g, b)
		if err != nil {
			return err
		}
	}

	return nil
}
