package utils

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"net"
	types "github.com/0187773933/MastersClosetRemoteDB/v1/types"
)

func ParseConfig( file_path string ) ( result types.ConfigFile ) {
	file_data , _ := ioutil.ReadFile( file_path )
	err := json.Unmarshal( file_data , &result )
	if err != nil { fmt.Println( err ) }
	return
}

// https://stackoverflow.com/a/28862477
func GetLocalIPAddresses() ( ip_addresses []string ) {
	host , _ := os.Hostname()
	addrs , _ := net.LookupIP( host )
	encountered := make( map[ string ]bool )
	for _ , addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			ip := ipv4.String()
			if !encountered[ ip ] {
				encountered[ ip ] = true
				ip_addresses = append( ip_addresses , ip )
			}
		}
	}
	return
}