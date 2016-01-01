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
    this.setState({
      isEditable: nextProps.isEditable
    });
  }

  handleOnFocus() {
    console.log('handleOnFocus');
    // TODO: somehow move the cursor to end
  }

  handleKeyDown(e) {
    console.log(e);
    var ENTER = 13;
    var SHIFT = 16;
    if (e.keyCode == ENTER && !e.shiftKey) {
      console.log('Enter pressed without shift', this);
      this.setState({
        isEditable: false
      });
    }
  }

  renderTextArea() {
    return (
      <textarea id="editable"
        autoFocus
        value={ this.props.children }
        onFocus={ this.handleOnFocus.bind(this) }
        onKeyDown={ this.handleKeyDown.bind(this) }
        onChange={ this.props.onEdit }>
      </textarea>
      );
  }

  render() {
    let tdProps = {
      className: this.props.className,
      style: this.props.style,
    };

    // Attach any properties on the column to this Td object to allow things like custom event handlers
    if (typeof (this.props.column) === 'object') {
      for (let key in this.props.column) {
        if (key !== 'key' && key !== 'name') {
          tdProps[key] = this.props.column[key];
        }
      }
    }

    let data = this.props.data;
    let pos = this.props.position;
    const rowCol = pos.rowId + '-' + pos.colId;
    if (typeof (this.props.children) !== 'undefined') {
      if (isReactComponent(this.props.children)) {
        data = this.props.children;
      } else if (
        typeof (this.props.data) === 'undefined' &&
        stringable(this.props.children)
      ) {
        data = this.props.children.toString();
      }

      if (isUnsafe(this.props.children)) {
        tdProps.dangerouslySetInnerHTML = {
          __html: this.props.children.toString()
        };
      } else {
        tdProps.children = data;
      }
      if (this.state.isEditable) {
        tdProps.children = this.renderTextArea();
      }
    }

    return (
      <td {...tdProps} data-custom-attribute={ rowCol }></td>
      );
  }
}
;
