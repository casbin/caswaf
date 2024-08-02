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
import {Button, Popconfirm, Table, Tag} from "antd";
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
    RuleBackend.getRules(this.props.account.name).then((res) => {
      this.setState({
        data: res.data,
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
    RuleBackend.deleteRule(this.state.data[i]).then((res) => {
      if (res.status === "error") {
        Setting.showMessage("error", `Failed to delete: ${res.msg}`);
      } else {
        Setting.showMessage("success", "Deleted successfully");
        this.setState({
          data: Setting.deleteRow(this.state.data, i),
        });
      }
    });
  }

  newRule() {
    const randomName = Setting.getRandomName();
    return {
      owner: this.props.account.name,
      name: `rule_${randomName}`,
      createdTime: moment().format(),
      type: "User-Agent",
      expressions: [],
      action: "Block",
      reason: "Your request is blocked.",
    };
  }

  renderTable(data) {
    const columns = [
      {
        title: i18next.t("general:Owner"),
        dataIndex: "owner",
        key: "owner",
        width: "200px",
        sorter: (a, b) => a.owner.localeCompare(b.owner),
      },
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
        title: i18next.t("general:Create time"),
        dataIndex: "createdTime",
        key: "createdTime",
        width: "200px",
        sorter: (a, b) => a.createdTime.localeCompare(b.createdTime),
        render: (text, rule, index) => {
          return Setting.getFormattedDate(text);
        },
      },
      {
        title: i18next.t("general:Update time"),
        dataIndex: "updatedTime",
        key: "updatedTime",
        width: "200px",
        sorter: (a, b) => a.updatedTime.localeCompare(b.updatedTime),
        render: (text, rule, index) => {
          return Setting.getFormattedDate(text);
        },
      },
      {
        title: i18next.t("general:Type"),
        dataIndex: "type",
        key: "type",
        sorter: (a, b) => a.type.localeCompare(b.type),
      },
      {
        title: i18next.t("general:Expressions"),
        dataIndex: "expressions",
        key: "expressions",
        sorter: (a, b) => a.expressions.localeCompare(b.expressions),
        render: (text, rule, index) => {
          return rule.expressions.map((expression, i) => {
            return (
              <Tag key={expression} color={"success"}>
                {expression.operator + " " + expression.value.slice(0, 20)}
              </Tag>
            );
          });
        },
      },
      {
        title: i18next.t("general:Action"),
        dataIndex: "",
        key: "op",
        render: (text, rule, index) => {
          return (
            <div>
              <Popconfirm
                title={`Sure to delete rule: ${rule.name} ?`}
                onConfirm={() => this.deleteRule(index)}
              >
                <Button style={{marginTop: "10px", marginBottom: "10px", marginRight: "10px"}} type="primary" onClick={() => this.props.history.push(`/rules/${rule.owner}/${rule.name}`)}>{i18next.t("general:Edit")}</Button>
                <Button type="danger">{i18next.t("general:Delete")}</Button>
              </Popconfirm>
            </div>
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
