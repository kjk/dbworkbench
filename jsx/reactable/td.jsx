import React from 'react';
import { stringable, isReactComponent } from './utils.jsx';
import { isUnsafe } from './unsafe.jsx';

export class Td extends React.Component {
    constructor(props, context) {
        super(props, context);
        this.state = {
            isEditable: false,
            value: 0,
        };
    }

    componentWillReceiveProps(nextProps) {
        this.setState({isEditable: nextProps.isEditable});
    }

    handleOnFocus() {
        console.log("handleOnFocus")
        // TODO: somehow move the cursor to end
    }

    handleClick(e){
        if (typeof this.props.onClick === 'function') {
            return this.props.onClick(e, this);
        }
    }

    handleKeyDown(e) {
        var ENTER = 13;
        if( e.keyCode == ENTER ) {
            console.log("Enter pressed", this)
            this.setState({isEditable: false});
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

        if (this.state.isEditable) {
            // console.log("Editable Cell", this.props, tdProps)
            return (
                <td {...tdProps}>
                    <textarea
                        id="editable"
                        autoFocus
                        value={this.props.children}
                        onFocus={this.handleOnFocus.bind(this)}
                        onKeyDown={this.handleKeyDown.bind(this)}
                        onChange={this.props.onEdit}>
                    </textarea>
                </td>
            );
        }

        return <td {...tdProps}></td>;
    }
};
