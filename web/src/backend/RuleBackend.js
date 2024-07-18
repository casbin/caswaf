// Copyright 2023 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import * as Setting from "../Setting";

export function getRules() {
  return fetch(`${Setting.ServerUrl}/api/get-rules`, {
    method: "GET",
    credentials: "include",
  }).then(res => res.json());
}

export function getRule(owner, name) {
  return fetch(`${Setting.ServerUrl}/api/get-rule?id=${owner}/${encodeURIComponent(name)}`, {
    method: "GET",
    credentials: "include",
  }).then(res => res.json());
}

export function addRule(rule) {
  return fetch(`${Setting.ServerUrl}/api/add-rule`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(rule),
  }).then(res => res.json());
}

export function updateRule(owner, name, rule) {
  return fetch(`${Setting.ServerUrl}/api/update-rule?id=${owner}/${encodeURIComponent(name)}`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(rule),
  }).then(res => res.json());
}

export function deleteRule(rule) {
  return fetch(`${Setting.ServerUrl}/api/delete-rule`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(rule),
  }).then(res => res.json());
}
