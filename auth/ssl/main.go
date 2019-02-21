/*
This is an example application to demonstrate parsing an ID Token.
*/
package main

import (
	"log"
        "fmt"

	oidc "github.com/coreos/go-oidc"

	"golang.org/x/net/context"
)

var rawIDToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjZiNDZmMTljZDhlYTRjN2ViZDI0MTZmYmMxMDA4NWEyIn0.eyJpc3MiOiJodHRwczovL2RleC5leGFtcGxlLmNvbTo4MDgwIiwic3ViIjoiQ2djNE9URXlOVFUzRWdabmFYUm9kV0kiLCJhdWQiOiJleGFtcGxlLWFwcCIsImV4cCI6MTU1MDcxMjgxMywiaWF0IjoxNTUwNzEyMjEzLCJlbWFpbCI6ImZodGpvYkBob3RtYWlsLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJncm91cHMiOlsiZGV2Il0sIm5hbWUiOiJmYW51eCJ9.OLrfnV0Ej9PwWCUpP3BferHJTqxigHYkMyfqvc3TFRKLWl-alH9MwA9Fo0T3Vj86kRxwHqu8BbOykKJJOnO3Agij8dInEGNu4ffb2mjo3WmJwUvWAcFuEV5GzI9UAkuK8QdnMiK-nBMy0givHwl4tVJergKEeABjUqX9WyRHjKDHPI-pGI4oZYTSek8kJpbGm1sBx9NIcWM6exICPWEKXI69heHNx76jgtLVV9pZdNJTw_TipL9svkrD3tEDeQHa47awkFNImXQGhom-d7QwbCmR2wXt4c79n2tG9R4-jffBX2iAoSuELBTAoVFFOeYhD2ZmEIBiS7ZpLqxmk_jnWw"


func main() {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "https://dex.example.com:8080")
	if err != nil {
		log.Fatal(err)
	}
	oidcConfig := &oidc.Config{
		ClientID: "test",
	}
	verifier := provider.Verifier(oidcConfig)

	idToken,err := verifier.Verify(ctx, rawIDToken)
        if err != nil {
		fmt.Println("err: ", err)
	}
		fmt.Println("idtoken: ", idToken)

}
