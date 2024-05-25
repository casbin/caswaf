import React from "react";
import {Link} from "react-router-dom";
import {Col, Row, Table} from "antd";
import * as Setting from "./Setting";
import * as SiteBackend from "./backend/SiteBackend";
import i18next from "i18next";

class DashboardPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      site: null,
    };
  }

  UNSAFE_componentWillMount() {
    this.getSite();
  }

  getSite() {
    SiteBackend.getSites(this.props.account.name)
      .then((res) => {
        if (res.status === "ok") {
          this.setState({
            site: res.data,
          });
        } else {
          Setting.showMessage("error", `Failed to get site: ${res.msg}`);
        }
      });
  }

  renderTable(site) {
    const columns = [
      {
        title: i18next.t("general:Owner"),
        dataIndex: "owner",
        key: "owner",
        width: "120px",
        sorter: (a, b) => a.owner.localeCompare(b.owner),
      },
      {
        title: i18next.t("general:Name"),
        dataIndex: "name",
        key: "name",
        width: "120px",
        sorter: (a, b) => a.name.localeCompare(b.name),
        render: (text, record, index) => {
          return (
            <Link to={`/dashboard/${record.owner}/${record.name}`}>
              {text}
            </Link>
          );
        },
      },
      {
        title: i18next.t("general:Create time"),
        dataIndex: "createdTime",
        key: "createdTime",
        width: "150px",
        sorter: (a, b) => a.createdTime.localeCompare(b.createdTime),
        render: (text, record, index) => {
          return Setting.getFormattedDate(text);
        },
      },
    ];

    return (
      <div>
        <Table columns={columns} dataSource={site} rowKey="name" size="middle" bordered pagination={{pageSize: 100}}
          title={() => (
            <div>
              {i18next.t("general:Dashboard")}&nbsp;&nbsp;&nbsp;&nbsp;
            </div>
          )}
          loading={site === null}
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
              this.renderTable(this.state.site)
            }
          </Col>
        </Row>
      </div>
    );
  }
}

export default DashboardPage;
