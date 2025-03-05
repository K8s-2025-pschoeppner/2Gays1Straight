# ctf
Capture the flag client and server for Kubernetes

# Repository
1. Create a personal access token (classic) in your developer settings
2. Use the docker do docker login ghcr.io
3. Tag your image with (docker tag ghcr.io/<user>/<image_name>)
4. Use docker push with your tag name

# Images
**Server:** ghcr.io/maytastico/2gays1straight

**Client:** ghcr.io/maytastico/2gays1straight-client

# Secret

To pull the secret add this to your kubefile

ghcr-secret-2gays-one-straight

```
spec:
  imagePullSecrets:
  - name: ghcr-secret-2gays-one-straight
```

# kubectl commands
```
create service account
kubectl --kubeconfig secret-secret.txt create serviceaccount straightmin

create role with read access, and add straightmin account to it
kubectl --kubeconfig secret-secret.txt create rolebinding straightmin-view --clusterrole=view --serviceaccount=2gays1straight:straightmin

```
