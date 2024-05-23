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
import {Button, Card, Col, Input, Row} from "antd";
import * as RecordBackend from "./backend/RecordBackend";
import * as Setting from "./Setting";
import i18next from "i18next";

// const {Option} = Select;

class RecordEditPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      owner: props.match.params.owner,
      id: props.match.params.id,
      record: null,
    };
  }

  UNSAFE_componentWillMount() {
    this.getRecord();
  }

  getRecord() {
    RecordBackend.getRecord(this.state.owner, this.state.id)
      .then((res) => {
        if (res.status === "ok") {
          this.setState({
            record: res.data,
          });
        } else {
          Setting.showMessage("error", `Failed to get record: ${res.msg}`);
        }
      });
  }

  updateRecordField(key, value) {
    const record = this.state.record;
    record[key] = value;
    this.setState({
      record: record,
    });
  }

  submitRecordEdit() {
    const record = Setting.deepCopy(this.state.record);
    RecordBackend.updateRecord(this.state.record.owner, this.state.id, record)
      .then((res) => {
        if (res.status === "error") {
          Setting.showMessage("error", `Failed to save: ${res.msg}`);
          this.updateRecordField("id", this.state.id);
        } else {
          Setting.showMessage("success", "Successfully saved");
          this.setState({
            id: this.state.record.id,
          });
          this.props.history.push(`/records/${this.state.record.owner}/${this.state.record.id}`);
          this.getRecord();
        }
      })
      .catch(error => {
        Setting.showMessage("error", `failed to save: ${error}`);
      });
  }

  renderRecord() {
    return (
      <Card size="small" title={
        <div>
          {i18next.t("general:Edit Record")}&nbsp;&nbsp;&nbsp;&nbsp;
          <Button type="primary" onClick={this.submitRecordEdit.bind(this)}>{i18next.t("general:Save")}</Button>
        </div>
      } style={{marginLeft: "5px"}} type="inner">
        <Row style={{marginTop: "10px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("general:Owner")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.record.owner} onChange={e => {
              this.updateRecordField("owner", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "10px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("general:CreatedTime")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.record.createdTime} onChange={e => {
              this.updateRecordField("createdTime", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "10px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("general:Method")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.record.method} onChange={e => {
              this.updateRecordField("method", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "10px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("general:Host")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.record.host} onChange={e => {
              this.updateRecordField("host", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "10px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("general:Path")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.record.path} onChange={e => {
              this.updateRecordField("path", e.target.value);
            }} />
          </Col>
        </Row>
        <Row style={{marginTop: "10px"}} >
          <Col style={{marginTop: "5px"}} span={2}>
            {i18next.t("general:UserAgent")}:
          </Col>
          <Col span={22} >
            <Input value={this.state.record.userAgent} onChange={e => {
              this.updateRecordField("userAgent", e.target.value);
            }} />
          </Col>
        </Row>
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
              this.state.record !== null ? this.renderRecord() : null
            }
          </Col>
          <Col span={1}>
          </Col>
        </Row>
        <Row style={{margin: 10}}>
          <Col span={2}>
          </Col>
          <Col span={18}>
            <Button type="primary" size="large" onClick={this.submitRecordEdit.bind(this)}>{i18next.t("general:Save")}</Button>
          </Col>
        </Row>
      </div>
    );
  }

}

export default RecordEditPage;
