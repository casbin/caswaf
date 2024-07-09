// Copyright 2021 The Casdoor Authors. All Rights Reserved.
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
import {Button, Popconfirm, Table} from "antd";
import moment from "moment";
import * as Setting from "./Setting";
import * as RuleBackend from "./backend/RuleBackend";
import i18next from "i18next";
import BaseListPage from "./BaseListPage";

class RuleListPage extends BaseListPage {
  UNSAFE_componentWillMount() {
    this.fetch();
  }

  fetch() {
    this.setState({
      loading: true,
    });
    RuleBackend.getRules().then((res) => {
      this.setState({
        data: res,
        loading: false,
      });
    });
  }

  addRule() {
    const newRule = this.newRule();
    RuleBackend.addRule(newRule).then((res) => {
      if (res.status === "error") {
        Setting.showMessage("error", `Failed to add: ${res.msg}`);
      } else {
        Setting.showMessage("success", "Rule added successfully");
        this.setState({
          data: Setting.prependRow(this.state.data, newRule),
        });
      }
    });
  }

  deleteRule(i) {
    // TODO: wait for backend implementation
    this.setState({
      data: Setting.deleteRow(this.state.data, i),
    });
  }

  newRule() {
    const randomName = Setting.getRandomName();
    return {
      owner: this.props.account.name,
      name: `rule_${randomName}`,
      createdTime: moment().format(),
      expression: "and key1 == value1 key2 == value2",
    };
  }

  renderTable(data) {
    const columns = [
      {
        title: i18next.t("general:Name"),
        dataIndex: "name",
        key: "name",
        width: "200px",
        sorter: (a, b) => a.name.localeCompare(b.name),
        render: (text, rule, index) => {
          return <a href={`/rules/${rule.owner}/${text}`}>{text}</a>;
        },
      },
      {
        title: i18next.t("general:Created Time"),
        dataIndex: "createdTime",
        key: "createdTime",
        width: "200px",
        sorter: (a, b) => a.createdTime.localeCompare(b.createdTime),
        render: (text, record, index) => {
          return Setting.getFormattedDate(text);
        },
      },
      {
        title: i18next.t("general:Updated Time"),
        dataIndex: "updatedTime",
        key: "updatedTime",
        width: "200px",
        sorter: (a, b) => a.updatedTime.localeCompare(b.updatedTime),
        render: (text, record, index) => {
          return Setting.getFormattedDate(text);
        },
      },
      {
        title: i18next.t("general:Expression"),
        dataIndex: "expression",
        key: "expression",
        sorter: (a, b) => a.expression.localeCompare(b.expression),
      },
      {
        title: i18next.t("general:Action"),
        dataIndex: "",
        key: "op",
        render: (text, rule, index) => {
          return (
            <Popconfirm
              title={`Sure to delete rule: ${rule.name} ?`}
              onConfirm={() => this.deleteRule(index)}
            >
              <Button style={{marginTop: "10px", marginBottom: "10px", marginRight: "10px"}} type="primary" onClick={() => this.props.history.push(`/rules/${rule.owner}/${rule.name}`)}>{i18next.t("general:Edit")}</Button>
              <Button type="danger">{i18next.t("general:Delete")}</Button>
            </Popconfirm>
          );
        },
      },
    ];

    return (
      <Table
        dataSource={data}
        columns={columns}
        rowKey="name"
        pagination={{pageSize: 1000}}
        loading={this.state.loading}
        onChange={this.handleTableChange}
        size="middle"
        bordered
        title={() => (
          <div>
            {i18next.t("general:Rule")}&nbsp;&nbsp;&nbsp;&nbsp;
            <Button type="primary" size="small" onClick={() => this.addRule()}>{i18next.t("general:Add")}</Button>
          </div>
        )}
      />
    );
  }
}

export default RuleListPage;
