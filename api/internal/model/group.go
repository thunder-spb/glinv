package model

// GroupName model
type Group struct {
	ID       int
	ParentID int
	Value    string
	Depth    int
}

// GVar model
type GVar struct {
	ID    int
	Name  string
	Value string
}
