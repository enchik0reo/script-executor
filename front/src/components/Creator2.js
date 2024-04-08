import React, { useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import axios from 'axios';
import 'react-toastify/dist/ReactToastify.css';

const baseurl = "/create/upload"

const Creator2 = () => {

    const [resetValues, setResetValues] = useState(false)
    const [ansver, setAnsver] = useState(false)
    const [id, setId] = useState(0)
    const [selectedFile, setSelectedFile] = useState(null)
    const [fileName, setFileName] = useState("")
    const [correctFile, setCorrectFile] = useState(null)
    const [chooseFile, setChooseFile] = useState("select")

    const handleFileChange = (event) => {
        event.preventDefault()
        if (event.target.files[0]) {
            setSelectedFile(event.target.files[0])
            setFileName(event.target.files[0].name)
            setChooseFile(event.target.files[0].name)
        }
    }

    const handleFileUpload = (event) => {
        event.preventDefault()
        if (selectedFile) {
            if (selectedFile.type === "application/x-shellscript") {
                setCorrectFile(selectedFile)
            } else {
                toast.warn("Required format: .sh")
            }
        } else {
            toast.warn("File not selected")
        }
    }

    const handleResetValues = () => {
        setResetValues(true)
    }

    useEffect(() => {
        if (resetValues) {
            setSelectedFile(null)
            setFileName(null)
            setCorrectFile(null)
            setChooseFile("select")
            setResetValues(false)
        }
    }, [resetValues])

    useEffect(() => {
        if (correctFile) {

            const formData = new FormData()
            formData.append("file", selectedFile)

            const config = {
                headers: {
                    "Content-Type": "multipart/form-data",
                },
            }

            axios.post(baseurl, formData, config).then((r) => {
                if (r.data.status === 500) {
                    toast.error(r.data.body.error)
                    handleResetValues()
                } else if (r.data.status === 400) {
                    toast.warn(r.data.body.error)
                    handleResetValues()
                } else {
                    setId(r.data.body.command_id)
                    setAnsver(true)
                    handleResetValues()
                }
            })
                .catch((error) => {
                    if (error) {
                        toast.error("Internal server error. Please, try later.")
                        console.error('Internal server error:', error)
                        handleResetValues()
                    }
                })
        }
    }, [correctFile, selectedFile])

    return (
        <div>
            <div className="create-wrapper2">
                <nav className="terminal-menu">
                    <ul>
                        <li>
                            <a className="menu-item" href="/" aria-current="page">write</a>
                        </li>
                        <li>
                            <a className="menu-item active" href="/upload">upload</a>
                        </li>
                    </ul>
                </nav>

                <form>

                    <div>
                        <label className="label2">Choose your file in format .sh: </label>
                        {!fileName ?
                            <label htmlFor="uploadFile" className="custom-file-upload">{chooseFile}</label>
                            :
                            <label htmlFor="uploadFile" className="custom-file-uploaded">{chooseFile}</label>
                        }

                        <input id="uploadFile" type="file" onChange={handleFileChange} style={{ display: "none" }} />
                    </div>

                    {ansver && <p className="ansver_id2">We got it! You can see details by id: {id}</p>}

                    <div>
                        {fileName ?
                            <button className="submit2" onClick={handleFileUpload}>
                                Execute
                            </button>
                            :
                            <button className="no-submit2" onClick={handleFileUpload}>
                                Execute
                            </button>
                        }
                    </div>
                </form>
            </div>
        </div>
    )
}

export default Creator2