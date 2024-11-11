import React, { useEffect, useState } from 'react';
import {deposit, getLige, getPonude, loginUser, passwordReset, registerUser} from './apiService';
import { format } from 'date-fns';
import './HomePage.css';

const HomePage = () => {
    const [lige, setLige] = useState([]);
    const [error] = useState(null);
    const [ponude, setPonude] = useState([]);
    const [username, setUsername] = useState('');
    const [AccountID, setID] = useState(0);
    const [isLoggedIn, setIsLoggedIn] = useState(false);
    const [funds, setFunds] = useState(0);
    const [showRegisterPopup, setShowRegisterPopup] = useState(false);
    const [showPasswordResetPopup, setShowPasswordResetPopup] = useState(false);

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
                const AccountID = response.id;
                const AccountBalance = response.account_balance;
                setUsername(username);
                setID(AccountID);
                setFunds(AccountBalance);
                setIsLoggedIn(true);
                console.log('Login successful');

            } else {
                alert('Invalid credentials');
            }
        } catch (err) {
            alert(`Login failed - ${err.message}`);
        }
    };
    const handleLogout = () => {
        setUsername('');
        setID(0);
        setFunds(0);
        setIsLoggedIn(false);
    };


    const handleRegisterClick = () => {
    setShowRegisterPopup(true);
};
    const handleRegister = async ()  => {

    const username = document.querySelector('.register-input[type="text"]').value;
    const password = document.querySelector('.register-input[type="password"]').value;
    const confirmPassword = document.querySelector('.register-input[type="password"]:nth-of-type(2)').value;

    if (password !== confirmPassword) {
        alert('Passwords do not match');
        return;
    }
    try{
        const response = await registerUser(username, password);
        if (!response.error) {
            alert('Registration successful');
            await loginUser(username, password);
        } else {
            alert('Error registering');
        }
    }catch (err) {
        alert(`Registration failed - ${err.message}`);
    }
    setShowRegisterPopup(false);
};

    const handlePasswordResetClick = () => {
        setShowPasswordResetPopup(true);
    }
    const handlePasswordReset = async () => {
        const username = document.querySelector('.register-input[type="text"]').value;
        const password = document.querySelector('.register-input[type="password"]').value;
        const confirmPassword = document.querySelector('.register-input[type="password"]:nth-of-type(2)').value;

        if (password !== confirmPassword) {
            alert('Passwords do not match');
            return;
        }
        try {
            const response = await passwordReset(username, password);
            if (!response.error) {
                alert('Password reset successful');
            } else {
                alert('Error resetting password');
            }
        } catch (err) {
            alert(`Password reset failed - ${err.message}`);
        }
        setShowPasswordResetPopup(false);
    }

    const handleAddFunds = async () => {
        const amount = prompt("Enter the amount to add:");
        if (amount) {
            try {
                const response = await deposit(AccountID, parseFloat(amount));
                if (!response.error) {
                    setFunds(funds + parseFloat(amount));
                    alert('Funds added successfully');
                } else {
                    alert('Error adding funds');
                }
            } catch (err) {
                alert(`Failed to add funds - ${err.message}`);
            }
        }
    };

    const getPonudeForLiga = (liga) => {
        return ponude
            .filter((ponuda) => liga.razrade.some((razrada) => razrada.ponude.includes(ponuda.id)))
            .sort((a, b) => new Date(a.vrijeme) - new Date(b.vrijeme));
    };

    if (error) return <div>{error}</div>;

    return (
        <div className={`home-page`}>
            {isLoggedIn ? (
                <div className="welcome-message">
                    Welcome, {username} (Balance: ${funds.toFixed(2)})
                    <button className="add-funds-button" onClick={handleAddFunds}>Add Funds</button>
                    <button className="logout-button" onClick={handleLogout}>Logout</button>
                </div>

            ) : (
                <div className="login-container">
                    <div className="login-form">
                        <input type="username" placeholder="nadimak" className="login-input"/>
                        <input type="password" placeholder="lozinka" className="login-input"/>
                        <button className="login-button" onClick={handleLogin}>PRIJAVA</button>
                    </div>
                    <div className="login-links">
                        <a href="#" className="register-link" onClick={handleRegisterClick}>Registracija</a> <a
                        href="#" className="forgot-password-link" onClick={handlePasswordResetClick}>Zaboravio sam lozinku!</a>
                    </div>
                </div>
            )}
            {showRegisterPopup && (
                <div className="popup">
                    <div className="popup-inner">
                        <h2>Register</h2>
                        <input type="text" placeholder="Username" className="register-input"/>
                        <input type="password" placeholder="Password" className="register-input"/>
                        <input type="password" placeholder="Confirm Password" className="register-input"/>
                        <button className="register-button" onClick={handleRegister}>Register</button>
                        <button className="close-button" onClick={() => setShowRegisterPopup(false)}>Close</button>
                    </div>
                </div>
            )}
            {showPasswordResetPopup && (
            <div className="popup">
                <div className="popup-inner">
                    <h2>Password Reset</h2>
                    <input type="text" placeholder="Username" className="register-input"/>
                    <input type="password" placeholder="New Password" className="register-input"/>
                    <input type="password" placeholder="Confirm Password" className="register-input"/>
                    <button className="register-button" onClick={handlePasswordReset}>Reset</button>
                    <button className="close-button" onClick={() => setShowPasswordResetPopup(false)}>Close</button>
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
