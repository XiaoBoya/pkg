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

package client

import (
	"context"
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"

	metav1alpha1 "github.com/katanomi/pkg/apis/meta/v1alpha1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

type ClientGitContent interface {
	Get(ctx context.Context, baseURL *duckv1.Addressable, option metav1alpha1.GitRepoFileOption, options ...OptionFunc) (*metav1alpha1.GitRepoFile, error)
	Create(ctx context.Context, baseURL *duckv1.Addressable, payload metav1alpha1.CreateRepoFilePayload, options ...OptionFunc) (*metav1alpha1.GitCommit, error)
}

type gitContent struct {
	client Client
	meta   Meta
	secret corev1.Secret
}

func newGitContent(client Client, meta Meta, secret corev1.Secret) ClientGitContent {
	return &gitContent{
		client: client,
		meta:   meta,
		secret: secret,
	}
}

func (g *gitContent) Get(ctx context.Context, baseURL *duckv1.Addressable, option metav1alpha1.GitRepoFileOption, options ...OptionFunc) (*metav1alpha1.GitRepoFile, error) {
	fileInfo := &metav1alpha1.GitRepoFile{}
	options = append(options, MetaOpts(g.meta), SecretOpts(g.secret), QueryOpts(map[string]string{"ref": option.Ref}), ResultOpts(fileInfo))
	if option.Repository == "" {
		return nil, errors.New("repo is empty string")
	} else if option.Path == "" {
		return nil, errors.New("file path is empty string")
	}
	uri := fmt.Sprintf("projects/%s/coderepositories/%s/content/%s", option.Project, option.Repository, option.Path)
	if err := g.client.Get(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}
	return fileInfo, nil
}

func (g *gitContent) Create(ctx context.Context, baseURL *duckv1.Addressable, payload metav1alpha1.CreateRepoFilePayload, options ...OptionFunc) (*metav1alpha1.GitCommit, error) {
	commitInfo := &metav1alpha1.GitCommit{}
	options = append(options, MetaOpts(g.meta), SecretOpts(g.secret), BodyOpts(payload.CreateRepoFileParams), ResultOpts(commitInfo))
	if payload.Repository == "" {
		return nil, errors.New("repo is empty string")
	}
	uri := fmt.Sprintf("projects/%s/coderepositories/%s/content", payload.Project, payload.Repository)
	if err := g.client.Post(ctx, baseURL, uri, options...); err != nil {
		return nil, err
	}

	return commitInfo, nil
}
