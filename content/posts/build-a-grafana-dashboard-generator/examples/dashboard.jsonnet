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