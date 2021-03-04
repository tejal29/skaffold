/*
Copyright 2020 The Skaffold Authors

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

package errors

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/instrumentation"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/runner/runcontext"
	"github.com/GoogleContainerTools/skaffold/proto/v1"
)

const (
	// These are phases in a Skaffolld
	Init        = Phase("Init")
	Build       = Phase("Build")
	Deploy      = Phase("Deploy")
	StatusCheck = Phase("StatusCheck")
	FileSync    = Phase("FileSync")
	DevInit     = Phase("DevInit")
	Cleanup     = Phase("Cleanup")

	// Report issue text
	reportIssueText = "If above error is unexpected, please open an issue https://github.com/GoogleContainerTools/skaffold/issues/new to report this error"
)

var (
	setRunContextOnce sync.Once
	runCtx            runcontext.RunContext

	reportIssueSuggestion = func(runcontext.RunContext) []*proto.Suggestion {
		return []*proto.Suggestion{{
			SuggestionCode: proto.SuggestionCode_OPEN_ISSUE,
			Action:         reportIssueText,
		}}
	}
)

type Phase string

// SetRunContext set Skaffold runCtx  once. This run context is used later to
// suggest actionable error messages based on skaffold command line options and run context
func SetRunContext(rc runcontext.RunContext) {
	setRunContextOnce.Do(func() {
		runCtx = rc
	})
}

// ActionableErr returns an actionable error message with suggestions
func ActionableErr(phase Phase, err error) *proto.ActionableErr {
	errCode, suggestions := getErrorCodeFromError(phase, err)
	return &proto.ActionableErr{
		ErrCode:     errCode,
		Message:     err.Error(),
		Suggestions: suggestions,
	}
}

func ShowAIError(err error) error {
	if IsSkaffoldErr(err) {
		instrumentation.SetErrorCode(err.(Error).StatusCode())
		return err
	}
	if p, ok := isProblem(err); ok {
		instrumentation.SetErrorCode(p.errCode)
		return p
	}

	var knownProblems = append(knownBuildProblems, knownDeployProblems...)
	for _, p := range append(knownProblems, knownInitProblems...) {
		if p.regexp.MatchString(err.Error()) {
			instrumentation.SetErrorCode(p.errCode)
			return p.withErr(err)
		}
	}
	return err
}

func IsOldImageManifestProblem(err error) (string, bool) {
	if err != nil && oldImageManifest.regexp.MatchString(err.Error()) {
		if s := oldImageManifest.suggestion(runCtx); s != nil {
			return fmt.Sprintf("%s. %s", oldImageManifest.description(err),
				concatSuggestions(oldImageManifest.suggestion(runCtx))), true
		}
		return "", true
	}
	return "", false
}

func getErrorCodeFromError(phase Phase, err error) (proto.StatusCode, []*proto.Suggestion) {
	var sErr Error
	if errors.As(err, &sErr) {
		return sErr.StatusCode(), sErr.Suggestions()
	}
	var p problem
	if errors.As(err, &p) {
		return p.errCode, p.suggestion(runCtx)
	}

	if problems, ok := allErrors[phase]; ok {
		for _, v := range problems {
			if v.regexp.MatchString(err.Error()) {
				return v.errCode, v.suggestion(runCtx)
			}
		}
	}
	return proto.StatusCode_UNKNOWN_ERROR, nil
}

func concatSuggestions(suggestions []*proto.Suggestion) string {
	var s strings.Builder
	for _, suggestion := range suggestions {
		if s.String() != "" {
			s.WriteString(" or ")
		}
		s.WriteString(suggestion.Action)
	}
	if s.String() == "" {
		return ""
	}
	s.WriteString(".")
	return s.String()
}

var allErrors = map[Phase][]problem{
	Build: append(knownBuildProblems, problem{
		regexp:     re(".*"),
		errCode:    proto.StatusCode_BUILD_UNKNOWN,
		suggestion: reportIssueSuggestion,
	}),
	Init: append(knownInitProblems, problem{
		regexp:     re(".*"),
		errCode:    proto.StatusCode_INIT_UNKNOWN,
		suggestion: reportIssueSuggestion,
	}),
	Deploy: append(knownDeployProblems, problem{
		regexp:     re(".*"),
		errCode:    proto.StatusCode_DEPLOY_UNKNOWN,
		suggestion: reportIssueSuggestion,
	}),
	StatusCheck: {{
		regexp:     re(".*"),
		errCode:    proto.StatusCode_STATUSCHECK_UNKNOWN,
		suggestion: reportIssueSuggestion,
	}},
	FileSync: {{
		regexp:     re(".*"),
		errCode:    proto.StatusCode_SYNC_UNKNOWN,
		suggestion: reportIssueSuggestion,
	}},
	DevInit: {oldImageManifest, {
		regexp:     re(".*"),
		errCode:    proto.StatusCode_DEVINIT_UNKNOWN,
		suggestion: reportIssueSuggestion,
	}},
	Cleanup: {{
		regexp:     re(".*"),
		errCode:    proto.StatusCode_CLEANUP_UNKNOWN,
		suggestion: reportIssueSuggestion,
	}},
}
