// Copyright 2020 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/jacobsa/gcloud/gcs"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Create token source from the JSON file at the supplide path.
func newTokenSourceFromPath(
	path string,
	scope string) (ts oauth2.TokenSource, err error) {
	// Read the file.
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("ReadFile(%q): %v", path, err)
		return
	}

	// Create a config struct based on its contents.
	jwtConfig, err := google.JWTConfigFromJSON(contents, scope)
	if err != nil {
		err = fmt.Errorf("JWTConfigFromJSON: %v", err)
		return
	}

	// Create the token source.
	ts = jwtConfig.TokenSource(context.Background())

	return
}

// GetTokenSource returns a TokenSource for GCS API given a key file, or
// with the default credentials.
func GetTokenSource(keyFile string) (tokenSrc oauth2.TokenSource, err error) {
	// Create the oauth2 token source.
	const scope = gcs.Scope_FullControl

	if keyFile != "" {
		tokenSrc, err = newTokenSourceFromPath(keyFile, scope)
		if err != nil {
			err = fmt.Errorf("newTokenSourceFromPath: %v", err)
			return
		}
	} else {
		tokenSrc, err = google.DefaultTokenSource(context.Background(), scope)
		if err != nil {
			err = fmt.Errorf("DefaultTokenSource: %v", err)
			return
		}
	}
	return
}
