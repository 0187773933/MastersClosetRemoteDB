package server

import (
	"fmt"
	"strings"
	fiber "github.com/gofiber/fiber/v2"
	// encryption "github.com/0187773933/encryption/v1/encryption"
)

func ( s *Server ) ValidateAPIKey() ( fiber.Handler ) {
	return func( c *fiber.Ctx ) error {
		key := c.Get( HEADER_API_KEY )
		key = strings.TrimSpace( key )
		if len( key ) > 256 { key = key[ :256 ] }
		if key == "" {
			fmt.Println( "missing api key" )
			return c.Status( fiber.StatusUnauthorized ).JSON( fiber.Map{
				"result": false ,
			})
		}
		if key != s.Config.ServerAPIKey {
			fmt.Println( "invalid api key:" , key )
			return c.Status( fiber.StatusUnauthorized ).JSON( fiber.Map{
				"result": false ,
			})
		}
		return c.Next()
	}
}