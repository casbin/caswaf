import * as Setting from "../Setting";

export function getGlobalCerts() {
  return fetch(`${Setting.ServerUrl}/api/get-global-certs`, {
    method: "GET",
    credentials: "include",
  }).then(res => res.json());
}

export function getCerts(owner) {
  return fetch(`${Setting.ServerUrl}/api/get-certs?owner=${owner}`, {
    method: "GET",
    credentials: "include",
  }).then(res => res.json());
}

export function getCert(owner, name) {
  return fetch(`${Setting.ServerUrl}/api/get-cert?id=${owner}/${encodeURIComponent(name)}`, {
    method: "GET",
    credentials: "include",
  }).then(res => res.json());
}

export function updateCert(owner, name, cert) {
  const newCert = Setting.deepCopy(cert);
  return fetch(`${Setting.ServerUrl}/api/update-cert?id=${owner}/${encodeURIComponent(name)}`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(newCert),
  }).then(res => res.json());
}

export function addCert(cert) {
  const newCert = Setting.deepCopy(cert);
  return fetch(`${Setting.ServerUrl}/api/add-cert`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(newCert),
  }).then(res => res.json());
}

export function deleteCert(cert) {
  const newCert = Setting.deepCopy(cert);
  return fetch(`${Setting.ServerUrl}/api/delete-cert`, {
    method: "POST",
    credentials: "include",
    body: JSON.stringify(newCert),
  }).then(res => res.json());
}
