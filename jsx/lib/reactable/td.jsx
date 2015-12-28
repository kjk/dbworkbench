import React from 'react';
import { isReactComponent } from './lib/is_react_component.jsx';
import { stringable } from './lib/stringable.jsx';
import { isUnsafe } from './unsafe.jsx';

export class Td extends React.Component {
    handleClick(e){
        if (typeof this.props.onClick === 'function') {
            return this.props.onClick(e, this);
        }
    }

    render() {
        var tdProps = {
            className: this.props.className,
            onClick: this.handleClick.bind(this),
            style: this.props.style,
        };

        // Attach any properties on the column to this Td object to allow things like custom event handlers
        if (typeof(this.props.column) === 'object') {
            for (var key in this.props.column) {
                if (key !== 'key' && key !== 'name') {
                    tdProps[key] = this.props.column[key];
                }
            }
        }

        var data = this.props.data;

        if (typeof(this.props.children) !== 'undefined') {
            if (isReactComponent(this.props.children)) {
                data = this.props.children;
            } else if (
                typeof(this.props.data) === 'undefined' &&
                    stringable(this.props.children)
            ) {
                data = this.props.children.toString();
            }

            if (isUnsafe(this.props.children)) {
                tdProps.dangerouslySetInnerHTML = { __html: this.props.children.toString() };
            } else {
                tdProps.children = data;
            }
        }

        if (this.props.isEditable) {
            // console.log("Editable Cell", this.props, tdProps)
            return <td {...tdProps}><textarea value={this.props.children} onChange={this.props.onEdit}></textarea></td>;
        }

        return <td {...tdProps}></td>;
    }
};
