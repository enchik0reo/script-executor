import React, { useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import axios from 'axios';
import 'react-toastify/dist/ReactToastify.css';
import limitValidation from './LimitValidation';

const baseurl = "/list"

const Commands = () => {

    const [values, setValues] = useState({
        limit: 5,
    })

    const [errors, setErrors] = useState({})
    const [dataIsCorrect, setDataIsCorrect] = useState(false)
    const [currentList, setCurrentList] = useState([])
    const [afterGet, setAfterGet] = useState(false)

    useEffect(() => {
        if (!afterGet) {
        const prevConfig = {
            params: {
                limit: 5,
            }
        }

        axios.get(baseurl, prevConfig).then((r) => {
            if (r.data.status === 500) {
                toast.error(r.data.body.error)
                setDataIsCorrect(false)
            } else if (r.data.status === 400) {
                toast.warn(r.data.body.error)
                setDataIsCorrect(false)
            } else {
                setCurrentList(r.data.body.commands)
                setAfterGet(true)
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
    }, [dataIsCorrect, afterGet])

    const handleChange = (event) => {
        setValues({
            ...values,
            [event.target.name]: event.target.value,
        })
    }

    const handleFormSubmit = (event) => {
        event.preventDefault()
        setErrors(limitValidation(values))
        setDataIsCorrect(true)
    }

    useEffect(() => {
        if (Object.keys(errors).length === 0 && dataIsCorrect) {

            const config = {
                params: {
                    limit: values.limit,
                }
            }

            axios.get(baseurl, config).then((r) => {
                if (r.data.status === 500) {
                    toast.error(r.data.body.error)
                    setDataIsCorrect(false)
                } else if (r.data.status === 400) {
                    toast.warn(r.data.body.error)
                    setDataIsCorrect(false)
                } else {
                    setCurrentList(r.data.body.commands)
                    setAfterGet(true)
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
            <div className="list-wrapper">
                <form className="form-wrapper">
                    <div className="limit">
                        <label className="label-limit">Show last</label>
                        <input
                            className="input-limit"
                            type="limit"
                            name="limit"
                            value={values.limit}
                            onChange={handleChange}
                        />

                        <label className="label-limit">commands</label>
                        <button className="show" onClick={handleFormSubmit}>
                            ok
                        </button>
                        {errors.limit && <p className="error-limit">{errors.limit}</p>}
                    </div>

                    <>
                        {afterGet ?
                            currentList && currentList.length > 0 ?
                                <table className="table">
                                    <thead>
                                        <tr>
                                            <th>id</th>
                                            <th>preview</th>
                                            <th>created at (UTC)</th>
                                            <th>status</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {currentList.map((row) => (
                                            <tr key={row.id}>
                                                <td>{row.id}</td>
                                                <td>{row.command_name}</td>
                                                <td>{row.created_at}</td>
                                                <td>{row.is_working ? "works" : "finished"}</td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                                :
                                <div>
                                    <p className="cmd-id">There are no working commands</p>
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

export default Commands