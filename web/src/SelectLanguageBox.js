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
