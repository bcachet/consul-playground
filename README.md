# Consul


Little playground to play with HAProxy auto-scaling backends via [server-template and Consul service discovery](https://developer.hashicorp.com/consul/docs/discover/load-balancer/ha)

We have setup an HAProxy server as load-balancer for our `backend` service (a Go HTTP server with and `/health` endpoint).

You can start the lab via `docker compose`:
```sh
docker compose up --detach --wait
curl localhost:8000/health
```

Visualization
* Consul UI is accessible through [http://localhost:8500/ui](http://localhost:8500/ui)
* HAProxy stats dashboard is accessible through [http://localhost:8404](http://localhost:8404)


We output the name of the host that handle the query in the response.

```sh
watch --precise --interval 1 curl localhost:8000/health
```

We can play with the backend services via

```sh
docker compose {stop,start} backend-{1,2}
```
