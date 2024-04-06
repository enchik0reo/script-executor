const LimitValidation = (values) => {

    let errors = {}

    const isNumber = !isNaN(values.id)

    if (!values.id) {
        errors.id = "ID is required"
    } else if (!isNumber) {
        errors.id = "ID can be number only"
    } else if (values.id < 1) {
        errors.id = "ID can't be less than 1"
    }

    return errors
}

export default LimitValidation