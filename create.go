package main

import (
  "fmt"
  "os"
  "bytes"
  "strings"
  "net/http"
  "github.com/spf13/viper"
  "github.com/spf13/cobra"
)

func create(cmd *cobra.Command, args []string) {
  server := viper.GetString("server")
  port := viper.GetString("port")
  user := viper.GetString("user")
  pass := viper.GetString("pass")
  target_repo := viper.GetString("target_repo")
  //fmt.Printf("[DEBUG] Config Determined:\nS:%v\nP:%v\nU:%v\nX:%v\n", server, port, user, pass)
  resp := GetArtifactoryRepos()
  var repos ArtifactoryResponse
  repos.Populate(resp)

  baseVers, archVers := repos.UniqueProperties()

  repo_entry := `[%v]
name=%v
baseurl=http://%v:%v/artifactory/%v
gpgcheck=0
`

  for baseVer,_ := range baseVers {
    for archVer,_ := range archVers {
      repo_file := fmt.Sprintf("%v-%v.repo", baseVer, archVer)
      var metadata_set bytes.Buffer
      f, err := os.Create(repo_file)
      defer f.Close()
      check(err)

      for _, repo := range repos.Repo {
        supported := IsSupported(repo, baseVer, archVer)
        if supported {
          endpoints_raw := repo.Properties.EndPoints
          if len(endpoints_raw) <= 0 {
            repo_name := strings.Replace(repo.Name, "-cache", "", 1)
            metadata_set.WriteString(fmt.Sprintf(repo_entry, repo_name, repo_name, server, port, repo_name))
            f.WriteString(fmt.Sprintf(repo_entry, repo_name, repo_name, server, port, repo_name))
          } else {
            endpoints := strings.Split(endpoints_raw[0], ",")
            //fmt.Printf("[DEBUG] EndPoints: (%v)\n", endpoints)
            for i, endpoint := range endpoints {
              //fmt.Printf("[DEBUG] EndPoint: (%v)\n", endpoint)
              repo_name := fmt.Sprintf("%v-%v", repo.Name, i)
              repo_name = strings.Replace(repo_name, "-cache", "", 1)
              repo_path := fmt.Sprintf("%v/%v", repo_name, endpoint)
              metadata_set.WriteString(fmt.Sprintf(repo_entry, repo_name, repo_name, server, port, repo_path))
              f.WriteString(fmt.Sprintf(repo_entry, repo_name, repo_name, server, port, repo_path))
            }
          }
          //metadata_set.WriteString(fmt.Sprintf("Blub %v\n", repo_file))
          //fmt.Printf("[DEBUG] %v: %v-%v.repo Supported!\n", repo.Name, baseVer, archVer)
        }
      }

      target_repo_uri := fmt.Sprintf("http://%v:%v@%v:%v/artifactory/%v/%v", user,
                                                                             pass,
                                                                             server,
                                                                             port,
                                                                             target_repo,
                                                                             repo_file)
      _, err = http.NewRequest("DELETE", target_repo_uri, nil)
      check(err)
      req, err := http.NewRequest("PUT", target_repo_uri, &metadata_set)
      check(err)
      req.Header.Set("Content-Type", "text/plain")
      client := &http.Client{}
      res, err := client.Do(req)
      check(err)
      defer res.Body.Close()
    }
  }
  fmt.Printf("Complete! Yum repository metadata has been pushed to the remote repository.\n")
}

