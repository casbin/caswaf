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
import {CheckCircleOutlined, DeleteOutlined, DownOutlined, MinusCircleOutlined, SyncOutlined, UpOutlined} from "@ant-design/icons";
import {Button, Col, Input, Row, Table, Tag, Tooltip} from "antd";
import * as Setting from "./Setting";

const {TextArea} = Input;

class NodeTable extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
    };
  }

  updateTable(table) {
    this.props.onUpdateTable(table);
  }

  parseField(key, value) {
    if (["no", "port", "processId"].includes(key)) {
      value = Setting.myParseInt(value);
    }
    return value;
  }

  updateField(table, index, key, value) {
    value = this.parseField(key, value);

    table[index][key] = value;
    this.updateTable(table);
  }

  addRow(table) {
    const row = {name: `New Node - ${table.length}`, version: "", diff: "", status: "", message: ""};
    if (table === undefined) {
      table = [];
    }

    if (table.length === 0) {
      row.name = this.props.account.hostname;
    }

    table = Setting.addRow(table, row);
    this.updateTable(table);
  }

  deleteRow(table, i) {
    table = Setting.deleteRow(table, i);
    this.updateTable(table);
  }

  upRow(table, i) {
    table = Setting.swapRow(table, i - 1, i);
    this.updateTable(table);
  }

  downRow(table, i) {
    table = Setting.swapRow(table, i, i + 1);
    this.updateTable(table);
  }

  renderTable(table) {
    const columns = [
      {
        title: "Name",
        dataIndex: "name",
        key: "name",
        width: "180px",
        render: (text, record, index) => {
          return (
            <Input value={text} onChange={e => {
              this.updateField(table, index, "name", e.target.value);
            }} />
          );
        },
      },
      {
        title: "Version",
        dataIndex: "version",
        key: "version",
        width: "300px",
        render: (text, record, index) => {
          if (text === "") {
            return null;
          }

          const versionInfo = JSON.parse(text);
          const link = versionInfo?.version !== "" ? `${Setting.getRepoUrl(this.props.siteName)}/releases/tag/${versionInfo?.version}` : "";
          let versionText = versionInfo?.version !== "" ? versionInfo?.version : "Unknown version";
          if (versionInfo?.commitOffset > 0) {
            versionText += ` (ahead+${versionInfo?.commitOffset})`;
          }

          return (
            <a target="_blank" rel="noreferrer" href={link}>{versionText}</a>
          );

          // return (
          //   <Input value={text} onChange={e => {
          //     this.updateField(table, index, "version", e.target.value);
          //   }} />
          // );
        },
      },
      {
        title: "Diff",
        dataIndex: "diff",
        key: "diff",
        render: (text, record, index) => {
          if (record.status === "") {
            return null;
          } else {
            return (
              <Tooltip title={
                <div style={{width: "800px"}}>
                  <TextArea autoSize={{minRows: 1, maxRows: 30}} value={text} />
                </div>
              }>
                {Setting.getShortText(text)}
              </Tooltip>
            );
          }
        },
      },
      {
        title: "Status",
        dataIndex: "status",
        key: "status",
        width: "150px",
        render: (text, record, index) => {
          if (record.status === "") {
            return null;
          } else if (record.status === "In Progress") {
            return (
              <Tag icon={<SyncOutlined spin />} color="processing">{text}</Tag>
            );
          } else if (record.status === "Running") {
            return (
              <Tag icon={<CheckCircleOutlined />} color="success">{text}</Tag>
            );
          } else if (record.status === "Stopped") {
            return (
              <Tag icon={<MinusCircleOutlined />} color="error">{text}</Tag>
            );
          } else {
            return text;
          }
        },
      },
      {
        title: "Message",
        dataIndex: "message",
        key: "message",
      },
      {
        title: "Action",
        key: "action",
        width: "100px",
        render: (text, record, index) => {
          return (
            <div>
              <Tooltip placement="bottomLeft" title={"Up"}>
                <Button style={{marginRight: "5px"}} disabled={index === 0} icon={<UpOutlined />} size="small" onClick={() => this.upRow(table, index)} />
              </Tooltip>
              <Tooltip placement="topLeft" title={"Down"}>
                <Button style={{marginRight: "5px"}} disabled={index === table.length - 1} icon={<DownOutlined />} size="small" onClick={() => this.downRow(table, index)} />
              </Tooltip>
              <Tooltip placement="topLeft" title={"Delete"}>
                <Button icon={<DeleteOutlined />} size="small" onClick={() => this.deleteRow(table, index)} />
              </Tooltip>
            </div>
          );
        },
      },
    ];

    return (
      <Table rowKey="index" columns={columns} dataSource={table} size="middle" bordered pagination={false}
        title={() => (
          <div>
            {this.props.title}&nbsp;&nbsp;&nbsp;&nbsp;
            <Button style={{marginRight: "5px"}} type="primary" size="small" onClick={() => this.addRow(table)}>{"Add"}</Button>
          </div>
        )}
      />
    );
  }

  render() {
    return (
      <div>
        <Row style={{marginTop: "20px"}} >
          <Col span={24}>
            {
              this.renderTable(this.props.table)
            }
          </Col>
        </Row>
      </div>
    );
  }
}

export default NodeTable;
