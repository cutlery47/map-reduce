# export
export $(grep -v '^#' .env | xargs)

# unset
unset $(grep -v '^#' .env | sed -E 's/(.*)=.*/\1/' | xargs)

# subst 
envsubst < deploy.yaml | kubectl apply -f -