import React from 'react';
import Stp from '../components/Stp';
import '../css/stop.css';

export default class Stop extends React.Component {
    render() {
        return (
            <>
                <h3 className="special-h3"><div className="terminal-prompt">Stop one command by id</div></h3>
                <Stp />
            </>
        )
    }
}