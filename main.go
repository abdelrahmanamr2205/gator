package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/AbdelrahmanAmr2205/gator/internal/config"
	"github.com/AbdelrahmanAmr2205/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", conf.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	st := &state{db: database.New(db), conf: &conf}
	cmds := commands{handlers: map[string]HandlerFunc{}}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", handlerFollow)
	cmds.register("following", handlerFollowing)

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("Too few arguments\nUsage: gator <command> <command_arguments>")
	}
	err = cmds.run(st, command{name: args[1], args: args[2:]})
	if err != nil {
		log.Fatal(err)
	}
}
