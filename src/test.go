package main

import (
    "fmt"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

type Person struct {
    Name  string
    Phone string
}


func Init_db() {
    var err error
    Session,err := mgo.Dial("")
    if err!=nil{
        panic(err)
    }

    Session.SetMode(mgo.Monotonic, true)
    _ = Session.DB(dbname).C("users")
}

func main() {
    session, err := mgo.Dial("")
    if err != nil {
        panic(err)
    }
    defer session.Close()

    // Optional. Switch the session to a monotonic behavior.
    session.SetMode(mgo.Monotonic, true)

    c := session.DB("test").C("people")
    err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
        &Person{"Cla", "+55 53 8402 8510"})
    if err != nil {
        panic(err)
    }

    result := Person{}
    err = c.Find(bson.M{"name": "Ale"}).One(&result)
    if err != nil {
        panic(err)
    }

    fmt.Println("Phone:", result.Phone)
}
