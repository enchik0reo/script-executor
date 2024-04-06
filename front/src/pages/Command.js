import React from 'react';
import Cmd from '../components/Cmd';
import '../css/cmd.css';

export default class Command extends React.Component {
    render() {
        return (
            <>
                <h3 className="special-h3"><div className="terminal-prompt">Check stdout of one command by id</div></h3>
                <Cmd />
            </>
        )
    }
}