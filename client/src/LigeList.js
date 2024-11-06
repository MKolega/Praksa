// LigeList.js
import React, { useEffect, useState } from 'react';
import { getLige } from './apiService';

const LigeList = () => {
    const [lige, setLige] = useState([]);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const data = await getLige();
                setLige(data);
            } catch (error) {
                console.error("Error fetching leagues:", error);
            }
        };
        fetchData();
    }, []);

    return (
        <div>
            <h2>Leagues</h2>
            <ul>
                {lige.map((liga) => (
                    <li key={liga.id}>{liga.name}</li>
                ))}
            </ul>
        </div>
    );
};

export default LigeList;
