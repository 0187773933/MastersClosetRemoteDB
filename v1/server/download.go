package server

import (
	"fmt"
	"bytes"
	// ulid "github.com/oklog/ulid/v2"
	bolt "github.com/boltdb/bolt"
	fiber "github.com/gofiber/fiber/v2"
	// utils "github.com/0187773933/MastersClosetRemoteDB/v1/utils"
)

type Change struct {
	ID string `json:"id"`
	UUID string `json:"uuid"`
}
func ( s *Server ) GetChangedUsersList( context *fiber.Ctx ) ( error ) {
	result := false
	client_id := context.Get( HEADER_CLIENT_ID )
	if client_id == "" { fmt.Println( "empty client_id" ); return context.Status( fiber.StatusBadRequest ).JSON( fiber.Map{ "result": result , } ) }
	sequence_id := context.Query( HEADER_SEQUENCE_ID )
	if sequence_id == "" { sequence_id = "0" }
	var changes_map map[string]string = make( map[string]string )
	s.DB.View( func( tx *bolt.Tx ) error {
		changed_bucket := tx.Bucket( []byte( "changed" ) )
		if changed_bucket == nil { return nil }
		c := changed_bucket.Cursor()
		last_k , _ := c.Last()
		for k , v := c.Seek( []byte( sequence_id ) ); k != nil && bytes.Compare( k , last_k ) <= 0; k, v = c.Next() {
			changes_map[ string( k ) ] = string( v )
		}
		return nil
	})
	changes := make( []Change , 0 )
	for k , v := range changes_map {
		change := Change{
			ID: k ,
			UUID: v ,
		}
		changes = append( changes , change )
	}
	return context.JSON( fiber.Map{
		"changes": changes ,
	})
}