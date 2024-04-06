const ScriptValidation = (values) => {

    let errors = {}

    if (!values.script) {
        errors.script = "Script is required"
    }

    return errors
}

export default ScriptValidation