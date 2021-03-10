package models

type Client struct {
	Id         int64
	Login      string
	FirstName  string
	LastName   string
	MiddleName *string
	Passport   string
	Birthday   int64
	Status     string
	CreatedAt  int64
}

type Card struct {
	Id        int64
	Number    string
	Balance   int64
	Issuer    string
	Holder    string
	OwnerId   int64
	Status    string
	CreatedAt int64
}

type Icon struct {
	Id    int64
	Title string
	Uri   string
}

type Mcc struct {
	Id   string
	Text string
}

type Transaction struct {
	Id          int64
	FromId      *int64
	ToId        *int64
	Sum         int64
	MccId       *string
	IconId      *int64
	Description *string
	CreatedAt   int64
}

type MostPopularSpending struct {
	Description string
	Count       int64
	IconUri     string
}

type MostExpensiveSpending struct {
	Description string
	Sum         int64
	IconUri     string
}
