import React, { useEffect, useState } from 'react';
import { getLige, getPonude } from './apiService';
import { format } from 'date-fns';
import './HomePage.css';

const HomePage = () => {
    const [lige, setLige] = useState([]);
    const [error, setError] = useState(null);
    const [ponude, setPonude] = useState([]);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const ligeData = await getLige();
                setLige(ligeData);

                const ponudeData = await getPonude();
                setPonude(ponudeData);

            } catch (err) {
                setError("Could not load data");
                console.error(err);
            }
        };
        fetchData();
    }, []);


    const getPonudeForLiga = (ligaId) => {
        return ponude
            .filter((ponuda) => ponuda.ligaId === ligaId)
            .sort((a, b) => new Date(a.vrijeme) - new Date(b.vrijeme));
    };

    if (error) return <div>{error}</div>;

    return (
        <div className="home-page">
            <h1>Sport</h1>
            {lige.map((liga) => (
                <div key={liga.id} className="liga-section">
                    <h2>{liga.naziv}</h2>
                    <table className="league-table">
                        <thead>
                        <tr>
                            <th>Match</th>
                            <th>Time</th>
                            <th>1</th>
                            <th>X</th>
                            <th>2</th>
                            <th>1X</th>
                            <th>X2</th>
                            <th>12</th>
                            <th>F+2</th>
                        </tr>
                        </thead>
                        <tbody>
                        {getPonudeForLiga(liga.id).map((ponuda) => (
                            <tr key={ponuda.id} className="ponuda-row">
                                <td>{ponuda.naziv}</td>
                                <td>{format(new Date(ponuda.vrijeme), 'dd/MM/yyyy HH:mm')}</td>
                                {['1', 'X', '2', '1X', 'X2', '12', 'f+2'].map((betType) => (
                                    <td key={betType}>
                                        {ponuda.tecajevi.find((tecaj) => tecaj.naziv === betType)?.tecaj || '-'}
                                    </td>
                                ))}
                            </tr>
                        ))}
                        </tbody>
                    </table>
                </div>
            ))}
        </div>
    );
};

export default HomePage;