# Deploying a multi container pod containing dind + golangApp images:

# In this demo, we provided the docker socket for our golang application via port tcp://[...]:2376
# and we used TLS/mTLS as verification and a way for securing the communication between the 2
# containers. 
