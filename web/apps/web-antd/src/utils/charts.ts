/**
 * 实时数据和图表组件
 */
import { ref, onMounted, onUnmounted } from 'vue';
import * as echarts from 'echarts';

// 图表主题配置
export const chartTheme = {
  color: [
    '#1890ff', '#52c41a', '#faad14', '#f5222d', '#722ed1',
    '#fa541c', '#13c2c2', '#eb2f96', '#a0d911', '#2f54eb'
  ],
  backgroundColor: 'rgba(0,0,0,0)',
  textStyle: {},
  title: {
    textStyle: {
      color: '#262626'
    },
    subtextStyle: {
      color: '#666'
    }
  },
  line: {
    itemStyle: {
      borderWidth: 1
    },
    lineStyle: {
      width: 2
    },
    symbolSize: 4,
    symbol: 'emptyCircle',
    smooth: false
  },
  radar: {
    itemStyle: {
      borderWidth: 1
    },
    lineStyle: {
      width: 2
    },
    symbolSize: 4,
    symbol: 'emptyCircle',
    smooth: false
  },
  bar: {
    itemStyle: {
      barBorderWidth: 0,
      barBorderColor: '#ccc'
    }
  },
  pie: {
    itemStyle: {
      borderWidth: 0,
      borderColor: '#ccc'
    }
  },
  scatter: {
    itemStyle: {
      borderWidth: 0,
      borderColor: '#ccc'
    }
  },
  boxplot: {
    itemStyle: {
      borderWidth: 0,
      borderColor: '#ccc'
    }
  },
  parallel: {
    itemStyle: {
      borderWidth: 0,
      borderColor: '#ccc'
    }
  },
  sankey: {
    itemStyle: {
      borderWidth: 0,
      borderColor: '#ccc'
    }
  },
  funnel: {
    itemStyle: {
      borderWidth: 0,
      borderColor: '#ccc'
    }
  },
  gauge: {
    itemStyle: {
      borderWidth: 0,
      borderColor: '#ccc'
    }
  },
  candlestick: {
    itemStyle: {
      color: '#52c41a',
      color0: '#ff4d4f',
      borderColor: undefined,
      borderColor0: undefined
    }
  },
  graph: {
    itemStyle: {
      borderWidth: 0,
      borderColor: '#ccc'
    },
    lineStyle: {
      width: 1,
      color: '#aaa'
    },
    symbolSize: 4,
    symbol: 'emptyCircle',
    smooth: false,
    color: [
      '#1890ff', '#52c41a', '#faad14', '#f5222d', '#722ed1',
      '#fa541c', '#13c2c2', '#eb2f96', '#a0d911', '#2f54eb'
    ],
    label: {
      color: '#ffffff'
    }
  },
  map: {
    itemStyle: {
      areaColor: '#eee',
      borderColor: '#444',
      borderWidth: 0.5
    },
    label: {
      color: '#000'
    },
    emphasis: {
      itemStyle: {
        areaColor: 'rgba(255,215,0,0.8)',
        borderColor: '#444',
        borderWidth: 1
      },
      label: {
        color: 'rgb(100,0,0)'
      }
    }
  },
  geo: {
    itemStyle: {
      areaColor: '#eee',
      borderColor: '#444',
      borderWidth: 0.5
    },
    label: {
      color: '#000'
    },
    emphasis: {
      itemStyle: {
        areaColor: 'rgba(255,215,0,0.8)',
        borderColor: '#444',
        borderWidth: 1
      },
      label: {
        color: 'rgb(100,0,0)'
      }
    }
  },
  categoryAxis: {
    axisLine: {
      show: true,
      lineStyle: {
        color: '#cccccc'
      }
    },
    axisTick: {
      show: true,
      lineStyle: {
        color: '#cccccc'
      }
    },
    axisLabel: {
      show: true,
      color: '#999999'
    },
    splitLine: {
      show: false,
      lineStyle: {
        color: ['#eeeeee']
      }
    },
    splitArea: {
      show: false,
      areaStyle: {
        color: ['rgba(250,250,250,0.3)', 'rgba(200,200,200,0.3)']
      }
    }
  },
  valueAxis: {
    axisLine: {
      show: true,
      lineStyle: {
        color: '#cccccc'
      }
    },
    axisTick: {
      show: true,
      lineStyle: {
        color: '#cccccc'
      }
    },
    axisLabel: {
      show: true,
      color: '#999999'
    },
    splitLine: {
      show: true,
      lineStyle: {
        color: ['#eeeeee']
      }
    },
    splitArea: {
      show: false,
      areaStyle: {
        color: ['rgba(250,250,250,0.3)', 'rgba(200,200,200,0.3)']
      }
    }
  },
  logAxis: {
    axisLine: {
      show: true,
      lineStyle: {
        color: '#cccccc'
      }
    },
    axisTick: {
      show: true,
      lineStyle: {
        color: '#cccccc'
      }
    },
    axisLabel: {
      show: true,
      color: '#999999'
    },
    splitLine: {
      show: true,
      lineStyle: {
        color: ['#eeeeee']
      }
    },
    splitArea: {
      show: false,
      areaStyle: {
        color: ['rgba(250,250,250,0.3)', 'rgba(200,200,200,0.3)']
      }
    }
  },
  timeAxis: {
    axisLine: {
      show: true,
      lineStyle: {
        color: '#cccccc'
      }
    },
    axisTick: {
      show: true,
      lineStyle: {
        color: '#cccccc'
      }
    },
    axisLabel: {
      show: true,
      color: '#999999'
    },
    splitLine: {
      show: true,
      lineStyle: {
        color: ['#eeeeee']
      }
    },
    splitArea: {
      show: false,
      areaStyle: {
        color: ['rgba(250,250,250,0.3)', 'rgba(200,200,200,0.3)']
      }
    }
  },
  toolbox: {
    iconStyle: {
      borderColor: '#999'
    },
    emphasis: {
      iconStyle: {
        borderColor: '#666'
      }
    }
  },
  legend: {
    textStyle: {
      color: '#999999'
    }
  },
  tooltip: {
    axisPointer: {
      lineStyle: {
        color: '#cccccc',
        width: 1
      },
      crossStyle: {
        color: '#cccccc',
        width: 1
      }
    }
  },
  timeline: {
    lineStyle: {
      color: '#DAE1F5',
      width: 2
    },
    itemStyle: {
      color: '#A4B1D7',
      borderWidth: 1
    },
    controlStyle: {
      color: '#A4B1D7',
      borderColor: '#A4B1D7',
      borderWidth: 1
    },
    checkpointStyle: {
      color: '#316bf3',
      borderColor: 'fff'
    },
    label: {
      color: '#A4B1D7'
    },
    emphasis: {
      itemStyle: {
        color: '#316bf3'
      },
      controlStyle: {
        color: '#316bf3',
        borderColor: '#316bf3',
        borderWidth: 2
      },
      label: {
        color: '#316bf3'
      }
    }
  },
  visualMap: {
    color: ['#bf444c', '#d88273', '#f6efa6']
  },
  dataZoom: {
    handleIcon: 'M10.7,11.9v-1.3H9.3v1.3c-4.9,0.3-8.8,4.4-8.8,9.4c0,5,3.9,9.1,8.8,9.4v1.3h1.3v-1.3c4.9-0.3,8.8-4.4,8.8-9.4C19.5,16.3,15.6,12.2,10.7,11.9z M13.3,24.4H6.7V23.1h6.6V24.4z M13.3,19.6H6.7v-1.4h6.6V19.6z',
    dataBackground: {
      lineStyle: {
        color: '#eeeeee',
        width: 1
      },
      areaStyle: {
        color: 'rgba(250,250,250,0.3)'
      }
    },
    dataBackgroundSelected: {
      lineStyle: {
        color: '#8392A5'
      },
      areaStyle: {
        color: 'rgba(90,90,90,0.3)'
      }
    },
    fillerColor: 'rgba(167,183,204,0.4)',
    handleColor: '#a7b7cc',
    handleStyle: {
      color: '#a7b7cc',
      borderColor: '#000'
    },
    moveHandleStyle: {
      color: '#A0B1C4',
      opacity: 0.3
    },
    selectedDataBackground: {
      lineStyle: {
        color: '#000000',
        width: 0.5
      },
      areaStyle: {
        color: 'rgba(0,0,0,0.3)'
      }
    }
  },
  markPoint: {
    label: {
      color: '#ffffff'
    },
    emphasis: {
      label: {
        color: '#ffffff'
      }
    }
  }
};

// 实时数据管理 Hook
export function useRealTimeData(apiFunction: Function, interval = 5000) {
  const data = ref<any>(null);
  const loading = ref(false);
  const error = ref<string | null>(null);
  let timer: NodeJS.Timeout | null = null;

  const fetchData = async () => {
    try {
      loading.value = true;
      error.value = null;
      const response = await apiFunction();
      data.value = response.data;
    } catch (err: any) {
      error.value = err.message || '数据加载失败';
      console.error('实时数据获取失败:', err);
    } finally {
      loading.value = false;
    }
  };

  const startPolling = () => {
    if (timer) return;
    
    fetchData(); // 立即获取一次数据
    timer = setInterval(fetchData, interval);
  };

  const stopPolling = () => {
    if (timer) {
      clearInterval(timer);
      timer = null;
    }
  };

  const refresh = () => {
    fetchData();
  };

  onMounted(() => {
    startPolling();
  });

  onUnmounted(() => {
    stopPolling();
  });

  return {
    data,
    loading,
    error,
    refresh,
    startPolling,
    stopPolling
  };
}

// 图表基础 Hook
export function useChart(containerRef: any, theme?: string) {
  const chart = ref<echarts.ECharts | null>(null);
  const loading = ref(false);

  const initChart = () => {
    if (!containerRef.value) return;
    
    chart.value = echarts.init(containerRef.value, theme || chartTheme);
    
    // 监听窗口大小变化
    const resizeHandler = () => {
      chart.value?.resize();
    };
    
    window.addEventListener('resize', resizeHandler);
    
    // 清理函数
    onUnmounted(() => {
      window.removeEventListener('resize', resizeHandler);
      chart.value?.dispose();
    });
  };

  const setOption = (option: any, notMerge = false) => {
    if (!chart.value) return;
    
    chart.value.setOption(option, notMerge);
  };

  const showLoading = () => {
    loading.value = true;
    chart.value?.showLoading();
  };

  const hideLoading = () => {
    loading.value = false;
    chart.value?.hideLoading();
  };

  const resize = () => {
    chart.value?.resize();
  };

  onMounted(() => {
    initChart();
  });

  return {
    chart,
    loading,
    setOption,
    showLoading,
    hideLoading,
    resize,
    initChart
  };
}

// 时间序列图表 Hook
export function useTimeSeriesChart(containerRef: any, options: {
  title?: string;
  yAxisName?: string;
  unit?: string;
  smooth?: boolean;
  showSymbol?: boolean;
} = {}) {
  const { chart, setOption, showLoading, hideLoading } = useChart(containerRef);
  
  const updateData = (data: {
    timestamps: number[];
    series: Array<{
      name: string;
      data: number[];
      color?: string;
      type?: 'line' | 'bar';
    }>;
  }) => {
    const option = {
      title: {
        text: options.title,
        left: 'center',
        textStyle: {
          fontSize: 14,
          fontWeight: 'normal'
        }
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'cross'
        },
        formatter: (params: any) => {
          if (!params || !Array.isArray(params) || params.length === 0) {
            return '';
          }
          let result = `${new Date(params[0].axisValue).toLocaleString()}<br/>`;
          params.forEach((param: any) => {
            result += `${param.marker}${param.seriesName}: ${param.value}${options.unit || ''}<br/>`;
          });
          return result;
        }
      },
      legend: {
        data: data.series.map(s => s.name),
        bottom: 0
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '15%',
        containLabel: true
      },
      xAxis: {
        type: 'time',
        data: data.timestamps,
        axisLabel: {
          formatter: (value: number) => {
            const date = new Date(value);
            return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`;
          }
        }
      },
      yAxis: {
        type: 'value',
        name: options.yAxisName,
        axisLabel: {
          formatter: `{value}${options.unit || ''}`
        }
      },
      series: data.series.map(series => ({
        name: series.name,
        type: series.type || 'line',
        data: series.data.map((value, index) => [data.timestamps[index], value]),
        smooth: options.smooth !== false,
        showSymbol: options.showSymbol === true,
        lineStyle: {
          width: 2
        },
        itemStyle: {
          color: series.color
        }
      }))
    };

    setOption(option);
  };

  return {
    chart,
    updateData,
    showLoading,
    hideLoading
  };
}

// 饼图 Hook
export function usePieChart(containerRef: any, options: {
  title?: string;
  radius?: string | string[];
  center?: string[];
  showPercent?: boolean;
} = {}) {
  const { chart, setOption, showLoading, hideLoading } = useChart(containerRef);
  
  const updateData = (data: Array<{ name: string; value: number; color?: string }>) => {
    const option = {
      title: {
        text: options.title,
        left: 'center',
        textStyle: {
          fontSize: 14,
          fontWeight: 'normal'
        }
      },
      tooltip: {
        trigger: 'item',
        formatter: '{a} <br/>{b}: {c} ({d}%)'
      },
      legend: {
        orient: 'vertical',
        left: 'left',
        data: data.map(item => item.name)
      },
      series: [
        {
          name: options.title || '数据分布',
          type: 'pie',
          radius: options.radius || '50%',
          center: options.center || ['50%', '60%'],
          data: data.map(item => ({
            ...item,
            itemStyle: {
              color: item.color
            }
          })),
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: 'rgba(0, 0, 0, 0.5)'
            }
          },
          label: {
            show: options.showPercent !== false,
            formatter: '{b}: {d}%'
          }
        }
      ]
    };

    setOption(option);
  };

  return {
    chart,
    updateData,
    showLoading,
    hideLoading
  };
}

// 仪表盘图表 Hook
export function useGaugeChart(containerRef: any, options: {
  title?: string;
  min?: number;
  max?: number;
  unit?: string;
  splitNumber?: number;
  thresholds?: Array<{ value: number; color: string }>;
} = {}) {
  const { chart, setOption, showLoading, hideLoading } = useChart(containerRef);
  
  const updateData = (value: number, title?: string) => {
    const thresholds = options.thresholds || [
      { value: 0, color: '#52c41a' },
      { value: 70, color: '#faad14' },
      { value: 90, color: '#ff4d4f' }
    ];

    // 根据阈值确定颜色
    let color = thresholds[0]?.color || '#52c41a';
    for (let i = thresholds.length - 1; i >= 0; i--) {
      const threshold = thresholds[i];
      if (threshold && value >= threshold.value) {
        color = threshold.color;
        break;
      }
    }

    const option = {
      title: {
        text: options.title,
        left: 'center',
        top: '75%',
        textStyle: {
          fontSize: 14,
          fontWeight: 'normal'
        }
      },
      series: [
        {
          name: title || '当前值',
          type: 'gauge',
          min: options.min || 0,
          max: options.max || 100,
          splitNumber: options.splitNumber || 10,
          radius: '80%',
          axisLine: {
            lineStyle: {
              width: 3,
              color: thresholds.map(t => [t.value / (options.max || 100), t.color])
            }
          },
          pointer: {
            itemStyle: {
              color: color
            }
          },
          axisTick: {
            distance: -30,
            length: 8,
            lineStyle: {
              color: '#fff',
              width: 2
            }
          },
          splitLine: {
            distance: -30,
            length: 30,
            lineStyle: {
              color: '#fff',
              width: 4
            }
          },
          axisLabel: {
            color: 'inherit',
            distance: 40,
            fontSize: 12,
            formatter: `{value}${options.unit || ''}`
          },
          detail: {
            valueAnimation: true,
            formatter: `{value}${options.unit || ''}`,
            color: 'inherit',
            fontSize: 20,
            offsetCenter: [0, '50%']
          },
          data: [
            {
              value: value,
              name: title || '当前值'
            }
          ]
        }
      ]
    };

    setOption(option);
  };

  return {
    chart,
    updateData,
    showLoading,
    hideLoading
  };
}

// 热力图 Hook
export function useHeatmapChart(containerRef: any, options: {
  title?: string;
  xAxisData?: string[];
  yAxisData?: string[];
  visualMap?: {
    min: number;
    max: number;
    calculable?: boolean;
    orient?: 'horizontal' | 'vertical';
    left?: string;
    bottom?: string;
  };
} = {}) {
  const { chart, setOption, showLoading, hideLoading } = useChart(containerRef);
  
  const updateData = (data: Array<[number, number, number]>) => {
    const option = {
      title: {
        text: options.title,
        left: 'center',
        textStyle: {
          fontSize: 14,
          fontWeight: 'normal'
        }
      },
      tooltip: {
        position: 'top',
        formatter: (params: any) => {
          if (!params || !params.data || !Array.isArray(params.data)) {
            return '';
          }
          const [x, y, value] = params.data;
          const xLabel = options.xAxisData?.[x] ?? x;
          const yLabel = options.yAxisData?.[y] ?? y;
          return `${xLabel} - ${yLabel}: ${value}`;
        }
      },
      grid: {
        height: '50%',
        top: '10%'
      },
      xAxis: {
        type: 'category',
        data: options.xAxisData || [],
        splitArea: {
          show: true
        }
      },
      yAxis: {
        type: 'category',
        data: options.yAxisData || [],
        splitArea: {
          show: true
        }
      },
      visualMap: {
        min: options.visualMap?.min || 0,
        max: options.visualMap?.max || 100,
        calculable: options.visualMap?.calculable !== false,
        orient: options.visualMap?.orient || 'horizontal',
        left: options.visualMap?.left || 'center',
        bottom: options.visualMap?.bottom || '15%'
      },
      series: [
        {
          name: options.title || '热力图',
          type: 'heatmap',
          data: data,
          label: {
            show: true
          },
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowColor: 'rgba(0, 0, 0, 0.5)'
            }
          }
        }
      ]
    };

    setOption(option);
  };

  return {
    chart,
    updateData,
    showLoading,
    hideLoading
  };
}

// 多维数据可视化 Hook
export function useRadarChart(containerRef: any, options: {
  title?: string;
  indicators: Array<{ name: string; max: number; min?: number }>;
} = { indicators: [] }) {
  const { chart, setOption, showLoading, hideLoading } = useChart(containerRef);
  
  const updateData = (data: Array<{ name: string; value: number[]; color?: string }>) => {
    const option = {
      title: {
        text: options.title,
        left: 'center',
        textStyle: {
          fontSize: 14,
          fontWeight: 'normal'
        }
      },
      tooltip: {
        trigger: 'item'
      },
      legend: {
        data: data.map(item => item.name),
        bottom: 0
      },
      radar: {
        indicator: options.indicators,
        radius: '60%'
      },
      series: [
        {
          name: options.title || '雷达图',
          type: 'radar',
          data: data.map(item => ({
            value: item.value,
            name: item.name,
            itemStyle: {
              color: item.color
            }
          }))
        }
      ]
    };

    setOption(option);
  };

  return {
    chart,
    updateData,
    showLoading,
    hideLoading
  };
}

// 数据格式化工具
export const formatters = {
  // 格式化文件大小
  fileSize: (bytes: number): string => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  },

  // 格式化百分比
  percentage: (value: number, decimals = 1): string => {
    return `${value.toFixed(decimals)}%`;
  },

  // 格式化数字（添加千分位分隔符）
  number: (value: number): string => {
    return value.toLocaleString();
  },

  // 格式化时间
  time: (timestamp: number, format = 'HH:mm'): string => {
    const date = new Date(timestamp);
    if (format === 'HH:mm') {
      return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`;
    }
    return date.toLocaleString();
  },

  // 格式化持续时间
  duration: (seconds: number): string => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;

    if (hours > 0) {
      return `${hours}小时${minutes}分钟`;
    } else if (minutes > 0) {
      return `${minutes}分钟${secs}秒`;
    } else {
      return `${secs}秒`;
    }
  }
};

// 颜色工具
export const colorUtils = {
  // 根据值获取颜色
  getColorByValue: (value: number, thresholds: Array<{ value: number; color: string }>): string => {
    let color = thresholds[0]?.color || '#52c41a';
    for (let i = thresholds.length - 1; i >= 0; i--) {
      const threshold = thresholds[i];
      if (threshold && value >= threshold.value) {
        color = threshold.color;
        break;
      }
    }
    return color;
  },

  // 生成渐变色
  generateGradient: (startColor: string, endColor: string, steps: number): string[] => {
    const colors = [];
    const startRGB = hexToRgb(startColor);
    const endRGB = hexToRgb(endColor);
    
    if (!startRGB || !endRGB) return [startColor, endColor];
    
    for (let i = 0; i < steps; i++) {
      const ratio = i / (steps - 1);
      const r = Math.round(startRGB.r + ratio * (endRGB.r - startRGB.r));
      const g = Math.round(startRGB.g + ratio * (endRGB.g - startRGB.g));
      const b = Math.round(startRGB.b + ratio * (endRGB.b - startRGB.b));
      colors.push(`rgb(${r}, ${g}, ${b})`);
    }
    
    return colors;
  }
};

// 辅助函数
function hexToRgb(hex: string): { r: number; g: number; b: number } | null {
  if (!hex || typeof hex !== 'string') {
    return null;
  }
  const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
  return result ? {
    r: parseInt(result[1] || '0', 16),
    g: parseInt(result[2] || '0', 16),
    b: parseInt(result[3] || '0', 16)
  } : null;
}

// 图表配置预设
export const chartPresets = {
  // CPU使用率时间序列图
  cpuUsage: {
    title: 'CPU使用率',
    yAxisName: '使用率',
    unit: '%',
    smooth: true,
    thresholds: [
      { value: 0, color: '#52c41a' },
      { value: 70, color: '#faad14' },
      { value: 90, color: '#ff4d4f' }
    ]
  },

  // 内存使用率时间序列图
  memoryUsage: {
    title: '内存使用率',
    yAxisName: '使用率',
    unit: '%',
    smooth: true,
    thresholds: [
      { value: 0, color: '#52c41a' },
      { value: 80, color: '#faad14' },
      { value: 95, color: '#ff4d4f' }
    ]
  },

  // GPU使用率仪表盘
  gpuGauge: {
    title: 'GPU使用率',
    min: 0,
    max: 100,
    unit: '%',
    thresholds: [
      { value: 0, color: '#52c41a' },
      { value: 70, color: '#faad14' },
      { value: 90, color: '#ff4d4f' }
    ]
  },

  // 任务状态饼图
  jobStatus: {
    title: '任务状态分布',
    radius: ['40%', '70%'],
    center: ['50%', '60%'],
    showPercent: true
  },

  // 模型类型分布
  modelTypes: {
    title: '模型类型分布',
    radius: '60%',
    center: ['50%', '50%'],
    showPercent: true
  }
};