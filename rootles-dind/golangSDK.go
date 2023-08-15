package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func main() {

	//BEGIN_PARTIE MTLS & CERTIFICATES
	// ca_certsDir := "./certs/ca"
	// client_certsDir := "./certs/client"
	// clientCerts, err := tlsconfig.Client(certprovider.NewFileProvider(client_certsDir, "cert.pem", "key.pem", ""))
	// if err != nil {
	// 	// Handle error
	// }

	// caCert, err := tlsconfig.CACertPool(certprovider.NewFileProvider(ca_certsDir, "cert.pem"))
	// if err != nil {
	// 	// Handle error
	// }

	// tlsConfig := &tlsconfig.Options{
	// 	InsecureSkipVerify: false, // Set to true if you want to skip server certificate verification (not recommended for production)
	// 	CAFile:             "",    // Not needed since we are using the caCert from above
	// 	CertFile:           "",    // Not needed since we are using clientCerts from above
	// 	KeyFile:            "",    // Not needed since we are using clientCerts from above
	// 	CABytes:            caCert,
	// 	CertBytes:          clientCerts.Certificate,
	// 	KeyBytes:           clientCerts.PrivateKey,
	// }
	// //END_PARTIE MTLS & CERTIFICATES
	// ctx := context.Background()

	// cli, err := client.NewClientWithOpts(
	// 	client.FromEnv, // Adjust the API version as needed
	// 	client.WithHost("tcp://localhost:2376"),
	// 	client.WithTLS(tlsConfig),
	// 	client.WithHTTPClient(&http.Client{}), // You can customize the HTTP client if needed
	// )
	// if err != nil {
	// 	// Handle error
	// }

	ctx := context.Background()
	// Create a new Docker client
	//cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	//cli, err := client.NewClientWithOpts(client.FromEnv, client.WithHost("tcp://localhost:2375"))
	//if err != nil {
	//	log.Fatal(err)
	//}
	// Load client TLS certificates and create a TLS configuration
	cert, err := tls.LoadX509KeyPair("./certs/client/cert.pem", "./certs/client/key.pem")
	if err != nil {
		log.Fatal("Failed to load client TLS certificates:", err)
	}

	caCert, err := ioutil.ReadFile("./certs/ca/cert.pem")
	if err != nil {
		log.Fatal("Failed to read CA certificate:", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	// Create a new HTTP client with mTLS configuration
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	// Create a new Docker client with custom httpClient using mTLS
	cli, err := client.NewClientWithOpts(
		client.WithHost("tcp://172.17.0.2:2376"),
		client.WithHTTPClient(httpClient),
	)

	if err != nil {
		log.Fatal("Failed to create Docker client:", err)
	}
	fmt.Println("Connected to docker socket with mtls on port 2376")
	// Pull the nginx image
	imageName := "nginx"
	resp, err := cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Close()

	// Print the status of the pull operation
	fmt.Println("Pulling nginx image...")
	_, err = io.Copy(os.Stdout, resp)
	if err != nil {
		log.Fatal(err)
	}
	//create a container from the pulled image
	resp1, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "nginx",
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}
	//start the previously created container
	if err := cli.ContainerStart(ctx, resp1.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	//pritnt the status of the container
	statusCh, errCh := cli.ContainerWait(ctx, resp1.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp1.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	// Check if the pull operation was successful
	imageExists, err := checkImageExists(cli, imageName)
	if err != nil {
		log.Fatal(err)
	}

	if imageExists {
		fmt.Println("Ubuntu image pulled successfully!")
	} else {
		fmt.Println("Failed to pull the Ubuntu image.")
	}
}

func checkImageExists(cli *client.Client, imageName string) (bool, error) {
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return false, err
	}

	for _, image := range images {
		for _, tag := range image.RepoTags {
			if tag == imageName {
				return true, nil
			}
		}
	}

	return false, nil
}
