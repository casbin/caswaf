import React from "react";
import {Card} from "antd";
import ReactECharts from "echarts-for-react";
import i18next from "i18next";

const BarChartCard = ({title, data}) => {
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
};

export default BarChartCard;
