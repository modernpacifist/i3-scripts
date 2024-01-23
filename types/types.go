package types

import (
	"go.i3wm.org/i3/v4"
)

type Config interface {
	Dump()
	ReadFile()
}

type Container struct {
	// TODO: add a node id and adress some of the containers by it <13-11-23, modernpacifist> //
	ID     i3.NodeID `json:"ID"`
	X      int64     `json:"X"`
	Y      int64     `json:"Y"`
	Width  int64     `json:"Width"`
	Height int64     `json:"Height"`
	Marks  []string  `json:"Marks"`
	// TODO: having a mark field here is extremely idiotic, since we have a map with keys of marks themselves <13-11-23, modernpacifist> //
}

type FocusMarkedContainerContainer struct {
	// TODO: add a node id and adress some of the containers by it <13-11-23, modernpacifist> //
	ID     i3.NodeID `json:"ID"`
	X      int64     `json:"X"`
	Y      int64     `json:"Y"`
	Width  int64     `json:"Width"`
	Height int64     `json:"Height"`
	Marks  []string  `json:"Marks"`
	Output string    `json:"Output"`
	// TODO: having a mark field here is extremely idiotic, since we have a map with keys of marks themselves <13-11-23, modernpacifist> //
}
