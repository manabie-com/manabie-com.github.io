+++
date = "2021-11-23T14:28:23+07:00"
author = "nvcnvn"
description = "Simulate Let’s Encrypt HTTPS certificate issuing with HTTP-01 challenge in your local k8s cluster"
title = "Simulate Let’s Encrypt certificate issuing in local Kubernetes"
categories = ["DevSecOps", "Infrastructure", "Security"]
tags = ["k8s", "cert-manager", "CoreDNS", "acme", "http-01"]
slug = "simulate-https-certificates-acme-k8s"
+++

We write test to make sure our code work as expected, no matter that a Go code or YAML config. This series we will 
show how we develop and do integration test using local k8s. The first part will show how we're simulate the HTTPs request 
flow to allow our Platform engineers test their config and allow our Front-end engineers connect their application to 
local server with a self-signed HTTPs certificate.  

#### Let’s Encrypt ACME & HTTP-01 challenge for dummy
Let’s Encrypt is a internet goodies, talk to them nicely and they will give an free HTTPS certificate for you. We're using the famous cert-manager to speak to them, the process can simplify as:
- **You**: buy and then point `example.com` to your sever.
- **Our cert-manager** (via API requests): Hello Mr. Let’s Encrypt, can I have a free HTTPS certificate do my `example.com` domain?
- **Let’s Encrypt** (via API responses): Sure thing! But first I need to know if you owned `example.com`. Here I have a secure-random-token, put this secure-random-token to this path `http://example.com/.well-known/acme-challenge/<SECURE-RANDOM-TOKEN>` and make sure I can access and view it.
- **Our cert-manager**: Its is done sir, can you check?
- **Let’s Encrypt**: Open `http://example.com/.well-known/acme-challenge/<SECURE-RANDOM-TOKEN>` in my browser, seem legit. OK, here is the key for your HTTPS certificate.

Let's Encrypt provide a better description or you can read the RFC 8555 for not so dummy:
- https://letsencrypt.org/docs/challenge-types/#http-01-challenge
- https://datatracker.ietf.org/doc/html/rfc8555#section-4

#### Setup your minikube
Cert-manager have a great tutorial here:
- https://cert-manager.io/docs/tutorials/acme/ingress/#step-2-deploy-the-nginx-ingress-controller
You should follow the post, they provide many explaination.
Their post assume you're testing on a real k8s cluster (GKE or any public cloud provider offer some free resource for testing). 
Our post is for minikube, some dependencies required:
- minikube (we're using v1.24.0)
- helm (v3.7.0)
```bash
minikube start # then wait
minikube addons enable ingress
# you can find this kuard.yaml in examples folder of this post
kubectl apply -f ./examples/kuard.yaml
kubectl apply -f ./examples/http-only-ingress.yaml
```
Everything should work as expected, note that we install everything in the same namespace for the simplicity of the post. 
Run `kubectl get pod -n default` and you should get everything ready and running like this:
```
kuard-5cd5556bc9-kxt6p                                 1/1     Running   0          14m
```
Checking the network services created by `kubectl get svc -n default`:
```
NAME         TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)                      AGE
kuard        ClusterIP      10.101.81.156   <none>        80/TCP                       22m
kubernetes   ClusterIP      10.96.0.1       <none>        443/TCP                      61m

```
Also, we should check for the ingress by `kubectl get svc -n ingress-nginx`:
```
NAME                                 TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
ingress-nginx-controller             NodePort    10.99.82.10      <none>        80:31563/TCP,443:30703/TCP   4m39s
ingress-nginx-controller-admission   ClusterIP   10.111.186.145   <none>        443/TCP                      4m39s
```
So we have:
- kuard: for the demo application
- ingress-nginx: play the role of network LoadBalancer for your cluster

If you're working with a cluster on GKE or EKS, they will create a real IP, a network LB and assign that to your cluster. 

See the new assigned NodePort `10.99.82.10`?
Then you can test thing out with `curl -H 'Host: example.example.com' 'http://10.98.223.158'` to receive a 404 page. 
Above command is equivalent with you modify your `hosts` file to access localhost via example.example.com. 
Next, let install cert-manager in the same namespace:
```bash
helm repo add jetstack https://charts.jetstack.io
helm install cert-manager jetstack/cert-manager --version v1.6.0 --set installCRDs=true
```
Check if stuff installed correctly:
```
NAME                                                   READY   STATUS    RESTARTS   AGE
cert-manager-6c576bddcf-hjdts                          1/1     Running   0          32m
cert-manager-cainjector-669c966b86-ggs9v               1/1     Running   0          32m
cert-manager-webhook-7d6cf57d55-mchqj                  1/1     Running   0          32m
kuard-5cd5556bc9-kxt6p                                 1/1     Running   0          46m
```
Then normally, if you are on a real cluster, you can create an `Issuer` object point to Let's Encrypt staging environment 
like this:
```yaml
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
    name: letsencrypt-staging
spec:
    acme:
    # The ACME server URL
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    # Email address used for ACME registration
    email: user@example.com
    # Name of a secret used to store the ACME account private key
    privateKeySecretRef:
        name: letsencrypt-staging
    # Enable the HTTP-01 challenge provider
    solvers:
    - http01:
        ingress:
            class:  nginx
```
Base on this configuration, cert-manager will speak to Let's Encrypt (staging in this case) to get the certificate. 
But Let's Encrypt cannot access to your sever since everything is just your local IP. Also, the staging environment still 
have some kind of rate-limit, you cannot use it if your CI/CD run really frequently. 
One solution for this, believe or not, is deploy your own fake Let's Encrypt to simulate the flow.

#### Introducing Pebble
https://github.com/letsencrypt/pebble
> A miniature version of Boulder, Pebble is a small ACME test server not suited for use as a production CA.
We have ready for you a minimal installation of Pebble ready for you already, let's install it in another namespace to 
simulate a different network :joy: 
```
kubectl create ns emulator
kubectl apply -f ./examples/pebble.yaml -n emulator
kubectl get pod -n emulator
```
and if everything run normally:
```
NAME                     READY   STATUS    RESTARTS   AGE
pebble-885bdd44c-cpv6r   1/1     Running   0          19s
```
Now go back to default namespace and install the issuer point to our local pebble:
```yaml
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: pebble-issuer
spec:
  acme:
    skipTLSVerify: true
    email: example@example.com
    server: https://pebble.emulator:14000/dir
    privateKeySecretRef:
      name: pk-pebble-issuer
    solvers:
      - selector:
        http01:
          ingress:
            class: nginx
```
Chicken and eggs issue here, our pebble itself don't have valid cert =]] so we `skipTLSVerify: true`. The `server` now 
pointed to `pebble` service an `emulator` namespace. We should check if the issuer config correctly by `kubectl get issuer`:
```
NAME            READY   AGE
pebble-issuer   True    10s
```
Now we should modify the installed http-only-ingress to have https: `kubectl apply -f ./examples/https-ingress.yaml`, you can see:
- kubectl get svc
```
NAME                        TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)                      AGE
cert-manager                ClusterIP      10.108.69.104   <none>          9402/TCP                     135m
cert-manager-webhook        ClusterIP      10.98.30.215    <none>          443/TCP                      135m
cm-acme-http-solver-8fm5v   NodePort       10.103.74.213   <none>          8089:30783/TCP               24m
kuard                       ClusterIP      10.101.81.156   <none>          80/TCP                       149m
```
cert-manager create a new `cm-acme-http-solver` service to handle the challenge verification, let check the challenge:
- kubectl get challenge
```
NAME                                           STATE     DOMAIN                AGE
pk-pebble-issuer-2n5gm-3378267180-1105783005   pending   example.example.com   3m43s
```
- kubectl describe challenge pk-pebble-issuer-2n5gm-3378267180-1105783005
```
Status:
  Presented:   true
  Processing:  true
  Reason:      Waiting for HTTP-01 challenge propagation: failed to perform self check GET request 'http://example.example.com/.well-known/acme-challenge/AfneYTxNeVkw25W2OPUcRPB0byhKfCwxisDFb9QJ9dw': Get "http://example.example.com/.well-known/acme-challenge/AfneYTxNeVkw25W2OPUcRPB0byhKfCwxisDFb9QJ9dw": dial tcp: lookup example.example.com on 10.96.0.10:53: no such host
  State:       pending
Events:
  Type    Reason     Age    From          Message
  ----    ------     ----   ----          -------
  Normal  Started    4m42s  cert-manager  Challenge scheduled for processing
  Normal  Presented  4m42s  cert-manager  Presented challenge using HTTP-01 challenge mechanism
```
Hmmm, I forgot, still the `example.example.com` is just a fake domain. Another hack needed for this blog post of hack. 
We can to modify the CoreDNS config so internal cluster can resolve to our internal IP, I call this a valid hack because 
its similar how the domain owner need to point the domain to the server's IP, and you don't need to keep a `minikube tunnel` open. 
I hope that when some one reading this post, k8s still use `CoreDNS` and installed it in the `kube-system` namespace: 
- kubectl -n kube-system describe configmap coredns
You will see the config file somewhat like this
```bash
ip=$(kubectl get svc ingress-nginx-controller --no-headers -n ingress-nginx | awk '{print$3}')
cat <<EOF | kubectl apply -f -
kind: ConfigMap
metadata:
  name: coredns
  namespace: kube-system
apiVersion: v1
data:
  Corefile: |
    .:53 {
        errors
        health {
           lameduck 5s
        }
        ready
        kubernetes cluster.local in-addr.arpa ip6.arpa {
           pods insecure
           fallthrough in-addr.arpa ip6.arpa
           ttl 30
        }
        prometheus :9153
        forward . 1.1.1.1
        cache 30
        loop
        reload
        loadbalance
    }
    example.example.com {
       hosts {
         $ip example.example.com
         fallthrough
       }
       whoami
    }
EOF
  kubectl delete pod -n kube-system --wait $(kubectl get pods -n kube-system | grep coredns | awk '{print$1}')
```
Delete everything (just for make sure):
```bash
kubectl delete -f ./examples/https-ingress.yaml
kubectl delete -f ./examples/pebble-issuer.yaml
kubectl delete -f ./examples/pebble.yaml -n emulator
```
and try again:
```bash
kubectl apply -f ./examples/pebble.yaml -n emulator
kubectl apply -f ./examples/pebble-issuer.yaml
kubectl apply -f ./examples/https-ingress.yaml
```
- kubectl get order
```
NAME                                      STATE   AGE
quickstart-example-tls-7k6zs-3378267180   valid   41s
```
- kubectl get cert
```
NAME                     READY   SECRET                   AGE
quickstart-example-tls   True    quickstart-example-tls   37s
```
Now let try to dial with HTTPS:
```bash
curl -H 'Host: example.example.com' https://$(minikube ip)
```
We should see some error:
```
curl: (60) SSL certificate problem: unable to get local issuer certificate
More details here: https://curl.haxx.se/docs/sslcerts.html

curl failed to verify the legitimacy of the server and therefore could not
establish a secure connection to it. To learn more about this situation and
how to fix it, please visit the web page mentioned above.
```
When you make HTTPS requests, the process invoke many more complex step, the OS or browser use a pre-installed certificate 
to validate your sever cert - one of them is Let's Encrypt root cert, and of course our pebble install have a test cert 
which no one trust (no one should trust the cert using in this example - if you're not lazy like me, generate your own). 
For now, just add option to ignore the validation `curl -k -H 'Host: example.example.com' https://$(minikube ip)` and 
you will see some HTML.

#### Testing things in your browser
For this, you need to modify your hosts files, to point the domain to `minikube ip`. 
We need to get the intermediate cert and use it to make the call success:
```
kubectl exec -n emulator deploy/pebble -- sh -c "apk add curl > /dev/null; curl -ksS https://localhost:15000/intermediates/0" > pebble.intermediate.pem.crt
curl --cacert pebble.intermediate.pem.crt -v https://example.example.com
```
And you should see some HTML. If you really want to, you can add the certificate to your chrome via:
- Settings > Privacy and security > Security > Manage certificates > Authorities tab
Choose the option for using it to verify the website, then you finally can access to https://example.example.com without 
a red warning, remember to remove that cert after you finished testing. 
#### Too much for a post already
I admit that a lot of effort just for preparing this, Pebble provide you tool to inject some chaos to the flow and 
designed for cert-manager and other smart teams to test their ACME client implementation. For our team, we need to 
spend even more time (mostly writing bash script) to make this process of setup new local cluster with only one command, 
we see this an opportunity to simulate the production environment, bring the development closer to production environment 
and it worth every single line of code.  

Almost forgot to add link to example here:
- https://github.com/manabie-com/manabie-com.github.io/tree/main/content/posts/simulate-https-certs-acme-local-k8s
