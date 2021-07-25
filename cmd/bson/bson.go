package main

import (
	"fmt"
	"github.com/sanity-io/litter"
	"go.mongodb.org/mongo-driver/bson"
	"myst/timestamp"
	"time"
)

func main() {
	type a struct {
		ID     string              `bson:"id"`
		Number timestamp.Timestamp `bson:"number"`
	}

	b := a{
		ID:     "asd",
		Number: timestamp.Timestamp{Time: time.Now()},
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
