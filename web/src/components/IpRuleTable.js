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
import {DeleteOutlined, DownOutlined, UpOutlined} from "@ant-design/icons";
import {Button, Col, Input, Row, Select, Table, Tooltip} from "antd";
import * as Setting from "../Setting";

const {Option} = Select;

class IpRuleTable extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      options: [],
    };
    for (let i = 0; i < this.props.table.length; i++) {
      const values = this.props.table[i].value.split(" ");
      const options = [];
      for (let j = 0; j < values.length; j++) {
        options[j] = {value: values[j], label: values[j]};
      }
      this.state.options.push(options);
    }
  }

  updateTable(table) {
    this.props.onUpdateTable(table);
  }

  updateField(table, index, key, value) {
    if (key === "value") {
      let v = "";
      for (let i = 0; i < value.length; i++) {
        v += value[i] + " ";
      }
      table[index][key] = v.trim();
    } else {
      table[index][key] = value;
    }
    this.updateTable(table);
  }

  addRow(table) {
    const row = {name: `New IP Rule - ${table.length}`, operator: "is in", value: "127.0.0.1"};
    if (table === undefined) {
      table = [];
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
    Setting.swapRow(this.state.options, i - 1, i);
    this.updateTable(table);
  }

  downRow(table, i) {
    table = Setting.swapRow(table, i, i + 1);
    Setting.swapRow(this.state.options, i, i + 1);
    this.updateTable(table);
  }

  renderTable(table) {
    const columns = [
      {
        title: "Name",
        dataIndex: "name",
        key: "name",
        width: "180px",
        render: (text, rule, index) => (
          <Input value={text} onChange={e => {
            this.updateField(table, index, "name", e.target.value);
          }} />
        ),
      },
      {
        title: "Operator",
        dataIndex: "operator",
        key: "operator",
        width: "180px",
        render: (text, rule, index) => (
          <Select value={text} virtual={false} style={{width: "100%"}} onChange={value => {
            this.updateField(table, index, "operator", value);
          }}>
            {
              [
                {value: "is in", text: "is in"},
                {value: "is not in", text: "is not in"},
              ].map((item, index) => <Option key={index} value={item.value}>{item.text}</Option>)
            }
          </Select>
        ),
      },
      {
        title: "Value",
        dataIndex: "value",
        key: "value",
        width: "100%",
        render: (text, rule, index) => (
          <Select
            mode="tags"
            style={{width: "100%"}}
            placeholder="Input IP Addresses"
            value={rule.value.split(" ")}
            onChange={value => this.updateField(table, index, "value", value)}
            options={this.state.options[index]}
          />
        ),
      },
      {
        title: "Action",
        key: "action",
        width: "100px",
        render: (text, rule, index) => (
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
        ),
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

export default IpRuleTable;
