
export function stringable(thing) {
    return thing !== null &&
        typeof(thing) !== 'undefined' &&
        typeof(thing.toString === 'function');
}

export function extractDataFrom(key, column) {
    var value;
    if (
        typeof(key) !== 'undefined' &&
            key !== null &&
                key.__reactableMeta === true
    ) {
        value = key.data[column];
    } else {
        value = key[column];
    }

    if (
        typeof(value) !== 'undefined' &&
            value !== null &&
                value.__reactableMeta === true
    ) {
        value = (typeof(value.props.value) !== 'undefined' && value.props.value !== null) ?
            value.props.value : value.value;
    }

    return (stringable(value) ? value : '');
}

const internalProps = {
    column: true,
    columns: true,
    sortable: true,
    filterable: true,
    sortBy: true,
    defaultSort: true,
    itemsPerPage: true,
    childNode: true,
    data: true,
    children: true
};

export function filterPropsFrom(baseProps) {
    baseProps = baseProps || {};
    var props = {};
    for (var key in baseProps) {
        if (!(key in internalProps)) {
            props[key] = baseProps[key];
        }
    }

    return props;
}

export function toArray(obj) {
    var ret = [];
    for (var attr in obj) {
        ret[attr] = obj;
    }

    return ret;
}

// this is a bit hacky - it'd be nice if React exposed an API for this
export function isReactComponent(thing) {
    return thing !== null && typeof(thing) === 'object' && typeof(thing.props) !== 'undefined';
}
