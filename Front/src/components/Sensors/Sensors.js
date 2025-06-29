import React, { useEffect, useState } from "react";
import { Link } from 'react-router-dom';
import "./Sensors.css"

function Sensors() {

    const [sensors, setSensors] = useState([

    ]);

    useEffect(() => {
        fetchSensors()
    }, []);

    const fetchSensors = async () => {
        try {
            const response = await fetch(`${process.env.REACT_APP_API_URL}/sensors`);
            const result = await response.json();

            setSensors(result)
        }
        catch (error) {
            console.error(error)
        }
    }

    return (
        <div className='container'>
            <div className='header1'>
                You can find the source code at:<br></br>
                <a href="https://github.com/kostinalex/RealTimeDashBoardExample" target='_blank' >https://github.com/kostinalex/RealTimeDashBoardExample</a>
            </div>

            <div className='header2'>Sensors</div>
            {
                sensors.map(sensor =>
                (<div key={sensor.id}><Link className="link" to={'/sensor/' + sensor.id}><div className="sensorWraper">
                    <div>{sensor.name}</div>
                    <div className={sensor.temperature > 22 ? "alert" : "temperature"}>{sensor.temperature}°C</div>
                </div></Link></div>)
                )
            }
        </div>
    );
}

export default Sensors;
