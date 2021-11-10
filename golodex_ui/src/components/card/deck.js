import React from 'react';
import {Button} from "@material-ui/core";
import AddBoxIcon from "@material-ui/icons/AddBox";
import RemoveIcon from "@material-ui/icons/RemoveCircle";

var arrayConstructor = [].constructor;
var objectConstructor = ({}).constructor;

export default class Deck extends React.Component {
    constructor(props) {
        super(props);
        this.getHeader = this.getHeader.bind(this);
        this.getRowsData = this.getRowsData.bind(this);
        this.getKeys = this.getKeys.bind(this);
        this.handleRemove = this.handleRemove.bind(this)
        this.handleEdit = this.handleEdit.bind(this)
        this.nestedEdit = this.nestedEdit.bind(this)
        this.handleAdd = this.handleAdd.bind(this)
        // Initial states
        this.state = {
            open: false,
            isEdit: false,
            disable: true,
            showConfirm: false
        }
        // const [open, setOpen] = React.useState(false);
        // const [isEdit, setEdit] = React.useState(false);
        // const [disable, setDisable] = React.useState(true);
        // const [showConfirm, setShowConfirm] = React.useState(false);
    }

    handleEdit = function (key, _event) {
        console.log("value changed: " + key + "-> " + _event.target.value)
        // this.props.setPage("card")
        // this.props.setCard(card)
        var _card = this.props.data[0]
        var _keys = key.replace("root.", "").split(".")
        // console.log("changing value: " + _key + "-> " + _card[_key])
        var _spliced = [..._keys]
        _spliced.splice(0, 1)
        _card[_keys[0]] = this.nestedEdit(_spliced, _event.target.value, _card[_keys[0]]);
        // console.log("changed: " + _card[_key])
        this.props.setRolodex(_card)
    }

    nestedEdit = function (_keys, _value, _data) {
        if (_keys.length === 1) {
            _data[_keys[0]] = _value
        } else {
            var _spliced = [..._keys]
            _spliced.splice(0, 1)
            _data[_keys[0]] = this.nestedEdit(_spliced, _value, _data[_keys[0]])
        }
        return _data
    }

    handleRemove = function (key, data) {
        console.log("trying to remove: " + key)
        // alert("removing entry: " + key)
        var _key = key.replace("root.", "")
        var _index = _key.split(".")
        //this should be setting the rolodex from Card.js which should be using useState...
        var _card = this.props.data[0]
        var _indexI = _index[_index.length - 1]
        _key = _key.replace("." + _indexI, "")
        console.log("-> " + _key + ":" + _indexI)
        // alert([..._card[_key]])
        var _spliced = [..._card[_key]]
        _spliced.splice(parseInt(_indexI), 1)
        // alert(_spliced)
        _card[_key] = _spliced
        this.props.setRolodex(_card)
    }

    handleAdd = function (key, data) {
        // alert("adding entry to: " + key)
        var _new = {}
        var keys = this.getKeys(data);
        keys.map((key, index) => {
            _new[key] = ""
        })
        var _key = key.replace("root.", "")
        //this should be setting the rolodex from Card.js which should be using useState...
        var _card = this.props.data[0]
        _card[_key] = [..._card[_key], _new]
        this.props.setRolodex(_card)
    }


    getKeys = function (data) {
        return Object.keys(data);
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
            return <tr key={index}><RenderRow key={index} data={row} keys={keys} handle={this}/></tr>
        })
    }

    render() {
        return (
            <div>
                <div style={{display: "flex", justifyContent: "space-between"}}>
                    <div>
                        {this.state.isEdit ? (
                            <div>
                                <Button onClick={this.handleAdd}>
                                    <AddBoxIcon onClick={this.handleAdd}/>
                                    ADD
                                </Button>
                            </div>
                        ) : (
                            <div>
                            </div>
                        )}
                    </div>
                </div>
                <table>
                    <tbody>
                    <RenderRows keyT={"root"} key={"root"} data={this.props.data[0]} handle={this}/>
                    </tbody>
                </table>
            </div>

        );
    }
}

const RenderRows = (props) => {
    console.log("trying to render rowS")
    console.log("-> " + props.keyT)
    var items = props.data;
    var keys = props.handle.getKeys(props.data);

    return keys.map((row, index) => {
        if (row.startsWith("_")) {
            return
        }
        return <tr key={props.keyT + "." + row}>
            <th class="form_header">{row}</th>
            <td class="form_data"><RenderData key={props.keyT + row} keyT={props.keyT + "." + row} data={items[row]} handle={props.handle} keys={keys}/></td>
        </tr>
    })
}

const RenderRow = (props) => {
    console.log("render row")
    console.log("-> " + props.keyT)
    return props.keys.map((keyT, index) => {
        if (keyT.startsWith("_")) {
            return
        }
        console.log("render row: " + keyT)
        return <RenderData isRemovable="true" data={props.data[keyT]} key={props.keyT + "." + index} keyT={props.keyT + "." + index} handle={props.handle} keys={props.handle.getKeys(props.data[keyT])}/>
    })
}

const RenderData = (props) => {
    console.log("trying to render data: " + props.keyT)
    if (props.data === arrayConstructor || (props.data !== undefined && props.data !== null && props.data.constructor === Array)) {
        console.log("gently down the stream")
        var items = props.data
        return (<td class="form_data">
            <table>
                <thead>
                <tr>
                    <th>
                        <Button onClick={() => props.handle.handleAdd(props.keyT, props.data[0])}>
                            <AddBoxIcon/>
                            ADD
                        </Button>
                    </th>
                </tr>
                </thead>
                <tbody>

                <RenderRow key={props.keyT} keyT={props.keyT} data={props.data} handle={props.handle} keys={props.handle.getKeys(props.data)}/>
                </tbody>
            </table>
        </td>)
    } else if (props.data === objectConstructor || (props.data !== undefined && props.data !== null && props.data.constructor === Object)) {
        console.log("object time")
        console.log("-> " + Object.keys(props))
        if (props.isRemovable == "true") {
            return (<td class="form_data" id={props.keyT}>
                <table>
                    <thead>
                    </thead>
                    <tbody>
                    <RenderRows key={props.keyT} keyT={props.keyT} data={props.data} handle={props.handle} keys={props.handle.getKeys(props.data)}/>
                    <tr>
                        <td>
                            <Button onClick={() => props.handle.handleRemove(props.keyT, props.data[0])}>
                                <RemoveIcon/>
                                Remove
                            </Button>
                        </td>
                    </tr>
                    </tbody>
                </table>
            </td>)
        } else {
            return (<td class="form_data" id={props.keyT}>
                <table>
                    <tbody>
                    <RenderRows key={props.keyT} keyT={props.keyT} data={props.data} handle={props.handle} keys={props.handle.getKeys(props.data)}/>
                    </tbody>
                </table>
            </td>)
        }
    } else {
        console.log("just a string..." + props.keyT)
        console.log("->" + props.data)
        if (props.keyT.includes("._")) {
            return <td key={props.keyT}>{props.data}</td>
        } else {
            return <td key={props.keyT} class="form_data"><input type="text" name={props.keyT} value={props.data || ""} onChange={e => props.handle.handleEdit(props.keyT, e)}/></td>
        }
    }
}