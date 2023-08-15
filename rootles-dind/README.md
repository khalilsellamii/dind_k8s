# Trying the rootless Docker in Docker image in kubernetes

## Architecture

<p align="center">
<img src="https://github.com/khalilsellamii/dind_k8s/blob/main/mtls_app/Screenshot%202023-08-15%20113901.png" alt="Alt text" width="650" height="400">
</p>  

## Documentation

<p align="center">
<img src="https://github.com/khalilsellamii/dind_k8s/blob/main/rootles-dind/image.png" alt="Alt text" width="650" height="400">
</p>  


## Implementation

Now, we create our deployment.yaml file again with new specifications namely :

```
 containers:
        - name: dind-daemon
          image: docker:dind-rootless
```

This what should the results look like ! 

<p align="center">
<img src="https://github.com/khalilsellamii/dind_k8s/blob/main/rootles-dind/e.png" alt="Alt text" width="750" height="200">
</p>  
