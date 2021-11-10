import React from 'react';
import Deck from "./deck";
import {Button} from "@material-ui/core";
import DoneIcon from "@material-ui/icons/Done";
import ClearIcon from "@material-ui/icons/Clear";
import DeleteIcon from "@material-ui/icons/Delete";
import Golodex from "../header/header";

export default class Card extends React.Component {
    constructor(props) {
        super(props);
        // this.getDeck = this.getDeck.bind(this);
        this.handleClose = this.handleClose.bind(this)
        this.handleDelete = this.handleDelete.bind(this)
        this.handleSave = this.handleSave.bind(this)
        this.setRolodex = this.setRolodex.bind(this)
        this.state = {
            error: null,
            isLoaded: false,
            save: "doing nothing",
            rolodex: {}
        };
    }

    handleClose = function () {
        this.props.setCard("")
        this.props.setPage("rolodex")
    }

    setRolodex = function (card) {
        this.setState({
            rolodex: [card]
        })
    }

    handleDelete = function () {
        // alert("editing card: " + this.props.card)
        if(this.props.card !== "_empty" && this.props.card !== "") {
            fetch(process.env.REACT_APP_GOLODEX_API_HOST + "/api/delete?card=" + this.props.card, {
                method: "DELETE",
                headers: {
                    "Content-Type": "application/json",
                    "X-GOLODEX-ID": this.props.getId(),
                    "X-GOLODEX-TOKEN": this.props.getToken(),
                }
            })
                .then(res => res.json())
                .then(
                    (result) => {
                        this.setState({
                            rolodex: [result],
                            isLoaded: true
                        });
                        this.props.setCard("_empty")
                        this.props.setPage("rolodex")
                    },
                    (error) => {
                        this.setState({
                            isLoaded: true,
                            error
                        });
                    }
                )
        } else {
            this.props.setCard("_empty")
            this.props.setPage("rolodex")
        }
        // window.location.reload()
    }

    handleSave = function () {
        // alert("editing card: " + this.props.card)
        if( this.state.save != "saving" ) {
            this.state.save = "saving"
            fetch(process.env.REACT_APP_GOLODEX_API_HOST + "/api/edit?card=" + this.props.card, {
                method: "PUT",
                body: JSON.stringify(this.state.rolodex[0]),
                headers: {
                    "Content-Type": "application/json",
                    "X-GOLODEX-ID": this.props.getId(),
                    "X-GOLODEX-TOKEN": this.props.getToken(),
                }
            })
                .then(res => res.json())
                .then(
                    (result) => {
                        this.setState({
                            rolodex: [result],
                            isLoaded: true,
                            save: "doing nothing",
                        });
                        this.props.setCard(result["_id"]);
                        this.state.save = "doing nothing"
                    },
                    (error) => {
                        this.setState({
                            isLoaded: true,
                            error
                        });
                    }
                )
        }
        // window.location.reload()
    }

    componentDidMount() {
        if (this.props.card !== "") {
            console.log("why am I fetching the card??? " + this.props.card)
            fetch(process.env.REACT_APP_GOLODEX_API_HOST + "/api/card?card=" + this.props.card, {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    "X-GOLODEX-ID": this.props.getId(),
                    "X-GOLODEX-TOKEN": this.props.getToken(),
                }
            })
                .then(res => res.json())
                .then(
                    (result) => {
                        this.setState({
                            rolodex: [result],
                            isLoaded: true
                        });
                    },
                    (error) => {
                        this.setState({
                            isLoaded: true,
                            error
                        });
                    }
                )
        } else {
            fetch(process.env.REACT_APP_GOLODEX_API_HOST + "/api/card?card=_empty", {
                method: "GET",
                headers: {
                    "Content-Type": "application/json",
                    "X-GOLODEX-ID": this.props.getId(),
                    "X-GOLODEX-TOKEN": this.props.getToken(),
                }
            })
                .then(res => res.json())
                .then(
                    (result) => {
                        this.setState({
                            rolodex: [result],
                            isLoaded: true
                        });
                    },
                    (error) => {
                        this.setState({
                            isLoaded: true,
                            error
                        });
                    }
                )
        }
    }

    render() {
        const {error, isLoaded, rolodex} = this.state;
        if (error) {
            return <div>Card Error: {error.message}</div>;
        } else if (!isLoaded) {
            return <div>Loading...</div>;
        } else {
            return (
                <div className="App">
                    <Golodex/>
                    <div className="content">
                        <div className="action_buttons">
                            <Button onClick={this.handleSave}>
                                <DoneIcon onClick={this.handleSave}/>
                                Save
                            </Button>
                            <Button onClick={this.handleDelete}>
                                <DeleteIcon onClick={this.handleDelete}/>
                                Delete
                            </Button>
                            <Button onClick={this.handleClose}>
                                <ClearIcon onClick={this.handleClose}/>
                                Close/Cancel
                            </Button>
                        </div>
                        <Deck data={rolodex} card={this.props.card} setPage={this.props.setPage} setCard={this.props.setCard} setRolodex={this.setRolodex}/>
                    </div>
                </div>
            );
        }
    }
}