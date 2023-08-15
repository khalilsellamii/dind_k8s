# Docker in Docker in kubernetes without mTLS

## Architecture

<p align="center">
<img src="https://github.com/khalilsellamii/dind_k8s/blob/main/my_dind_app_without_mtls/image.png" alt="Alt text" width="650" height="400">
</p>  

## Implementation
we need to modify our Golang code to configure and set up our client, and create a new dockerSDK client with the desired specifications
```
- // Create a new Docker client
-	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

+ //Create a new docker client
+ cli, err := client.NewClientWithOpts(client.FromEnv, 
  client.WithHost("tcp://172.17.0.2:2376"))
```

Now, we create our deployment.yaml file again with ne specifications namely :

`+` env variable DOCKER_TLS_CERTDIR =`""`   
⇒ This variable set to “” will allow us to disable the auto generation of certificates and access the docker socket without mtls process.  
`+` env variable docker_host to seek `2375`  

This what should the results look like ! 

<p align="center">
<img src="https://github.com/khalilsellamii/dind_k8s/blob/main/my_dind_app_without_mtls/a.png" alt="Alt text" width="750" height="400">
</p>  
