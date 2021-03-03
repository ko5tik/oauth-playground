
# Deploying to kuberenetes

as kubernetes tries to pull docker images from public dockerhub ( we do not like to have it )
we have to push our image to local minikube registry

  
````shell
eval $(minikube -p minikube docker-env)

docker image build -t client-credentials-server .

kubectl create -f deployment.yaml

````

Now   **in this shell only**  docker will push images into kubernetes registry and everything shall work
