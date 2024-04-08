const LimitValidation = (values) => {

    let errors = {}

    const isNumber = !isNaN(values.limit)

    if (!values.limit) {
        errors.limit = "Limit is required"
    } else if (values.limit > 1000) {
        errors.limit = "Limit can be max 1000"
    } else if (values.limit < 5) {
        errors.limit = "Limit can be min 5"
    } else if (!isNumber) {
        errors.limit = "Limit can be number only"
    }

    

    return errors
}

export default LimitValidation