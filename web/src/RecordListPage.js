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
import {Button, Col, Popconfirm, Row, Table} from "antd";
import * as Setting from "./Setting";
import * as RecordBackend from "./backend/RecordBackend";
import i18next from "i18next";

class RecordListPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      records: null,
    };
  }

  UNSAFE_componentWillMount() {
    this.getRecords();
  }

  getRecords() {
    RecordBackend.getRecords(this.props.account.name)
      .then((res) => {
        if (res.status === "ok") {
          this.setState({
            records: res.data,
          });
        } else {
          Setting.showMessage("error", `Failed to get records: ${res.msg}`);
        }
      });
  }

  deleteRecord(i) {
    RecordBackend.deleteRecord(this.state.records[i])
      .then((res) => {
        if (res.status === "error") {
          Setting.showMessage("error", `Failed to delete: ${res.msg}`);
        } else {
          Setting.showMessage("success", "Record deleted successfully");
          this.setState({
            records: Setting.deleteRow(this.state.records, i),
          });
        }
      }
      )
      .catch(error => {
        Setting.showMessage("error", `Record failed to delete: ${error}`);
      });
  }

  renderTable(records) {

    const columns = [
      {
        title: i18next.t("general:Id"),
        dataIndex: "id",
        key: "id",
        width: "40px",
        sorter: (a, b) => a.id - b.id,
      },
      {
        title: i18next.t("general:Owner"),
        dataIndex: "owner",
        key: "owner",
        width: "40px",
        sorter: (a, b) => a.owner.localeCompare(b.owner),
      },
      {
        title: i18next.t("general:CreatedTime"),
        dataIndex: "createdTime",
        key: "createdTime",
        width: "70px",
        sorter: (a, b) => a.createdTime.localeCompare(b.createdTime),
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
        title: i18next.t("general:RequestURI"),
        dataIndex: "requestURI",
        key: "requestURI",
        width: "100px",
        sorter: (a, b) => a.requestURI.localeCompare(b.requestURI),
      },
      {
        title: i18next.t("general:UserAgent"),
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
        <Table columns={columns} dataSource={records} rowKey="name" size="middle" bordered pagination={{pageSize: 1000}}
          title={() => (
            <div>
              {i18next.t("general:Records")}&nbsp;&nbsp;&nbsp;&nbsp;
            </div>
          )}
          loading={records === null}
        />
      </div>
    );
  }

  render() {
    return (
      <div>
        <Row style={{width: "100%"}}>
          <Col span={24}>
            {
              this.renderTable(this.state.records)
            }
          </Col>
        </Row>
      </div>
    );
  }
}

export default RecordListPage;
