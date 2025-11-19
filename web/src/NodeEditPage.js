// Copyright 2025 The casbin Authors. All Rights Reserved.
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
import {Button, Card, Col, Input, Row} from "antd";
import {LinkOutlined} from "@ant-design/icons";
import * as NodeBackend from "./backend/NodeBackend";
import * as Setting from "./Setting";
import i18next from "i18next";

const {TextArea} = Input;

class NodeEditPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      owner: props.match.params.owner,
      nodeName: props.match.params.nodeName,
      node: null,
    };
  }

  UNSAFE_componentWillMount() {
    this.getNode();
  }

  getNode() {
    NodeBackend.getNode(this.state.owner, this.state.nodeName)
      .then((res) => {
        if (res.status === "ok") {
          this.setState({
            node: res.data,
          });
        } else {
          Setting.showMessage("error", `Failed to get node: ${res.msg}`);
        }
      });
  }

  parseNodeField(key, value) {
    return value;
  }

  updateNodeField(key, value) {
    value = this.parseNodeField(key, value);

    const node = this.state.node;
    node[key] = value;
    this.setState({
      node: node,
    });
  }

  submitNodeEdit(willExist) {
    const node = Setting.deepCopy(this.state.node);
    NodeBackend.updateNode(this.state.owner, this.state.nodeName, node)
      .then((res) => {
        if (res.status === "ok") {
          if (res.data) {
            Setting.showMessage("success", "Successfully saved");
            this.setState({
              nodeName: this.state.node.name,
            });
            if (willExist) {
              this.props.history.push("/nodes");
            } else {
              this.props.history.push(`/nodes/${this.state.node.owner}/${this.state.node.name}`);
            }
          } else {
            Setting.showMessage("error", "failed to save: server side failure");
            this.updateNodeField("name", this.state.nodeName);
          }
        } else {
          Setting.showMessage("error", `failed to save: ${res.msg}`);
        }
      })
      .catch(error => {
        Setting.showMessage("error", `failed to save: ${error}`);
      });
  }

  render() {
    return (
      <div>
        <Row style={{width: "100%"}}>
          <Col span={1}>
          </Col>
          <Col span={22}>
            {
              this.state.node !== null ? this.renderNode() : null
            }
          </Col>
          <Col span={1}>
          </Col>
        </Row>
        <Row style={{margin: 10}}>
          <Col span={2}>
          </Col>
          <Col span={18}>
            <Button type="primary" size="large" onClick={() => this.submitNodeEdit(false)}>{i18next.t("general:Save")}</Button>
            <Button style={{marginLeft: "20px"}} type="primary" size="large" onClick={() => this.submitNodeEdit(true)}>{i18next.t("general:Save & Exit")}</Button>
          </Col>
        </Row>
      </div>
    );
  }

  renderNode() {
    return (
      <Card size="small" title={
        <div>
          {i18next.t("node:Edit Node")}&nbsp;&nbsp;&nbsp;&nbsp;
          <Button onClick={() => this.submitNodeEdit(false)}>{i18next.t("general:Save")}</Button>
          <Button style={{marginLeft: "20px"}} onClick={() => this.submitNodeEdit(true)}>{i18next.t("general:Save & Exit")}</Button>
        </div>
      } style={{marginLeft: "5px"}} type="inner">
        <Row style={{marginTop: "10px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("general:Name")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.node.name} onChange={e => {
              this.updateNodeField("name", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("general:Display name")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.node.displayName} onChange={e => {
              this.updateNodeField("displayName", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("general:Tag")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.node.tag} onChange={e => {
              this.updateNodeField("tag", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("node:Hostname")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.node.hostname} onChange={e => {
              this.updateNodeField("hostname", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("node:IP address")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.node.ipAddress} onChange={e => {
              this.updateNodeField("ipAddress", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "20px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("node:Description")}:
          </Col>
          <Col span={22} >
            <TextArea autoSize={{minRows: 3, maxRows: 100}} value={this.state.node.description} onChange={e => {
              this.updateNodeField("description", e.target.value);
            }} />
          </Col>
        </Row>
      </Card>
    );
  }
}

export default NodeEditPage;
