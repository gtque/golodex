import React from 'react';
import {Button} from "@material-ui/core";
import AddBoxIcon from "@material-ui/icons/AddBox";
import SearchIcon from "@material-ui/icons/Search";

var arrayConstructor = [].constructor;
var objectConstructor = ({}).constructor;

export default class Deck extends React.Component {
    constructor(props) {
        super(props);
        this.getHeader = this.getHeader.bind(this);
        this.getRowsData = this.getRowsData.bind(this);
        this.getKeys = this.getKeys.bind(this);
        this.handleEdit = this.handleEdit.bind(this)
        // Initial states
        this.state = {
            open: false,
            isEdit: false,
            disable: true,
            showConfirm: false,
            search: ""
        }
    }

    handleEdit = function (card) {
        console.log("edit clicked...: " + card)
        this.props.setPage("card")
        this.props.setCard(card)
    }

    getKeys = function (data) {
        if(data !== null && data !== undefined) {
            return Object.keys(data);
        }
        return {}
    }

    getHeader = function (data) {
        var keys = this.getKeys(data);
        return keys.map((key, index) => {
            if (key.startsWith("_")) {
                return
            }
            return <th key={key}>{key.toUpperCase()}</th>
        })
    }

    getRowsData = function (data) {
        var items = data;
        var keys = this.getKeys(data);
        return items.map((row, index) => {
            return <tr key={index}><RenderRow key={index} data={row} keys={keys} handleEdit={this.handleEdit} getKeys={this.getKeys} getHeader={this.getHeader} getRowsData={this.getRowsData}/></tr>
        })
    }

    render() {
        if(this.props.data && this.props.data.length > 0) {
            return (
                <table>
                    <thead>
                    <tr>{this.getHeader(this.props.data[0])}</tr>
                    </thead>
                    <tbody>
                    <RenderRows keyT={"root"} key={"root"} data={this.props.data} handleEdit={this.handleEdit} getKeys={this.getKeys} getHeader={this.getHeader} getRowsData={this.getRowsData}/>
                    </tbody>
                </table>
            );
        } else {
            return (
                <table/>
            )
        }
    }
}

const RenderRows = (props) => {
    console.log("trying to render rowS")
    console.log("-> " + props.keyT)
    var items = props.data;
    var keys = props.getKeys(props.data[0]);

    return items.map((row, index) => {
        if (props.keyT === "root") {
            return <tr id={row._id} onClick={() => props.handleEdit(row._id)} class="full_row" key={props.keyT + index}><RenderRow key={props.keyT + index} keyT={props.keyT + index} data={row} handleEdit={props.handleEdit} keys={keys}
                                                                                                                                   getKeys={props.getKeys} getHeader={props.getHeader} getRowsData={props.getRowsData}/></tr>
        }
        return <tr key={props.keyT + index} class="nested"><RenderRow key={props.keyT + index} keyT={props.keyT + index} data={row} handleEdit={props.handleEdit} keys={keys} getKeys={props.getKeys} getHeader={props.getHeader}
                                                                      getRowsData={props.getRowsData}/></tr>
    })
}

const RenderRow = (props) => {
    console.log("trying to render row")
    console.log("-> " + props.keyT)
    return props.keys.map((keyT, index) => {
        if (keyT.startsWith("_")) {
            return
        }
        console.log("render row: " + keyT)
        // return <td key={props.data[key]}>{props.data[key]}</td>
        return <RenderData data={props.data[keyT]} key={props.keyT + keyT} keyT={props.keyT + keyT} handleEdit={props.handleEdit} getKeys={props.getKeys} getHeader={props.getHeader} getRowsData={props.getRowsData}/>
    })
}

const RenderData = (props) => {
    console.log("trying to render data: " + props.keyT)
    if (props.data === arrayConstructor || (props.data !== undefined && props.data !== null && props.data.constructor === Array)) {
        console.log("gently down the stream")
        return (<td class="nested_rows">
            <table>
                {/*<thead>*/}
                {/*<tr>{props.getHeader(props.data[0])}</tr>*/}
                {/*</thead>*/}
                <tbody>
                <RenderRows key={props.keyT} keyT={props.keyT} data={props.data} keys={props.getKeys(props.data)} handleEdit={props.handleEdit} getKeys={props.getKeys} getHeader={props.getHeader} getRowsData={props.getRowsData}/>
                </tbody>
            </table>
        </td>)
        // var items = props.data;
        // var keys = props.getKeys(props.data[0]);
        // console.log("row,row,row")
        // return items.map((row, index) => {
        //     return <RenderRow key={props.keyT+index} keyT={props.keyT+index} data={row} keys={keys} getKeys={props.getKeys} getHeader={props.getHeader} getRowsData={props.getRowsData}/>
        // })
        // return <RenderRows key={props.keyT} keyT={props.keyT} data={props.data} keys={keys} getKeys={props.getKeys} getHeader={props.getHeader} getRowsData={props.getRowsData}/>
        // return items.map((row, index) => {
        //     console.log("row: " + index)
        //     return <tr key={props.keyT+index}><RenderRow key={props.keyT+index} keyT={props.keyT+index} data={row} keys={keys} getKeys={props.getKeys} getHeader={props.getHeader} getRowsData={props.getRowsData}/></tr>
        // })
    } else if (props.data === objectConstructor || (props.data !== undefined && props.data !== null && props.data.constructor === Object)) {
        console.log("object time")
        console.log("-> " + Object.keys(props))
        // return <td><table><thead><tr>boo</tr></thead><tbody><tr><td key={props.keyT}>smile</td></tr></tbody></table></td>
        return (<td class="nested_obj">
            <table>
                {/*<thead>*/}
                {/*<tr>{props.getHeader(props.data)}</tr>*/}
                {/*</thead>*/}
                <tbody>
                <tr><RenderRow key={props.keyT} keyT={props.keyT} data={props.data} keys={props.getKeys(props.data)} handleEdit={props.handleEdit} getKeys={props.getKeys} getHeader={props.getHeader} getRowsData={props.getRowsData}/></tr>
                </tbody>
            </table>
        </td>)
    } else {
        console.log("just a string..." + props.keyT)
        console.log("->" + props.data)
        return <td key={props.keyT}>{props.data}</td>
    }
}