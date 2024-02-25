package main

import (
	initModule "github.com/NeptuneYeh/simplebank/init"
)

//const (
//	dbDriver      = "postgres"
//	dbSource      = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
//	serverAddress = "0.0.0.0:8080"
//)

func main() {
	// noob style
	//conn, err := sql.Open(dbDriver, dbSource)
	//if err != nil {
	//	log.Fatal("cannot connect to db: ", err)
	//}
	//
	//// 依賴建立
	//store := postgresdb.NewStore(conn)
	//server := api.NewServer(store)
	//
	//err = server.Start(serverAddress)
	//if err != nil {
	//	log.Fatal("cannot start server: ", err)
	//}
	initProcess := initModule.NewMainInitProcess()
	initProcess.Run()
}
