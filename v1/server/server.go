package server

import (
	"fmt"
	bolt "github.com/boltdb/bolt"
	fiber "github.com/gofiber/fiber/v2"
	fiber_cookie "github.com/gofiber/fiber/v2/middleware/encryptcookie"
	fiber_cors "github.com/gofiber/fiber/v2/middleware/cors"
	types "github.com/0187773933/MastersClosetRemoteDB/v1/types"
	utils "github.com/0187773933/MastersClosetRemoteDB/v1/utils"
)

type Server struct {
	FiberApp *fiber.App `json:"fiber_app"`
	Config types.ConfigFile `json:"config"`
	DB *bolt.DB `json:"_"`
}

func New( config types.ConfigFile , db *bolt.DB ) ( server Server ) {

	server.FiberApp = fiber.New()
	server.Config = config
	server.DB = db

	server.FiberApp.Use( server.LogRequest )
	// server.FiberApp.Use( favicon.New( favicon.Config{
	// 	File: "./v1/server/cdn/favicon.ico" ,
	// }))

	// temp_key := fiber_cookie.GenerateKey()
	// fmt.Println( temp_key )
	server.FiberApp.Use( fiber_cookie.New( fiber_cookie.Config{
		Key: server.Config.ServerCookieSecret ,
	}))

	allow_origins_string := fmt.Sprintf( "%s, %s" , server.Config.ServerBaseUrl , server.Config.ServerLiveUrl )
	fmt.Println( "Using Origins:" , allow_origins_string )
	server.FiberApp.Use( fiber_cors.New( fiber_cors.Config{
		AllowOrigins: allow_origins_string ,
		AllowHeaders:  "Origin, Content-Type, Accept, key" ,
		AllowCredentials: true ,
	}))

	server.SetupRoutes()
	// server.FiberApp.Get( "/*" , func( context *fiber.Ctx ) ( error ) { return context.Redirect( "/" ) } )
	return
}

func ( s *Server ) LogRequest( context *fiber.Ctx ) ( error ) {
	ip_address := context.Get( "x-forwarded-for" )
	if ip_address == "" { ip_address = context.IP() }
	// log_message := fmt.Sprintf( "%s === %s === %s === %s === %s" , time_string , GlobalConfig.FingerPrint , ip_address , context.Method() , context.Path() )
	// log_message := fmt.Sprintf( "%s === %s === %s === %s" , time_string , GlobalConfig.FingerPrint , context.Method() , context.Path() )
	// log_message := fmt.Sprintf( "%s === %s" , context.Method() , context.Path() )
	log_message := fmt.Sprintf( "%s === %s === %s" , ip_address , context.Method() , context.Path() );
	// log.Println( log_message )
	fmt.Println( log_message )
	return context.Next()
}

func ( s *Server ) SetupRoutes() {

}

func ( s *Server ) Start() {
	fmt.Println( "\n" )
	local_ips := utils.GetLocalIPAddresses()
	for _ , x_ip := range local_ips {
		fmt.Printf( "Listening on http://%s:%s\n" , x_ip , s.Config.ServerPort )
	}
	fmt.Println( "Listening on http://localhost:%s\n" , s.Config.ServerPort )
	fmt.Printf( "Admin Login @ http://localhost:%s/admin/login\n" , s.Config.ServerPort )
	fmt.Printf( "Admin Login @ %s/admin/login\n" , s.Config.ServerLiveUrl )
	fmt.Printf( "Admin Username === %s\n" , s.Config.AdminUsername )
	fmt.Printf( "Admin Password === %s\n" , s.Config.AdminPassword )
	s.FiberApp.Listen( fmt.Sprintf( ":%s" , s.Config.ServerPort ) )
}