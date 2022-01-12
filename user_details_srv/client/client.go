package main

import (
	"context"
	"fmt"
	"log"

	"github.com/mauricioww/user_microsrv/user_details_srv/detailspb"
	"google.golang.org/grpc"
)

func main() {

	var grpc_client *grpc.ClientConn
	var grpc_err error
	{
		grpc_client, grpc_err = grpc.Dial("0.0.0.0:50052", grpc.WithInsecure())
		if grpc_err != nil {
			log.Fatal(grpc_err)
		}
	}

	client := detailspb.NewUserDetailsServiceClient(grpc_client)

	set(client, 0)
}

func set(c detailspb.UserDetailsServiceClient, id int) {
	req := &detailspb.SetUserDetailsRequest{
		// UserId:       uint32(id),
		// Country:      "Mexico",
		// City:         "CDMX",
		// MobileNumber: "0000000001",
		// Married:      true,
		// Height:       1.75,
		// Weigth:       76.0,
	}

	res, err := c.SetUserDetails(context.TODO(), req)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
