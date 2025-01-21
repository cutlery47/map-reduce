export $(grep -v '^#' .env | xargs)

minikube start

minikube mount $LOCAL_FILE_DIRECTORY:$MINIKUBE_FILE_DIRECTORY

unset $(grep -v '^#' .env | sed -E 's/(.*)=.*/\1/' | xargs)