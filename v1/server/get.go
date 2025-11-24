package server

import (
	"fmt"
	"encoding/json"
	bolt "github.com/boltdb/bolt"
	fiber "github.com/gofiber/fiber/v2"
	// utils "github.com/0187773933/MastersClosetRemoteDB/v1/utils"
	encryption "github.com/0187773933/encryption/v1/encryption"
	user "github.com/0187773933/MastersCloset/v1/user"
)

func ( s *Server ) GetUser( ctx *fiber.Ctx ) ( error ) {
	uuid := ctx.Get( HEADER_UUID_KEY )
	if uuid == "" { fmt.Println( "empty uuid" ); return ctx.Status( fiber.StatusBadRequest ).JSON( fiber.Map{ "result": false , } ) }
	var viewed_user user.User
	s.DB.View( func( tx *bolt.Tx ) error {
		users_bucket := tx.Bucket( []byte( "users" ) )
		user_encrypted_bytes := users_bucket.Get( []byte( uuid ) )
		if len( user_encrypted_bytes ) < 1 {
			return ctx.Status( 200 ).JSON( fiber.Map{
				"user": viewed_user ,
			})
		}
		user_bytes := encryption.ChaChaDecryptBytes( s.Config.BoltDBEncryptionKeyClients , user_encrypted_bytes )
		json.Unmarshal( user_bytes , &viewed_user )
		return nil
	})
	return ctx.Status( 200 ).JSON( fiber.Map{
		"user": viewed_user ,
	})
}