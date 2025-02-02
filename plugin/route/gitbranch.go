/*
Copyright 2021 The Katanomi Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package route

import (
	"net/http"

	kerrors "github.com/katanomi/pkg/errors"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	"github.com/katanomi/pkg/plugin/client"
)

type gitBranchLister struct {
	impl client.GitBranchLister
	tags []string
}

// NewGitBranchLister create a git branch lister route with plugin client
func NewGitBranchLister(impl client.GitBranchLister) Route {
	return &gitBranchLister{
		tags: []string{"git", "repositories", "branch"},
		impl: impl,
	}
}

// Register route
func (a *gitBranchLister) Register(ws *restful.WebService) {
	repositoryParam := ws.PathParameter("repository", "branch belong to repository")
	projectParam := ws.PathParameter("project", "repository belong to project")
	ws.Route(
		ws.GET("/projects/{project}/coderepositories/{repository}/branches").To(a.ListBranch).
			Doc("ListBranch").Param(projectParam).Param(repositoryParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitBranchList{}),
	)
}

// ListBranch list branch by repo
func (a *gitBranchLister) ListBranch(request *restful.Request, response *restful.Response) {
	option := GetListOptionsFromRequest(request)
	repo := request.PathParameter("repository")
	project := request.PathParameter("project")
	branchList, err := a.impl.ListGitBranch(request.Request.Context(), metav1alpha1.GitRepo{Repository: repo, Project: project}, option)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, branchList)
}

type gitBranchCreator struct {
	impl client.GitBranchCreator
	tags []string
}

// NewGitBranchCreator create a git branch create route with plugin client
func NewGitBranchCreator(impl client.GitBranchCreator) Route {
	return &gitBranchCreator{
		tags: []string{"git", "repositories", "branch"},
		impl: impl,
	}
}

// Register route
func (a *gitBranchCreator) Register(ws *restful.WebService) {
	repositoryParam := ws.PathParameter("repository", "branch belong to repository")
	projectParam := ws.PathParameter("project", "repository belong to project")
	ws.Route(
		ws.POST("/projects/{project}/coderepositories/{repository}/branches").To(a.CreateBranch).
			Doc("CreateBranch").Param(projectParam).Param(repositoryParam).
			Metadata(restfulspec.KeyOpenAPITags, a.tags).
			Returns(http.StatusOK, "OK", metav1alpha1.GitBranch{}),
	)
}

// CreateBranch create branch
func (a *gitBranchCreator) CreateBranch(request *restful.Request, response *restful.Response) {
	repo := request.PathParameter("repository")
	project := request.PathParameter("project")
	var params metav1alpha1.CreateBranchParams
	if err := request.ReadEntity(&params); err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	payload := metav1alpha1.CreateBranchPayload{GitRepo: metav1alpha1.GitRepo{Repository: repo, Project: project}, CreateBranchParams: params}
	gitBranchObj, err := a.impl.CreateGitBranch(request.Request.Context(), payload)
	if err != nil {
		kerrors.HandleError(request, response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, gitBranchObj)
}
