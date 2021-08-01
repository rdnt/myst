package main

import (
	"fmt"
	"time"

	timestamp2 "myst/pkg/timestamp"

	"github.com/sanity-io/litter"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	type a struct {
		ID     string               `bson:"id"`
		Number timestamp2.Timestamp `bson:"number"`
	}

	b := a{
		ID:     "asd",
		Number: timestamp2.Timestamp{Time: time.Now()},
	}
	var d a

	c, err := bson.Marshal(b)
	if err != nil {
		panic(err)
	}

	fmt.Printf("marshaled %X\n", c)

	err = bson.Unmarshal(c, &d)
	if err != nil {
		panic(err)
	}
	fmt.Println(litter.Sdump(b.Number.Unix(), d.Number.Unix()))
}
