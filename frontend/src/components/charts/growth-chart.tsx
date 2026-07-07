"use client"

import { useRef, useEffect } from "react"
import * as echarts from "echarts/core"
import { LineChart } from "echarts/charts"
import { GridComponent, TooltipComponent, TitleComponent } from "echarts/components"
import { CanvasRenderer } from "echarts/renderers"

echarts.use([LineChart, GridComponent, TooltipComponent, TitleComponent, CanvasRenderer])

interface DataPoint {
  date: string
  score: number
}

interface GrowthChartProps {
  data: DataPoint[]
  color?: string
  loading?: boolean
}

export function GrowthChart({ data, color = "#00ff41", loading }: GrowthChartProps) {
  const chartRef = useRef<HTMLDivElement>(null)
  const instanceRef = useRef<echarts.ECharts | null>(null)

  useEffect(() => {
    if (!chartRef.current || loading) return

    if (!instanceRef.current) {
      instanceRef.current = echarts.init(chartRef.current, undefined, {
        renderer: "canvas",
      })
    }

    const isDark = document.documentElement.classList.contains("dark")
    const textColor = isDark ? "#888899" : "#6b7280"
    const borderColor = isDark ? "#1f1f23" : "#e5e7eb"

    instanceRef.current.setOption({
      tooltip: {
        trigger: "axis",
        backgroundColor: isDark ? "#111114" : "#fff",
        borderColor,
        textStyle: { color: isDark ? "#ededed" : "#111" },
      },
      grid: { left: 40, right: 20, top: 20, bottom: 30 },
      xAxis: {
        type: "category",
        data: data.map((d) => d.date),
        axisLine: { lineStyle: { color: borderColor } },
        axisLabel: { color: textColor, fontSize: 11 },
      },
      yAxis: {
        type: "value",
        splitLine: { lineStyle: { color: borderColor } },
        axisLabel: { color: textColor, fontSize: 11 },
      },
      series: [
        {
          data: data.map((d) => d.score),
          type: "line",
          smooth: true,
          symbol: "circle",
          symbolSize: 6,
          lineStyle: { color, width: 2 },
          itemStyle: { color },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: color + "33" },
              { offset: 1, color: color + "00" },
            ]),
          },
        },
      ],
    })

    const handler = () => {
      setTimeout(() => instanceRef.current?.resize(), 0)
    }
    window.addEventListener("resize", handler)
    return () => window.removeEventListener("resize", handler)
  }, [data, color, loading])

  useEffect(() => {
    return () => instanceRef.current?.dispose()
  }, [])

  if (loading) {
    return <div className="h-64 glass rounded-xl animate-pulse flex items-center justify-center text-muted-foreground">Loading chart...</div>
  }

  return <div ref={chartRef} className="h-64 w-full" />
}
