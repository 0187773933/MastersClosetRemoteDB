package server

import (
	"fmt"
	ulid "github.com/oklog/ulid/v2"
	bolt "github.com/boltdb/bolt"
	fiber "github.com/gofiber/fiber/v2"
	// utils "github.com/0187773933/MastersClosetRemoteDB/v1/utils"
)

func ( s *Server ) ImportUser( context *fiber.Ctx ) ( error ) {
	result := false
	uuid := context.Get( HEADER_UUID_KEY )
	if uuid == "" { fmt.Println( "empty uuid" ); return context.Status( fiber.StatusBadRequest ).JSON( fiber.Map{ "result": result , } ) }
	client_id := context.Get( HEADER_CLIENT_ID )
	if client_id == "" { fmt.Println( "empty client_id" ); return context.Status( fiber.StatusBadRequest ).JSON( fiber.Map{ "result": result , } ) }
	// clients_changed_bucket_name := fmt.Sprintf( "changed-client-%s" , client_id )
	changed_bucket_name := fmt.Sprintf( "changed" )
	body := context.Body()
	if len( body ) == 0 { fmt.Println( "empty body" ); return context.JSON( fiber.Map{ "result": result , } ) }
	sequence := ""
	db_result := s.DB.Update( func( tx *bolt.Tx ) error {
		users_bucket , users_bucket_err := tx.CreateBucketIfNotExists( []byte( "users" ) )
		if users_bucket_err != nil { fmt.Println( users_bucket_err ); return users_bucket_err }
		user_store_result := users_bucket.Put( []byte( uuid ) , body )
		if user_store_result != nil { fmt.Println( user_store_result ); return user_store_result }
		changed_bucket , changed_bucket_err := tx.CreateBucketIfNotExists( []byte( changed_bucket_name ) )
		if changed_bucket_err != nil { fmt.Println( changed_bucket_err ); return changed_bucket_err }
		// sequence , _ = changed_bucket.NextSequence()
		sequence = ulid.Make().String()
		changed_bucket_store_result := changed_bucket.Put( []byte( sequence ) , []byte( uuid ) )
		if changed_bucket_store_result != nil { fmt.Println( changed_bucket_store_result ); return changed_bucket_store_result }
		if changed_bucket.Stats().KeyN > s.Config.MaxTrackedChanges {
			c := changed_bucket.Cursor()
			for k , _ := c.First(); k != nil && changed_bucket.Stats().KeyN > s.Config.MaxTrackedChanges; k , _ = c.Next() {
				changed_bucket.Delete( k )
			}
		}
		fmt.Println( "ImportUser - tracking change for user:" , client_id , "::" , uuid , " with sequence:" , sequence )
		return nil
	})
	if db_result != nil { fmt.Println( db_result ); return context.Status( 500 ).JSON( fiber.Map{ "result": result } ) }
	result = true
	return context.Status( 200 ).JSON( fiber.Map{
		"result": result ,
		"sequence": sequence ,
	})
}