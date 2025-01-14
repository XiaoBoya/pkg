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
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/emicklei/go-restful/v3"
	"github.com/katanomi/pkg/errors"
)

const (
	// PluginMetaHeader header to store metadata for the plugin
	PluginMetaHeader = "X-Plugin-Meta"
)

type metaContextKey struct{}

// Meta Plugin meta with base url and version info, for calling plugin api
type Meta struct {
	Version string `json:"version,omitempty"`
	BaseURL string `json:"baseURL,omitempty"`
}

// WithContext returns a copy of parent include with the plugin meta
func (p *Meta) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, metaContextKey{}, p)
}

// ExtraMeta extract meta from a specific context
func ExtraMeta(ctx context.Context) *Meta {
	value := ctx.Value(metaContextKey{})
	if v, ok := value.(*Meta); ok {
		return v
	}

	return nil
}

// MetaFilter meta filter for go restful, parsing plugin meta
func MetaFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	encodedMeta := req.HeaderParameter(PluginMetaHeader)
	if encodedMeta == "" {
		chain.ProcessFilter(req, resp)
		return
	}

	decodedMeta, err := base64.StdEncoding.DecodeString(encodedMeta)
	if err != nil {
		errors.HandleError(req, resp, fmt.Errorf("decode meta error: %s", err.Error()))
		return
	}

	meta := &Meta{}
	if err = json.Unmarshal(decodedMeta, meta); err != nil {
		errors.HandleError(req, resp, fmt.Errorf("decode meta error: %s", err.Error()))
		return
	}

	ctx := req.Request.Context()
	req.Request = req.Request.WithContext(meta.WithContext(ctx))

	chain.ProcessFilter(req, resp)
}
