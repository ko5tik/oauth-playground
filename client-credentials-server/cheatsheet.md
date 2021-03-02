
# Deploying to cuberenetes

as kubernetes tries to pull docker images from public dockerhub ( we do not like to have it )
we have toi push our image to local minikube registry

  
````shell
eval $(minikube -p minikube docker-env)

docker image build -t client-credentials-server .

kubectl create -f deployment.yaml

````

Now   __in this shell only__  docker will push images into kubernetes registry and everything shall work
