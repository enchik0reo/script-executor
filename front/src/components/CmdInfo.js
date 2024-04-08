import React from 'react';

export default class CmdInfo extends React.Component {
    render() {

        return (
            <div className="cmd">
                <p className="cmd-id">Command id: {this.props.cmd.id}</p>
                <p className="cmd-id">Short description: {this.props.cmd.command_name}</p>
                <p className="cmd-id">Created at: {this.props.cmd.created_at}</p>
                <p className="cmd-id">Status: {this.props.cmd.is_working ? "works" : "finished"}</p>
                <div className="cmd-wrapper">
                    Output history:
                    {this.props.cmd.output ?
                        this.props.cmd.output.map((row) => (
                            <tr>
                                <td>{row}</td>
                            </tr>
                        ))
                        :
                        <div></div>
                    }
                </div>
            </div>
        )
    }
}