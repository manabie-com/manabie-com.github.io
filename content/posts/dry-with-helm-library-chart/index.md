+++
date = "2022-08-25T10:00:00+07:00"
author = "vync"
description = "Using library chart to avoid repeating helm code"
title = "DRY with library chart"
categories = ["Helm", "Library Chart"]
tags = ["Helm", "Library Chart"]
slug = "dry-with-helm-library-chart" 
+++
***Helm Charts help us define, install, and upgrade even the most complex Kubernetes application. However, when we have many services, we can easily repeat the helm code. In this blog, we use library chart to avoid that.***
#### **Why we need library chart?**
The more services we have, the more times helm code may be repeated. For examples, we have two services ***dog*** and ***cat*** and we need to define **ConfigMap** for them as below.

*dog/templates/configmap.yaml*
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Release.Name | printf "%s-%s" .Chart.Name }}
data:
  speak: "bow wow"
```

*cat/templates/configmap.yaml*
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Release.Name | printf "%s-%s" .Chart.Name }}
data:
  speak: "meow meow"
```
As we can see the *ConfigMap* of them have the same structure and **data** field is different only. So we need the method to create the common structure for their configs.
#### **How we create library chart?**
We create a common chart for ***dog*** and ***cat*** and name it as ***animal***.
```bash
helm create animal
```
Go to *animal/Chart.yaml* and change type from *application* to *library*.
```yaml
apiVersion: v2
name: animal
description: A Helm chart for Kubernetes
type: application => library
version: 0.1.0
appVersion: "1.16.0"
```
Now, we can carry out creating common *ConfigMap* for both ***dog*** and ***cat*** services by creating *animal/templates/_configmap.yaml*. By convention, the common files must be beginning with an underscore(_). In this file, we define **animal.configmap.tpl** as below.
```yaml
{{- define "animal.configmap.tpl" -}}
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Release.Name | printf "%s-%s" .Chart.Name }}
data: {}
{{- end -}}
```
After defining common *ConfigMap* in **animal library**, we need define a function to merge it with specific config in each service. So we continue to define **animal.util.merge** in *animal/templates/_util.yaml*.
```yaml
{{- define "animal.util.merge" -}}
{{- $top := first . -}}
{{- $overrides := fromYaml (include (index . 1) $top) | default (dict ) -}}
{{- $tpl := fromYaml (include (index . 2) $top) | default (dict ) -}}
{{- toYaml (merge $overrides $tpl) -}}
{{- end -}}
```
Then we call the merge function in *animal/templates/_configmap.yaml* by adding 
```yaml
{{- define "animal.configmap" -}}
{{- include "animal.util.merge" (append . "animal.configmap.tpl") -}}
{{- end -}}
```
Therefore, *animal/templates/_configmap.yaml* file will become
```yaml
{{- define "animal.configmap.tpl" -}}
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Release.Name | printf "%s-%s" .Chart.Name }}
data: {}
{{- end -}}
{{- define "animal.configmap" -}}
{{- include "animal.util.merge" (append . "animal.configmap.tpl") -}}
{{- end -}}
```
From now, we have already completed the definition of the library **animal**.
#### **Using the library chart to specific chart**
First of all, when we want to use other things, we need import them. So in *dog/Chart.yaml* and *cat/Chart.yaml* must be added **animal dependency** as follow.
```yaml
dependencies:
- name: animal
  version: 0.1.0
  repository: file://../animal
```
After importing it, we need to run *helm dependency update* in each service.
```bash
helm dependency update dog/
helm dependency update cat/
```
Then we can use the common *ConfigMap* in *dog/templates/configmap.yaml* and *cat/templates/configmap.yaml* as below.

*dog/templates/configmap.yaml*
```yaml
{{- include "animal.configmap" (list . "dog.configmap") -}}
{{- define "dog.configmap" -}}
data:
  speak: "bow wow"
{{- end -}}
```
*cat/templates/configmap.yaml*
```yaml
{{- include "animal.configmap" (list . "cat.configmap") -}}
{{- define "cat.configmap" -}}
data:
  speak: "meow meow"
{{- end -}}
```
Finally, we use *helm install* to test them.
#### **Summary**
A library chart is a type of Helm chart that defines chart primitives or definitions which can be shared by Helm templates in other charts. This allows users to share snippets of code that can be re-used across charts, avoiding repetition and keeping charts DRY. Thank you for your reading.
