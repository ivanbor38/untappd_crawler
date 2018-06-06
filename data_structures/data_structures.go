package data_structures


type BeerTable struct {
	Id                string
	UserRating        float64
	GlobalRating      float64
	BeerName          string
	BreweryName       string
	Style             string
	ABV               float64
	IBU               int
	DegustationNumber int
}

type VenueTable struct {
	Id        string
	VenueName string
	Category  string
	Address   string
	Checkins  int
}

type BadgeTable struct {
	Id         string
	BadgeName  string
	BadgeLevel int
}

type FriendsTable struct {
	Id         string
	Friend  string
	FriendName string

}

type User struct {
	Id string
	Name string
	Distance int

}
