## WIP Game

1. Redis pub/sub for game server websockets
2. Postgres for player info etc.
3. Some UI? (htmx/react/?)
4. helm, redis, postgres

### Redis (https://artifacthub.io/packages/helm/bitnami/redis)
1. `helm install my-release oci://registry-1.docker.io/bitnamicharts/redis`
2. `export REDIS_PASSWORD=$(kubectl get secret --namespace default my-release-redis -o jsonpath="{.data.redis-password}" | base64 -d)`
3. `kubectl port-forward --namespace default svc/my-release-redis-master 6379:6379 &`
4. `REDISCLI_AUTH="$REDIS_PASSWORD" redis-cli -h 127.0.0.1 -p 6379`
