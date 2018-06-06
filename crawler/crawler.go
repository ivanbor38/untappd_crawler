package crawler

import (
	//"time"
	"os"
	"database/sql"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
	//"github.com/foize/go.fifo"
	"strconv"
	"strings"
	"encoding/csv"
	//"github.com/visheratin/ico-crawler/misc"
	. "data_structures"
)





func Spyder(mainPageLink, userId string, depth int) ([]FriendsTable, []User) {
	queue := []string{}
	queue = append(queue, userId)
	var visited map[string]int
	var dist map[string]int
	visited = make(map[string]int)
	dist = make(map[string]int)
	dist[userId] = 1

	edges := []FriendsTable{}
	users := []User{}
	start, _ := GetUser(mainPageLink, userId)
	users = append(users, start)
	var currName string

	for {
		if len(queue) == 0 {
			break
		}

		currName, queue = queue[0], queue[1:]
		if dist[currName] > depth {
			continue
		}
		if visited[currName] == 1 {
			continue
		}
		visited[currName] = 1

		friends, _ := GetFriends(mainPageLink, currName)

		edges = append(edges, friends...)
		for k := 0; k < len(friends); k++ {
			tmpName := friends[k].Friend
			if dist[tmpName] == 0 {
				dist[tmpName] = dist[currName] + 1
			} else {
				if dist[tmpName] > dist[currName] + 1 {
					dist[tmpName] = dist[currName] + 1
				}
			}
			if visited[tmpName] == 0 {
				queue = append(queue, tmpName)
			}
		}
	}
	for z := 0; z < len(edges); z++ {
		kk := User{}
		e := edges[z]
		kk.Id = e.Friend
		kk.Name = e.FriendName
		kk.Distance = dist[e.Friend] - 1
		users = append(users, kk)
	}

	//fmt.Println(edges)
	//fmt.Println(users)
	return edges, users

}




func GetUser(mainPageLink, userId string) (User, error) {

	mainPage := mainPageLink + userId

	doc, _ := goquery.NewDocument(mainPage)



	s := doc.Find("h1").Text()

	result := User{}
	result.Id = userId
	result.Name = s

	//fmt.Println(result)

	return result, nil
}

func GetFriends(mainPageLink, userId string) ([]FriendsTable, error) {

	FriendsLink := mainPageLink + userId + "/friends"
	doc, err := goquery.NewDocument(FriendsLink)
	if err != nil {
		return nil, err
	}

	result := []FriendsTable{}

	doc.Find("div.info").Each(func(i int, s *goquery.Selection) {
		tmp := FriendsTable{}

		friendnick := s.Find("span.username").Text()
		friendname := s.Find("h1").Text()
		//fmt.Println(friendname)
		tmp.Id = userId
		tmp.Friend = friendnick
		tmp.FriendName = friendname
		if friendnick != userId {
			result = append(result, tmp)
		}


	})
	//fmt.Println(result)
	return result, nil
}


func GetBadges(mainPageLink, userId string) ([]BadgeTable, error) {

	BadgesLink := mainPageLink + userId + "/badges"
	doc, err := goquery.NewDocument(BadgesLink)
	if err != nil {
		return nil, err
	}

	result := []BadgeTable{}

	doc.Find("div.item.badge-item.not-retired.level").Each(func(i int, s *goquery.Selection) {
		tmp := BadgeTable{}

		badgename := s.Find("p.name").Text()

		badgename = strings.Split(badgename, "(")[0]
		badgename = strings.Trim(badgename, " ")

		level := s.Find("div.level-box").Text()

		tmp.Id = userId
		tmp.BadgeName = badgename
		cc := 1
		if len(level) > 0 {
			cc, _ = strconv.Atoi(level)
		}
		tmp.BadgeLevel = cc

		result = append(result, tmp)

	})
	return result, nil
}

func GetVenues(mainPageLink, userId string) ([]VenueTable, error) {

	venuesLink := mainPageLink + userId + "/venues?sort=highest_checkin"
	doc, err := goquery.NewDocument(venuesLink)
	if err != nil {
		return nil, err
	}
	result := []VenueTable{}
	doc.Find("div.venue-item").Each(func(i int, s *goquery.Selection) {
		tmp := VenueTable{}
		name := s.Find("p.name").Text()
		category := s.Find("p.category").Text()
		address := s.Find("p.address").Text()
		checkins := s.Find("p.check-ins").Text()

		jj := strings.LastIndex(checkins, ":")
		checkins_str := strings.Trim(checkins[jj+1:], " ")

		address = strings.Trim(address, "\n")
		address = strings.Trim(address, " ")

		tmp.Id = userId
		tmp.VenueName = name
		tmp.Category = category
		tmp.Address = strings.Trim(address, "\n")
		tmp.Checkins, _ = strconv.Atoi(checkins_str)

		result = append(result, tmp)

	})
	return result, nil
}

func GetBeers(mainPageLink, userId string) ([]BeerTable, error) {
	beerLink := mainPageLink + userId + "/beers?sort=highest_rated_their"
	// !!! venuesLink = mainPageLink + "/venues?sort=highest_checkin"
	doc, err := goquery.NewDocument(beerLink)
	if err != nil {
		return nil, err
	}
	result := []BeerTable{}
	doc.Find("div.beer-item").Each(func(i int, s *goquery.Selection) {
		tmp := BeerTable{}
		brewery := s.Find("p.brewery").Text()
		name := s.Find("p.name").Text()
		style := s.Find("p.style").Text()
		ratings_raw := s.Find("div.you").Text()

		//ratings bullshit

		h := strings.Index(ratings_raw, "(")
		j := strings.Index(ratings_raw, ")")
		k := ratings_raw[h+1 : j]

		m := strings.LastIndex(ratings_raw, "(")
		n := strings.LastIndex(ratings_raw, ")")
		o := ratings_raw[m+1 : n]

		abv_raw := s.Find("p.abv").Text()
		hh := strings.Index(abv_raw, "%")
		abv_str := strings.Trim(abv_raw[:hh], "\n")
		//fmt.Printf(abv_str)

		total := s.Find("p.check-ins").Text()
		jj := strings.LastIndex(total, ":")
		total_str := strings.Trim(total[jj+1:], " ")
		//fmt.Printf(total_str)

		ibu_raw := s.Find("p.ibu").Text()
		kk := strings.Index(ibu_raw, " IBU")
		ibu_str := strings.Trim(ibu_raw[:kk], "\n")

		//fmt.Printf(ibu_raw + "\n")

		if ibu_str != "No" {
			tmp.IBU, _ = strconv.Atoi(ibu_str)
		}

		tmp.Id = userId
		tmp.DegustationNumber, _ = strconv.Atoi(total_str)
		tmp.ABV, _ = strconv.ParseFloat(abv_str, 64)
		tmp.UserRating, _ = strconv.ParseFloat(k, 64)
		tmp.GlobalRating, _ = strconv.ParseFloat(o, 64)
		tmp.BreweryName = brewery
		tmp.BeerName = name
		tmp.Style = style

		result = append(result, tmp)

	})
	return result, nil
}

func BeerToSqlite(mainPageLink, userId string) {
	tbl, _ := GetBeers(mainPageLink, userId)
	if len(tbl) == 0 {
		return
	}
	database, _ := sql.Open("sqlite3", "Data.db")
	statement, _ := database.Prepare(`CREATE TABLE IF NOT EXISTS Beer (

		Id 				TEXT,
		UserRating        REAL,
		GlobalRating      REAL,
		BeerName          TEXT,
		BreweryName       TEXT,
		Style             TEXT,
		ABV               REAL,
		IBU               INTEGER,
		DegustationNumber INTEGER,
		unique(Id, BeerName, BreweryName)
		)`)
	statement.Exec()

	SqlStr := `INSERT INTO Beer (Id, UserRating, GlobalRating, BeerName, BreweryName, Style, ABV, IBU,DegustationNumber) VALUES `

	vals := []interface{}{}

	for num := 0; num < len(tbl); num++ {
		SqlStr += "(?, ?, ?, ?, ?, ?, ?, ?, ?), "
		vals = append(vals, tbl[num].Id, tbl[num].UserRating, tbl[num].GlobalRating, tbl[num].BeerName, tbl[num].BreweryName, tbl[num].Style, tbl[num].ABV, tbl[num].IBU, tbl[num].DegustationNumber)

	}
	SqlStr = SqlStr[0 : len(SqlStr)-2]
	statement, _ = database.Prepare(SqlStr)
	statement.Exec(vals...)

}

func VenuesToSqlite(mainPageLink, userId string) {
	vens, _ := GetVenues(mainPageLink, userId)
	if len(vens) == 0 {
		return
	}
	database, _ := sql.Open("sqlite3", "Data.db")
	statement, _ := database.Prepare(`CREATE TABLE IF NOT EXISTS Venues (
		Id        TEXT,
		VenueName TEXT,
		Category  TEXT,
		Address   TEXT,
		Checkins  INTEGER,
		unique(Id, VenueName)
		)`)

	statement.Exec()

	SqlStr := `INSERT INTO Venues (Id, VenueName, Category, Address, Checkins) VALUES `

	vals := []interface{}{}

	for num := 0; num < len(vens); num++ {
		SqlStr += "(?, ?, ?, ?, ?), "
		vals = append(vals, vens[num].Id, vens[num].VenueName, vens[num].Category, vens[num].Address, vens[num].Checkins)

	}
	SqlStr = SqlStr[0 : len(SqlStr)-2]
	statement, _ = database.Prepare(SqlStr)

	statement.Exec(vals...)

}

func BadgesToSqlite(mainPageLink, userId string) {
	badges, _ := GetBadges(mainPageLink, userId)
	if len(badges) == 0 {
		return
	}
	database, _ := sql.Open("sqlite3", "Data.db")
	statement, _ := database.Prepare(`CREATE TABLE IF NOT EXISTS Badges (
		Id         TEXT,
		BadgeName  TEXT,
		BadgeLevel INTEGER,
		unique(Id, BadgeName)
		)`)

	statement.Exec()

	SqlStr := `INSERT INTO Badges (Id, BadgeName, BadgeLevel) VALUES `

	vals := []interface{}{}

	for num := 0; num < len(badges); num++ {
		SqlStr += "(?, ?, ?), "
		vals = append(vals, badges[num].Id, badges[num].BadgeName, badges[num].BadgeLevel)
	}

	SqlStr = SqlStr[0 : len(SqlStr)-2]
	statement, _ = database.Prepare(SqlStr)

	statement.Exec(vals...)

}

func UsersToCsv(arr []User) {


	headers := []string{"Id", "Name", "Distance"}

	f, _ := os.OpenFile("users.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	defer f.Close()
	//my_struct := arr[0]
	w := csv.NewWriter(f)
	defer w.Flush()
	values := make([][]string, 0)
	for _, v := range(arr) {
		c := []string{v.Id, v.Name, strconv.Itoa(v.Distance)}
		values = append(values, c)
	}

	if err := w.Write(headers); err != nil {
		fmt.Println(headers)
		//write failed do something
	}

	if err := w.WriteAll(values); err != nil {
		fmt.Println(values)
		//write failed do something
	}




}

func EdgesToCsv(arr []FriendsTable) {
	//fmt.Println(arr)
	headers := []string{"Id", "Friend", "FriendName"}
	dupl := make(map[string]int)
	f, _ := os.OpenFile("edges.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	defer f.Close()
	//my_struct := arr[0]
	w := csv.NewWriter(f)
	defer w.Flush()
	values := make([][]string, 0)
	for _, v := range(arr) {

		c := []string{v.Id, v.Friend, v.FriendName}
		if dupl[v.Id] == 1 {
			continue
		}
		values = append(values, c)
		dupl[v.Id] = 1
	}

	if err := w.Write(headers); err != nil {
		fmt.Println(headers)
	    //write failed do something
	}

	if err := w.WriteAll(values); err != nil {
		fmt.Println(values)
	    //write failed do something
	}

}


func WriteToSqlite(mainPageLink, userId string) {
	BadgesToSqlite(mainPageLink, userId)
	VenuesToSqlite(mainPageLink, userId)
	BeerToSqlite(mainPageLink, userId)
}



func Total(mainPageLink, userId string, depth int) {
	edges, users := Spyder(mainPageLink, userId, depth)
	//fmt.Println(edges)
	EdgesToCsv(edges)
	UsersToCsv(users)
	var unique map[string]int
	unique = make(map[string]int)


	for us := 0; us < len(users); us++ {
		unique[users[us].Id] = 1
	}

	for u := range unique {
		//fmt.Println(u)

		go WriteToSqlite(mainPageLink, u)
		}



	fmt.Println("cycle ended")
	fmt.Scanln()

}
