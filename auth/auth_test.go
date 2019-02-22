package main

/*
	fmt.Println("=========veriry========\n\n")
	object, err := jose.ParseSigned(idToken)
	if err != nil {
		fmt.Printf("parse signed failed: %s", err)
	}

	// Now we can verify the signature on the payload. An error here would
	// indicate the the message failed to verify, e.g. because the signature was
	// broken or the message was tampered with.
	priv, ok := Priv.Key.(*rsa.PrivateKey)
	if !ok {
		fmt.Println("to rsa private key failed")
	}
	output, err := object.Verify(priv.Public())
	if err != nil {
		fmt.Printf("Verify failed: %s", err)
	}

	fmt.Printf(string(output))
	fmt.Println("=========veriry========\n\n")

	fmt.Println("=========oidc verify========\n\n")
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "https://fist.sealyun.svc.cluster.local:8080")
	if err != nil {
		fmt.Println("new provider failed: ", err)
	}
	oidcConfig := &oidc.Config{
		ClientID: "sealyun-fist",
	}
	verifier := provider.Verifier(oidcConfig)
	token, err := verifier.Verify(ctx, idToken)
	if err != nil {
		fmt.Println("valify failed: ", err)
	}
	fmt.Println("token: ", token)
	fmt.Println("=========oidc verify========\n\n")

*/
