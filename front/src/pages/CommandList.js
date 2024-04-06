import React from 'react';
import Commands from '../components/Commands';
import '../css/list.css';

export default class Create extends React.Component {
    render() {
        return (
            <>
                <h3 className="special-h3"><div className="terminal-prompt">Check latest commands</div></h3>
                <Commands />
            </>
        )
    }
}