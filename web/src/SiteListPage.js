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
import * as SiteBackend from "./backend/SiteBackend";
import i18next from "i18next";

class SiteListPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      sites: null,
    };
  }

  UNSAFE_componentWillMount() {
    this.getSites();
  }

  getSites() {
    SiteBackend.getSites(this.props.account.name)
      .then((res) => {
        this.setState({
          sites: res,
        });
      });
  }

  newSite() {
    const randomName = Setting.getRandomName();
    return {
      owner: this.props.account.name,
      name: `site_${randomName}`,
      createdTime: moment().format(),
      displayName: `New Site - ${randomName}`,
      domain: "door.casdoor.com",
      host: "http://localhost:8000",
      sslMode: "HTTP",
      sslCert: "cert_casdoor.com",
    };
  }

  addSite() {
    const newSite = this.newSite();
    SiteBackend.addSite(newSite)
      .then((res) => {
        Setting.showMessage("success", "Site added successfully");
        this.setState({
          sites: Setting.prependRow(this.state.sites, newSite),
        });
      }
      )
      .catch(error => {
        Setting.showMessage("error", `Site failed to add: ${error}`);
      });
  }

  deleteSite(i) {
    SiteBackend.deleteSite(this.state.sites[i])
      .then((res) => {
        Setting.showMessage("success", "Site deleted successfully");
        this.setState({
          sites: Setting.deleteRow(this.state.sites, i),
        });
      }
      )
      .catch(error => {
        Setting.showMessage("error", `Site failed to delete: ${error}`);
      });
  }

  renderTable(sites) {
    const columns = [
      {
        title: i18next.t("general:Name"),
        dataIndex: "name",
        key: "name",
        width: "120px",
        sorter: (a, b) => a.name.localeCompare(b.name),
        render: (text, record, index) => {
          return (
            <Link to={`/sites/${text}`}>
              {text}
            </Link>
          );
        },
      },
      // {
      //   title: i18next.t("general:Create time"),
      //   dataIndex: "createdTime",
      //   key: "createdTime",
      //   width: "180px",
      //   sorter: (a, b) => a.createdTime.localeCompare(b.createdTime),
      //   render: (text, record, index) => {
      //     return Setting.getFormattedDate(text);
      //   },
      // },
      {
        title: i18next.t("general:Display name"),
        dataIndex: "displayName",
        key: "displayName",
        // width: "200px",
        sorter: (a, b) => a.displayName.localeCompare(b.displayName),
      },
      {
        title: i18next.t("site:Domain"),
        dataIndex: "domain",
        key: "domain",
        width: "150px",
        sorter: (a, b) => a.domain.localeCompare(b.domain),
        render: (text, record, index) => {
          return (
            <a target="_blank" rel="noreferrer" href={`https://${text}`}>
              {text}
            </a>
          );
        },
      },
      {
        title: i18next.t("site:Host"),
        dataIndex: "host",
        key: "host",
        width: "180px",
        sorter: (a, b) => a.host.localeCompare(b.host),
        render: (text, record, index) => {
          return (
            <a target="_blank" rel="noreferrer" href={text}>
              {text}
            </a>
          );
        },
      },
      {
        title: i18next.t("site:Public IP"),
        dataIndex: "publicIp",
        key: "publicIp",
        width: "150px",
        sorter: (a, b) => a.publicIp.localeCompare(b.publicIp),
      },
      {
        title: i18next.t("site:Node"),
        dataIndex: "node",
        key: "node",
        width: "150px",
        sorter: (a, b) => a.node.localeCompare(b.node),
      },
      {
        title: i18next.t("site:SSL mode"),
        dataIndex: "sslMode",
        key: "sslMode",
        width: "150px",
        sorter: (a, b) => a.sslMode.localeCompare(b.sslMode),
      },
      {
        title: i18next.t("site:SSL cert"),
        dataIndex: "sslCert",
        key: "sslCert",
        width: "180px",
        sorter: (a, b) => a.sslCert.localeCompare(b.sslCert),
        render: (text, record, index) => {
          return (
            <Link to={`/certs/${text}`}>
              {text}
            </Link>
          );
        },
      },
      {
        title: i18next.t("site:Casdoor endpoint"),
        dataIndex: "casdoorEndpoint",
        key: "casdoorEndpoint",
        width: "200px",
        sorter: (a, b) => a.host.localeCompare(b.host),
        render: (text, record, index) => {
          return (
            <a target="_blank" rel="noreferrer" href={text}>
              {text}
            </a>
          );
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
              <Button style={{marginTop: "10px", marginBottom: "10px", marginRight: "10px"}} type="primary" onClick={() => this.props.history.push(`/sites/${record.name}`)}>{i18next.t("general:Edit")}</Button>
              <Popconfirm
                title={`Sure to delete site: ${record.name} ?`}
                onConfirm={() => this.deleteSite(index)}
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
        <Table columns={columns} dataSource={sites} rowKey="name" size="middle" bordered pagination={{pageSize: 100}}
          title={() => (
            <div>
              {i18next.t("general:Sites")}&nbsp;&nbsp;&nbsp;&nbsp;
              <Button type="primary" size="small" onClick={this.addSite.bind(this)}>{i18next.t("general:Add")}</Button>
            </div>
          )}
          loading={sites === null}
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
              this.renderTable(this.state.sites)
            }
          </Col>
          <Col span={1}>
          </Col>
        </Row>
      </div>
    );
  }
}

export default SiteListPage;
