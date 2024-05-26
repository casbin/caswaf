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
import ReactECharts from "echarts-for-react";
import {Card, Col, Radio, Row, Table} from "antd";
import * as Setting from "./Setting";
import * as DashboardBackend from "./backend/DashboardBackend";
import i18next from "i18next";

class DashboardDetailPage extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      classes: props,
      site: null,
      userAgents: [{}],
      sites: [{}],
      uniqueIPCount: 0,
      totleRequestCount: 0,
      ipAddresses: [{}],
      requestCountOverTime: [{}],
      rangeType: "All",
    };
  }

  UNSAFE_componentWillMount() {
    this.getAllData(this.state.rangeType);
  }

  async getMetric(type, rangeType, top) {
    rangeType = rangeType === "All" ? "month" : rangeType.toLowerCase();
    const count = this.getRangeValue(rangeType);
    if (type === "UserAgent" || type === "IPAddress") {
      top = 10;
    }
    return DashboardBackend.getMetric(type, rangeType, count, top).then((res) => {
      if (res.status === "ok") {
        return res;
      } else {
        Setting.showMessage("error", res.msg);
      }
    });
  }

  async getMetricOverTime(rangeType) {
    rangeType = rangeType === "All" ? "week" : rangeType.toLowerCase();
    const count = this.getRangeValue(rangeType);
    const timeType = this.getGranularity(rangeType);
    return DashboardBackend.getMetricOverTime(rangeType, count, timeType).then((res) => {
      if (res.status === "ok") {
        return res;
      } else {
        Setting.showMessage("error", res.msg);
      }
    });
  }

  getAllData(rangeType) {
    this.getUserAgents(rangeType);
    this.getIPAddresses(rangeType);
    this.getSites(rangeType);
    this.getRequestCount(rangeType);
  }

  getRangeValue(rangeType) {
    switch (rangeType) {
    case "hour":
      return 72;
    case "day":
      return 7;
    case "week":
      return 12;
    case "month":
      return 12;
    default:
      return 7;
    }
  }

  getGranularity(rangeType) {
    switch (rangeType) {
    case "hour":
      return "hour";
    case "day":
      return "hour";
    case "week":
      return "day";
    case "month":
      return "month";
    default:
      return "day";
    }
  }

  getUserAgents(rangeType) {
    this.getMetric("userAgent", rangeType, 10).then(res => {
      this.setState({
        userAgents: res.data,
      });
    });
  }

  getIPAddresses(rangeType) {
    this.getMetric("ip", rangeType).then((res) => {
      this.setState({
        ipAddresses: res.data.slice(0, 10),
        uniqueIPCount: res.data.length,
      });
    });
  }

  getRequestCount(rangeType) {
    this.getMetricOverTime(rangeType).then((res) => {
      this.setState({
        requestCountOverTime: res.data,
        totleRequestCount: res.data2,
      });
    });
  }

  getSites(rangeType) {
    this.getMetric("site", rangeType).then((res) => {
      this.setState({
        sites: res.data,
      });
    });
  }

  renderUserAgentsTable() {
    const columns = [
      {
        title: i18next.t("general:User-Agent"),
        dataIndex: "data",
        key: "data",
        width: "440px",
      },
      {
        title: i18next.t("general:Count"),
        dataIndex: "count",
        key: "count",
        width: "40px",
        sorter: (a, b) => a.count - b.count,
      },
    ];

    return (
      <Card title={i18next.t("general:Top 10 User-Agents")}>
        <Table
          columns={columns}
          dataSource={this.state.userAgents}
          rowKey="userAgent"
          size="small"
          pagination={{hideOnSinglePage: true}}
        />
      </Card>
    );

  }

  renderIPAddressTable() {
    const columns = [
      {
        title: i18next.t("general:IP Address"),
        dataIndex: "data",
        key: "data",
        width: "140px",
      },
      {
        title: i18next.t("general:Count"),
        dataIndex: "count",
        key: "count",
        width: "20px",
        sorter: (a, b) => a.count - b.count,
      },
    ];

    return (
      <Card title={i18next.t("general:Top 10 IP Addresses")}>
        <Table
          columns={columns}
          dataSource={this.state.ipAddresses}
          rowKey="ipAddress"
          size="small"
          pagination={{hideOnSinglePage: true}}
        />
      </Card>
    );

  }

  renderSitesPieChart() {
    return this.renderPieChart("Sites", this.state.sites);
  }

  renderPieChart(title, data) {
    const d = data.map((item) => {
      return {value: item.count, name: item.data};
    });
    const option = {
      tooltip: {
        trigger: "item",
      },
      legend: {
        top: "5%",
        left: "right",
        orient: "vertical",
      },
      series: [
        {
          name: title,
          type: "pie",
          radius: ["40%", "70%"],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: "#fff",
            borderWidth: 2,
          },
          label: {
            show: false,
            position: "center",
          },
          emphasis: {
            label: {
              show: true,
              fontSize: 30,
              fontWeight: "bold",
            },
          },
          labelLine: {
            show: false,
          },
          data: d,
        },
      ],
    };
    return (
      <Card title={i18next.t(`general:${title}`)}>
        <ReactECharts option={option}></ReactECharts>
      </Card>
    );
  }

  renderTotalRequestCountStatistic() {
    return this.renderStatistic(i18next.t("general:Total Request Count"), this.state.totleRequestCount);
  }

  renderUniqueIPCountStatistic() {
    return this.renderStatistic(i18next.t("general:Unique IP Count"), this.state.uniqueIPCount);
  }

  renderStatistic(title, value) {
    const option = {
      series: [
        {
          type: "scatter",
          data: [[0, 0]],
          symbolSize: 1,
          label: {
            show: true,
            formatter: [
              value,
            ].join("\n"),
            color: "#000",
            fontSize: 64,
          },
        },
      ],
      xAxis: {
        axisLabel: {show: false},
        axisLine: {show: false},
        splitLine: {show: false},
        axisTick: {show: false},
        min: -1,
        max: 1,
      },
      yAxis: {
        axisLabel: {show: false},
        axisLine: {show: false},
        splitLine: {show: false},
        axisTick: {show: false},
        min: -1,
        max: 1,
      },
    };
    return (
      <Card title={i18next.t(`general:${title}`)}>
        <ReactECharts option={option}></ReactECharts>
      </Card>
    );
  }

  renderBarChart(title, data) {
    const option = {
      tooltip: {
        trigger: "axis",
        axisPointer: {
          type: "shadow",
        },
      },
      grid: {
        left: "3%",
        right: "4%",
        bottom: "3%",
        containLabel: true,
      },
      xAxis: {
        type: "category",
        data: data.map((item) => item.data),
        axisTick: {
          alignWithLabel: true,
        },
      },
      yAxis: {
        type: "value",
      },
      series: [
        {
          name: title,
          type: "bar",
          barWidth: "60%",
          data: data.map((item) => item.count),
        },
      ],
    };
    return (
      <Card title={i18next.t(`general:${title}`)}>
        <ReactECharts option={option}></ReactECharts>
      </Card>
    );
  }

  renderRadio() {
    return (
      <div style={{margin: "10px", float: "right"}}>
        <Radio.Group style={{marginBottom: "10px"}} buttonStyle="solid" value={this.state.rangeType} onChange={e => {
          const rangeType = e.target.value;
          this.getAllData(rangeType);
          this.setState({
            rangeType: rangeType,
          });
        }}>
          <Radio.Button value={"All"}>{i18next.t("usage:All")}</Radio.Button>
          <Radio.Button value={"Hour"}>{i18next.t("usage:Hour")}</Radio.Button>
          <Radio.Button value={"Day"}>{i18next.t("usage:Day")}</Radio.Button>
          <Radio.Button value={"Week"}>{i18next.t("usage:Week")}</Radio.Button>
          <Radio.Button value={"Month"}>{i18next.t("usage:Month")}</Radio.Button>
        </Radio.Group>
      </div>
    );
  }

  render() {
    return (
      <div>
        {this.renderRadio()}
        <Row style={{width: "100%"}}>
          <Col span={4}>
            {
              this.renderTotalRequestCountStatistic()
            }
          </Col>
          <Col span={20}>
            {
              this.renderBarChart("Request Count Over Time", this.state.requestCountOverTime)
            }
          </Col>
        </Row>
        <Row style={{width: "100%"}}>
          <Col span={20}>
            {this.renderSitesPieChart()}
          </Col>
          {/* <Col span={10}>
            {this.renderHTTPVersionPieChart()}
          </Col> */}
          <Col span={4}>
            {
              this.renderUniqueIPCountStatistic()
            }
          </Col>
        </Row>
        <Row>
          <Col span={8}>
            {
              this.renderIPAddressTable()
            }
          </Col>
          <Col span={16}>
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
