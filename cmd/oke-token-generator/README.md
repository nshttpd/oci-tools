#### oke-token-generator

OKE utilizes the version 2.0.0 of the kubectl config. This means that every time a kubectl command
is executed against the masters a token needs to be generated via the OCI CLI. This is very slow
and painful. This little command line utility can be used in place of the OCI CLI to generate this
token for speedier kubectl access.

