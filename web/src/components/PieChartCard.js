import React from "react";
import {Card} from "antd";
import ReactECharts from "echarts-for-react";
import i18next from "i18next";

const PieChartCard = ({title, data}) => {
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
        data: data,
      },
    ],
  };

  return (
    <Card title={i18next.t(`general:${title}`)}>
      <ReactECharts option={option}></ReactECharts>
    </Card>
  );
};

export default PieChartCard;
