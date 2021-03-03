
# Deploying to kuberenetes

as kubernetes tries to pull docker images from public dockerhub ( we do not like to have it )
we have to push our image to local minikube registry

  
````shell
eval $(minikube -p minikube docker-env)

docker image build -t client-credentials-server .

kubectl create -f deployment.yaml

````

Now   **in this shell only**  docker will push images into kubernetes registry and everything shall work


### expose internal port 

kubectl expose deployment/client-credentials-server --type="NodePort" --port 3846
````shell
ko5tik@firefly ~/customer/vw/oauth-playground/client-credentials-server $ kubectl get services
NAME                        TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
client-credentials-server   NodePort    10.106.212.30   <none>        3846:30245/TCP   6s
kubernetes                  ClusterIP   10.96.0.1       <none>        443/TCP          47h

````

Now is it reachable at 

ko5tik@firefly ~ $ minikube ip
192.168.49.2

