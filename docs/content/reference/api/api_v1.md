---
title: "Autopilot Config Files"
weight: 5
---

<a name="top"></a>


## Table of Contents

- [autopilot.proto](#autopilot.proto)
    - [AutopilotProject](#autopilot.AutopilotProject)
    - [MetricsQuery](#autopilot.MetricsQuery)
    - [Parameter](#autopilot.Parameter)
    - [Phase](#autopilot.Phase)
  
  
  
  

- [autopilot-operator.proto](#autopilot-operator.proto)
    - [AutopilotOperator](#autopilot.AutopilotOperator)
  
    - [MeshProvider](#autopilot.MeshProvider)
  
  
  




<a name="autopilot.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## autopilot.proto
The following Schema defines the structure of the `autopilot.yaml` configuration file.

This file is used to generate and re-generate the project structure, as well
as execute tasks related to build and deployment. It can be consumed
both via the `ap` CLI as well as in `codegen` packages.


<a name="autopilot.AutopilotProject"></a>

### AutopilotProject
The AutopilotProject file is the root configuration file for the project itself.

This file will be used to build and deploy the autopilot operator.
It is loaded automatically by the autopilot CLI. Its
default location is 'autopilot.yaml'


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| kind | [string](#string) |  | the name (kubernetes Kind) of the top-level CRD for the operator Specified via the `ap init <Kind>` command |
| apiVersion | [string](#string) |  | the ApiVersion of the top-level CRD for the operator |
| operatorName | [string](#string) |  | the name of the Operator this is used to name and label loggers, k8s resources, and metrics exposed by the operator. Should be [valid Kube resource names](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names). |
| phases | [][Phase](#autopilot.Phase) | repeated | Each phase represents a different stage in the lifecycle of the CRD (e.g. Pending/Succeeded/Failed). <br> Each phase specifies a unique name and its own set of inputs and outputs. |
| enableFinalizer | [bool](#bool) |  | enable use of a Finalizer to handle object deletion |
| customParameters | [][Parameter](#autopilot.Parameter) | repeated | custom Parameters which extend Autopilot's builtin types |
| queries | [][MetricsQuery](#autopilot.MetricsQuery) | repeated | custom Queries which extend Autopilot's metrics queries |






<a name="autopilot.MetricsQuery"></a>

### MetricsQuery
MetricsQueries extend the query options available to workers.
MetricsQueries are accessible to workers via generated client code
that lives in <project root>/pkg/metrics


The following MetricsQuery:

```
name: success-rate
parameters:
- Name
- Namespace
- Interval
queryTemplate: |
    sum(
        rate(
            envoy_cluster_upstream_rq{
                kubernetes_namespace="{{ .Namespace }}",
                kubernetes_pod_name=~"{{ .Name }}-[0-9a-zA-Z]+(-[0-9a-zA-Z]+)",
                envoy_response_code!~"5.*"
            }[{{ .Interval }}]
        )
    )
    /
    sum(
        rate(
            envoy_cluster_upstream_rq{
                kubernetes_namespace="{{ .Namespace }}",
                kubernetes_pod_name=~"{{ .Name }}-[0-9a-zA-Z]+(-[0-9a-zA-Z]+)"
            }[{{ .Interval }}]
        )
    )
    * 100
```

would produce the following `metrics` Interface:

```go
type CanaryDeploymentMetrics interface {
    metrics.Client
    GetIstioSuccessRate(ctx context.Context, Namespace, Name, Interval string) (*metrics.QueryResult, error)
    GetIstioRequestDuration(ctx context.Context, Namespace, Name, Interval string) (*metrics.QueryResult, error)
    GetEnvoySuccessRate(ctx context.Context, Namespace, Name, Interval string) (*metrics.QueryResult, error)
    GetEnvoyRequestDuration(ctx context.Context, Namespace, Name, Interval string) (*metrics.QueryResult, error)
}
```


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| queryTemplate | [string](#string) |  |  |
| parameters | [][string](#string) | repeated |  |






<a name="autopilot.Parameter"></a>

### Parameter
Custom Parameters allow code to be generated
for inputs/outputs that are not built-in to Autopilot.
These types must be Kubernetes-compatible Go structs.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| lowerName | [string](#string) |  | the fully lower-case name of this resource e.g. "pods", "services", "replicasets", "configmaps" |
| singleName | [string](#string) |  | the singular CamelCased name of the resource equivalent to Kind |
| pluralName | [string](#string) |  | the plural CamelCased name of the resource equivalent to the pluralized form of Kind |
| importPrefix | [string](#string) |  | import prefix used by generated code |
| package | [string](#string) |  | go package (import path) to the go struct for the resource |
| apiGroup | [string](#string) |  | Kubernetes API group for the resource e.g. "networking.istio.io" |
| isCrd | [bool](#bool) |  | indicates whether the resource is a CRD if true, the Resource will be added to the operator's runtime.Scheme |






<a name="autopilot.Phase"></a>

### Phase
MeshProviders provide an interface to monitoring and managing a specific
mesh.

Autopilot does not abstract the mesh API - Autopilot developers must
still reason able about Provider-specific CRDs. Autopilot's job is to
abstract operational concerns such as discovering control plane configuration
and monitoring metrics.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | name of the phase. must be unique |
| description | [string](#string) |  | description of the phase. used for comments and docs |
| initial | [bool](#bool) |  | indicates whether this is the initial phase of the system. exactly one phase must be the initial phase |
| final | [bool](#bool) |  | indicates whether this is a "final" or "resting" phase of the system. when the CRD is in the final phase, no more processing will be done on it |
| inputs | [][string](#string) | repeated | the set of inputs for this phase the inputs will be retrieved by the scheduler and passed to the worker as input parameters

custom inputs can be defined in the autopilot.yaml |
| outputs | [][string](#string) | repeated | the set of outputs for this phase the inputs will be propagated to k8s storage (etcd) by the scheduler.

custom outputs can be defined in the autopilot.yaml |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="autopilot-operator.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## autopilot-operator.proto
autopilot-operator.proto defines the API Schema for the autopilot-operator.yaml configuration file.
this file provides the bootstrap configuration that is loaded to the
operator at boot-time/runtime


<a name="autopilot.AutopilotOperator"></a>

### AutopilotOperator
The AutopilotOperator file is the bootstrap
Configuration file for the Operator.
It is stored and mounted to the operator as a Kubernetes ConfigMap.
The Operator will hot-reload when the configuration file changes.
Default name is 'autopilot-operator.yaml' and should be stored in the project root.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| version | [string](#string) |  | version of the operator used for logging and metrics default is "0.0.1" |
| meshProvider | [MeshProvider](#autopilot.MeshProvider) |  | meshProvider determines how the operator will connect to a service mesh Default is "SMI" |
| controlPlaneNs | [string](#string) |  | controlPlaneNs is the namespace the control plane lives in Default is "istio-system" |
| workInterval | [google.protobuf.Duration](#google.protobuf.Duration) |  | workInterval to sets the interval at which CRD workers resync. Default is 5s |
| metricsAddr | [string](#string) |  | Serve metrics on this address. Set to empty string to disable metrics defaults to ":9091" |
| enableLeaderElection | [bool](#bool) |  | Enable leader election. This will prevent more than one operator from running at a time defaults to true |
| watchNamespace | [string](#string) |  | if non-empty, watchNamespace will restrict the Operator to watching resources in a single namespace if empty (default), the Operator must have Cluster-scope RBAC permissions (ClusterRole/Binding) can also be set via the WATCH_NAMESPACE environment variable |
| leaderElectionNamespace | [string](#string) |  | The namespace to use for Leader Election (requires read/write ConfigMap permissions) defaults to the watchNamespace |
| logLevel | [google.protobuf.UInt32Value](#google.protobuf.UInt32Value) |  | Log level for the operator's logger values: 0 - Debug 1 - Info 2 - Warn 3 - Error 4 - DPanic 5 - Panic 6 - Fatal Defaults to Info |





 <!-- end messages -->


<a name="autopilot.MeshProvider"></a>

### MeshProvider
MeshProviders provide an interface to monitoring and managing a specific
mesh.
Autopilot does not abstract the mesh API - Autopilot developers must
still reason able about Provider-specific CRDs. Autopilot's job is to
abstract operational concerns such as discovering control plane configuration
and monitoring metrics.

| Name | Number | Description |
| ---- | ------ | ----------- |
| Istio | 0 | the Operator will utilize Istio mesh for metrics and configuration |
| Custom | 1 | the Operator will utilize a locally deployed Prometheus instance for metrics (Currently unimplemented) |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->


