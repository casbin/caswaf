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

import React from "react";
import * as Setting from "./Setting";
import {Dropdown, Menu} from "antd";
import {createFromIconfontCN} from "@ant-design/icons";
import "./App.less";

const IconFont = createFromIconfontCN({
  scriptUrl: "//at.alicdn.com/t/font_2680620_ffij16fkwdg.js",
});

const LanguageItems = [
  {lang: "en", label: "English", icon: "icon-en"},
  {lang: "zh", label: "中文", icon: "icon-zh"},
];

class SelectLanguageBox extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
    };
  }

  render() {
    return <Dropdown overlay={<Menu>{LanguageItems.map(({lang, label, icon}) => <Menu.Item key={lang} onClick={() => Setting.changeLanguage(lang)}><IconFont type={icon} />{label}</Menu.Item>)}</Menu>}>
      <div className="language-box"></div>
    </Dropdown>;
  }
}

export default SelectLanguageBox;
