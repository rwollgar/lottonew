import React, {useState, useEffect} from "react";

import dynamic from 'next/dynamic';
const Chart = dynamic(() => import('react-apexcharts'),{ssr: false});


const LnChart = (props) => {

    const [chartState, setChartState] = useState({
        options: {
            chart: {
                id: "basic-bar"
            },
            xaxis: {
                categories: [1991, 1992, 1993, 1994, 1995, 1996, 1997, 1998, 1999]
            }
        },
        series: [
            {
                name: "series-1",
                data: [30, 40, 45, 50, 49, 60, 70, 91]
            }
        ]
    });


    return (
        
        <div>
            <Chart
                options={chartState.options}
                series={chartState.series}
                type="bar"
                width="90%"
            />            
        </div>

    )

}

export default LnChart;