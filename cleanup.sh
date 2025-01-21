export $(grep -v '^#' .env | xargs)

kubectl delete pod/master-pod
kubectl delete deployment/worker-deployment
kubectl delete service/$MASTER_HOST

minikube delete configmap/env-map

minikube stop

unset $(grep -v '^#' .env | sed -E 's/(.*)=.*/\1/' | xargs)