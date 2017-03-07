// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/cihub/seelog"
)

const registryCacheVersion = "1.0"

type RegistryCache struct {
	Registries map[string]*AuthEntry
	Version    string
}

type fileCredentialCache struct {
	path           string
	filename       string
	cachePrefixKey string
}

func newRegistryCache() *RegistryCache {
	return &RegistryCache{
		Registries: make(map[string]*AuthEntry),
		Version:    registryCacheVersion,
	}
}

// NewFileCredentialsCache returns a new file credentials cache.
//
// path is used for temporary files during save, and filename should be a relative filename
// in the same directory where the cache is serialized and deserialized.
//
// cachePrefixKey is used for scoping credentials for a given credential cache (i.e. region and
// accessKey).
func NewFileCredentialsCache(path string, filename string, cachePrefixKey string) CredentialsCache {
	if _, err := os.Stat(path); err != nil {
		os.MkdirAll(path, 0700)
	}
	return &fileCredentialCache{path: path, filename: filename, cachePrefixKey: cachePrefixKey}
}

func (f *fileCredentialCache) Get(registry string) *AuthEntry {
	log.Debugf("Checking file cache for %s", registry)
	registryCache := f.init()
	return registryCache.Registries[f.cachePrefixKey+registry]
}

func (f *fileCredentialCache) Set(registry string, entry *AuthEntry) {
	log.Debugf("Saving credentials to file cache for %s", registry)
	registryCache := f.init()

	registryCache.Registries[f.cachePrefixKey+registry] = entry

	err := f.save(registryCache)
	if err != nil {
		log.Infof("Could not save cache: %s", err)
	}
}

// List returns all of the available AuthEntries (regardless of prefix)
func (f *fileCredentialCache) List() []*AuthEntry {
	registryCache := f.init()

	// optimize allocation for copy
	entries := make([]*AuthEntry, 0, len(registryCache.Registries))

	for _, entry := range registryCache.Registries {
		entries = append(entries, entry)
	}

	return entries
}

func (f *fileCredentialCache) Clear() {
	err := os.Remove(f.fullFilePath())
	if err != nil {
		log.Infof("Could not clear cache: %s")
	}
}

func (f *fileCredentialCache) fullFilePath() string {
	return filepath.Join(f.path, f.filename)
}

// Saves credential cache to disk. This writes to a temporary file first, then moves the file to the config location.
// This elminates from reading partially written credential files, and reduces (but does not eliminate) concurrent
// file access. There is not guarantee here for handling multiple writes at once since there is no out of process locking.
func (f *fileCredentialCache) save(registryCache *RegistryCache) error {
	defer log.Flush()

	file, err := ioutil.TempFile(f.path, ".config.json.tmp")
	if err != nil {
		return err
	}

	buff, err := json.MarshalIndent(registryCache, "", "  ")
	if err != nil {
		file.Close()
		os.Remove(file.Name())
		return err
	}

	_, err = file.Write(buff)

	if err != nil {
		file.Close()
		os.Remove(file.Name())
		return err
	}

	file.Close()
	// note this is only atomic when relying on linux syscalls
	os.Rename(file.Name(), f.fullFilePath())
	return err
}

func (f *fileCredentialCache) init() *RegistryCache {
	registryCache, err := f.load()
	if err != nil {
		log.Infof("Could not load existing cache: %v", err)
		f.Clear()
		registryCache = newRegistryCache()
	}
	return registryCache
}

// Loading a cache from disk will return errors for malformed or incompatible cache files.
func (f *fileCredentialCache) load() (*RegistryCache, error) {
	registryCache := newRegistryCache()

	file, err := os.Open(f.fullFilePath())
	if os.IsNotExist(err) {
		return registryCache, nil
	}

	if err != nil {
		return nil, err
	}

	defer file.Close()

	if err = json.NewDecoder(file).Decode(&registryCache); err != nil {
		return nil, err
	}

	if registryCache.Version != registryCacheVersion {
		return nil, fmt.Errorf("Registry cache version %#v is not compatible with %#v. Ignoring existing cache.",
			registryCache.Version,
			registryCacheVersion)
	}

	return registryCache, nil
}
