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
