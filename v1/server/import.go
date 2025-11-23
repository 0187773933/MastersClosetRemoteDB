package server

import (
	"fmt"
	"time"
	"bytes"
	"net/http"
	"io"
	"encoding/json"
	ulid "github.com/oklog/ulid/v2"
	bolt "github.com/boltdb/bolt"
	fiber "github.com/gofiber/fiber/v2"
	// utils "github.com/0187773933/MastersClosetRemoteDB/v1/utils"
)

func ( s *Server ) ImportUser( ctx *fiber.Ctx ) ( error ) {
	result := false
	uuid := ctx.Get( HEADER_UUID_KEY )
	if uuid == "" { fmt.Println( "empty uuid" ); return ctx.Status( fiber.StatusBadRequest ).JSON( fiber.Map{ "result": result , } ) }
	client_id := ctx.Get( HEADER_CLIENT_ID )
	if client_id == "" { fmt.Println( "empty client_id" ); return ctx.Status( fiber.StatusBadRequest ).JSON( fiber.Map{ "result": result , } ) }
	// clients_changed_bucket_name := fmt.Sprintf( "changed-client-%s" , client_id )
	changed_bucket_name := fmt.Sprintf( "changed" )
	body := ctx.Body()
	if len( body ) == 0 { fmt.Println( "empty body" ); return ctx.JSON( fiber.Map{ "result": result , } ) }
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
		// if changed_bucket.Stats().KeyN > s.Config.MaxTrackedChanges {
		// 	c := changed_bucket.Cursor()
		// 	for k , _ := c.First(); k != nil && changed_bucket.Stats().KeyN > s.Config.MaxTrackedChanges; k , _ = c.Next() {
		// 		changed_bucket.Delete( k )
		// 	}
		// }
		fmt.Println( "ImportUser - tracking change for user:" , client_id , "::" , uuid , " with sequence:" , sequence )
		return nil
	})
	if db_result != nil { fmt.Println( db_result ); return ctx.Status( 500 ).JSON( fiber.Map{ "result": result } ) }
	result = true
	if s.Config.MirrorToGlobal == true {
		go func() {
			timer := time.AfterFunc( 30*time.Second , func() {
				fmt.Println( "MirrorToGlobal still running after 30s — giving up" )
			})
			defer timer.Stop()
			mtgr := s.MirrorToGlobal( uuid , &body )
			fmt.Println( "MirrorToGlobal result for user:" , uuid , "::" , mtgr )
			timer.Stop()
		}()
	}
	return ctx.Status( 200 ).JSON( fiber.Map{
		"result": result ,
		"sequence": sequence ,
	})
}

type UploadResult struct {
	Result bool `json:"result"`
	Sequence string `json:"sequence"`
}
func ( s *Server ) MirrorToGlobal( uuid string , u_bytes *[]byte ) (  result UploadResult ) {
	http_client := &http.Client{ Timeout: 10 * time.Second }
	req , err := http.NewRequest( "POST" , s.Config.MirrorHostUrl + "/admin/user/import" , bytes.NewReader( *u_bytes ) )
	if err != nil {
		fmt.Println( err )
		return
	}
	req.Header.Set( "Content-Type" , "application/json" )
	q := req.URL.Query()
	q.Add( "k" , s.Config.MirrorHostAPIKey )
	q.Add( "uuid" , uuid )
	req.URL.RawQuery = q.Encode()
	// req.Header.Set( fmt.Sprintf( "%s-CLIENT-ID" , rs.CONFIG.RemoteHostHeaderPrefix ) , rs.CONFIG.RemoteHostClientID )
	// req.Header.Set( fmt.Sprintf( "%s-UUID" , rs.CONFIG.RemoteHostHeaderPrefix ) , string( *uuid ) )
	// req.Header.Set( fmt.Sprintf( "%s-API-KEY" , rs.CONFIG.RemoteHostHeaderPrefix ) , rs.CONFIG.RemoteHostAPIKey )
	resp , err := http_client.Do( req )
	if err != nil {
		fmt.Println( err )
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		snippet , _ := io.ReadAll( io.LimitReader( resp.Body , 2048 ) )
		x := fmt.Errorf( "http %d: %s", resp.StatusCode , string( snippet ) )
		fmt.Println( x )
		return
	}
	body_bytes , _ := io.ReadAll( io.LimitReader( resp.Body , 1<<20 ) )
	if err := json.Unmarshal( body_bytes , &result ); err != nil {
		x := fmt.Errorf( "decode response: %w" , err )
		fmt.Println( x )
		return
	}
	return
}