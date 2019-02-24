package main

import (
	"encoding/json"
	"errors"
	"github.com/prometheus/common/log"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"sync"
	"time"
)

// Takes validated results and stores them as the task specifies, either locally, remotely, or both
func StoreResults(finalResultChan <-chan FinalMIDAResult, monitoringChan chan<- TaskStats, retryChan chan<- SanitizedMIDATask, storageWG *sync.WaitGroup, pipelineWG *sync.WaitGroup) {
	for r := range finalResultChan {

		r.Stats.Timing.BeginStorage = time.Now()

		if !r.SanitizedTask.TaskFailed {
			// Store results here from a successfully completed task
			outputPathURL, err := url.Parse(r.SanitizedTask.OutputPath)
			if err != nil {
				Log.Error(err)
			} else {
				if outputPathURL.Host == "" {
					err = StoreResultsLocalFS(r)
					if err != nil {
						log.Error("Failed to store results: ", err)
					}
				} else {
					// Remote storage not yet implemented
					Log.Info("Remote storage not yet implemented")
				}
			}

		}

		// Remove all data from crawl
		// TODO: Add ability to save user data directory (without saving crawl data inside it)
		err := os.RemoveAll(r.SanitizedTask.UserDataDirectory)
		if err != nil {
			Log.Fatal(err)
		}

		if r.SanitizedTask.TaskFailed {
			if r.SanitizedTask.CurrentAttempt >= r.SanitizedTask.MaxAttempts {
				// We are abandoning trying this task. Too bad.
				Log.Error("Task failed after ", r.SanitizedTask.MaxAttempts, " attempts.")

			} else {
				// "Squash" task results and put the task back at the beginning of the pipeline
				Log.Debug("Retrying task...")
				r.SanitizedTask.CurrentAttempt++
				pipelineWG.Add(1)
				retryChan <- r.SanitizedTask
			}
		}

		r.Stats.Timing.EndStorage = time.Now()

		// Send stats to Prometheus
		if viper.GetBool("EnableMonitoring") {
			r.Stats.Timing.EndStorage = time.Now()
			monitoringChan <- r.Stats
		}

		pipelineWG.Done()
	}

	storageWG.Done()
}

// Given a valid FinalMIDAResult, stores it according to the output
// path specified in the sanitized task within the result
func StoreResultsLocalFS(r FinalMIDAResult) error {
	outpath := path.Join(r.SanitizedTask.OutputPath, r.SanitizedTask.RandomIdentifier)
	_, err := os.Stat(outpath)
	if err != nil {
		err = os.MkdirAll(outpath, 0755)
		if err != nil {
			Log.Error("Failed to create local output directory")
			return errors.New("failed to create local output directory")
		}
	} else {
		Log.Error("Output directory for task already exists")
		return errors.New("output directory for task already exists")
	}

	// Store resource metadata from crawl (DevTools requestWillBeSent and responseReceived data)
	if r.SanitizedTask.ResourceMetadata {
		data, err := json.Marshal(r.ResourceMetadata)
		if err != nil {
			Log.Error(err)
		}
		err = ioutil.WriteFile(path.Join(r.SanitizedTask.OutputPath, r.SanitizedTask.RandomIdentifier,
			DefaultResourceMetadataFile), data, 0644)
		if err != nil {
			Log.Error(err)
		}

	}

	// Store raw resources downloaded during crawl (named for their request IDs)
	if r.SanitizedTask.AllResources {
		_, err = os.Stat(path.Join(r.SanitizedTask.UserDataDirectory, r.SanitizedTask.RandomIdentifier, DefaultFileSubdir))
		if err != nil {
			Log.Error("AllResources requested but no files directory exists within temporary results directory")
			Log.Error("Files will not be stored")
			return errors.New("files temporary directory does not exist")
		} else {
			err = os.Rename(path.Join(r.SanitizedTask.UserDataDirectory, r.SanitizedTask.RandomIdentifier, DefaultFileSubdir),
				path.Join(outpath, DefaultFileSubdir))
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if r.SanitizedTask.ScriptMetadata {
		data, err := json.Marshal(r.ScriptMetadata)
		if err != nil {
			Log.Error(err)
		}
		err = ioutil.WriteFile(path.Join(r.SanitizedTask.OutputPath, r.SanitizedTask.RandomIdentifier,
			DefaultScriptMetadataFile), data, 0644)
		if err != nil {
			Log.Error(err)
		}
	}

	// Store raw scripts parsed by the browser during crawl (named by hashes)
	if r.SanitizedTask.AllScripts {
		_, err = os.Stat(path.Join(r.SanitizedTask.UserDataDirectory, r.SanitizedTask.RandomIdentifier, DefaultScriptSubdir))
		if err != nil {
			Log.Error("AllScripts requested but no files directory exists within temporary results directory")
			Log.Error("Scripts will not be stored")
			return errors.New("scripts temporary directory does not exist")
		} else {
			err = os.Rename(path.Join(r.SanitizedTask.UserDataDirectory, r.SanitizedTask.RandomIdentifier, DefaultScriptSubdir),
				path.Join(outpath, DefaultScriptSubdir))
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return nil
}