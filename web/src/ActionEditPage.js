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
import {Button, Card, Col, Input, InputNumber, Row, Select} from "antd";
import * as ActionBackend from "./backend/ActionBackend";
import * as Setting from "./Setting";
import i18next from "i18next";

const {Option} = Select;

class ActionEditPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      owner: props.match.params.owner,
      actionName: props.match.params.actionName,
      action: null,
    };
  }

  UNSAFE_componentWillMount() {
    this.getAction();
  }

  getAction() {
    ActionBackend.getAction(this.state.owner, this.state.actionName)
      .then((res) => {
        if (res.status === "ok") {
          this.setState({
            action: res.data,
          });
        } else {
          Setting.showMessage("error", `Failed to get action: ${res.msg}`);
        }
      });
  }

  updateActionField(key, value) {
    const action = this.state.action;
    action[key] = value;
    this.setState({
      action: action,
    });
  }

  submitActionEdit() {
    const action = Setting.deepCopy(this.state.action);
    ActionBackend.updateAction(this.state.action.owner, this.state.actionName, action)
      .then((res) => {
        if (res.status === "error") {
          Setting.showMessage("error", `Failed to save: ${res.msg}`);
          this.setState({
            action: action,
          });
        } else {
          Setting.showMessage("success", "Successfully saved");
          this.setState({
            actioName: this.state.action.name,
          });
          this.props.history.push(`/actions/${this.state.action.owner}/${this.state.action.name}`);
          this.getAction();
        }
      })
      .catch(error => {
        Setting.showMessage("error", `failed to save: ${error}`);
      });
  }

  renderAction() {
    return (
      <Card size="small" title={
        <div>
          {i18next.t("general:Edit Action")}&nbsp;&nbsp;&nbsp;&nbsp;
          <Button type="primary" onClick={this.submitActionEdit.bind(this)}>{i18next.t("general:Save")}</Button>
        </div>
      } style={{marginTop: "10px"}} type="inner">
        <Row style={{marginTop: "20px"}}>
          <Col span={2} style={{marginTop: "5px"}}>
            {i18next.t("general:Name")}:
          </Col>
          <Col span={22}>
            <Input value={this.state.action.name} disabled={true} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("rule:Type")}:
          </Col>
          <Col span={22} >
            <Select virtual={false} value={this.state.action.type} defaultValue={"Block"} style={{width: "100%"}} onChange={value => {
              this.updateActionField("type", value);
            }}>
              {
                [
                  {value: "Allow", text: i18next.t("rule:Allow")},
                  // {value: "redirect", text: "Redirect"},
                  {value: "Block", text: i18next.t("rule:Block")},
                  // {value: "drop", text: "Drop"},
                  {value: "CAPTCHA", text: i18next.t("rule:Captcha")},
                ].map((item, index) => <Option key={index} value={item.value}>{item.text}</Option>)
              }
            </Select>
          </Col>
        </Row>
        {
          (this.state.action.type === "CAPTCHA") ? (
            <Row style={{marginTop: "20px"}} >
              <Col style={{marginTop: "5px"}} span={2}>
                {i18next.t("rule:Immunity times")}:
              </Col>
              <Col span={22} >
                <InputNumber value={this.state.action.immunityTimes} addonAfter={i18next.t("usage:minutes")} onChange={e => {
                  this.updateActionField("immunityTimes", e);
                }} />
              </Col>
            </Row>
          ) : null
        }
        {
          this.state.action.type === "Allow" || this.state.action.type === "Block" ? (
            <Row style={{marginTop: "20px"}} >
              <Col style={{marginTop: "5px"}} span={2}>
                {i18next.t("rule:Status Code")}:
              </Col>
              <Col span={22} >
                <InputNumber value={this.state.action.statusCode} min={100} max={599} onChange={e => {
                  this.updateActionField("statusCode", e);
                }} />
              </Col>
            </Row>
          ) : null
        }
      </Card>
    );
  }

  render() {
    return (
      <div>
        <Row style={{width: "100%"}}>
          <Col span={1}>
          </Col>
          <Col span={22}>
            {
              this.state.action !== null ? this.renderAction() : null
            }
          </Col>
          <Col span={1}>
          </Col>
        </Row>
        <Row style={{margin: 10}}>
          <Col span={2}>
          </Col>
          <Col span={18}>
            <Button type="primary" size="large" onClick={this.submitActionEdit.bind(this)}>{i18next.t("general:Save")}</Button>
          </Col>
        </Row>
      </div>
    );
  }
}

export default ActionEditPage;
