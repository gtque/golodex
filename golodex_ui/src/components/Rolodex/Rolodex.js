import React from 'react';
import Deck from "./deck";
import Golodex from "../header/header"
import {Button} from "@material-ui/core";
import AddBoxIcon from "@material-ui/icons/AddBox";
import SearchIcon from "@material-ui/icons/Search";
import SortIcon from "@material-ui/icons/SortByAlphaOutlined";

export default class Rolodex extends React.Component {
    constructor(props) {
        super(props);
        // this.getDeck = this.getDeck.bind(this);
        this.getSort = this.getSort.bind(this);
        this.getSearch = this.getSearch.bind(this);
        this.setSort = this.setSort.bind(this);
        this.setSearch = this.setSearch.bind(this);
        this.handleSearch = this.handleSearch.bind(this);
        this.handleAdd = this.handleAdd.bind(this)
        this.handleSort = this.handleSort.bind(this)
        this.state = {
            error: null,
            isLoaded: false,
            sort: "ascending",
            isReloaded: "",
            search: "",
            rolodex: []
        };
    }

    handleAdd = function () {
        this.props.setPage("card")
        this.props.setCard("_empty")
    }

    handleSort = function () {
        // alert("sort: " + this.state.sort)
        switch (this.state.sort) {
            case "ascending":
                // alert("setting sort to descending")
                this.setSort("descending")
                break
            case "descending":
                // alert("setting sort to ascending")
                this.setSort("ascending")
                break
            default:
                // alert("setting sort to default")
                this.setSort("descending")
                break
        }
        this.handleSearch()
    }

    getSort = function () {
        return this.state.sort
    }

    getSearch = function () {
        // const {search} = this.state;
        // alert("getting state: " + Object.keys(this.state))
        return this.state.search
    }

    setSort = function (sortType) {
        // this.setState({
        //     sort: sortType
        // })
        this.state.sort = sortType
    }

    setSearch = function (searchTerm) {
        this.setState({
            search: searchTerm.target.value
        })
    }

    handleSearch = function () {
        // alert("searching for: " + this.getSearch())
        //window.location.reload()
        // alert("searching sort order: " + this.state.sort)
        // console.log("searching........")
        this.setState({
            isReloaded: "searching..." + this.getSearch()
        })
        fetch(process.env.REACT_APP_GOLODEX_API_HOST + "/api/rolodex?sort=" + this.getSort() + "&search=" + this.getSearch(), {
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
                        rolodex: result,
                        isLoaded: true,
                        isReloaded: ""
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

    componentDidMount() {
        this.handleSearch()
        // return {
        //     "cards":
        //         {
        //             "_id": "0419055f7fe9956fc3bc3fb00f0005bd",
        //             "_rev": "4-35abba1ef3d411d0fe20800741406e19",
        //             "name": "Athan",
        //             "phone": "123456789",
        //             "address": "123 Sesame Street"
        //         }
        // }
        // return {
        //     "cards": [
        //         {
        //             "_id": "0419055f7fe9956fc3bc3fb00f0005bd",
        //             "_rev": "4-35abba1ef3d411d0fe20800741406e19",
        //             "name": {"first":"Athan","last":"Peanut"},
        //             "phone": "123456789",
        //             "address": "123 Sesame Street"
        //         }
        //     ]
        // }
        // return {
        //     "cards": [
        //         {
        //             "_id": "0419055f7fe9956fc3bc3fb00f0005bd",
        //             "_rev": "4-35abba1ef3d411d0fe20800741406e19",
        //             "name": {
        //                 "first": "Catniss",
        //                 "middle": "Everdeen",
        //                 "last": "Angeli"
        //             },
        //             "phone_numbers": [
        //                 {
        //                     "number": "1-234-555-1811",
        //                     "type": "home"
        //                 }
        //             ],
        //             "addresses": [
        //                 {
        //                     "street_1": "123 Sesame Street",
        //                     "street_2": "",
        //                     "city": "Alpharetta",
        //                     "state": "GA",
        //                     "zip_code": "30009",
        //                     "type": "home"
        //                 },
        //                 {
        //                     "street_1": "415 S Broad St",
        //                     "street_2": "",
        //                     "city": "Alpharetta",
        //                     "state": "GA",
        //                     "zip_code": "30009",
        //                     "type": "business"
        //                 },
        //                 {
        //                     "street_1": "1130 Mayfield Manor Drive",
        //                     "street_2": "",
        //                     "city": "Alpharetta",
        //                     "state": "GA",
        //                     "zip_code": "30009",
        //                     "type": "vacation"
        //                 }
        //             ]
        //         },
        //         {
        //             "_id": "0419055f7fe9956fc3bc3fb00f0005bd",
        //             "_rev": "4-35abba1ef3d411d0fe20800741406e19",
        //             "name": {
        //                 "first": "Catniss",
        //                 "middle": "Everdeen",
        //                 "last": "Angeli"
        //             },
        //             "phone_numbers": [
        //                 {
        //                     "number": "1-234-555-1811",
        //                     "type": "home"
        //                 }
        //             ],
        //             "addresses": [
        //                 {
        //                     "street_1": "123 Sesame Street",
        //                     "street_2": "",
        //                     "city": "Alpharetta",
        //                     "state": "GA",
        //                     "zip_code": "30009",
        //                     "type": "home"
        //                 },
        //                 {
        //                     "street_1": "415 S Broad St",
        //                     "street_2": "",
        //                     "city": "Alpharetta",
        //                     "state": "GA",
        //                     "zip_code": "30009",
        //                     "type": "business"
        //                 },
        //                 {
        //                     "street_1": "1130 Mayfield Manor Drive",
        //                     "street_2": "",
        //                     "city": "Alpharetta",
        //                     "state": "GA",
        //                     "zip_code": "30009",
        //                     "type": "vacation"
        //                 }
        //             ]
        //         }
        //     ]
        // }
    }

    render() {
        const {error, isLoaded, rolodex, isReloaded} = this.state;
        if (error) {
            return <div>Error: {error.message}</div>;
        } else if (!isLoaded) {
            return <div>Loading...</div>;
        } else {
            console.log("reloaded: " + isReloaded)
            return (
                <div className="App">
                    <Golodex/>
                    <div className="content">
                        <div style={{display: "flex", justifyContent: "space-between"}}>
                            <div>
                                <div>
                                    <Button onClick={this.handleAdd}>
                                        <AddBoxIcon onClick={this.handleAdd}/>
                                        ADD
                                    </Button>
                                    <Button onClick={this.handleSort}>
                                        <SortIcon onClick={this.handleSort}/>
                                        Sort
                                    </Button>
                                    <input type="text" name="search" value={this.getSearch() || ""} id="search" onChange={e => this.setSearch(e)}/>
                                    <Button onClick={this.handleSearch}>
                                        <SearchIcon onClick={this.handleSearch}/>
                                        Search
                                    </Button>
                                </div>
                            </div>
                        </div>
                        <Deck data={rolodex.cards} setPage={this.props.setPage} setCard={this.props.setCard}/>
                    </div>
                </div>
            );
        }
    }
}