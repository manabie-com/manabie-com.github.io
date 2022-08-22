+++
title: "Build a Grafana dashboard using Grafonnet"
date: 2022-08-22T16:19:32+07:00
author = "buivuanh"
description = "This tutorial helps you build Grafana dashboard config by using Grafonnet and golang."
categories = ["Monitoring", "Automation"]
tags = ["Grafana", "configure-Grafana-as-code", "Grafonnet"]
slug = "build-a-grafana-dashboard-using-grafonnet"
+++

## Overview

- Currently, if we want to create a Grafana dashboard, we often use UI and copy json config into our repo.
- And We have to manager folder which contains dashboard configs JSON files.

=> Maintain a set of dashboards for people with conflicting preferences and goals. Or you’re dealing with confusion
because one graph shows errors in one color and a new one uses a different color. Or there’s a new feature in Grafana
and you need to change 50 dashboards to use it. As a result, making sure your dashboards work and look good is no longer
a simple process.

=> That’s why we want to use your dashboard as code.

## [Jsonnet](https://github.com/google/go-jsonnet)

A data templating language for app and tool developers

- Generate config data
- Side-effect free
- Organize, simplify, unify
- Manage sprawling config

The example:

- You define a local variable and then that same local variable is referenced later:

```
local greeting = "hello world!";

{
 foo: "bar",
 dict: {
    nested: greeting
 },
} 
```

=>

```
{
 "foo": "bar",
 "dict": {
    "nested": "hello world!"
 }
} 
```

- Or, use functions:

```
// A function that returns an object.
local Person(name='Alice') = {
  name: name,
  welcome: 'Hello ' + name + '!',
};
{
  person1: Person(),
  person2: Person('Bob'),
}
```

=>

```
{
  "person1": {
    "name": "Alice",
    "welcome": "Hello Alice!"
  },
  "person2": {
    "name": "Bob",
    "welcome": "Hello Bob!"
  }
}
```

- Patches: You call your Jsonnet function and append to it the snippet of JSON, and it simply adds — overwrites: like
  so:

```
dashboard: {
  new(title, uid): {...}
}

dashboard.new(...) + {
 schemaVersion: 22,
}
```

=>

```
{
  "title": "super dash",
  "uid": "my-dash",
  "Panels": [],
  "schemaVersion": 22
}
```

- Imports: Not only can you create functions, but you also can put those functions into files.

```
{
  alertlist:: import 'alertlist.libsonnet',
  dashboard:: import 'dashboard.libsonnet',
  gauge:: error 'gauge is removed, migrate to gaugePanel',
  gaugePanel:: import 'gauge_panel.libsonnet',
}
```

## [Grafonnet](https://github.com/grafana/grafonnet-lib)

- Is a very simple library that provides you with the basics: creating a dashboard, creating a panel, creating a single
  stat panel, and so on.
- See more examples: https://github.com/grafana/grafonnet-lib/tree/master/examples
- E.g: we have one `dashboard.jsonnet` jsonnet file, maybe call it is main file which will evaluates it to json file

```
local grafana = import 'grafonnet/grafana.libsonnet';
local dashboard = grafana.dashboard;
local annotation = grafana.annotation;
local template = grafana.template;
local graphPanel = grafana.graphPanel;
local prometheus = grafana.prometheus;

dashboard.new(
  'title',
  editable=true,
  refresh='5s',
  time_from='now-6h',
  time_to='now',
  timepicker={},
  schemaVersion=27,
  uid='uid',
)
.addAnnotation(annotation.default)
.addTemplate(
  template.datasource(
    name='cluster',
    query='prometheus',
    current='Thanos',
    hide='',
  )
)
.addTemplate(
  template.new(
    name='namespace',
    datasource='${cluster}',
    query={
        query: 'label_values(grpc_io_server_completed_rpcs, namespace)',
        refId: 'StandardVariableQuery'
    },
    label='Namespace',
    hide='',
    refresh='load',
    definition='label_values(grpc_io_server_completed_rpcs, namespace)'
  )
)
.addPanel(
  graphPanel.new(
    title='Number of requests',
    datasource='${cluster}',
    fill=1,
    legend_show=true,
    lines=true,
    linewidth=1,
    pointradius=2,
    stack=true,
    shared_tooltip=true,
    value_type='individual',
  ).resetYaxes().
  addYaxis(
    format='none',
  ).addYaxis(
    format='short',
  ).addTarget(
    prometheus.custom_target(
        expr='expr',
        legendFormat='legendFormat',
    )
  ), gridPos={
    x: 0,
    y: 1,
    w: 18,
    h: 10,
  }
)
```

- And now, we are going to make config from above jsonnet file.

```
func MakeConfig() error {
	vm := jsonnet.MakeVM()
	jsonData, err := vm.EvaluateFile(`dashboard.jsonnet`)
	if err != nil {
		return err
	}
	fmt.Println(jsonData)
	return nil
}
```

- Result will look like:

```
{
   "__inputs": [ ],
   "__requires": [ ],
   "annotations": {
      "list": [
         {
            "builtIn": 1,
            "datasource": "-- Grafana --",
            "enable": true,
            "hide": true,
            "iconColor": "rgba(0, 211, 255, 1)",
            "name": "Annotations & Alerts",
            "type": "dashboard"
         }
      ]
   },
   "editable": true,
   "gnetId": null,
   "graphTooltip": 0,
   "hideControls": false,
   "id": null,
   "links": [ ],
   "panels": [
      {
         "aliasColors": { },
         "bars": false,
         "dashLength": 10,
         "dashes": false,
         "datasource": "${cluster}",
         "fill": 1,
         "fillGradient": 0,
         "gridPos": {
            "h": 10,
            "w": 18,
            "x": 0,
            "y": 1
         },
         "hiddenSeries": false,
         "id": 2,
         "legend": {
            "alignAsTable": false,
            "avg": false,
            "current": false,
            "max": false,
            "min": false,
            "rightSide": false,
            "show": true,
            "sideWidth": null,
            "total": false,
            "values": false
         },
         "lines": true,
         "linewidth": 1,
         "links": [ ],
         "nullPointMode": "null",
         "percentage": false,
         "pointradius": 2,
         "points": false,
         "renderer": "flot",
         "repeat": null,
         "seriesOverrides": [ ],
         "spaceLength": 10,
         "stack": true,
         "steppedLine": false,
         "targets": [
            {
               "expr": "expr",
               "legendFormat": "legendFormat",
               "refId": "A"
            }
         ],
         "thresholds": [ ],
         "timeFrom": null,
         "timeShift": null,
         "title": "Number of requests",
         "tooltip": {
            "shared": true,
            "sort": 0,
            "value_type": "individual"
         },
         "type": "graph",
         "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": [ ]
         },
         "yaxes": [
            {
               "format": "none",
               "label": null,
               "logBase": 1,
               "max": null,
               "min": null,
               "show": true
            },
            {
               "format": "short",
               "label": null,
               "logBase": 1,
               "max": null,
               "min": null,
               "show": true
            }
         ],
         "yaxis": {
            "align": false,
            "alignLevel": null
         }
      }
   ],
   "refresh": "5s",
   "rows": [ ],
   "schemaVersion": 27,
   "style": "dark",
   "tags": [ ],
   "templating": {
      "list": [
         {
            "current": {
               "text": "Thanos",
               "value": "Thanos"
            },
            "hide": 0,
            "label": null,
            "name": "cluster",
            "options": [ ],
            "query": "prometheus",
            "refresh": 1,
            "regex": "",
            "type": "datasource"
         },
         {
            "allValue": null,
            "current": { },
            "datasource": "${cluster}",
            "definition": "label_values(grpc_io_server_completed_rpcs, namespace)",
            "hide": 0,
            "includeAll": false,
            "label": "Namespace",
            "multi": false,
            "name": "namespace",
            "options": [ ],
            "query": {
               "query": "label_values(grpc_io_server_completed_rpcs, namespace)",
               "refId": "StandardVariableQuery"
            },
            "refresh": 1,
            "regex": "",
            "sort": 0,
            "tagValuesQuery": "",
            "tags": [ ],
            "tagsQuery": "",
            "type": "query",
            "useTags": false
         }
      ]
   },
   "time": {
      "from": "now-6h",
      "to": "now"
   },
   "timepicker": { },
   "timezone": "browser",
   "title": "title",
   "uid": "uid",
   "version": 0
}
```

## More
When you are able to config dashboard as code, in particular: using golang to make config base Grafonnet lib, you can totally add more feature such as :
- Create and manage dashboard config files with jsonnet format instead of raw json files.
- Can reuse one or more pannels for multiple dashboards.
- Build jsonnet template file which you can fill automatically params such metric names, grpc methods ....