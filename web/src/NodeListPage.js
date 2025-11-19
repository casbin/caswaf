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
import {Button, Popconfirm, Table} from "antd";
import moment from "moment";
import * as Setting from "./Setting";
import * as NodeBackend from "./backend/NodeBackend";
import i18next from "i18next";
import BaseListPage from "./BaseListPage";

class NodeListPage extends BaseListPage {
  UNSAFE_componentWillMount() {
    this.setState({
      pagination: {
        ...this.state.pagination,
        current: 1,
        pageSize: 10,
      },
    });
    this.fetch({pagination: this.state.pagination});
  }

  fetch = (params = {}) => {
    const sortField = params.sortField, sortOrder = params.sortOrder;
    if (!params.pagination) {
      params.pagination = {current: 1, pageSize: 10};
    }
    this.setState({
      loading: true,
    });
    NodeBackend.getNodes(this.props.account.name, params.pagination.current, params.pagination.pageSize, "", "", sortField, sortOrder).then((res) => {
      this.setState({
        loading: false,
      });
      if (res.status === "ok") {
        this.setState({
          data: res.data,
          pagination: {
            ...params.pagination,
            total: res.data2,
          },
        });
      } else {
        this.setState({loading: false});
      }
    });
  };

  addNode() {
    const newNode = this.newNode();
    NodeBackend.addNode(newNode).then((res) => {
      if (res.status === "error") {
        Setting.showMessage("error", `Failed to add: ${res.msg}`);
      } else {
        Setting.showMessage("success", "Node added successfully");
        this.setState({
          data: Setting.prependRow(this.state.data, newNode),
        });
        this.fetch();
      }
    });
  }

  deleteNode(i) {
    NodeBackend.deleteNode(this.state.data[i]).then((res) => {
      if (res.status === "error") {
        Setting.showMessage("error", `Failed to delete: ${res.msg}`);
      } else {
        Setting.showMessage("success", "Deleted successfully");
        this.fetch({
          pagination: {
            ...this.state.pagination,
            current: this.state.pagination.current > 1 && this.state.data.length === 1 ? this.state.pagination.current - 1 : this.state.pagination.current,
          },
        });
      }
    });
  }

  newNode() {
    const randomName = Setting.getRandomName();
    return {
      owner: this.props.account.name,
      name: `node_${randomName}`,
      createdTime: moment().format(),
      displayName: `Node ${randomName}`,
      tag: "",
      hostname: "",
      ipAddress: "",
      description: "",
    };
  }

  renderTable(nodes) {
    const columns = [
      {
        title: i18next.t("general:Name"),
        dataIndex: "name",
        key: "name",
        width: "120px",
        sorter: (a, b) => a.name.localeCompare(b.name),
        render: (text, record, index) => {
          return (
            <a href={`/nodes/${record.owner}/${record.name}`}>{text}</a>
          );
        },
      },
      {
        title: i18next.t("general:Display name"),
        dataIndex: "displayName",
        key: "displayName",
        width: "150px",
        sorter: (a, b) => a.displayName.localeCompare(b.displayName),
      },
      {
        title: i18next.t("general:Tag"),
        dataIndex: "tag",
        key: "tag",
        width: "120px",
        sorter: (a, b) => a.tag.localeCompare(b.tag),
      },
      {
        title: i18next.t("node:Hostname"),
        dataIndex: "hostname",
        key: "hostname",
        width: "150px",
        sorter: (a, b) => a.hostname.localeCompare(b.hostname),
      },
      {
        title: i18next.t("node:IP address"),
        dataIndex: "ipAddress",
        key: "ipAddress",
        width: "130px",
      },
      {
        title: i18next.t("node:Description"),
        dataIndex: "description",
        key: "description",
      },
      {
        title: i18next.t("general:Created time"),
        dataIndex: "createdTime",
        key: "createdTime",
        width: "160px",
        sorter: (a, b) => a.createdTime.localeCompare(b.createdTime),
        render: (text, record, index) => {
          return Setting.getFormattedDate(text);
        },
      },
      {
        title: i18next.t("general:Action"),
        dataIndex: "",
        key: "op",
        width: "170px",
        fixed: (Setting.isMobile()) ? "false" : "right",
        render: (text, record, index) => {
          return (
            <div>
              <Button style={{marginTop: "10px", marginBottom: "10px", marginRight: "10px"}} type="primary" onClick={() => this.props.history.push(`/nodes/${record.owner}/${record.name}`)}>{i18next.t("general:Edit")}</Button>
              <Popconfirm
                title={`Sure to delete: ${record.name} ?`}
                onConfirm={() => this.deleteNode(index)}
                disabled={!Setting.isLocalAdminUser(this.props.account)}
              >
                <Button style={{marginBottom: "10px"}} type="danger" disabled={!Setting.isLocalAdminUser(this.props.account)}>{i18next.t("general:Delete")}</Button>
              </Popconfirm>
            </div>
          );
        },
      },
    ];

    return (
      <div>
        <Table columns={columns} dataSource={nodes} rowKey="name" size="middle" bordered pagination={this.state.pagination} loading={this.state.loading}
          onChange={this.fetch}
          title={() => (
            <div>
              {i18next.t("general:Nodes")}&nbsp;&nbsp;&nbsp;&nbsp;
              <Button type="primary" size="small" disabled={!Setting.isLocalAdminUser(this.props.account)} onClick={this.addNode.bind(this)}>{i18next.t("general:Add")}</Button>
            </div>
          )}
        />
      </div>
    );
  }

  render() {
    return (
      <div>
        {
          this.renderTable(this.state.data)
        }
      </div>
    );
  }
}

export default NodeListPage;
