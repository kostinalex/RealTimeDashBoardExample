import './App.css';
import { Routes, Route, Navigate } from 'react-router-dom';


import Sensors from './components/Sensors/Sensors';
import Sensor from './components/Sensor/Sensor';

function App() {
  return (
    <div>
      <Routes>
        <Route path="/" element={<Navigate to="/sensors" replace />} />
        <Route path="/sensors" element={<Sensors />} />
        <Route path="/sensor/:id" element={<Sensor />} />
      </Routes>
    </div>
  );
}

export default App;
