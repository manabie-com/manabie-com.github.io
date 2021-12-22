+++
date = "2021-12-15T14:28:23+07:00"
author = "sanglh"
description = "How to work with helm chart and something about helm subchart/dependency feature"
title = "Deploying application on Kubenetes with Helm"
categories = ["DevSecOps", "Helm", "K8S"]
tags = ["k8s", "helm", "kubernetes"]
+++

## Why do we need to use helm chart?

Working with k8s, we have many k8s resources: ConfigMap, Secret, Service, Deployment, Pod,... K8s resources can be created by creating the yaml file, and use command kubectl apply -f <link_to_the_file>. For example, we create a file with name is deployment.yaml:

```yaml

  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: nginx-deployment
  spec:
    replicas: 2
    selector:
      matchLabels:
        app: nginx-deployment
    template:
      metadata:
        labels:
          app: nginx-deployment
      spec:
        containers:
        - name: nginx-deployment
          image: nginx
          ports:
          - containerPort: 8080

```
To create the k8s deployment, we use the command:
```bash

  kubectl apply -f deployment.yaml

```

With each resource, we can create a file (deployment.yaml, service.yaml, configmap.yaml, secret.yaml,...) and use the kubectl command to apply them. But we have a problem. If we work with microservices and we have multiple services => the number of yaml files is very large, and makes it difficult for us to manage them.

We need a solution to manage all of them in a common template => Helm can solve this problem

Before going to next part, you need to download and install minikube on your local device.

## Introduction to helm chart
**Helm is a tool that streamlines installing and managing Kubernetes applications.** To work with helm charts, we need to install the helm command. This link is about how to install: <a>https://helm.sh/docs/intro/install</a>

I will introduce you to the structure of the helm chart. To create a chart, we can use my example command:
```bash

  helm create mynginx

```
That command will create a mynginx folder. Let's take a look in the mynginx folder where we have some folders and files:
- charts: this folder will contain another subchart which our chart will depend on. I will talk more about that subchart in the next section
- templates: this folder contains yaml files about the k8s resource. We can add more, remove or modify the files in this folder
- Chart.yaml: this file is the description about the chart information
- _helpers.tpl: this file defines some yaml key-value which we can import to yaml files in templates folders

Let take a look about the file template/deployments.yaml:

```yaml

  apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: {{ include "mynginx.fullname" . }}
    labels:
      {{- include "mynginx.labels" . | nindent 4 }}
  spec:
    {{- if not .Values.autoscaling.enabled }}
    replicas: {{ .Values.replicaCount }}
    {{- end }}
    selector:
      matchLabels:
        {{- include "mynginx.selectorLabels" . | nindent 6 }}
    template:
      metadata:
        {{- with .Values.podAnnotations }}
        annotations:
          {{- toYaml . | nindent 8 }}
        {{- end }}
        labels:
          {{- include "mynginx.selectorLabels" . | nindent 8 }}
      spec:
        {{- with .Values.imagePullSecrets }}
        imagePullSecrets:
          {{- toYaml . | nindent 8 }}
        {{- end }}
        serviceAccountName: {{ include "mynginx.serviceAccountName" . }}
        securityContext:
          {{- toYaml .Values.podSecurityContext | nindent 8 }}
        containers:
          - name: {{ .Chart.Name }}
            securityContext:
              {{- toYaml .Values.securityContext | nindent 12 }}
            image: "{{ .Values.image.repository }}:{{ .Values.image.tag | default .Chart.AppVersion }}"
            imagePullPolicy: {{ .Values.image.pullPolicy }}
  ..........................................................

```

This yaml file looks like the template of k8s resource, except we see alot of this {{ }} syntax. The statement in {{ }} can be a value, operator, expression, etc. Let me explain some popular uses:

- {{ include … }}: this syntax helps us to import the yaml template that is defined in file **_helpers.tpl**
- {{ .Values.&#60;key&#62; }}: gets the **value** at the position **key** of files values.yaml and fill that value to the k8s template yaml file
- {{ -if &#60;expression &#62; }} &#60;some yaml line&#62; {{- end }}: if the expression returns true, the &#60;some yaml line&#62; will be added to the template file
- {{ with .Values.key }} {{- toYaml . }}   {{- end }}, {{- toYaml .Values.key }}: get the **value** at the position **key** of files **values.yaml** and convert it to yaml format.
- indent &#60;n&#62;: indents the content with n tabs
- nindent &#60;n&#62;: same with indent but the yaml code is created in the new line 

Look at the **_helper.tpl**, there are some basic things we need to know:
```yaml

  {{- define "mynginx.name" -}}
  {{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
  {{- end }}

```
This script is about the value of **mynginx.name**. The first thing is the **default** operator. The value of **.Chart.Name** is defined in file **Chart.yaml**. If the value of key nameOverride in file values.yaml is not empty, it is chosen, if not, .Chart.Name is chosen. Trunc 63 is the function to get the first 63 characters and trimSuffix “-” in here is used to remove all the characters “-” at the end of the string. You can read more about helm chart function at here: https://helm.sh/docs/chart_template_guide/function_list/

``` yaml

  {{- define "mynginx.selectorLabels" -}}
  app.kubernetes.io/name: {{ include "mynginx.name" . }}
  app.kubernetes.io/instance: {{ .Release.Name }}
  {{- end }}

```
This script defines the template of “mynginx.selectorLabels”. This returns some yaml code. **{{ .Release.Name }}** is the value which we pass in the helm install command. I will introduce the helm command install now.

## Install/Upgrade helm chart
To install all the services that are defined in the helm chart into k8s cluster, we use helm install/upgrade command. The format I usually use to install/upgrade is:

```bash

    helm upgrade --install <ReleaseName> <path to the chart> \
    --create-namespace --namespace= “specific namespace to install service” \
    --values “path to the values file” \
    --set=<key1>=<value1> \
    --set=<key2>=<value2> \
    ......
    --set=<keyN>=<valueN>
    
```

- The &#60;ReleaseName&#62; will become the value of {{ .Release.Name }} in file **_heplers.tpl** I mentioned above
- “path to the values file": The default values for the helm chart are defined in file values.yaml. But if we want to use value from another file, we can specify the path to that file by using --values flag
- --set=&#60;key&#62;=&#60;value&#62;: in the case the value is dynamic, we can set value for the key by using --set flag


With current chart, you can apply and create nginx service by running (make sure to install and start minikube in your local device)

```bash

  helm upgrade --install mynginx ./mynginx

```

Now you can run: _kubectl get deployment_

``` bash

    NAME                    READY   UP-TO-DATE   AVAILABLE   AGE
    mynginx                 1/1     1            1           3m29s
    
```

And _kubectl get service_

``` bash

    NAME                    TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
    kubernetes              ClusterIP   10.96.0.1       <none>        443/TCP   9m36s
    mynginx                 ClusterIP   10.107.22.190   <none>        80/TCP    4m15s

```

You can see that all things are created for you. 

More about helm install command, you can refer at: https://helm.sh/docs/helm/helm_install/

## Helm subchart and global values
We have two use-case that need to use this feature of helm:
- To import some template that common and dynamic by value of each helm chart
- To install all of our microservices in one helm command.

### Import some template that common and dynamic by value of each helm chart 

Imagine that you have many nginx apps and you need to monitor them. The solution is you can use a sidecar container for each app. 
For monitoring nginx, each services you need add some config like this:
- A config map for file nginx.conf
- Some yaml code for mapping the config map into volume of container
- Some yaml code for create sidecar container with image nginx/nginx-prometheus-exporter:0.9.0

The code above is repeated for each service. So we can create a helm chart with the type “library" and import it to each service.

You can follow these steps:

- Create an library helm, I named it is “metrics"

```bash

  helm create metrics

```
- In file Chart.yaml, change type from **application** to **library**
- In the metrics/template folder, you delete all files and create a file _metrics.yaml with the content:
```yaml

  {{- define "metrics.volume.tpl" -}}
  name: {{ .Chart.Name }}-config
  configMap:
    name: {{ .Release.Name }}
    items:
    - key: nginx.conf
      path: nginx.conf
  {{- end -}}
  
  {{- define "metrics.volume-mount.tpl" -}}
  name: {{ .Chart.Name }}-config
  mountPath: /etc/nginx/conf.d/nginx.conf
  subPath: nginx.conf
  {{- end -}}
  
  {{- define "metrics.container.tpl" -}}
  name: {{ .Chart.Name }}-exporter
  image: nginx/nginx-prometheus-exporter:0.9.0
  imagePullPolicy: {{ .Values.image.pullPolicy }}
  ports:
    - name: exporter
      containerPort: 9113
      protocol: TCP
  args:
    - -nginx.scrape-uri=http://localhost:8080/stub_status
  resources:
    {{- toYaml .Values.resources | nindent 12 }}
  {{- end -}}
  
  {{- define "metrics.configmap.tpl" -}}
  apiVersion: v1
  kind: ConfigMap
  metadata:
    name: {{ .Chart.Name }}
  data:
    nginx.conf: |
      server {
          listen       8080;
          server_name  localhost;
          location /stub_status {
              stub_status;
          }
      }
  {{- end -}}

```
- Go to folder mynginx,  and add this block to file Chart.yaml:

```yaml

    dependencies:
    - name: metrics
      version: 0.1.0
      repository: file://../metrics

```
- Run the command: <i>helm dependency update</i><br>
After you run the command, it will create the .tgz file in charts folder
- Now we need to add config for file nginx.conf by using configmap. Create configmap.yaml file in templates folder and just import the template from subchart:

```yaml

  {{ include "metrics.configmap.tpl" . }}

```
- In the templates/deployment.yaml, import volume template, volume mount template and sidecar container template like this:

```yaml

  apiVersion: apps/v1
  kind: Deployment
  spec:
    …
    template:
       spec:
         volumes:
         - {{ include "metrics.volume.tpl" . | nindent 8 }}
         containers:
         - name: {{ .Chart.Name }}
           ....
           volumeMounts:
            - {{- include "metrics.volume-mount.tpl" . | nindent 12 }}
         - {{- include "metrics.container.tpl" . | nindent 10 }}

```
- Expose port 9113 in the service.yaml:

```yaml

    - port: 9113
      targetPort: exporter
      protocol: TCP
      name: exporter

```
- Now we can upgrade helm chart:

```bash

    helm upgrade --install mynginx ./mynginx
    
```
- Check the service by run __kubectl get services__

```bash

    NAME                    TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)           AGE
    kubernetes              ClusterIP   10.96.0.1       <none>        443/TCP           6m31s
    mynginx                 ClusterIP   10.106.193.81   <none>        80/TCP,9113/TCP   2m59s

```

Now you have the prometheus exporter working at port 9113. You can run: __kubectl port-forward svc/mynginx 9113:9113__ and call to localhost:9113/metrics to get metrics data
```bash

  sanglh@sanglh-G5-5500:~$ curl localhost:9113/metrics
  # HELP nginx_connections_accepted Accepted client connections
  # TYPE nginx_connections_accepted counter
  nginx_connections_accepted 6123
  # HELP nginx_connections_active Active client connections
  # TYPE nginx_connections_active gauge
  nginx_connections_active 1
  # HELP nginx_connections_handled Handled client connections
  # TYPE nginx_connections_handled counter
  nginx_connections_handled 6123
  # HELP nginx_connections_reading Connections where NGINX is reading the request header
  # TYPE nginx_connections_reading gauge
  nginx_connections_reading 0
  # HELP nginx_connections_waiting Idle client connections
  # TYPE nginx_connections_waiting gauge
  nginx_connections_waiting 0
  # HELP nginx_connections_writing Connections where NGINX is writing the response back to the client
  # TYPE nginx_connections_writing gauge
  nginx_connections_writing 1
  # HELP nginx_http_requests_total Total http requests
  # TYPE nginx_http_requests_total counter
  nginx_http_requests_total 6124
  # HELP nginx_up Status of the last metric scrape
  # TYPE nginx_up gauge
  nginx_up 1
  # HELP nginxexporter_build_info Exporter build information
  # TYPE nginxexporter_build_info gauge
  nginxexporter_build_info{commit="5f88afbd906baae02edfbab4f5715e06d88538a0",date="2021-03-22T20:16:09Z",version="0.9.0"} 1

```

### Install all of our microservices in one helm command

Currently each service has 1 helm chart. So that if you have N service, you need to run N command. Imaging that if you want to deploy all in one command, you can do with 2 ways:

- create many of deployment, service, configmap,... template in one helm chart
- use subchart


In this blog I will use subchart. I created another chart with the name “another-nginx". You can read at … And I created the chart name “all-in-one".

You can add dependencies into all-in-one/Chart.yaml like this and run helm dependency update:

```yaml

    dependencies:
    - name: another-nginx
      version: 0.1.0
      repository: file://../another-nginx
    - name: mynginx
      version: 0.1.0
      repository: file://../mynginx

```

But this way is only effective when dependencies are stable. When you update something in another-nginx or mynginx, you need to run helm dependency update to update the chart all-in-one. So the efficient way we can do is copy mynginx and another nginx folder to all-in-one/charts. When you run helm upgrade --install all-in-one ./all-in-one, it will return the same result with using helm dependency.

Check the deployment result by run __kubectl get deploy__:

```bash

    NAME                       READY   UP-TO-DATE   AVAILABLE   AGE
    all-in-one-another-nginx   1/1     1            1           3m51s
    all-in-one-mynginx         1/1     1            1           3m51s
    
```

You can see the deployment name contains “all-in-one”. You can remove the “all-in-one” prefix by change the nameOverride and fullnameOverride in file values.yaml of all-in-one/chart/another-nginx and all-in-one/chart/mynginx to these name:

- File all-in-one/charts/mynginx/values.yaml

```yaml

    nameOverride: "mynginx"
    fullnameOverride: "mynginx"

```
- File all-in-one/charts/another-nginx/values.yaml

```yaml

    nameOverride: "another-nginx"
    fullnameOverride: "another-nginx"
    
```

Save file and run __helm upgrade --install all-in-one ./all-in-one__ again. After that recheck deploy and you can see the name is changed:

```bash

    NAME            READY   UP-TO-DATE   AVAILABLE   AGE
    another-nginx   1/1     1            1           13s
    mynginx         1/1     1            1           15s

```

The final thing you need to know about helm subchart is Global Chart Values. This helps us so much to control the value of each subchart or common value of the all-in-one chart. Read more at: https://helm.sh/docs/chart_template_guide/subcharts_and_globals/ 


The source code for this blog is posted at:
https://github.com/manabie-com/manabie-com.github.io/tree/main/content/posts/introduce-about-helm-chart