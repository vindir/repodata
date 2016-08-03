package main

import (
  "github.com/spf13/viper"
  "gopkg.in/resty.v0"
  "encoding/json"
  "fmt"
)

type RepoProperties struct {
  ArchSupport []string `json:"yumrepo.arch_support"`
  BaseVerSupport []string `json:"yumrepo.basever_support"`
  AutoSync []string `json:"yumrepo.auto_sync"`
  EndPoints []string `json:"yumrepo.endpoints"`
}

type Repo struct {
  Name string `json:"repo"`
  Uri string `json:"uri"`
  Path string `json:"path"`
  Properties RepoProperties `json:"properties"`
}

type ArtifactoryResponse struct {
   Repo []Repo `json:"results"`
}

func (r *ArtifactoryResponse) Populate(resp *resty.Response) ArtifactoryResponse {
  err := json.Unmarshal(resp.Body(), &r)
  check(err)
  return *r
}

func (r *ArtifactoryResponse) UniqueProperties() (map[string]struct{}, map[string]struct{}) {
  baseVers := make(map[string]struct{})
  archVers := make(map[string]struct{})
  for _, repo := range r.Repo {
    for _, baseVer := range repo.Properties.BaseVerSupport {
      baseVers[baseVer] = struct{}{}
    }
    for _, archVer := range repo.Properties.ArchSupport {
      archVers[archVer] = struct{}{}
    }
    //fmt.Printf("[DEBUG] i: (%v) -- repo: (%v)\n", i, repo)
  }
  return baseVers, archVers
}

func IsSupported(repo Repo, baseVer, archVer string) (bool) {
  supported := false
  for _, repoBaseVal := range repo.Properties.BaseVerSupport {
    for _, repoArchVal := range repo.Properties.ArchSupport {
      if repoBaseVal == baseVer &&
      repoArchVal == archVer {
        supported = true
      }
    }
  }
  return supported
}

func GetArtifactoryRepos() (*resty.Response) {
  resp, err := resty.R().
               SetHeader("Accept", "application/json").
               SetHeader("X-Result-Detail", "info, properties").
               Get(fmt.Sprintf("http://%v:%v/artifactory/api/search/prop?yumrepo.auto_sync", viper.GetString("server"),
                                                                                             viper.GetString("port")))

  check(err)
  //fmt.Printf("[DEBUG] Response Type: %T\n", resp)
  return resp
}

