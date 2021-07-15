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

package validation

import (
	"net/url"

	"k8s.io/apimachinery/pkg/util/validation/field"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// ValidateAddressable validates an addressable for non-optional addresses
func ValidateAddressable(addr duckv1.Addressable, optional bool, fld *field.Path) field.ErrorList {
	errs := field.ErrorList{}

	if (addr.URL == nil || addr.URL.String() == "") && !optional {
		errs = append(errs, field.Required(fld.Child("url"), "url is required"))
	} else if addr.URL != nil && !optional {
		if _, err := url.ParseRequestURI(addr.URL.String()); err != nil {
			errs = append(errs, field.Invalid(fld.Child("url"), addr.URL.String(), err.Error()))
		}
	}
	return errs
}
