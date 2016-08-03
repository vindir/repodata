package main

import (
  "github.com/spf13/viper"
  "github.com/spf13/cobra"
)

func check(e error) {
  if e != nil {
    panic(e)
  }
}

var repodata = &cobra.Command{
    Use:   "repodata",
    Short: "Generate metadata using Artifactory repositories",
    Long:  `Generate metadata using Artifactory repositories.  Also pushes data back
up to a configured generic repo in artifactory.`,
    Run: create,
}

var server, port, user, pass, target_repo string

func main() {
  repodata.Flags().StringVarP(&server, "server", "s", "artifactory.server.missing", "Artifactory server FQDN")
  repodata.Flags().StringVarP(&port, "port", "p", "8081", "Artifactory server port to connect to")
  repodata.Flags().StringVarP(&user, "user", "u", "", "Artifactory user to deploy metadata with")
  repodata.Flags().StringVarP(&pass, "pass", "x", "", "Artifactory password to deploy metadata with")
  repodata.Flags().StringVarP(&target_repo, "target_repo", "t", "yum-repository-metadata", "Target generic repository to push yum metadata to")
  viper.BindPFlag("server", repodata.Flags().Lookup("server"))
  viper.BindPFlag("port", repodata.Flags().Lookup("port"))
  viper.BindPFlag("user", repodata.Flags().Lookup("user"))
  viper.BindPFlag("pass", repodata.Flags().Lookup("pass"))
  viper.BindPFlag("target_repo", repodata.Flags().Lookup("target_repo"))
  viper.SetEnvPrefix("repodata")
  viper.AutomaticEnv()
  viper.SetConfigName("repodata")
  viper.AddConfigPath("/etc")
  viper.AddConfigPath("$HOME/.appname")
  viper.AddConfigPath(".")
  err := viper.ReadInConfig()
  if err != nil {
    //panic(fmt.Errorf("Fatal error config file: %s \n", err))
  }
  repodata.Execute()
}


