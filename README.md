# SCMT

[![Build Status](https://travis-ci.org/eeayiaia/scmt.svg?branch=master)](https://travis-ci.org/eeayiaia/scmt)

This is the git repository for the bachelor thesis "Design and development of single-board supercomputers" given at Chalmers University of Technology 2016

## Building scmt
Welcome!

Requirements:
* Go [here](https://golang.org/doc/install) for instruction to setup your golang environment.

### Cloning the repository
There is two ways to setup the repository using golang. Either by [cloning the repo directly into your golang environment](#Cloning\ into\ the\ environment), or by [linking from the outside to your golang environment](#Linking\ into\ the\ environment).

#### Cloning into the environment
```bash
$> cd $GOPATH
$> git clone git@github.com:eeayiaia/scmt.git src/github.com/eeyiaia/scmt
```

go to [building](#Building) in order to build the project.

#### Linking into the environment
```bash
$> git clone git@github.com:eeayiaia/scmt.git
$> pwd
$HOME/git
```

The repository will be placed in `$HOME/git/scmt`.

```bash
$> ln -sf $HOME/git/scmt $GOPATH/src/github.com/eeyiaia/scmt
```

will create a hard link from `$HOME/git/scmt` to `$GOPATH/src/github.com/eeyiaia/scmt`.

### Building
You can build the code from any path, because the golang compiler will always try to look for your code in your `$GOPATH`. 

```bash
$> cd $HOME/github.com/eeyiaia/scmt   # or $HOME/git/scmt if you prefer
$> go build
```

will produce a binary in your current working directory!

```bash
go install github.com/eeayiaia/scmt
```

will produce a binary in your `$GOPATH/bin` directory.

### Setting up the repository
Run `setup-repo.sh` in order to setup the proper git hooks.

```bash
$> ./setup-repo.sh
```

## Using scmt
scmt is fairly easy to use to create your own cluster, but certain requirements are needed before-hand.

Requirements:
* master-node with two network interfaces (external & internal)
* a switch to connect the devices together

### Connecting the master-node
The master-node should be connected both to an external interface (to access the outside-world) and an internal interface. The internal interface should be connected to the switch.

The first step is to find your master-node on the external interface (eg. your router). You can run a simple network-scan using nmap to do this.

```bash
$> nmap -F 192.168.0.1-255
```

Replace `192.168.0.1-255` with your routers subnet. Then remotely access your master-node in order to install and build scmt.

Make sure you know the name of the external and internal interfaces, as this is needed later. 

The next step is to build scmt - see above; once scmt is built proceed to the next section.

### Setting up the master-node
All you need to do in order to configure the master-node properly is to run scmt. In order to run scmt requires the environment variable `$SCMT_ROOT` to be set to the root of the repository. You can choose whether to use our wrapper-script `run-scmt.sh` or by setting the environment variable by yourself.

Run `scmt` to generate a configuration file. If you wish to change any configuration parameters later you can edit the file `resources/scmt.json`.

### Running scmt
You can either run scmt in your shell, or start scmt as a daemon process. 

In order to run scmt and follow the log-output (to see if it works), run the following:
```bash
$> scmt -d
$> tail -f $SCMT_ROOT/resources/scmt.log
```

If scmt successfully starts, your cluster is ready to be used! You can now connect unconfigured single board devices to the internal interface (via the switch) and use them once configured.

## Customising scmt
Customising scmt can be done in two ways, either by adding features or changing the core functionality. We do not - however - recommend changing the core functionality of scmt as this can break your cluster. If you wish to extend the functionality we recommend creating a new plugin.

### Event-Action system
scmt uses a very simple event-action system, which is designed to be modular by default. An event trigger certain types of actions - which is basically a set of shell scripts that modify the cluster.

Three types of events are currently supported:
* Connect - when a new device connects
* Reconnect - when an already existing device reconnects (after disconnection)
* Disconnect - when an already existing device disconnects

### Action structure
Actions are really a set of scripts, stored in a unix-fashion in `resources/scripts.d`. 

Currently there are two types of actions:
* Device - actions being run on a node
* Master - actions being run on the master-node

As you can see in the folder structure below, it is pretty straight-forward:

```bash
$> tree resources/scripts.d
resources/scripts.d
├── device.init.d
│   ├── 00-set-hostname.sh
│   ├── 10-set-hosts.sh
│   ├── 15-setup-nfs.client.sh
│   ├── 20-resizefs.sh
│   ├── 30-setup-approx.sh
│   ├── 35-apt-get-update.sh
│   └── 50-uninstall-packages.sh
├── master.init.d
│   ├── 00-apt-get-update.sh
│   ├── 05-install-dependencies.sh
│   ├── 10-create-database.sh
│   ├── 15-setup-nfs.sh
│   ├── 20-setup-approx.sh
│   ├── 40-setup-dhcpd.sh
│   ├── 45-symlink-usr.bin.sh
│   ├── 50-config-apparmor.sh
│   ├── 60-setup-network.sh
│   ├── 65-iptables.sh
│   ├── 70-setup-init.sh
│   └── resources
│       ├── baseDHCPD.conf
│       ├── create_database.sql
│       └── scmt
├── master.newdevice.d
│   ├── 00-add-static-ip.sh
│   └── 10-add-to-hosts.sh
├── master.removedevice.d
│   ├── 00-remove-static-ip.sh
│   └── 10-remove-from-hosts.sh
└── utils.sh

5 directories, 26 files
```

Upon connection and reconnection, the `master.newdevice.d` and `device.init.d` actions is being run.
Upon disconnection, the `master.removedevice.d` actions is being run.

*Note: If you are to edit any of the actions - or add new actions - remember that some actions depend on each other.*

### Plugins
Plugins follow the very same event-action system as described above, but in an isolated folder structure.

```bash
$> tree resources/plugins.d
resources/plugins.d
├── ganglia
│   ├── device.init.d
│   │   ├── 00-install.sh
│   │   └── helpscript
│   │       └── regex.py
│   ├── master.init.d
│   │   ├── 00-install.sh
│   │   └── helpscript
│   │       └── regex.py
│   ├── master.newdevice.d
│   │   └── 00-ganglia-add-node.sh
│   └── master.removedevice.d
│       ├── 00-ganglia-remove-node.sh
│       └── helpscript
│           └── regex.py
```
