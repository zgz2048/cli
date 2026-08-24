// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package base

import (
	"encoding/json"
	"testing"

	"github.com/larksuite/cli/internal/citation"
	"github.com/larksuite/cli/internal/envvars"
	"github.com/larksuite/cli/internal/httpmock"
)

func TestBaseGetCitationEnvelope(t *testing.T) {
	t.Setenv(envvars.CliCitation, "1")
	factory, stdout, reg := newExecuteFactory(t)
	reg.Register(&httpmock.Stub{
		Method: "GET",
		URL:    "/open-apis/base/v3/bases/app_x",
		Body: map[string]interface{}{
			"code": 0,
			"data": map[string]interface{}{
				"base_token": "app_x",
				"name":       "Demo Base",
				"url":        "https://example.larkoffice.com/base/app_x",
			},
		},
	})

	if err := runShortcut(t, BaseBaseGet, []string{"+base-get", "--base-token", "app_x"}, factory, stdout); err != nil {
		t.Fatalf("err=%v", err)
	}
	var envelope struct {
		Citations []citation.Citation `json:"citations"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &envelope); err != nil {
		t.Fatalf("decode output: %v\nraw=%s", err, stdout.String())
	}
	if len(envelope.Citations) != 1 {
		t.Fatalf("citations = %#v, want 1 entry", envelope.Citations)
	}
	got := envelope.Citations[0]
	if got.SourceType != citation.SourceBase || got.URL != "https://example.larkoffice.com/base/app_x" {
		t.Errorf("source_type/url = %d %q", got.SourceType, got.URL)
	}
	if got.Title != "Demo Base" || got.ResourceID != "app_x" {
		t.Errorf("title/resource_id = %q %q", got.Title, got.ResourceID)
	}
}

func TestBaseGetCitationDoesNotGuessMissingURL(t *testing.T) {
	t.Setenv(envvars.CliCitation, "1")
	factory, stdout, reg := newExecuteFactory(t)
	reg.Register(&httpmock.Stub{
		Method: "GET",
		URL:    "/open-apis/base/v3/bases/app_x",
		Body: map[string]interface{}{
			"code": 0,
			"data": map[string]interface{}{"base_token": "app_x", "name": "Demo Base"},
		},
	})

	if err := runShortcut(t, BaseBaseGet, []string{"+base-get", "--base-token", "app_x"}, factory, stdout); err != nil {
		t.Fatalf("err=%v", err)
	}
	var envelope map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &envelope); err != nil {
		t.Fatalf("decode output: %v\nraw=%s", err, stdout.String())
	}
	if _, ok := envelope["citations"]; ok {
		t.Fatalf("missing native URL must omit citations: %s", stdout.String())
	}
}
