## yumfactory
Generate yum metadata for artifactory hosted yum repos

The metadata generated is pushed back up to an artifactory repository
where it can be remotely pulled in by yum clients using include directives
in their yum.conf

### Building

The yumfactory has a few external depencencies.  If you aren't familiar with 
setting up your go environment please read through [this go introduction](https://howistart.org/posts/go/1).

```
# Once you've got a properly configured go environment these can be pulled down with go get
go get github.com/spf13/viper
go get github.com/spf13/cobra
go get gopkg.in/resty.v0
```

Once dependencies are in place.  Build the project as normal
```bash
go get github.com/vindir/yumfactory
cd $GO_HOME/src/github.com/vindir/yumfactory
go build
```

### Configuring

yumfactory uses the [viper package](https://github.com/spf13/viper) and supports
configuration files in yaml or json format with a properly set extension. The
filename should be `repodata` as in `repodata.yaml`

Example:
```yaml
$ cat /etc/repodata.yaml
---
server: neovpartifactory1.neo.vocalocity.com
port: 8081
user: repodata
pass: automation
target_repo: yum-repository-metadata
```

### Using

```bash
$ ./repodata --help
Generate metadata using Artifactory repositories.  Also pushes data back
up to a configured generic repo in artifactory.

Usage:
  repodata [flags]

Flags:
  -x, --pass string          Artifactory password to deploy metadata with
  -p, --port string          Artifactory server port to connect to (default "8081")
  -s, --server string        Artifactory server FQDN (default "artifactory.server.missing")
  -t, --target_repo string   Target generic repository to push yum metadata to (default "yum-repository-metadata")
  -u, --user string          Artifactory user to deploy metadata with
  ```
  
There are no required arguments. If the configuration file sets a working set of options, then repodata will run run fine without any arguments.

```bash
$ ./repodata
Complete! Yum repository metadata has been pushed to the remote repository.
```

