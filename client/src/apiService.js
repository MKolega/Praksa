// apiService.js
const BASE_URL = 'http://localhost:8080/api'; // Update to match your server

export const getLige = async () => {
    const response = await fetch(`${BASE_URL}/lige`);
    if (!response.ok) {
        throw new Error(`Error fetching leagues: ${response.statusText}`);
    }
    return response.json();
};

export const getPonude = async () => {
    const response = await fetch(`${BASE_URL}/ponude`);
    if (!response.ok) {
        throw new Error(`Error fetching ponude: ${response.statusText}`);
    }
    return response.json();
}

export const getPlayers = async () => {
    const response = await fetch(`${BASE_URL}/players`);
    if (!response.ok) {
        throw new Error(`Error fetching players: ${response.statusText}`);
    }
    return response.json();
};

export const getPlayerById = async (id) => {
    const response = await fetch(`${BASE_URL}/players/${id}`);
    if (!response.ok) {
        throw new Error(`Error fetching player ${id}: ${response.statusText}`);
    }
    return response.json();
};

export const createPonuda = async (data) => {
    const response = await fetch(`${BASE_URL}/ponude`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
    });
    if (!response.ok) {
        throw new Error(`Error creating ponuda: ${response.statusText}`);
    }
    return response.json();
};

export const deposit = async (id, amount) => {
    const response = await fetch(`${BASE_URL}/deposit/${id}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ amount }),
    });
    if (!response.ok) {
        throw new Error(`Error making deposit for player ${id}: ${response.statusText}`);
    }
    return response.json();
};

export const uplata = async (id, data) => {
    const response = await fetch(`${BASE_URL}/uplata/${id}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
    });
    if (!response.ok) {
        throw new Error(`Error creating uplata for player ${id}: ${response.statusText}`);
    }
    return response.json();
};