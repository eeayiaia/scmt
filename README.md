# SCMT

This is the git repository for the bachelor thesis "Design and development of single-board supercomputers" given at Chalmers University of Technology 2016

## Setting up golang env
Welcome!

Requirements:
* Go [here](https://golang.org/doc/install) for instruction to setup your golang environment.

### Cloning the repository
There is two ways to setup the repository using golang. Either by [cloning the repo directly into your golang environment](#Cloning\ into\ the\ environment), or by [linking from the outside to your golang environment](#Linking\ into\ the\ environment).

#### Cloning into the environment
```bash
$> cd $GOPATH
$> git clone git@github.com:eeayiaia/scmt.git src/scmt
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
$> ln -sf $HOME/git/scmt $GOPATH/src/scmt
```

will create a hard link from `$HOME/git/scmt` to `$GOPATH/src/scmt`.

### Building
You can build the code from any path, because the golang compiler will always try to look for your code in your `$GOPATH`. 

```bash
go build scmt
```

will produce a binary in your current working directory!

```bash
go install scmt
```

will produce a binary in your `$GOPATH/bin` directory.
