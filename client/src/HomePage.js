import React, { useEffect, useState } from 'react';
import { getLige, getPonude, loginUser} from './apiService';
import { format } from 'date-fns';
import './HomePage.css';

const HomePage = () => {
    const [lige, setLige] = useState([]);
    const [error, setError] = useState(null);
    const [ponude, setPonude] = useState([]);
    const [username, setUsername] = useState('');
    const [isLoggedIn, setIsLoggedIn] = useState(false);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const ligeData = await getLige();
                setLige(ligeData);

                const ponudeData = await getPonude();
                setPonude(ponudeData);

            } catch (err) {
                alert("Could not load data");
            }
        };
        fetchData();
    }, []);

    const handleLogin = async () => {
        const username = document.querySelector('.login-input[type="username"]').value;
        const password = document.querySelector('.login-input[type="password"]').value;

        try {
            const response = await loginUser(username, password);
            if (!response.error) {
                setUsername(username);
                setIsLoggedIn(true);
                console.log('Login successful');
            } else {
                alert('Invalid credentials');
            }
        } catch (err) {
            alert(`Login failed - ${err.message}`);
        }
    };

    const getPonudeForLiga = (liga) => {
        return ponude
            .filter((ponuda) => liga.razrade.some((razrada) => razrada.ponude.includes(ponuda.id)))
            .sort((a, b) => new Date(a.vrijeme) - new Date(b.vrijeme));
    };

    if (error) return <div>{error}</div>;

    return (
        <div className="home-page">
            {isLoggedIn ? (
                <div className="welcome-message">Welcome, {username}</div>
            ) : (
                <div className="login-container">
                    <div className="login-links">
                        <a href="/register" className="register-link">Registracija</a>
                        <a href="/forgot-password" className="forgot-password-link">Zaboravio sam lozinku!</a>
                    </div>
                    <div className="login-form">
                        <input type="username" placeholder="nadimak" className="login-input"/>
                        <input type="password" placeholder="lozinka" className="login-input"/>
                        <button className="login-button" onClick={handleLogin}>PRIJAVI ME</button>
                    </div>
                </div>
            )}
            <h1 className="sport">Sport</h1>
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
                        {getPonudeForLiga(liga).map((ponuda) => (
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
