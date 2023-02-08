package models

type ApiKeys struct {
	Key         string   `bson:key`
	Holder      string   `bson:holder`
	Scope       []string `bson:scope`
	Application string   `bson:application`
}

type SmellyServer struct {
	Address       string   `bson:address`
	Version       string   `bson:version`
	Players       []string `bson:players`
	OnlinePlayers int      `bson:"online_players"`
	DateCreated   string   `bson:datecreated`
	DateUpdated   string   `bson:dateupdated`
}

type PlayerServerHistory struct {
	Player  string   `bson:player`
	Servers []string `bson:servers`
}
