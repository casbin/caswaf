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
import * as Setting from "./Setting";
import * as RuleBackend from "./backend/RuleBackend";
import i18next from "i18next";
import WafRuleTable from "./components/WafRuleTable";

const {Option} = Select;

class RuleEditPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      owner: props.match.params.owner,
      ruleName: props.match.params.ruleName,
      rule: null,
    };
  }

  UNSAFE_componentWillMount() {
    this.getRule();
  }

  getRule() {
    RuleBackend.getRule(this.state.owner, this.state.ruleName).then((res) => {
      this.setState({
        rule: res.data,
      });
    });
  }

  updateRuleField(key, value) {
    const rule = Setting.deepCopy(this.state.rule);
    rule[key] = value;
    this.setState({
      rule: rule,
    });
  }

  renderRule() {
    return (
      <Card size="small" title={
        <div>
          Edit Rule&nbsp;&nbsp;&nbsp;&nbsp;
          <Button type="primary" onClick={this.submitRuleEdit.bind(this)}>Save</Button>
        </div>
      } style={{marginTop: 10}} type="inner">
        <Row style={{marginTop: "20px"}}>
          <Col span={2} style={{marginTop: "5px"}}>
            {i18next.t("rule:Name")}:
          </Col>
          <Col span={22}>
            <Input value={this.state.rule.name} disabled={true} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}}>
          <Col span={2} style={{marginTop: "5px"}}>
            {i18next.t("rule:Type")}:
          </Col>
          <Col span={22}>
            <Select virtual={false} value={this.state.rule.type} style={{width: "100%"}} onChange={value => {
              this.updateRuleField("type", value);
            }}>
              {
                [
                  {value: "waf", text: "WAF"},
                  {value: "ip", text: "IP"},
                  {value: "ua", text: "UA"},
                  {value: "frequency", text: "Frequency"},
                  {value: "complex", text: "Complex"},
                ].map((item, index) => <Option key={index} value={item.value}>{item.text}</Option>)
              }
            </Select>
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            Expressions:
          </Col>
          <Col span={22} >
            {
              this.state.rule.type === "waf" ? (
                <WafRuleTable
                  title={"Expressions"}
                  table={this.state.rule.expressions}
                  ruleName={this.state.rule.name}
                  account={this.props.account}
                  onUpdateTable={(value) => {this.updateRuleField("expressions", value);}}
                />
              ) : null
            }
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}}>
          <Col span={2} style={{marginTop: "5px"}}>
            {i18next.t("rule:Disruptive Action")}:
          </Col>
          <Col span={22}>
            <Select virtual={false} value={this.state.rule.disruptiveAction} style={{width: "100%"}} onChange={value => {
              this.updateRuleField("disruptiveAction", value);
            }}>
              {
                [
                  {value: "allow", text: "Allow"},
                  {value: "redirect", text: "Redirect"},
                  {value: "deny", text: "Deny"},
                  {value: "drop", text: "Drop"},
                ].map((item, index) => <Option key={index} value={item.value}>{item.text}</Option>)
              }
            </Select>
          </Col>
        </Row>
        {
          this.state.rule.disruptiveAction === "redirect" && (
            <Row style={{marginTop: "20px"}}>
              <Col span={2} style={{marginTop: "5px"}}>
                {i18next.t("rule:Redirect URL")}:
              </Col>
              <Col span={22}>
                <Input value={this.state.rule.redirectUrl} onChange={e => {
                  this.updateRuleField("redirectUrl", e.target.value);
                }} />
              </Col>
            </Row>
          )
        }
        {
          this.state.rule.disruptiveAction === "deny" && (
            <Row style={{marginTop: "20px"}}>
              <Col span={2} style={{marginTop: "5px"}}>
                {i18next.t("rule:Status Code")}:
              </Col>
              <Col span={22}>
                <InputNumber value={this.state.rule.statusCode} defaultValue={403} onChange={e => {
                  this.updateRuleField("statusCode", e.target.value);
                }} />
              </Col>
            </Row>
          )
        }
        <Row style={{marginTop: "20px"}}>
          <Col span={2} style={{marginTop: "5px"}}>
            {i18next.t("rule:LogAction")}:
          </Col>
          <Col span={22}>
            <Select virtual={false} value={this.state.rule.logAction} style={{width: "100%"}} onChange={value => {
              this.updateRuleField("logAction", value);
            }}>
              {
                [
                  {value: "log", text: "Log"},
                  {value: "nolog", text: "Not-log"},
                ].map((item, index) => <Option key={index} value={item.value}>{item.text}</Option>)
              }
            </Select>
          </Col>
        </Row>
        {this.state.rule.logAction === "log" && (
          <Row style={{marginTop: "20px"}}>
            <Col span={2} style={{marginTop: "5px"}}>
              {i18next.t("rule:Log Message")}:
            </Col>
            <Col span={22}>
              <Input value={this.state.rule.logMessage} onChange={e => {
                this.updateRuleField("logMessage", e.target.value);
              }} />
            </Col>
          </Row>
        )}
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
              this.state.rule !== null ? this.renderRule() : null
            }
          </Col>
          <Col span={1}>
          </Col>
        </Row>
        <Row style={{margin: 10}}>
          <Col span={2}>
          </Col>
          <Col span={18}>
            <Button type="primary" size="large" onClick={this.submitRuleEdit.bind(this)}>{i18next.t("general:Save")}</Button>
          </Col>
        </Row>
      </div>
    );
  }

  submitRuleEdit() {
    const rule = Setting.deepCopy(this.state.rule);
    RuleBackend.updateRule(this.state.owner, this.state.ruleName, rule)
      .then((res) => {
        if (res.status !== "error") {
          Setting.showMessage("success", "Rule updated successfully");
          this.setState({
            rule: rule,
          });
        } else {
          Setting.showMessage("error", `Rule failed to update: ${res.msg}`);
          this.setState({
            ruleName: this.state.rule.name,
          });
          this.props.history.push(`/rules/${this.state.rule.owner}/${this.state.rule.name}`);
          this.getRule();
        }
      });
  }
}

export default RuleEditPage;
