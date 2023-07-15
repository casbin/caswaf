import * as Setting from "../Setting";

export function getGlobalSites() {
  return fetch(`${Setting.ServerUrl}/api/get-global-sites`, {
    method: "GET",
    credentials: "include",
  }).then(res => res.json());
}

export function getSites(owner) {
  return fetch(`${Setting.ServerUrl}/api/get-sites?owner=${owner}`, {
    method: "GET",
    credentials: "include",
  }).then(res => res.json());
}

export function getSite(owner, name) {
  return fetch(`${Setting.ServerUrl}/api/get-site?id=${owner}/${encodeURIComponent(name)}`, {
    method: "GET",
    credentials: "include",
  }).then(res => res.json());
}

export function updateSite(owner, name, site) {
  const newSite = Setting.deepCopy(site);
  return fetch(`${Setting.ServerUrl}/api/update-site?id=${owner}/${encodeURIComponent(name)}`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(newSite),
  }).then(res => res.json());
}

export function addSite(site) {
  const newSite = Setting.deepCopy(site);
  return fetch(`${Setting.ServerUrl}/api/add-site`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(newSite),
  }).then(res => res.json());
}

export function deleteSite(site) {
  const newSite = Setting.deepCopy(site);
  return fetch(`${Setting.ServerUrl}/api/delete-site`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(newSite),
  }).then(res => res.json());
}
