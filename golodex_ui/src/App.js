import React from 'react';
import logo from './logo.svg';
import './App.css';
import {GoogleLogin} from 'react-google-login';
import Rolodex from "./components/Rolodex/Rolodex";
import Card from "./components/card/Card";

const handleLogin = async googleData => {
    const res = await fetch(process.env.REACT_APP_GOLODEX_API_HOST + "/api/login_google", {
        method: "POST",
        body: JSON.stringify({
            Token: googleData.tokenId,
            Id: "google-" + googleData.googleId,
            Email: googleData.profileObj.email,
            Name: googleData.profileObj.name,
        }),
        headers: {
            "Content-Type": "application/json"
        }
    })
    const data = await res.json()
    console.log(data);
    // store returned user somehow
    setToken(data);
    window.location.reload()
}

function setToken(loggedIn) {
    sessionStorage.setItem('token', JSON.stringify({
        token: loggedIn?.token,
        id: loggedIn?.id
    }));
}

function getToken() {
    const tokenString = sessionStorage.getItem('token');
    const userToken = JSON.parse(tokenString);
    return userToken?.token
}

function getId() {
    const tokenString = sessionStorage.getItem('token');
    const userToken = JSON.parse(tokenString);
    return userToken?.id
}

// function setCard(card) {
//     sessionStorage.setItem("card", card)
// }

// function getCard() {
//     return sessionStorage.getItem('card');
// }

const responseGoogle = (response) => {
    console.log(response);
}

function App() {

    const [page, setPage] = React.useState("login");
    // const card = getCard();
    const [card, setCard] = React.useState({})
    const token = getToken();

    if(page === "login" && token) {
        setPage("rolodex")
    }
    // if (!token) {
    switch (page) {
        case "login":
            return (
                <div className="App">
                    <header className="App-header">
                        <img src={logo} className="App-logo" alt="logo"/>
                        <p>
                            Welcome to Golodex!
                        </p>
                        <GoogleLogin
                            clientId={process.env.REACT_APP_GOOGLE_CLIENT_ID}
                            buttonText="Log in with Google"
                            onSuccess={handleLogin}
                            onFailure={responseGoogle}
                            cookiePolicy={'single_host_origin'}
                        />
                    </header>
                </div>
            );
            // } else {
        case "rolodex":
            return (
                <Rolodex key="the_rolodex" getId={getId} getToken={getToken} setPage={setPage} setCard={setCard}/>
            );
        case "card":
            if(card === "") {
                return (
                    <Card card="" getId={getId} getToken={getToken} setPage={setPage} setCard={setCard}/>
                );
            } else {
                // log.console("session card: " + cardId)
                return <Card card={card} getId={getId} getToken={getToken} setPage={setPage} setCard={setCard}/>
            }
        default:
            setPage("rolodex")
            break;
    }
}

export default App;
