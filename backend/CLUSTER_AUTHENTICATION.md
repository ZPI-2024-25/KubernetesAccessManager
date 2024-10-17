# Out of Cluster Authentication

To run the backend outside of the cluster, use:

```
go run main.go
```

It will use the `~/.kube/config` file to authenticate with the cluster.

If you have non standard kubeconfig file, you can specify it this way:

```
go run main.go -kubeconfig=/kubeconfig/path
```

# In Cluster Authentication

To run the backend inside the cluster, use:

```
go run main.go -in-cluster
```