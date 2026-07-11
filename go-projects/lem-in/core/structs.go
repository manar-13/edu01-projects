package core

type Room struct {
	Name       string
	X, Y       int
	Links      []*Room
	Parent     *Room
	FlowParent *Room
	Flow       []*Room
	FlowFrom   *Room
	Occupied   bool
}

type PathsSet struct {
	Paths       [][]*Room
	Lengths     []int
	PathsAmount int
}

type Result struct {
	Finished   int
	AntNum     int
	Moves      int
	Left       int
	FirstPrint bool
}

type Data struct {
	Ants      int
	Rooms     map[string]*Room
	Input     []string
	Start     *Room
	End       *Room
	BestSet   *PathsSet
	BestSpeed int
	RoomOrder []*Room
}
