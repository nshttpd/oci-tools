#### oke-token-generator

OKE utilizes the version 2.0.0 of the kubectl config. This means that every time a kubectl command
is executed against the masters a token needs to be generated via the OCI CLI. This is very slow
and painful. This little command line utility can be used in place of the OCI CLI to generate this
token for speedier kubectl access.

It uses the same config file that the OCI CLI uses for credentials to interact with the OKE API.

In the generated kube config file for the cluster you'll see :

```
user:
  exec:
    apiVersion: client.authentication.k8s.io/v1beta1
    args:
    - ce
    - cluster
    - generate-token
    - --cluster-id
    - <cluster ocid>
    command: oci
    env: []
```

to use the quicker token generator you will replace that with :

```
user:
  exec:
    apiVersion: client.authentication.k8s.io/v1beta1
    args:
    - --cluster-id
    - <cluster ocid>
    - --region
    - <region of cluster>
    command: oke-token-generator
    env: []
```

If the tenancy isn't under `[DEFAULT]`, but a different profile you add another parameter called `--profile`
to the above configuration.

```
user:
  exec:
    apiVersion: client.authentication.k8s.io/v1beta1
    args:
    - --cluster-id
    - <cluster ocid>
    - --region
    - <region of cluster>
    - --profile
    - <profile>
    command: oke-token-generator
    env: []
```

Before with OCI CLI :

```
$ time kubectl get nodes
NAME         STATUS   ROLES   AGE   VERSION

real	0m1.609s
user	0m0.776s
sys	0m0.436s
```

With oke-token-generator :

```
$ time kubectl get nodes
NAME         STATUS   ROLES   AGE   VERSION

real	0m0.926s
user	0m0.108s
sys	0m0.178s
```
