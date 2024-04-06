import React, { useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import axios from 'axios';
import 'react-toastify/dist/ReactToastify.css';
import idValidation from './IDValidation';

const baseurl = "/stop"

const Stp = () => {

    const [values, setValues] = useState({
        id: 0,
    })

    const [errors, setErrors] = useState({})
    const [dataIsCorrect, setDataIsCorrect] = useState(false)
    const [currentAnsver, setCurrentAnsver] = useState("")

    const handleChange = (event) => {
        setValues({
            ...values,
            [event.target.name]: event.target.value,
        })
    }

    const handleFormSubmit = (event) => {
        event.preventDefault()
        setErrors(idValidation(values))
        setDataIsCorrect(true)
    }

    useEffect(() => {
        if (Object.keys(errors).length === 0 && dataIsCorrect) {

            const jsonData = {
                id: values.id
            }

            axios.put(baseurl, jsonData, {}).then((r) => {
                if (r.data.status === 202) {
                    setCurrentAnsver(r.data)
                    setDataIsCorrect(false)
                } else if (r.data.status === 304) {
                    setCurrentAnsver(r.data)
                    setDataIsCorrect(false)
                } else if (r.data.status === 500) {
                    toast.error(r.data.body.error)
                    setDataIsCorrect(false)
                } else if (r.data.status === 400) {
                    toast.warn(r.data.body.error)
                    setDataIsCorrect(false)
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
            <div className="cmd-list-wrapper">
                <form className="form-wrapper">
                    <div className="limit">
                        <label className="label-id">Interrupt command execution with id</label>
                        <input
                            className="input-id"
                            type="id"
                            name="id"
                            value={values.id}
                            onChange={handleChange}
                        />
                        <button className="show-id" onClick={handleFormSubmit}>
                            ok
                        </button>
                        {errors.id && <p className="error-limit">{errors.id}</p>}
                    </div>

                    <>
                        {Object.keys(currentAnsver).length > 0 ?
                            currentAnsver.status === 202 ?
                                <div>
                                    <p className="stop-info">The command stopped sucesfully.</p>
                                </div>
                                :
                                <div>
                                    <p className="stop-info">The command with current id isn't working.</p>
                                </div>
                            :
                            <div></div>
                        }
                    </>
                </form>
            </div>
        </div>
    )
}

export default Stp