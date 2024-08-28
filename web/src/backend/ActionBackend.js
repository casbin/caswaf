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

export function getActions(owner) {
  return fetch(`${Setting.ServerUrl}/api/get-actions?owner=${owner}`, {
    method: "GET",
    credentials: "include",
  }).then(res => res.json());
}

export function getAction(owner, name) {
  return fetch(`${Setting.ServerUrl}/api/get-action?id=${owner}/${encodeURIComponent(name)}`, {
    method: "GET",
    credentials: "include",
  }).then(res => res.json());
}

export function addAction(action) {
  return fetch(`${Setting.ServerUrl}/api/add-action`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(action),
  }).then(res => res.json());
}

export function updateAction(owner, name, action) {
  return fetch(`${Setting.ServerUrl}/api/update-action?id=${owner}/${encodeURIComponent(name)}`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(action),
  }).then(res => res.json());
}

export function deleteAction(action) {
  return fetch(`${Setting.ServerUrl}/api/delete-action`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(action),
  }).then(res => res.json());
}
