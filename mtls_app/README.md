# Dockerin DOcker in kubernetes with mTLS

## Architecture

<p align="center">
<img src="https://github.com/khalilsellamii/dind_k8s/blob/main/mtls_app/Screenshot%202023-08-15%20113901.png" alt="Alt text" width="650" height="400">
</p>  

## Implementation
we need to modify our Golang code to configure and set up our mtls client, and recreate a new dockerSDK client with the desired specifications
```
+ cert, err := tls.LoadX509KeyPair("./certs/client/cert.pem", "./certs/client/key.pem")
+	 if err != nil {
+		log.Fatal("Failed to load client TLS certificates:", err)
+ 	}

+ 	caCert, err := ioutil.ReadFile("./certs/ca/cert.pem")
+ 	if err != nil {
+ 		log.Fatal("Failed to read CA certificate:", err)
+	 }

+ 	caCertPool := x509.NewCertPool()
+ 	caCertPool.AppendCertsFromPEM(caCert)

+ 	tlsConfig := &tls.Config{
+ 		Certificates: []tls.Certificate{cert},
+ 		RootCAs:      caCertPool,
+ 	}
+ 	// Create a new HTTP client with mTLS configuration
+ 	httpClient := &http.Client{
+ 		Transport: &http.Transport{
+ 			TLSClientConfig: tlsConfig,
+ 		},
+ 	}

+ 	// Create a new Docker client with custom httpClient using mTLS
+ 	cli, err := client.NewClientWithOpts(
-     client.WithHost("tcp://172.17.0.2:2375")
+ 		client.WithHost("tcp://172.17.0.2:2376"),
+ 		client.WithHTTPClient(httpClient),
+ 	)
```

Now, we create our deployment.yaml file again with ne specifications namely :

`+` env variable DOCKER_TLS_CERTDIR =`/certs`   
`+` the volume mounts for the certificates to be shared between the containers in the mtls handshake process (2 volumes ⇒ `/certs/ca` and  `/certs/client`)  
`+` env variable docker_host to seek `2376` not `2375`  

This what should the results look like ! 

<p align="center">
<img src="https://github.com/khalilsellamii/dind_k8s/blob/main/mtls_app/image.png" alt="Alt text" width="1000" height="400">
</p>  
