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
import {Card} from "antd";
import ReactECharts from "echarts-for-react";
import i18next from "i18next";

const StatisticCard = ({title, value}) => {
  const option = {
    series: [
      {
        type: "scatter",
        data: [[0, 0]],
        symbolSize: 1,
        label: {
          show: true,
          formatter: [value].join("\n"),
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
};

export default StatisticCard;
