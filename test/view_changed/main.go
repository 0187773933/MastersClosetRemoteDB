package main

import (
	"fmt"
	bolt "github.com/boltdb/bolt"
)

func main() {
	db , db_err := bolt.Open( "mct_remote.db" , 0600 , &bolt.Options{} )
	if db_err != nil { panic( db_err.Error() ) }
	db.View( func( tx *bolt.Tx ) error {
		b := tx.Bucket( []byte( "changed" ) )
		if b == nil { fmt.Println( "no changed bucket" ); return nil }

		c := b.Cursor()
		for k , v := c.First(); k != nil; k, v = c.Next() {
			fmt.Println( "changed:" , string( k ) , "->" , string( v ) )
		}
		return nil
	})
}