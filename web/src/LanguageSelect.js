// Copyright 2024 The Casbin Authors. All Rights Reserved.
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

import React from "react";
import * as Setting from "./Setting";
import {Dropdown, Menu} from "antd";
import {GlobalOutlined} from "@ant-design/icons";

function flagIcon(country, alt) {
  return <img src={`${Setting.StaticBaseUrl}/flag-icons/${country}.svg`} alt={alt} style={{marginRight: 8, width: 24}} />;
}

class LanguageSelect extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      languages: Setting.Countries.map(item => item.key),
    };

    // Preload flag icons
    Setting.Countries.forEach((country) => {
      new Image().src = `${Setting.StaticBaseUrl}/flag-icons/${country.country}.svg`;
    });
  }

  items = Setting.Countries.map((country) => ({
    key: country.key,
    label: (
      <span>
        {flagIcon(country.country, country.alt)}
        {country.label}
      </span>
    ),
  }));

  getLanguages(languages) {
    const select = [];
    for (const language of languages) {
      this.items.forEach((item) => item.key === language ? select.push(item) : null);
    }
    return select;
  }

  render() {
    const languageItems = this.getLanguages(this.state.languages);

    const menu = (
      <Menu onClick={(e) => Setting.setLanguage(e.key)}>
        {languageItems.map(item => (
          <Menu.Item key={item.key}>
            {item.label}
          </Menu.Item>
        ))}
      </Menu>
    );

    return (
      <Dropdown overlay={menu}>
        <div className="select-box" style={{display: languageItems.length === 0 ? "none" : null, ...this.props.style}}>
          <GlobalOutlined style={{fontSize: "24px"}} />
        </div>
      </Dropdown>
    );
  }
}

export default LanguageSelect;
