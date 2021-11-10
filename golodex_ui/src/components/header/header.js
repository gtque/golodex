import React from 'react';

import logo from "../../images/logo_transparent.png";

export default class Deck extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (
            <header className="golodex-header">
                <div className="logo_label">
                    <img src={logo} className="logo" alt="logo"/>
                    <p className="label">
                        Golodex,
                        Golodex<br/>
                        oh, how you rolodex.
                    </p>
                </div>
            </header>
        );
    }
}
