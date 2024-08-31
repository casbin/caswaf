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
import * as ActionBackend from "./backend/ActionBackend";
import i18next from "i18next";
import BaseListPage from "./BaseListPage";

class ActionListPage extends BaseListPage {
  UNSAFE_componentWillMount() {
    this.fetch();
  }

  fetch = (params = {}) => {
    this.setState({
      loading: true,
    });
    ActionBackend.getActions(this.props.account.name).then((res) => {
      this.setState({
        loading: false,
      });
      if (res.status === "ok") {
        this.setState({
          data: res.data,
        });
      } else {
        this.setState({loading: false});
      }
    });
  };

  addAction() {
    const newAction = this.newAction();
    ActionBackend.addAction(newAction).then((res) => {
      if (res.status === "error") {
        Setting.showMessage("error", `Failed to add: ${res.msg}`);
      } else {
        Setting.showMessage("success", "Action added successfully");
        this.setState({
          data: Setting.prependRow(this.state.data, newAction),
        });
        this.fetch();
      }
    });
  }

  deleteAction(i) {
    ActionBackend.deleteAction(this.state.data[i]).then((res) => {
      if (res.status === "error") {
        Setting.showMessage("error", `Failed to delete: ${res.msg}`);
      } else {
        Setting.showMessage("success", "Deleted successfully");
        this.fetch();
      }
    });
  }

  newAction() {
    const randomName = Setting.getRandomName();
    return {
      owner: this.props.account.name,
      name: `action_${randomName}`,
      createdTime: moment().format(),
      type: "CAPTCHA",
      statusCode: 302,
      immunityTimes: 30,
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
        render: (text, action, index) => {
          return <a href={`/actions/${action.owner}/${text}`}>{text}</a>;
        },
      },
      {
        title: i18next.t("general:Create time"),
        dataIndex: "createdTime",
        key: "createdTime",
        width: "200px",
        sorter: (a, b) => a.createdTime.localeCompare(b.createdTime),
        render: (text, action, index) => {
          return Setting.getFormattedDate(text);
        },
      },
      {
        title: i18next.t("rule:Type"),
        dataIndex: "type",
        key: "type",
        sorter: (a, b) => a.type.localeCompare(b.type),
        render: (text, action, index) => {
          return (
            <Tag color="blue">
              {i18next.t(`action:${text}`)}
            </Tag>
          );
        },
      },
      {
        title: i18next.t("general:Action"),
        dataIndex: "",
        key: "op",
        render: (text, action, index) => {
          return (
            <div>
              <Popconfirm
                title={`Sure to delete action: ${action.name} ?`}
                onConfirm={() => this.deleteAction(index)}
              >
                <Button style={{marginTop: "10px", marginBottom: "10px", marginRight: "10px"}} type="primary" onClick={() => this.props.history.push(`/actions/${action.owner}/${action.name}`)}>{i18next.t("general:Edit")}</Button>
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
        loading={this.state.loading}
        onChange={this.handleTableChange}
        size="middle"
        bordered
        title={() => (
          <div>
            {i18next.t("general:Actions")}&nbsp;&nbsp;&nbsp;&nbsp;
            <Button type="primary" size="small" onClick={() => this.addAction()}>{i18next.t("general:Add")}</Button>
          </div>
        )}
      />
    );
  }
}

export default ActionListPage;
