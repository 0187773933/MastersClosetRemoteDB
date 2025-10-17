package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"path/filepath"
	bolt "github.com/boltdb/bolt"
	utils "github.com/0187773933/MastersClosetRemoteDB/v1/utils"
	server "github.com/0187773933/MastersClosetRemoteDB/v1/server"
)

var s server.Server
var db *bolt.DB

func SetupCloseHandler() {
	c := make( chan os.Signal )
	signal.Notify( c , os.Interrupt , syscall.SIGTERM , syscall.SIGINT )
	go func() {
		<-c
		fmt.Println( "\r- Ctrl+C pressed in Terminal" )
		fmt.Println( "Shutting Down Master's Closet Remote DB Server" )
		s.FiberApp.Shutdown()
		db.Close()
		os.Exit( 0 )
	}()
}

func main() {
	SetupCloseHandler()
	config_file_path , _ := filepath.Abs( os.Args[ 1 ] )
	config := utils.ParseConfig( config_file_path )
	db , db_err := bolt.Open( config.BoltDBPath , 0600 , &bolt.Options{} )
	if db_err != nil { panic( db_err.Error() ) }
	s = server.New( config , db )
	s.Start()
}