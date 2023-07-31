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
import {Link} from "react-router-dom";
import {Button, Col, Popconfirm, Row, Table} from "antd";
import moment from "moment";
import * as Setting from "./Setting";
import * as CertBackend from "./backend/CertBackend";
import i18next from "i18next";

class CertListPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      certs: null,
    };
  }

  UNSAFE_componentWillMount() {
    this.getCerts();
  }

  getCerts() {
    CertBackend.getCerts(this.props.account.name)
      .then((res) => {
        this.setState({
          certs: res,
        });
      });
  }

  newCert() {
    const randomName = Setting.getRandomName();
    return {
      owner: this.props.account.name,
      name: `cert_${randomName}`,
      createdTime: moment().format(),
      displayName: `New Cert - ${randomName}`,
      type: "SSL",
      cryptoAlgorithm: "RSA",
      expireTime: "",
      certificate: "",
      privateKey: "",
    };
  }

  addCert() {
    const newCert = this.newCert();
    CertBackend.addCert(newCert)
      .then((res) => {
        Setting.showMessage("success", "Cert added successfully");
        this.setState({
          certs: Setting.prependRow(this.state.certs, newCert),
        });
      }
      )
      .catch(error => {
        Setting.showMessage("error", `Cert failed to add: ${error}`);
      });
  }

  deleteCert(i) {
    CertBackend.deleteCert(this.state.certs[i])
      .then((res) => {
        Setting.showMessage("success", "Cert deleted successfully");
        this.setState({
          certs: Setting.deleteRow(this.state.certs, i),
        });
      }
      )
      .catch(error => {
        Setting.showMessage("error", `Cert failed to delete: ${error}`);
      });
  }

  renderTable(certs) {
    const columns = [
      {
        title: i18next.t("general:Name"),
        dataIndex: "name",
        key: "name",
        width: "120px",
        sorter: (a, b) => a.name.localeCompare(b.name),
        render: (text, record, index) => {
          return (
            <Link to={`/certs/${text}`}>
              {text}
            </Link>
          );
        },
      },
      {
        title: i18next.t("general:Create time"),
        dataIndex: "createdTime",
        key: "createdTime",
        width: "180px",
        sorter: (a, b) => a.createdTime.localeCompare(b.createdTime),
        render: (text, record, index) => {
          return Setting.getFormattedDate(text);
        },
      },
      {
        title: i18next.t("general:Display name"),
        dataIndex: "displayName",
        key: "displayName",
        // width: "200px",
        sorter: (a, b) => a.displayName.localeCompare(b.displayName),
      },
      {
        title: i18next.t("cert:Type"),
        dataIndex: "type",
        key: "type",
        width: "180px",
        sorter: (a, b) => a.type.localeCompare(b.type),
      },
      {
        title: i18next.t("cert:Crypto algorithm"),
        dataIndex: "cryptoAlgorithm",
        key: "cryptoAlgorithm",
        width: "180px",
        sorter: (a, b) => a.cryptoAlgorithm.localeCompare(b.cryptoAlgorithm),
      },
      {
        title: i18next.t("cert:Expire time"),
        dataIndex: "expireTime",
        key: "expireTime",
        width: "180px",
        sorter: (a, b) => a.expireTime.localeCompare(b.expireTime),
        render: (text, record, index) => {
          return Setting.getFormattedDate(text);
        },
      },
      {
        title: i18next.t("cert:Certificate"),
        dataIndex: "certificate",
        key: "certificate",
        width: "180px",
        sorter: (a, b) => a.certificate.localeCompare(b.certificate),
        render: (text, record, index) => {
          return Setting.getShortText(text);
        },
      },
      {
        title: i18next.t("cert:Private key"),
        dataIndex: "privateKey",
        key: "privateKey",
        width: "180px",
        sorter: (a, b) => a.privateKey.localeCompare(b.privateKey),
        render: (text, record, index) => {
          return Setting.getShortText(text);
        },
      },
      {
        title: i18next.t("general:Action"),
        dataIndex: "action",
        key: "action",
        width: "180px",
        render: (text, record, index) => {
          return (
            <div>
              <Button style={{marginTop: "10px", marginBottom: "10px", marginRight: "10px"}} type="primary" onClick={() => this.props.history.push(`/certs/${record.name}`)}>{i18next.t("general:Edit")}</Button>
              <Popconfirm
                title={`Sure to delete cert: ${record.name} ?`}
                onConfirm={() => this.deleteCert(index)}
                okText="OK"
                cancelText="Cancel"
              >
                <Button style={{marginBottom: "10px"}} type="danger">{i18next.t("general:Delete")}</Button>
              </Popconfirm>
            </div>
          );
        },
      },
    ];

    return (
      <div>
        <Table columns={columns} dataSource={certs} rowKey="name" size="middle" bordered pagination={{pageSize: 100}}
          title={() => (
            <div>
              {i18next.t("general:Certs")}&nbsp;&nbsp;&nbsp;&nbsp;
              <Button type="primary" size="small" onClick={this.addCert.bind(this)}>{i18next.t("general:Add")}</Button>
            </div>
          )}
          loading={certs === null}
        />
      </div>
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
              this.renderTable(this.state.certs)
            }
          </Col>
          <Col span={1}>
          </Col>
        </Row>
      </div>
    );
  }
}

export default CertListPage;
