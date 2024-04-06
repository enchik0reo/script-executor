import React from 'react';
import Creator2 from '../components/Creator2';
import '../css/create.css';

export default class Create2 extends React.Component {
    render() {
        return (
            <>
                <h3 className="special-h3"><div className="terminal-prompt">Run your bash-script on remote server</div></h3>
                <Creator2 />
            </>
        )
    }
}