# Yorc Plugin example

[![Build Status](https://travis-ci.org/ystia/yorc-plugin-example.svg?branch=master)](https://travis-ci.org/ystia/yorc-plugin-example) [![Go Report Card](https://goreportcard.com/badge/github.com/ystia/yorc-plugin-example)](https://goreportcard.com/report/github.com/ystia/yorc-plugin-example)

This repository provides an example of plugin extending the Ystia orchestrator, as described in [Yorc documentation on plugins](https://yorc.readthedocs.io/en/latest/plugins.html).

The plugin example implemented here provides :

* a new TOSCA definition for a Compute Instance to create on demand, in file [tosca/mycustom-types.yaml](https://github.com/ystia/yorc-plugin-example/blob/master/tosca/mycustom-types.yaml)
* an example of application TOSCA topology template using this new definition, in file [tosca/topology.yaml](https://github.com/ystia/yorc-plugin-example/blob/master/tosca/topology.yaml)
* a [delegate executor](https://github.com/ystia/yorc-plugin-example/blob/master/src/delegate.go) that will manage the provisioning of such compute instance (here it just prints logs and send events)
* an [operation executor](https://github.com/ystia/yorc-plugin-example/blob/master/src/operation.go) allowing to execute operations (here it just prints logs and send events)
* This plugin expects an infrastructure `myinfra` property `myprop` to be defined (in a real case, it could be a URL and credentials to access the service allowing to manage the infrastructure).

## Build

On a linux host, install [go 1.11](https://golang.org/dl/) or a newer version.
Then build the plugin running :

```bash
$ make
```

The plugin will be available at `bin/my-plugin`

## Test the plugin in a development environment

You can quickly setup a development environment that will allow you to test your plugin using Yorc docker images.

Download the latest Yorc docker image:

```bash
$ docker pull ystia-docker.jfrog.io/ystia/yorc:latest
```

Run Yorc mounting the directory `bin` in your host on the container directory `/var/yorc/plugins` (default path where Yorc expects to find plugins),
and mounting as well the directory `tosca` where an example of TOSCA deployment topology is provided.

Define as well the infrastructure `myinfra` property `myprop` used by the example plugin. It can be defined in Yorc configuration file, here it is passed as an environment variable `YORC_INFRA_MYINFRA_MYPROP`:

```bash
$ docker run -d --rm \
    -e 'YORC_INFRA_MYINFRA_MYPROP=myvalue' \
    -e 'YORC_LOG=1' \
    --mount "type=bind,src=$PWD/bin,dst=/var/yorc/plugins" \
    --mount "type=bind,src=$PWD/tosca,dst=/var/yorc/topology" \
    --name yorc \
    ystia-docker.jfrog.io/ystia/yorc:latest
```

Now that a Yorc server is running, it has automatically loaded the plugin available through the mount in the default plugins directory `/var/yorc/plugins`, and delegate/operation executors able to manage our new types are now registered.

We can deploy the example application topology making use of these newly defined types:

```bash
docker exec -it yorc sh -c "yorc d deploy --id my-test-app /var/yorc/topology/topology.yaml"
```

You can then check the deployment logs, running :

```bash
docker exec -it yorc sh -c "yorc d logs --from-beginning  my-test-app"
```

You should see logs of the installation workflow describing state changes as well as the event sent by our plugin delegate executor implementation:

```
[INFO][install][...][Compute][][delegate][install][]**********Provisioning node "Compute" of type "mytosca.types.Compute"
```

as well as the event sent by our plugin operation executor:
```
[INFO][my-test-app][install][...][Soft][][standard][create][]******Executing operation "standard.create" on node "Soft"
```

## Going further

You can check one of the Yorc implementation for OpenStack, AWS, Google Cloud, hosts pool, or SLURM.
For example:

* the SLURM implementation of [ExecDelegate](https://github.com/ystia/yorc/blob/v3.2.0-M4/prov/slurm/executor.go#L84) which is calling SLURM command `salloc` to allocate resources
* the SLURM implementation of [ExecAsyncOperation](https://github.com/ystia/yorc/blob/v3.2.0-M4/prov/slurm/executor.go#L64) which is calling SLURM `sbatch` to submit a job.
* the [registration](https://github.com/ystia/yorc/blob/v3.2.0-M4/prov/slurm/init.go#L28) of these delegate and operation executors for Yorc [TOSCA SLURM types](https://github.com/ystia/yorc/blob/v3.2.0-M4/data/tosca/yorc-slurm-types.yml).
