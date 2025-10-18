package server

import (
	"fmt"
	// "bytes"
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
	changes := make( []Change , 0 )
	client_id := context.Get( HEADER_CLIENT_ID )
	if client_id == "" { fmt.Println( "empty client_id" ); return context.Status( fiber.StatusBadRequest ).JSON( fiber.Map{ "changes": changes , } ) }
	sequence_id := context.Get( HEADER_SEQUENCE_ID )
	if sequence_id == "" { sequence_id = "0" }

	var changes_map map[string]string = make( map[string]string )

	s.DB.View( func( tx *bolt.Tx ) error {
		changed_bucket := tx.Bucket( []byte( "changed" ) )
		if changed_bucket == nil { return nil }
		c := changed_bucket.Cursor()

		if sequence_id == "0" {
			for k , v := c.First(); k != nil; k, v = c.Next() {
				fmt.Println( "  found changed user - sequence:" , string( k ) , ":: uuid:" , string( v ) )
				changes_map[ string( k ) ] = string( v )
			}
		} else {
			c.Seek( []byte( sequence_id ) )
			for k , v := c.Next(); k != nil; k, v = c.Next() {
				fmt.Println( "  found changed user - sequence:" , string( k ) , ":: uuid:" , string( v ) )
				changes_map[ string( k ) ] = string( v )
			}
		}

		return nil
	})

	for k , v := range changes_map {
		change := Change{
			ID: k ,
			UUID: v ,
		}
		changes = append( changes , change )
	}

	return context.Status( 200 ).JSON( fiber.Map{
		"changes": changes ,
	})
}

func ( s *Server ) DownloadUser( context *fiber.Ctx ) ( error ) {
	uuid := context.Get( HEADER_UUID_KEY )
	fmt.Println( "DownloadUser - uuid:" , uuid )
	var user_bytes []byte
	s.DB.View( func( tx *bolt.Tx ) error {
		users_bucket := tx.Bucket( []byte( "users" ) )
		user_bytes = users_bucket.Get( []byte( uuid ) )
		return nil
	})
	return context.Status( 200 ).JSON( fiber.Map{
		"uuid": uuid ,
		"user_bytes": user_bytes ,
	})
}