import React, { useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import axios from 'axios';
import 'react-toastify/dist/ReactToastify.css';
import idValidation from './IDValidation';
import CmdInfo from './CmdInfo';

const baseurl = "/cmd"

const Cmd = () => {

    const [values, setValues] = useState({
        id: 0,
    })

    const [errors, setErrors] = useState({})
    const [dataIsCorrect, setDataIsCorrect] = useState(false)
    const [currentCmd, setCurrentCmd] = useState({})

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

            const config = {
                params: {
                    id: values.id,
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
                    setCurrentCmd(r.data.body)
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
                        <label className="label-id">Command id</label>
                        <input
                            className="input-id"
                            type="id"
                            name="id"
                            value={values.id}
                            onChange={handleChange}
                        />
                        <button className="show-id" onClick={handleFormSubmit}>
                            show
                        </button>
                        {errors.id && <p className="error-limit">{errors.id}</p>}
                    </div>

                    <>
                        {Object.keys(currentCmd).length > 0 ?
                            currentCmd.command.id > 0 ?
                                <CmdInfo cmd={currentCmd.command} />
                                :
                                <div>
                                    <p className="cmd-id">There is no stdout history for this id</p>
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

export default Cmd