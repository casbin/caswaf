import React from "react";
import {Col, Row, Table} from "antd";
import i18next from "i18next";

class DashboardDetailPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      owner: props.match.params.owner,
      siteName: props.match.params.siteName,
      site: null,
      userAgents: [{}],
    };
  }

  UNSAFE_componentWillMount() {
    this.getUserAgents();
  }

  getUserAgents() {
    const userAgents = [
      {
        userAgent: "Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Mobile Safari/537.36",
        count: "1",
      },
      {
        userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246",
        count: "35",
      },
      {
        userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9",
        count: "13",
      },
      {
        userAgent: "Mozilla/5.0 (Windows NT 6.1; Trident/7.0; rv:11.0) like Gecko",
        count: "3",
      },
    ];
    this.setState({
      userAgents: userAgents,
    });
  }

  renderUserAgentsTable() {
    const columns = [
      {
        title: i18next.t("general:User Agent"),
        dataIndex: "userAgent",
        key: "userAgent",
        width: "140px",
      },
      {
        title: i18next.t("general:Count"),
        dataIndex: "count",
        key: "count",
        width: "20px",
      },
    ];

    return (
      <div>
        <Table
          columns={columns}
          dataSource={this.state.userAgents}
          rowKey="userAgent"
          size="small"
          bordered
          pagination={{hideOnSinglePage: true}}
        />
      </div>
    );

  }

  render() {
    return (
      <div>
        <h2>Dashboard Detail Page</h2>
        <Row style={{width: "100%"}}>
          <Col span={4}>
            {
              this.renderUserAgentsTable()
            }
          </Col>
        </Row>
      </div>
    );
  }
}

export default DashboardDetailPage;
