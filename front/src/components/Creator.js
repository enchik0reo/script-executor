import React, { useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import axios from 'axios';
import scriptValidation from './ScriptValidation';
import 'react-toastify/dist/ReactToastify.css';

const baseurl = "/create"

const Creator = () => {

    const [values, setValues] = useState({
        script: "",
    })

    const [errors, setErrors] = useState({})
    const [dataIsCorrect, setDataIsCorrect] = useState(false)
    const [resetValues, setResetValues] = useState(false)
    const [ansver, setAnsver] = useState(false)
    const [id, setId] = useState(0)

    const handleChange = (event) => {
        setValues({
            ...values,
            [event.target.name]: event.target.value,
        })
    }

    const handleFormSubmit = (event) => {
        setAnsver(false)
        event.preventDefault()
        setErrors(scriptValidation(values))
        setDataIsCorrect(true)
    }

    const handleResetValues = () => {
        setResetValues(true)
    }

    useEffect(() => {
        if (resetValues) {
            setValues({
                script: "",
            })

            setResetValues(false)
        }
    }, [resetValues, values])

    useEffect(() => {
        if (Object.keys(errors).length === 0 && dataIsCorrect) {

            const jsonData = {
                script: values.script,
            }

            axios.post(baseurl, jsonData, {}).then((r) => {
                if (r.data.status === 500) {
                    toast.error(r.data.body.error)
                    setDataIsCorrect(false)
                } else if (r.data.status === 400) {
                    toast.warn(r.data.body.error)
                    setDataIsCorrect(false)
                } else {
                    setId(r.data.body.command_id)
                    setAnsver(true)
                    setDataIsCorrect(false)
                    handleResetValues()
                }
            })
                .catch((error) => {
                    if (error) {
                        toast.error("Internal server error. Please, try later.")
                        console.error('Internal server error:', error)
                        setDataIsCorrect(false)
                    }
                })
        }
    }, [errors, dataIsCorrect, values])

    return (
        <div>
            <div className="create-wrapper">
                <nav className="terminal-menu">
                    <ul>
                        <li>
                            <a className="menu-item active" href="/" aria-current="page">write</a>
                        </li>
                        <li>
                            <a className="menu-item" href="/upload">upload</a>
                        </li>
                    </ul>
                </nav>

                <form className="form-wrapper">
                    <div className="password">
                        <label className="label">Write script here</label>
                        <textarea className="script-text"
                            type="script"
                            name="script"
                            value={values.script}
                            onChange={handleChange} />
                        {errors.script && <p className="error">{errors.script}</p>}
                        {ansver && <p className="ansver_id">We got it! You can see details by id: {id}</p>}
                    </div>

                    <div>
                        <button className="submit" onClick={handleFormSubmit}>
                            Execute
                        </button>
                    </div>
                </form>
            </div>
        </div>
    )
}

export default Creator