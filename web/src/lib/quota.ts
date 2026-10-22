// Copyright 2026 The casbin Authors. All Rights Reserved.
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

import i18next from "i18next";

import type {QuotaConfig, QuotaWindow} from "@/types";

const currencySymbols: Record<string, string> = {USD: "$", CNY: "¥", EUR: "€", GBP: "£", JPY: "¥"};

/** An amount as the vendor bills it. A currency with a symbol gets it; anything
 * else keeps its own name, because a made-up symbol would be a guess. */
export function formatQuota(amount: number, unit: string) {
  const rounded = Math.abs(amount) < 1 && amount !== 0 ? amount.toFixed(4) : amount.toFixed(2);
  const symbol = currencySymbols[unit.toUpperCase()];
  if (symbol) {
    return amount < 0 ? `-${symbol}${rounded.slice(1)}` : `${symbol}${rounded}`;
  }
  return unit === "" ? rounded : `${rounded} ${unit}`;
}

export function emptyQuotaConfig(): QuotaConfig {
  return {
    url: "",
    headers: {},
    token: "",
    remaining: "",
    used: "",
    total: "",
    unit: "USD",
    scale: 0,
    manual: false,
    initial: 0,
    since: "",
  };
}

/** A ready-made endpoint for a family of sites, so the paths do not have to be
 * worked out from the vendor's API documentation. */
export interface QuotaPreset {
  key: string;
  label: string;
  /** The site, or the image, the project's own mark comes from. */
  icon: string;
  config: QuotaConfig;
}

export const quotaPresets: QuotaPreset[] = [
  {
    key: "newapi",
    label: "New API",
    icon: "newapi.pro",
    config: {
      ...emptyQuotaConfig(),
      url: "/api/user/self",
      headers: {Authorization: "Bearer {{token}}", "New-Api-User": ""},
      remaining: "data.quota",
      used: "data.used_quota",
      // New API counts in units of 1/500000 USD.
      scale: 500000,
    },
  },
  {
    key: "oneapi",
    label: "One API",
    // One API has no site of its own, and the logo in its repository is
    // OpenAI's, so it is marked by where the project lives instead.
    icon: "github.com",
    config: {
      ...emptyQuotaConfig(),
      url: "/api/user/self",
      headers: {Authorization: "Bearer {{token}}"},
      remaining: "data.quota",
      used: "data.used_quota",
      scale: 500000,
    },
  },
];

/** What a plan limit is called on the page. */
export function windowLabel(window: QuotaWindow) {
  switch (window.name) {
  case "week":
    return i18next.t("provider:Week");
  case "premium":
    return i18next.t("provider:Premium requests");
  case "chat":
    return i18next.t("provider:Chat");
  default:
    return window.name;
  }
}

/** How long, in minutes, until a window that still has most of its allowance
 *  resets and loses it. Zero when that is not about to happen: the allowance a
 *  plan gives for one window is not carried into the next. */
export function minutesToUnusedReset(window: QuotaWindow, now = Date.now()) {
  if ((window.name !== "5h" && window.name !== "week") || window.remainingPercent < 50 || !window.resetTime) {
    return 0;
  }
  const minutes = Math.ceil((new Date(window.resetTime).getTime() - now) / 60000);
  const soon = window.name === "week" ? 24 * 60 : 60;
  return minutes > 0 && minutes <= soon ? minutes : 0;
}
