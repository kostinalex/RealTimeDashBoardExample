import React, { useEffect, useState, useRef } from "react";
import { useParams } from "react-router-dom";
import { Line } from "react-chartjs-2";
import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
} from "chart.js";

import "./Sensor.css"
import Loading from "../Loading/Loading";


ChartJS.register(
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend
);

const options = {
    responsive: true,
    plugins: {
        legend: {
            position: "top",
        },
        title: {
            display: false,
        },
    },
    scales: {
        x: {
            ticks: {
                padding: 15,
                display: true,
                align: 'center',
                font: {
                    size: 10,
                    family: 'Poppins',
                },
                minRotation: 0,
                maxRotation: 0,
                maxTicksLimit: 6,
            },
            grid: {
                display: false,
            },
            border: {
                display: false,
            },
        },
        y: {
            max: 45,
            min: 5,
            ticks: {
                padding: 15,
                stepSize: 5,
                callback: (value) => value + '°C',
                font: {
                    size: 10,
                    family: 'Poppins',
                },
            },
            grid: {
                display: true,
                lineWidth: 1.5,
                color: (context) => {
                    if (
                        context.tick.value === 0
                    ) {
                        return 'rgba(0,0,0,0)';
                    }
                    return '#eeeeee';
                },
                drawTicks: false,
                drawOnChartArea: true,
            },
            border: {
                display: false,
                dash: [5, 2],
            },
        },
    },

};

function Sensor() {
    const { id } = useParams();
    const [sensorId, setSensorId] = useState(null);
    const [loading, setLoading] = useState(false);
    const [start, setStart] = useState(new Date(new Date().getTime() - 6 * 60 * 60 * 1000).toISOString().substring(0, 19));
    const [end, setEnd] = useState(new Date().toISOString().substring(0, 19));


    const [chartData, setChartData] = useState({
        labels: ['No data 1', 'No data 2', 'No data 3'],
        datasets: [
            {
                label: "Sensor Reading",
                data: [1, 2, 3],
                borderColor: "rgba(75, 192, 192, 1)",
                backgroundColor: "rgba(75, 192, 192, 0.2)",
                tension: 0.4,
                pointRadius: 1,
                pointHoverRadius: 1,
            },
        ],
    });

    const intervalRef = useRef(null);

    useEffect(() => {
        setSensorId(id);
        fetchSensorData()
    }, [id]);

    const fetchSensorData = async () => {

        try {
            setLoading(true)
            const response = await fetch(`${process.env.REACT_APP_API_URL}/readings/${id}/${start}Z/${end}Z`);
            const result = await response.json();
            const readings = result.readings
            const name = result.name

            setChartData({
                labels: readings.map(c => c.date),
                datasets: [
                    {
                        label: "Threshold",
                        data: readings.map(c => 30),
                        borderColor: "red",
                        backgroundColor: "red",
                        tension: 0.4,
                        pointRadius: 1,
                        pointHoverRadius: 1,
                    },
                    {
                        label: `Sensor "${name}"`,
                        data: readings.map(c => c.temperature),
                        borderColor: "rgba(75, 192, 192, 1)",
                        backgroundColor: "rgba(75, 192, 192, 1)",
                        tension: 0.4,
                        pointRadius: 1,
                        pointHoverRadius: 1,
                    },


                ],
            });

            setTimeout(() => { setLoading(false) }, 500)

        } catch (error) {
            console.error("Error fetching sensor data:", error);
        }
    };

    const changeStart = (e) => {
        setStart(e.target.value)
    }

    const changeEnd = (e) => {
        setEnd(e.target.value)
    }

    return (
        <div className="container">
            <div className="controlsWraper">
                <div><input value={start} onChange={changeStart} type="datetime-local" /></div>
                <div><input value={end} onChange={changeEnd} type="datetime-local" /></div>
                <div><button onClick={fetchSensorData}>Apply</button></div>
            </div>
            <div className="chartWraper">
                {
                    loading ? (<Loading />) : (
                        <div>

                            <div >
                                <Line data={chartData} options={options} />
                            </div>
                        </div>
                    )
                }
            </div>
        </div>
    );
}

export default Sensor;
