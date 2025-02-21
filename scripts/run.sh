export $(grep -v '^#' .env | xargs)

kubectl create configmap env-map --from-env-file=.env

envsubst < deploy.yaml | kubectl apply -f -

unset $(grep -v '^#' .env | sed -E 's/(.*)=.*/\1/' | xargs)