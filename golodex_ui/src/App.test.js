import React from 'react';
import {render, screen} from '@testing-library/react';
import App from './App'
import {getId, getToken, setToken} from './App';
import Rolodex from "./components/Rolodex/Rolodex"
import {mount, shallow} from 'enzyme'
import Enzyme from 'enzyme';
import Adapter from '@wojtekmaj/enzyme-adapter-react-17';//'enzyme-adapter-react-16';
Enzyme.configure({ adapter: new Adapter() });

beforeEach(() => {
    fetch.resetMocks();
});

test('renders welcome', () => {
    render(<App/>);
    const linkElement = screen.getByText(/Welcome to Golodex!/i);
    expect(linkElement).toBeInTheDocument();
});

test('react version', () => {
    console.log(React.version);
});

test('set token', () => {
    setToken({
        "token": "test",
        "id": "tester"
    })
    expect(getToken()).toBe("test")
    expect(getId()).toBe("tester")
});

test('renders deck', async () => {
    fetch.mockResponseOnce(JSON.stringify(
        {
            "cards": [
                {
                    "_id": "0419055f7fe9956fc3bc3fb00f0005bd",
                    "_rev": "4-35abba1ef3d411d0fe20800741406e19",
                    "name": {
                        "first": "Catniss",
                        "middle": "Everdeen",
                        "last": "Angeli"
                    },
                    "phone_numbers": [
                        {
                            "number": "1-234-555-1811",
                            "type": "home"
                        },
                        {
                            "number": "9-876-555-1811",
                            "type": "work"
                        }
                    ],
                    "addresses": [
                        {
                            "street_1": "123 Sesame Street",
                            "street_2": "",
                            "city": "Alpharetta",
                            "state": "GA",
                            "zip_code": "30009",
                            "type": "home"
                        },
                        {
                            "street_1": "415 S Broad St",
                            "street_2": "",
                            "city": "Alpharetta",
                            "state": "GA",
                            "zip_code": "30009",
                            "type": "business"
                        },
                        {
                            "street_1": "1130 Mayfield Manor Drive",
                            "street_2": "",
                            "city": "Alpharetta",
                            "state": "GA",
                            "zip_code": "30009",
                            "type": "vacation"
                        }
                    ]
                }
            ]
        }))
    const allOver = () => new Promise((resolve) => setImmediate(resolve));
    const setPage = jest.fn();
    const setCard = jest.fn();
    const wrapper = mount(<Rolodex key="the_rolodex" getId={getId} getToken={getToken} setPage={setPage} setCard={setCard}/>);
    await allOver()
    expect(wrapper.instance().state.isReloaded).toBe("");
    expect(wrapper.text().includes('123 Sesame Street')).toBe(true);
    expect(wrapper.text().includes('vacation')).toBe(true);
    expect(wrapper.text().includes('tr')).toBe(true);

    // console.log(wrapper.text().text)
    const the_doc = wrapper.render();
    const table = the_doc.find('table');
    const rows = the_doc.find('tr');
    expect(table).toHaveLength(4);
    expect(rows).toHaveLength(8);

});