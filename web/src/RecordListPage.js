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
import * as RecordBackend from "./backend/RecordBackend";
import i18next from "i18next";
import BaseListPage from "./BaseListPage";

class RecordListPage extends BaseListPage {

  UNSAFE_componentWillMount() {
    this.fetch();
  }

  fetch = (params = {}) => {
    const sortField = params.sortField, sortOrder = params.sortOrder;
    if (!params.pagination) {
      params.pagination = {current: 1, pageSize: 10};
    }
    this.setState({loading: true});
    RecordBackend.getRecords(this.props.account.name, params.pagination.current, params.pagination.pageSize, sortField, sortOrder)
      .then((res) => {
        this.setState({
          loading: false,
        });
        if (res.status === "ok") {
          this.setState({
            data: res.data,
          });
        } else {
          Setting.showMessage("error", `Failed to get records: ${res.msg}`);
        }
      });
  };

  newRecord() {
    const randomName = Setting.getRandomName();
    return {
      owner: this.props.account.name,
      name: `record_${randomName}`,
      createdTime: moment().format(),
      method: "GET",
      host: "door.casdoor.com",
      path: "/",
      userAgent: "",
    };
  }

  addRecord() {
    const newRecord = this.newRecord();
    RecordBackend.addRecord(newRecord)
      .then((res) => {
        if (res.status === "error") {
          Setting.showMessage("error", `Failed to add: ${res.msg}`);
        } else {
          Setting.showMessage("success", "Record added successfully");
          this.setState({
            data: Setting.addRow(this.state.data, res.data),
          });
          this.fetch();
        }
      }
      )
      .catch(error => {
        Setting.showMessage("error", `Record failed to add: ${error}`);
      });
  }

  deleteRecord(i) {
    RecordBackend.deleteRecord(this.state.data[i])
      .then((res) => {
        if (res.status === "error") {
          Setting.showMessage("error", `Failed to delete: ${res.msg}`);
        } else {
          Setting.showMessage("success", "Record deleted successfully");
          this.setState({
            data: Setting.deleteRow(this.state.data, i),
          });
        }
      }
      )
      .catch(error => {
        Setting.showMessage("error", `Record failed to delete: ${error}`);
      });
  }

  renderTable(data) {
    const columns = [
      {
        title: i18next.t("general:ID"),
        dataIndex: "id",
        key: "id",
        width: "30px",
        sorter: (a, b) => a.id - b.id,
      },
      {
        title: i18next.t("general:Owner"),
        dataIndex: "owner",
        key: "owner",
        width: "30px",
        sorter: (a, b) => a.owner.localeCompare(b.owner),
      },
      {
        title: i18next.t("general:Created time"),
        dataIndex: "createdTime",
        key: "createdTime",
        width: "70px",
        sorter: (a, b) => a.createdTime.localeCompare(b.createdTime),
        render: (text, record, index) => {
          return Setting.getFormattedDate(text);
        },
      },
      {
        title: i18next.t("general:Method"),
        dataIndex: "method",
        key: "method",
        width: "30px",
        sorter: (a, b) => a.method.localeCompare(b.method),
      },
      {
        title: i18next.t("general:Host"),
        dataIndex: "host",
        key: "host",
        width: "50px",
        sorter: (a, b) => a.host.localeCompare(b.host),
      },
      {
        title: i18next.t("general:Path"),
        dataIndex: "path",
        key: "path",
        width: "100px",
        sorter: (a, b) => a.path.localeCompare(b.path),
      },
      {
        title: i18next.t("general:Client ip"),
        dataIndex: "clientIp",
        key: "clientIp",
        width: "100px",
        sorter: (a, b) => a.clientIp.localeCompare(b.clientIp),
      },
      {
        title: i18next.t("general:User-Agent"),
        dataIndex: "userAgent",
        key: "userAgent",
        width: "240px",
        sorter: (a, b) => a.userAgent.localeCompare(b.userAgent),
      },
      {
        title: i18next.t("general:Action"),
        dataIndex: "action",
        key: "action",
        width: "180px",
        render: (text, record, index) => {
          return (
            <div>
              <Button style={{marginTop: "10px", marginBottom: "10px", marginRight: "10px"}} type="primary" onClick={() => this.props.history.push(`/records/${record.owner}/${record.id}`)}>{i18next.t("general:Edit")}</Button>
              <Popconfirm
                title={"Sure to delete record?"}
                onConfirm={() => this.deleteRecord(index)}
                okText="OK"
                cancelText="Cancel"
              >
                <Button style={{marginTop: "10px", marginBottom: "10px", marginRight: "10px"}} type="danger">{i18next.t("general:Delete")}</Button>
              </Popconfirm>
            </div>
          );
        },
      },
    ];

    return (
      <div>
        <Table columns={columns} dataSource={data} rowKey="name" size="middle" bordered pagination={{pageSize: 1000}}
          title={() => (
            <div>
              {i18next.t("general:Records")}&nbsp;&nbsp;&nbsp;&nbsp;
              <Button type="primary" size="small" onClick={this.addRecord.bind(this)}>{i18next.t("general:Add")}</Button>
            </div>
          )}
          loading={this.state.loading}
          onChange={this.handleTableChange}
        />
      </div>
    );
  }

}

export default RecordListPage;
